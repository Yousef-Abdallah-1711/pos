const vulnerableNextIntlVersion = '4.15.1'
const expectedWatcherRange = '^2.4.1'

function readPackage(pkg) {
  if (pkg.name !== 'next-intl' || pkg.version !== vulnerableNextIntlVersion) {
    return pkg
  }

  const watcherRange = pkg.dependencies?.['@parcel/watcher']
  if (watcherRange === undefined) {
    return pkg
  }
  if (watcherRange !== expectedWatcherRange) {
    throw new Error(
      `Review next-intl@${pkg.version}: expected unused @parcel/watcher ${expectedWatcherRange}, found ${watcherRange ?? 'no dependency'}`,
    )
  }

  // next-intl@4.15.1 ships no imports of @parcel/watcher. Removing this
  // unused dependency also removes the vulnerable production micromatch path.
  delete pkg.dependencies['@parcel/watcher']
  return pkg
}

module.exports = {
  hooks: { readPackage },
}
