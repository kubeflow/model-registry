const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');
const { setupWebpackDotenvFilesForEnv } = require('./dotenv');
const { name } = require('../package.json');

const { moduleFederationPlugins } = require('./moduleFederation');

const RELATIVE_DIRNAME = process.env._RELATIVE_DIRNAME;
const IS_PROJECT_ROOT_DIR = process.env._IS_PROJECT_ROOT_DIR;
const IMAGES_DIRNAME = process.env._IMAGES_DIRNAME;
const PUBLIC_PATH = process.env._PUBLIC_PATH;
const SRC_DIR = process.env._SRC_DIR;
const COMMON_DIR = process.env._COMMON_DIR;
const DIST_DIR = process.env._DIST_DIR;
const OUTPUT_ONLY = process.env._OUTPUT_ONLY;
const FAVICON = process.env.FAVICON;
const PRODUCT_NAME = process.env.PRODUCT_NAME;
const COVERAGE = process.env.COVERAGE;
const DEPLOYMENT_MODE = process.env._DEPLOYMENT_MODE;
const BASE_PATH = DEPLOYMENT_MODE === 'integrated' ? '/model-registry/' : PUBLIC_PATH;

if (OUTPUT_ONLY !== 'true') {
  console.info(
    `\nPrepping files...\n  SRC DIR: ${SRC_DIR}\n  OUTPUT DIR: ${DIST_DIR}\n  PUBLIC PATH: ${PUBLIC_PATH}\n  BASE_PATH: ${BASE_PATH}\n`,
  );
  if (COVERAGE === 'true') {
    console.info('\nAdding code coverage instrumentation.\n');
  }
}

module.exports = (env) => ({
  entry: {
    app: path.join(SRC_DIR, 'index.ts'),
  },
  module: {
    rules: [
      {
        test: /\.(tsx|ts|jsx|js)?$/,
        exclude: [/node_modules/, /__tests__/, /__mocks__/],
        include: [SRC_DIR, COMMON_DIR],
        use: [
          COVERAGE === 'true' && '@jsdevtools/coverage-istanbul-loader',
          env === 'development'
            ? { loader: 'swc-loader' }
            : {
                loader: 'ts-loader',
                options: {
                  transpileOnly: true,
                },
              },
        ],
      },
      {
        test: /\.(svg|ttf|eot|woff|woff2)$/,
        // only process modules with this loader
        // if they live under a 'fonts' or 'pficon' directory
        include: [
          path.resolve(RELATIVE_DIRNAME, 'node_modules/patternfly/dist/fonts'),
          path.resolve(
            RELATIVE_DIRNAME,
            'node_modules/@patternfly/react-core/dist/styles/assets/fonts',
          ),
          path.resolve(
            RELATIVE_DIRNAME,
            'node_modules/@patternfly/react-core/dist/styles/assets/pficon',
          ),
          path.resolve(RELATIVE_DIRNAME, 'node_modules/@patternfly/patternfly/assets/fonts'),
          path.resolve(RELATIVE_DIRNAME, 'node_modules/@patternfly/patternfly/assets/pficon'),
        ],
        use: {
          loader: 'file-loader',
          options: {
            // Limit at 50k. larger files emitted into separate files
            limit: 5000,
            outputPath: 'fonts',
            name: '[name].[ext]',
          },
        },
      },
      {
        test: /\.svg$/,
        include: (input) => input.indexOf('background-filter.svg') > 1,
        use: [
          {
            loader: 'url-loader',
            options: {
              limit: 5000,
              outputPath: 'svgs',
              name: '[name].[ext]',
            },
          },
        ],
      },
      {
        test: /\.svg$/,
        // only process SVG modules with this loader if they live under a 'bgimages' directory
        // this is primarily useful when applying a CSS background using an SVG
        include: (input) => input.indexOf(IMAGES_DIRNAME) > -1,
        use: {
          loader: 'svg-url-loader',
          options: {
            limit: 10000,
          },
        },
      },
      {
        test: /\.svg$/,
        // only process SVG modules with this loader when they don't live under a 'bgimages',
        // 'fonts', or 'pficon' directory, those are handled with other loaders
        include: (input) =>
          input.indexOf(IMAGES_DIRNAME) === -1 &&
          input.indexOf('fonts') === -1 &&
          input.indexOf('background-filter') === -1 &&
          input.indexOf('pficon') === -1,
        use: {
          loader: 'raw-loader',
          options: {},
        },
      },
      {
        test: /\.(jpg|jpeg|png|gif)$/i,
        include: [
          SRC_DIR,
          COMMON_DIR,
          path.resolve(RELATIVE_DIRNAME, 'node_modules/patternfly'),
          path.resolve(RELATIVE_DIRNAME, 'node_modules/@patternfly/patternfly/assets/images'),
          path.resolve(RELATIVE_DIRNAME, 'node_modules/@patternfly/react-styles/css/assets/images'),
          path.resolve(
            RELATIVE_DIRNAME,
            'node_modules/@patternfly/react-core/dist/styles/assets/images',
          ),
          path.resolve(
            RELATIVE_DIRNAME,
            'node_modules/@patternfly/react-core/node_modules/@patternfly/react-styles/css/assets/images',
          ),
          path.resolve(
            RELATIVE_DIRNAME,
            'node_modules/@patternfly/react-table/node_modules/@patternfly/react-styles/css/assets/images',
          ),
          path.resolve(
            RELATIVE_DIRNAME,
            'node_modules/@patternfly/react-inline-edit-extension/node_modules/@patternfly/react-styles/css/assets/images',
          ),
        ],
        use: [
          {
            loader: 'url-loader',
            options: {
              limit: 5000,
              outputPath: 'images',
              name: '[name].[ext]',
            },
          },
        ],
      },
      {
        test: /\.s[ac]ss$/i,
        use: [
          // Creates `style` nodes from JS strings
          'style-loader',
          // Translates CSS into CommonJS
          'css-loader',
          // Compiles Sass to CSS
          'sass-loader',
        ],
      },
      {
        test: /\.ya?ml$/,
        use: 'js-yaml-loader',
      },
    ],
  },
  output: {
    filename: '[name].bundle.js',
    path: DIST_DIR,
    publicPath: 'auto',
    uniqueName: name,
  },
  plugins: [
    ...moduleFederationPlugins,
    ...setupWebpackDotenvFilesForEnv({
      directory: RELATIVE_DIRNAME,
      isRoot: IS_PROJECT_ROOT_DIR,
    }),
    new HtmlWebpackPlugin({
      template: path.join(SRC_DIR, 'index.html'),
      title: PRODUCT_NAME,
      favicon: path.join(SRC_DIR, 'images', FAVICON),
      publicPath: BASE_PATH,
      base: {
        href: BASE_PATH,
      },
      chunks: ['app'],
    }),
    new CopyPlugin({
      patterns: [
        {
          from: path.join(SRC_DIR, 'locales'),
          to: path.join(DIST_DIR, 'locales'),
          noErrorOnMissing: true,
        },
        {
          from: path.join(SRC_DIR, 'favicons'),
          to: path.join(DIST_DIR, 'favicons'),
          noErrorOnMissing: true,
        },
        {
          from: path.join(SRC_DIR, 'images'),
          to: path.join(DIST_DIR, 'images'),
          noErrorOnMissing: true,
        },
        {
          from: path.join(SRC_DIR, 'favicon.ico'),
          to: path.join(DIST_DIR),
          noErrorOnMissing: true,
        },
        {
          from: path.join(SRC_DIR, 'favicon.png'),
          to: path.join(DIST_DIR),
          noErrorOnMissing: true,
        },
        {
          from: path.join(SRC_DIR, 'manifest.json'),
          to: path.join(DIST_DIR),
          noErrorOnMissing: true,
        },
        {
          from: path.join(SRC_DIR, 'robots.txt'),
          to: path.join(DIST_DIR),
          noErrorOnMissing: true,
        },
      ],
    }),
  ],
  resolve: {
    extensions: ['.js', '.ts', '.tsx', '.jsx'],
    alias: {
      '~': path.resolve(SRC_DIR),
    },
    symlinks: false,
    cacheWithContext: false,
  },
});
