import assert from "node:assert/strict";
import { mkdtemp, mkdir, rm, symlink } from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import test from "node:test";
import { resolveOutsideRepository } from "./path-boundary.mjs";

test("requires fixture and output paths to resolve outside the repository", async (t) => {
  const temporaryRoot = await mkdtemp(path.join(os.tmpdir(), "servio-financial-path-boundary-"));
  t.after(() => rm(temporaryRoot, { recursive: true, force: true }));
  const repositoryRoot = path.join(temporaryRoot, "repo");
  await mkdir(repositoryRoot);

  await assert.rejects(
    resolveOutsideRepository(path.join(repositoryRoot, "captures", "page.png"), repositoryRoot, "output"),
    /must resolve outside the repository/
  );
  assert.equal(
    await resolveOutsideRepository(path.join(temporaryRoot, "fixtures", "synthetic.json"), repositoryRoot, "fixture"),
    path.join(temporaryRoot, "fixtures", "synthetic.json")
  );
});

test("rejects a path that escapes into the repository through a symlink", async (t) => {
  const temporaryRoot = await mkdtemp(path.join(os.tmpdir(), "servio-financial-path-symlink-"));
  t.after(() => rm(temporaryRoot, { recursive: true, force: true }));
  const repositoryRoot = path.join(temporaryRoot, "repo");
  const outsideRoot = path.join(temporaryRoot, "outside");
  const link = path.join(outsideRoot, "repo-link");
  await Promise.all([mkdir(repositoryRoot), mkdir(outsideRoot)]);

  try {
    await symlink(repositoryRoot, link, process.platform === "win32" ? "junction" : "dir");
  } catch (error) {
    if (["EPERM", "EACCES", "ENOTSUP"].includes(error.code)) {
      t.skip(`Symlink creation is unavailable: ${error.code}`);
      return;
    }
    throw error;
  }

  await assert.rejects(
    resolveOutsideRepository(path.join(link, "captures", "page.png"), repositoryRoot, "output"),
    /must resolve outside the repository/
  );
});
