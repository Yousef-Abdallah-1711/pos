#!/usr/bin/env node

const { execSync } = require('node:child_process');
const fs = require('node:fs');
const path = require('node:path');

if (process.env.POS_SKIP_DEV_PREFLIGHT === '1') {
  console.log('[frontend:predev] Skipped (POS_SKIP_DEV_PREFLIGHT=1).');
  process.exit(0);
}

const isWindows = process.platform === 'win32';
const rootDir = path.resolve(__dirname, '..');

const appConfigs = [
  { name: 'landing', port: 3001, trace: path.join(rootDir, 'apps', 'landing', '.next', 'trace') },
  { name: 'dashboard', port: 3002, trace: path.join(rootDir, 'apps', 'dashboard', '.next', 'trace') },
  {
    name: 'restaurant-website',
    port: 3003,
    trace: path.join(rootDir, 'apps', 'restaurant-website', '.next', 'trace'),
  },
  { name: 'super-admin', port: 3004, trace: path.join(rootDir, 'apps', 'super-admin', '.next', 'trace') },
];

function getListeningPids(port) {
  try {
    if (isWindows) {
      const output = execSync(`netstat -ano -p tcp | findstr :${port}`, {
        encoding: 'utf8',
        stdio: ['ignore', 'pipe', 'ignore'],
      });

      return Array.from(
        new Set(
          output
            .split(/\r?\n/)
            .map((line) => line.trim())
            .filter((line) => line && line.includes('LISTENING'))
            .map((line) => line.split(/\s+/).pop())
            .filter((pid) => pid && Number(pid) > 0)
        )
      );
    }

    const output = execSync(`lsof -ti tcp:${port}`, {
      encoding: 'utf8',
      stdio: ['ignore', 'pipe', 'ignore'],
    });

    return Array.from(
      new Set(
        output
          .split(/\r?\n/)
          .map((line) => line.trim())
          .filter((pid) => pid && Number(pid) > 0)
      )
    );
  } catch {
    return [];
  }
}

function killPid(pid) {
  try {
    if (isWindows) {
      execSync(`taskkill /PID ${pid} /F /T`, {
        stdio: ['ignore', 'ignore', 'ignore'],
      });
      return true;
    }

    process.kill(Number(pid), 'SIGKILL');
    return true;
  } catch {
    return false;
  }
}

for (const app of appConfigs) {
  const pids = getListeningPids(app.port).filter((pid) => Number(pid) !== process.pid);
  if (pids.length === 0) {
    console.log(`[frontend:predev] Port ${app.port} (${app.name}) is free.`);
  } else {
    console.log(
      `[frontend:predev] Freeing port ${app.port} (${app.name}). PID(s): ${pids.join(', ')}`
    );
    for (const pid of pids) {
      const ok = killPid(pid);
      if (ok) {
        console.log(`[frontend:predev] Killed PID ${pid}`);
      } else {
        console.warn(`[frontend:predev] Failed to kill PID ${pid}`);
      }
    }
  }

  try {
    fs.rmSync(app.trace, { force: true });
    console.log(`[frontend:predev] Cleared trace: ${app.trace}`);
  } catch (error) {
    console.warn(`[frontend:predev] Could not clear trace (${app.trace}): ${error.message}`);
  }
}
