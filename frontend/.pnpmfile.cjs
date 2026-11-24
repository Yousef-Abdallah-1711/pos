'use strict';

const TARGET_PACKAGE = 'next-intl';
const TARGET_VERSION = '4.15.1';
const OMITTED_DEPENDENCY = '@parcel/watcher';

module.exports = {
  hooks: {
    readPackage(pkg) {
      if (pkg.name !== TARGET_PACKAGE || pkg.version !== TARGET_VERSION) {
        return pkg;
      }

      if (pkg.dependencies) {
        delete pkg.dependencies[OMITTED_DEPENDENCY];
      }

      return pkg;
    },
  },
};
