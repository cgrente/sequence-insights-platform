"""
A tiny Python client for the Sequence Insights API.

This is intentionally lightweight (requests + dataclasses-style dicts) to keep it easy to copy/paste.
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, List, Optional
import requests

@dataclass
class SequenceIngestResponse:
    base_url: str = ""
    api_key: Optional[str] = None
    timeout_seconds: int = 30

    def _headers(self) -> Dict[str, str]:
        headers = {"Content-Type": "application/json"}
        if self.api_key:
            headers["Authorization"] = f"Bearer {self.api_key}"
        return headers
    
    def health(self) -> Dict[str, Any]:
        response = requests.get(
            f"{self.base_url.rstrip('/')}/v1/health",
            headers=self._headers(),
            timeout=self.timeout_seconds
        )
        response.raise_for_status()
        return response.json()
    
    def ingest_sequence(self, values: List[int]) -> Dict[str, Any]:
        payload = {"values": values}
        response = requests.post(
            f"{self.base_url.rstrip('/')}/v1/sequences/ingest",
            json=payload,
            headers=self._headers(),
            timeout=self.timeout_seconds
        )
        response.raise_for_status()
        return response.json()
    
    def get_sequence(self, sequence_id: str) -> Dict[str, Any]:
        response = requests.get(
            f"{self.base_url.rstrip('/')}/v1/sequences/{sequence_id}",
            headers=self._headers(),
            timeout=self.timeout_seconds
        )
        response.raise_for_status()
        return response.json()