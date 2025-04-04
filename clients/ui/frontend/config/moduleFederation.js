const path = require('path');
const { ModuleFederationPlugin } = require('@module-federation/enhanced/webpack');

const deps = require('../package.json').dependencies;

const moduleFederationConfig = {
  name: 'modelRegistry',
  filename: 'remoteEntry.js',

  shared: {
    react: { singleton: true, eager: true, requiredVersion: deps.react },
    'react-dom': { singleton: true, eager: true, requiredVersion: deps['react-dom'] },
    'react-router': { singleton: true, eager: true, requiredVersion: deps['react-router'] },
    'react-router-dom': { singleton: true, eager: true, requiredVersion: deps['react-router-dom'] },
  },
  exposes: {
    // TODO expose api. eg:
    // './index': './src/plugin/index.tsx',
    // './plugin': './src/plugin/index.tsx',
  },
  // For module federation to work when optimization.runtimeChunk="single":
  // See https://github.com/webpack/webpack/issues/18810
  runtime: false,
  // TODO generate types when exposing api
  dts: false,
};

module.exports = {
  moduleFederationPlugins: [new ModuleFederationPlugin(moduleFederationConfig)],
};
