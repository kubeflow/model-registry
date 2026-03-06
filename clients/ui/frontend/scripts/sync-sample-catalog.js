/**
 * Syncs sample-catalog.yaml from manifests/ (repo root) into the frontend so the UI
 * shows the same content. Run before build when building from the full model-registry repo.
 * If manifests/ is not present (e.g. frontend used as subtree), this script no-ops.
 * To use a different YAML at build time, set SAMPLE_CATALOG_YAML_PATH in webpack build.
 */
const fs = require('fs');
const path = require('path');

const source = path.join(
  __dirname,
  '../../../manifests/kustomize/options/catalog/base/sample-catalog.yaml',
);
const dest = path.join(
  __dirname,
  '../src/app/pages/modelCatalogSettings/sample-catalog.yaml',
);

if (!fs.existsSync(source)) {
  process.exit(0);
}
fs.copyFileSync(source, dest);
