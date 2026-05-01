"""FastAPI server fronting the local Qwen3.5-4B classifier via Ollama."""
from __future__ import annotations

import hashlib
import json
import os
import time
from pathlib import Path

import httpx
from fastapi import FastAPI
from pydantic import BaseModel

from .cache import TTLCache


MODEL = os.environ.get("JUDGE_MODEL", "qwen3.5-4b")
OLLAMA = os.environ.get("OLLAMA_HOST", "http://ollama:11434")
CACHE_TTL = int(os.environ.get("JUDGE_CACHE_TTL", "300"))

PROMPT_DIR = Path(__file__).parent / "prompts"

app = FastAPI(title="ClawFirm Policy Judge", version="0.1.0")
cache: TTLCache = TTLCache(ttl=CACHE_TTL)


class JudgeRequest(BaseModel):
    tool: str | None = None
    args: dict | None = None
    identity: str
    text: str | None = None
    rubric: str = "routing"   # "routing" | "sensitivity" | "complexity"


class JudgeReply(BaseModel):
    rubric: str
    cls: str
    confidence: float
    reason: str
    cached: bool = False


def _cache_key(req: JudgeRequest) -> str:
    canonical = json.dumps(req.model_dump(), sort_keys=True, separators=(",", ":"))
    return hashlib.sha256(canonical.encode("utf-8")).hexdigest()


def _load_rubric(name: str) -> str:
    return (PROMPT_DIR / f"{name}.txt").read_text(encoding="utf-8")


@app.get("/healthz")
def healthz():
    return {"status": "ok", "model": MODEL}


@app.post("/v1/judge", response_model=JudgeReply)
async def judge(req: JudgeRequest):
    key = _cache_key(req)
    if hit := cache.get(key):
        return JudgeReply(**hit, cached=True)

    rubric = _load_rubric(req.rubric)
    prompt = (
        f"{rubric}\n\n"
        f"INPUT:\ntool={req.tool}\nargs={json.dumps(req.args or {})}\n"
        f"identity={req.identity}\ntext={req.text or ''}\n\n"
        f"OUTPUT JSON ONLY:"
    )
    started = time.time()
    async with httpx.AsyncClient(timeout=20) as client:
        r = await client.post(
            f"{OLLAMA}/api/generate",
            json={"model": MODEL, "prompt": prompt, "stream": False, "format": "json"},
        )
        r.raise_for_status()
        raw = r.json().get("response", "{}")
    try:
        out = json.loads(raw)
        reply = JudgeReply(
            rubric=req.rubric,
            cls=str(out.get("class", "unknown")),
            confidence=float(out.get("confidence", 0.0)),
            reason=str(out.get("reason", ""))[:512],
        )
    except (json.JSONDecodeError, ValueError, TypeError):
        reply = JudgeReply(rubric=req.rubric, cls="unknown", confidence=0.0, reason="bad_json")

    cache.set(key, reply.model_dump(exclude={"cached"}))
    print(f"judge rubric={req.rubric} cls={reply.cls} conf={reply.confidence:.2f} ms={int(1000*(time.time()-started))}")
    return reply


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
