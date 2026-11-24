import { realpath } from "node:fs/promises";
import path from "node:path";

export async function resolveOutsideRepository(candidate, repositoryRoot, label) {
  const realRepositoryRoot = await realpath(repositoryRoot);
  const requestedPath = path.resolve(candidate);
  const missingSegments = [];
  let existingAncestor = requestedPath;
  let resolvedAncestor;

  while (!resolvedAncestor) {
    try {
      resolvedAncestor = await realpath(existingAncestor);
    } catch (error) {
      if (error.code !== "ENOENT" && error.code !== "ENOTDIR") throw error;
      const parent = path.dirname(existingAncestor);
      if (parent === existingAncestor) throw error;
      missingSegments.unshift(path.basename(existingAncestor));
      existingAncestor = parent;
    }
  }

  const resolvedPath = path.resolve(resolvedAncestor, ...missingSegments);
  const relativePath = path.relative(realRepositoryRoot, resolvedPath);
  const isInsideRepository = relativePath === "" || (
    !path.isAbsolute(relativePath) && relativePath !== ".." && !relativePath.startsWith(`..${path.sep}`)
  );
  if (isInsideRepository) throw new Error(`${label} must resolve outside the repository.`);
  return resolvedPath;
}
