"""Canonical, conservative size of an image's uncompressed Docker-save layers."""

import json
import tarfile
from pathlib import Path


def uncompressed_layer_bytes(archive: Path) -> int:
    """Return exact layer tar member bytes from a single-image docker-save archive.

    This is stable across Docker Engine versions and conservatively accounts for
    the space needed to materialize every layer. Docker's inspect Size can vary
    between engines even when the image config and all layer contents match.
    """
    with tarfile.open(archive, 'r:gz') as saved:
        manifest_stream = saved.extractfile('manifest.json')
        if manifest_stream is None:
            raise RuntimeError('Docker archive is missing its manifest')
        entries = json.load(manifest_stream)
        if len(entries) != 1:
            raise RuntimeError('Expected one image in Docker archive')
        layers = entries[0].get('Layers')
        if not isinstance(layers, list) or not layers or len(layers) > 256:
            raise RuntimeError('Docker archive has an invalid layer list')
        total = 0
        for name in layers:
            if not isinstance(name, str) or not name or name.startswith('/') or '..' in Path(name).parts:
                raise RuntimeError('Docker archive has an unsafe layer path')
            try:
                member = saved.getmember(name)
            except KeyError as error:
                raise RuntimeError('Docker archive layer is missing') from error
            if not member.isfile() or member.size <= 0:
                raise RuntimeError('Docker archive layer is empty or invalid')
            total += member.size
        return total
