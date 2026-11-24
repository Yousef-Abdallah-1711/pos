'use strict'

module.exports.hooks = {
  readPackage(pkg) {
    if (pkg.name === 'next-intl' && pkg.version === '4.15.1') {
      if (pkg.dependencies) {
        delete pkg.dependencies['@parcel/watcher']
      }
    }
    return pkg
  },
}
