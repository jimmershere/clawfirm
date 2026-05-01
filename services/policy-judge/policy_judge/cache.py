"""Tiny in-process TTL cache. Replace with Redis at SMB+ tiers."""
import time


class TTLCache:
    def __init__(self, ttl: int):
        self.ttl = ttl
        self._store: dict[str, tuple[float, dict]] = {}

    def get(self, key: str) -> dict | None:
        if key not in self._store:
            return None
        ts, val = self._store[key]
        if time.time() - ts > self.ttl:
            del self._store[key]
            return None
        return val

    def set(self, key: str, value: dict) -> None:
        self._store[key] = (time.time(), value)
