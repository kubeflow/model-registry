const { execSync } = require('child_process');
const path = require('path');
const { merge } = require('webpack-merge');
const { setupWebpackDotenvFilesForEnv, setupDotenvFilesForEnv } = require('./dotenv');
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
const ReactRefreshWebpackPlugin = require('@pmmmwh/react-refresh-webpack-plugin');
const SpeedMeasurePlugin = require('speed-measure-webpack-plugin');

const smp = new SpeedMeasurePlugin({ disable: !process.env.MEASURE });

setupDotenvFilesForEnv({ env: 'development' });
const webpackCommon = require('./webpack.common.js');

const RELATIVE_DIRNAME = process.env._RELATIVE_DIRNAME;
const IS_PROJECT_ROOT_DIR = process.env._IS_PROJECT_ROOT_DIR;
const SRC_DIR = process.env._SRC_DIR;
const COMMON_DIR = process.env._COMMON_DIR;
const PUBLIC_PATH = process.env._PUBLIC_PATH;
const DIST_DIR = process.env._DIST_DIR;
const HOST = process.env._HOST;
const PORT = process.env._PORT;
const PROXY_PROTOCOL = process.env._PROXY_PROTOCOL;
const PROXY_HOST = process.env._PROXY_HOST;
const PROXY_PORT = process.env._PROXY_PORT;
const DEPLOYMENT_MODE = process.env._DEPLOYMENT_MODE;
const AUTH_METHOD = process.env._AUTH_METHOD;
const BASE_PATH = DEPLOYMENT_MODE === 'kubeflow' ? '/model-registry/' : PUBLIC_PATH;

// Function to generate headers based on deployment mode
const getProxyHeaders = () => {
  if (AUTH_METHOD === 'internal') {
    return {
      'kubeflow-userid': 'user@example.com',
    };
  }
  if (AUTH_METHOD === 'user_token') {
    try {
      const token = execSync(
        "kubectl config view --raw --minify --flatten -o jsonpath='{.users[].user.token}'",
      )
        .toString()
        .trim();
      const username = execSync("kubectl auth whoami -o jsonpath='{.status.userInfo.username}'")
        .toString()
        .trim();
      // eslint-disable-next-line no-console
      console.info('Logged in as user:', username);
      return {
        Authorization: `Bearer ${token}`,
        'x-forwarded-access-token': token,
      };
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Failed to get Kubernetes token:', error.message);
      return {};
    }
  }
  return {};
};

module.exports = smp.wrap(
  merge(
    {
      plugins: [
        ...setupWebpackDotenvFilesForEnv({
          directory: RELATIVE_DIRNAME,
          env: 'development',
          isRoot: IS_PROJECT_ROOT_DIR,
        }),
      ],
    },
    webpackCommon('development'),
    {
      mode: 'development',
      devtool: 'eval-source-map',
      optimization: {
        runtimeChunk: 'single',
        removeEmptyChunks: true,
      },
      devServer: {
        host: HOST,
        port: PORT,
        compress: true,
        historyApiFallback: true,
        hot: true,
        open: false,
        proxy: [
          {
            context: ['/api', '/model-registry/api'],
            target: {
              host: PROXY_HOST,
              protocol: PROXY_PROTOCOL,
              port: PROXY_PORT,
            },
            changeOrigin: true,
            headers: getProxyHeaders(),
          },
        ],
        devMiddleware: {
          stats: 'errors-only',
        },
        client: {
          overlay: false,
        },
        static: {
          directory: DIST_DIR,
          publicPath: BASE_PATH,
        },
        onListening: (devServer) => {
          if (devServer) {
            // eslint-disable-next-line no-console
            console.log(
              `\x1b[32m✓ Dashboard available at: \x1b[4mhttp://localhost:${
                devServer.server.address().port
              }\x1b[0m`,
            );
          }
        },
      },
      module: {
        rules: [
          {
            test: /\.css$/,
            include: [
              SRC_DIR,
              COMMON_DIR,
              path.resolve(RELATIVE_DIRNAME, 'node_modules/@patternfly'),
              path.resolve(
                RELATIVE_DIRNAME,
                'node_modules/mod-arch-shared/node_modules/@patternfly',
              ),
            ],
            use: ['style-loader', 'css-loader'],
          },
        ],
      },
      plugins: [
        new ForkTsCheckerWebpackPlugin(),
        new ReactRefreshWebpackPlugin({ overlay: false }),
      ],
    },
  ),
);
