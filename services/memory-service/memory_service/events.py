"""Append-only event log with hash chain."""
from __future__ import annotations

import hashlib
import json
import time
from dataclasses import asdict, dataclass, field
from pathlib import Path
from typing import Any


@dataclass
class Event:
    kind: str               # "append" | "tombstone" | "snapshot"
    user: str
    scope: str              # "public" | "tenant" | "user" | "session" | "secret"
    payload: dict[str, Any]
    prev_hash: str = ""
    hash: str = ""
    ts: float = field(default_factory=time.time)


class EventLog:
    """File-backed append-only log. One JSON line per event."""

    def __init__(self, path: str | Path, *, hash_chain: bool = True):
        self.path = Path(path)
        self.path.mkdir(parents=True, exist_ok=True)
        self.file = self.path / "events.jsonl"
        self.hash_chain = hash_chain
        self._prev_hash = self._load_prev_hash()

    def append(self, event: Event) -> Event:
        if self.hash_chain:
            event.prev_hash = self._prev_hash
            event.hash = self._chain(event)
            self._prev_hash = event.hash
        with self.file.open("a", encoding="utf-8") as f:
            f.write(json.dumps(asdict(event), separators=(",", ":")) + "\n")
        return event

    def _chain(self, event: Event) -> str:
        canonical = json.dumps(
            {
                "prev_hash": event.prev_hash,
                "kind": event.kind,
                "user": event.user,
                "scope": event.scope,
                "payload": event.payload,
            },
            sort_keys=True,
            separators=(",", ":"),
        )
        return hashlib.sha256(canonical.encode("utf-8")).hexdigest()

    def _load_prev_hash(self) -> str:
        if not self.file.exists():
            return ""
        last = ""
        with self.file.open("r", encoding="utf-8") as f:
            for line in f:
                last = line
        if not last:
            return ""
        return json.loads(last).get("hash", "")
