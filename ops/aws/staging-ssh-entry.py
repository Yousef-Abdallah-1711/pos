#!/usr/bin/python3
"""Forced-command SSH entrypoint; no shell, arbitrary command or forwarding access."""
import os
import re
import sys

match = re.fullmatch(r'servio-deploy ([0-9a-f]{40})', os.environ.get('SSH_ORIGINAL_COMMAND', ''))
if not match:
    raise SystemExit('Only an immutable SERVIO staging deployment is permitted')
os.execv('/usr/bin/sudo', ['sudo', '-n', '/usr/local/sbin/servio-staging-release', match.group(1)])
