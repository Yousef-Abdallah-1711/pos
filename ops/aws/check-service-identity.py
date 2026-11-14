#!/usr/bin/env python3
"""Exercise the service's real Linux identity drop without touching deployment state."""
import os
from pathlib import Path
import pwd
import subprocess
import uuid


def main():
    if os.geteuid()==0:
        raise SystemExit('Run as the unprivileged test user with noninteractive sudo')
    user=pwd.getpwuid(os.getuid()).pw_name
    keys={'NoNewPrivileges','PrivateTmp','ProtectSystem','ProtectHome','AmbientCapabilities'}
    unit=Path(__file__).with_name('servio-deploy.service').read_text()
    properties=[line for line in unit.splitlines() if '=' in line and line.split('=',1)[0] in keys]
    child="""import os
assert os.geteuid()!=0, 'acceptance must run without root'
status=dict(line.split(':',1) for line in open('/proc/self/status') if ':' in line)
assert status['NoNewPrivs'].strip()=='1', 'NNP must survive the identity drop'
assert int(status['CapEff'],16)==0 and int(status['CapAmb'],16)==0, 'identity capabilities leaked'
try:
    os.setuid(0)
except PermissionError:
    print('PASS unprivileged identity, zero capabilities, NNP and denied root regain')
else:
    raise SystemExit('Unprivileged acceptance regained root')
"""
    parent="import subprocess,sys; subprocess.run(['runuser','-u',sys.argv[1],'--','/usr/bin/python3','-c',sys.argv[2]],check=True)"
    command=['sudo','-n','systemd-run','--unit=servio-identity-check-'+uuid.uuid4().hex[:12],
             '--wait','--collect','--pipe','-pUser=root','-pType=oneshot','-pTimeoutStartSec=30']
    command.extend('-p'+item for item in properties)
    subprocess.run(command+['/usr/bin/python3','-c',parent,user,child],check=True,timeout=45)


if __name__=='__main__':main()
