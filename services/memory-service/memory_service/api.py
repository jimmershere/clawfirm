"""FastAPI front + gRPC server for the Memory Service.

MVP scaffold: only the REST surface is wired. gRPC server lands in
Weeks 9-10 of the 90-day roadmap.
"""
from __future__ import annotations

import os

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

from .events import Event, EventLog


app = FastAPI(title="ClawFirm Memory Service", version="0.1.0")
log = EventLog(os.environ.get("MEMORY_EVENT_LOG", "/var/lib/memory/events"))


class AppendRequest(BaseModel):
    user: str
    scope: str
    payload: dict


class RecallRequest(BaseModel):
    user: str
    query: str
    top_k: int = 5


@app.get("/healthz")
def healthz():
    return {"status": "ok"}


@app.post("/v1/append")
def append(req: AppendRequest):
    if req.scope not in {"public", "tenant", "user", "session", "secret"}:
        raise HTTPException(400, f"invalid scope: {req.scope}")
    event = log.append(Event(kind="append", user=req.user, scope=req.scope, payload=req.payload))
    return {"hash": event.hash, "prev_hash": event.prev_hash, "ts": event.ts}


@app.post("/v1/recall")
def recall(req: RecallRequest):
    # MVP scaffold returns an empty result set. The actual retrieval lands
    # once the per-store backends in stores/* are wired up.
    return {"hits": []}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
