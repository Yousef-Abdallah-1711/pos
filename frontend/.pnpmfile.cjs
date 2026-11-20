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

  // Keep watcher available as a root development dependency for next-intl's
  // config plugin, but remove it from next-intl's production dependency graph
  // so the production audit excludes this build-only package.
  delete pkg.dependencies['@parcel/watcher']
  return pkg
}

module.exports = {
  hooks: { readPackage },
}
