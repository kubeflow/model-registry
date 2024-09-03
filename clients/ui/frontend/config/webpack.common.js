/* eslint-disable @typescript-eslint/no-var-requires */

const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');
const TsconfigPathsPlugin = require('tsconfig-paths-webpack-plugin');
const Dotenv = require('dotenv-webpack');
const ASSET_PATH = process.env.ASSET_PATH || '/';
const IMAGES_DIRNAME = 'images';
const relativeDir = path.resolve(__dirname, '..');
module.exports = (env) => {
  return {
    module: {
      rules: [
        {
          test: /\.(tsx|ts|jsx)?$/,
          use: [
            {
              loader: 'ts-loader',
              options: {
                transpileOnly: true,
                experimentalWatchApi: true
              }
            }
          ]
        },
        {
          test: /\.(svg|ttf|eot|woff|woff2)$/,
          type: 'asset/resource',
          // only process modules with this loader
          // if they live under a 'fonts' or 'pficon' directory
          include: [
            path.resolve(relativeDir, 'node_modules/patternfly/dist/fonts'),
            path.resolve(relativeDir, 'node_modules/@patternfly/react-core/dist/styles/assets/fonts'),
            path.resolve(relativeDir, 'node_modules/@patternfly/react-core/dist/styles/assets/pficon'),
            path.resolve(relativeDir, 'node_modules/@patternfly/patternfly/assets/fonts'),
            path.resolve(relativeDir, 'node_modules/@patternfly/patternfly/assets/pficon')
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
            path.resolve(relativeDir, 'src'),
            path.resolve(relativeDir, 'node_modules/patternfly'),
            path.resolve(relativeDir, 'node_modules/@patternfly/patternfly/assets/images'),
            path.resolve(relativeDir, 'node_modules/@patternfly/react-styles/css/assets/images'),
            path.resolve(relativeDir, 'node_modules/@patternfly/react-core/dist/styles/assets/images'),
            path.resolve(
              relativeDir,
              'node_modules/@patternfly/react-core/node_modules/@patternfly/react-styles/css/assets/images'
            ),
            path.resolve(
              relativeDir,
              'node_modules/@patternfly/react-table/node_modules/@patternfly/react-styles/css/assets/images'
            ),
            path.resolve(
              relativeDir,
              'node_modules/@patternfly/react-inline-edit-extension/node_modules/@patternfly/react-styles/css/assets/images'
            )
          ],
          type: 'asset/inline',
          use: [
            {
              options: {
                limit: 5000,
                outputPath: 'images',
                name: '[name].[ext]'
              }
            }
          ]
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
        }
      ]
    },
    output: {
      filename: '[name].bundle.js',
      path: path.resolve(relativeDir, 'dist'),
      publicPath: ASSET_PATH
    },
    plugins: [
      new HtmlWebpackPlugin({
        template: path.resolve(relativeDir, 'src', 'index.html')
      }),
      new Dotenv({
        systemvars: true,
        silent: true
      }),
      new CopyPlugin({
        patterns: [{ from: './src/images', to: 'images' }]
      })
    ],
    resolve: {
      extensions: ['.js', '.ts', '.tsx', '.jsx'],
      plugins: [
        new TsconfigPathsPlugin({
          configFile: path.resolve(relativeDir, './tsconfig.json')
        })
      ],
      symlinks: false,
      cacheWithContext: false
    }
  };
};
