"""Synchronize the landing Docker context with the canonical backend offer."""
from pathlib import Path
root = Path(__file__).resolve().parents[1]
source = root / 'backend/internal/commerce/offer.json'
target = root / 'frontend/apps/landing/src/content/commercial-offer.json'
target.write_bytes(source.read_bytes())
