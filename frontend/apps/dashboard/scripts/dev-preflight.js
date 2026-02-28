#!/usr/bin/env node

const { execSync } = require('node:child_process');
const fs = require('node:fs');
const path = require('node:path');

const DEFAULT_PORT = 3002;
const port = Number(process.env.PORT || DEFAULT_PORT);
const projectRoot = path.resolve(__dirname, '..');
const traceFilePath = path.join(projectRoot, '.next', 'trace');

const isWindows = process.platform === 'win32';

function getListeningPids(targetPort) {
  try {
    if (isWindows) {
      const output = execSync(`netstat -ano -p tcp | findstr :${targetPort}`, {
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

    const output = execSync(`lsof -ti tcp:${targetPort}`, {
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

function cleanupTraceFile() {
  try {
    fs.rmSync(traceFilePath, { force: true });
    console.log(`[dashboard:predev] Cleared stale trace file: ${traceFilePath}`);
  } catch (error) {
    console.warn(
      `[dashboard:predev] Could not clear trace file (${traceFilePath}): ${error.message}`
    );
  }
}

function main() {
  const pids = getListeningPids(port).filter((pid) => Number(pid) !== process.pid);

  if (pids.length === 0) {
    console.log(`[dashboard:predev] Port ${port} is free.`);
  } else {
    console.log(
      `[dashboard:predev] Freeing port ${port}. Stale listener PID(s): ${pids.join(', ')}`
    );

    for (const pid of pids) {
      const killed = killPid(pid);
      if (killed) {
        console.log(`[dashboard:predev] Killed PID ${pid}`);
      } else {
        console.warn(`[dashboard:predev] Failed to kill PID ${pid}`);
      }
    }
  }

  cleanupTraceFile();
}

main();
