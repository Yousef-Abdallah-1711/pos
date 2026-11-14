"""Root-reviewed migration compatibility and explicit mixed-image recovery state."""
import copy
import json
from pathlib import Path
import re

POLICY=Path('/etc/servio/migration-rollbacks.json')

def requires_compatible_backend(pending, previous):
    if previous.get('recovery_required'):return True
    if not POLICY.exists():return False
    info=POLICY.stat()
    if info.st_uid!=0 or info.st_mode&0o022:raise RuntimeError('Migration rollback policy must be root-owned and not writable by others')
    policy=json.loads(POLICY.read_text())
    if any(value not in ['compatible','keep_backend'] for value in policy.values()):raise RuntimeError('Unknown migration rollback policy')
    return any(policy.get(name)=='keep_backend' for name in pending or [])

def compatible_recovery(previous, attempt):
    state=copy.deepcopy(previous)
    manifest=attempt['manifest'];runtime=attempt['runtime']
    if attempt['backend_ready']:
        for field in ['images','image_ids']:
            state[field]=dict(state.get(field,{}),backend=manifest[field]['backend'])
        state['runtime_image_ids']=dict(state.get('runtime_image_ids',{}),backend=runtime['backend'])
        state['application_shas']=dict(state.get('application_shas',{}),backend=manifest['sha'])
    backend=state.get('runtime_image_ids',{}).get('backend')
    if not isinstance(backend,str) or not re.fullmatch('sha256:[0-9a-f]{64}',backend):raise RuntimeError('Missing compatible recovery backend identity')
    old=Path(previous['compose']);config=json.loads(old.read_text())
    for name in ['backend','migrate']:
        config['services'][name]['image']=backend
        config['services'][name]['pull_policy']='never'
    # Keep relative mounts anchored beside the previous Compose file.
    fallback=old.with_name('recovery-'+manifest['sha']+'.json')
    temporary=fallback.with_suffix('.tmp');temporary.write_text(json.dumps(config,indent=2)+'\n');temporary.chmod(0o600);temporary.replace(fallback)
    state.update(compose=str(fallback),recovery_required=True,attempted_sha=manifest['sha'])
    return state,fallback
