const path = require('path');
const relativeDir = path.resolve(__dirname, '..');
module.exports = {
  stylePaths: [
    path.resolve(relativeDir, 'src'),
    path.resolve(relativeDir, 'node_modules/patternfly'),
    path.resolve(relativeDir, 'node_modules/@patternfly/patternfly'),
    path.resolve(relativeDir, 'node_modules/@patternfly/react-styles/css'),
    path.resolve(relativeDir, 'node_modules/@patternfly/react-core/dist/styles/base.css'),
    path.resolve(relativeDir, 'node_modules/@patternfly/react-core/dist/esm/@patternfly/patternfly'),
    path.resolve(relativeDir, 'node_modules/@patternfly/react-core/node_modules/@patternfly/react-styles/css'),
    path.resolve(relativeDir, 'node_modules/@patternfly/react-table/node_modules/@patternfly/react-styles/css'),
    path.resolve(relativeDir, 'node_modules/@patternfly/react-inline-edit-extension/node_modules/@patternfly/react-styles/css')
  ]
};
