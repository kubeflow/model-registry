/* eslint-disable @typescript-eslint/no-var-requires */

const path = require('path');
const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');
const { stylePaths } = require('./stylePaths');
const HOST = process.env.HOST || 'localhost';
const PORT = process.env.PORT || '9000';
const PROXY_HOST = process.env.PROXY_HOST || 'localhost';
const PROXY_PORT = process.env.PROXY_PORT || '4000';
const PROXY_PROTOCOL = process.env.PROXY_PROTOCOL || 'http:';
const relativeDir = path.resolve(__dirname, '..');
const PROXY_TARGET = `${PROXY_PROTOCOL}//${PROXY_HOST}:${PROXY_PORT}`;

module.exports = merge(common('development'), {
  mode: 'development',
  devtool: 'eval-source-map',
  devServer: {
    host: HOST,
    port: PORT,
    historyApiFallback: true,
    open: true,
    static: {
      directory: path.resolve(relativeDir, 'dist')
    },
    client: {
      overlay: true
    },
    proxy: [
      {
        context: ['/api'],
        target: PROXY_TARGET,
        host: PROXY_HOST,
        protocol: PROXY_PROTOCOL,
        port: PROXY_PORT,
        secure: false
      },
    ],
  },
  module: {
    rules: [
      {
        test: /\.css$/,
        include: [...stylePaths],
        use: ['style-loader', 'css-loader']
      }
    ]
  }
});
