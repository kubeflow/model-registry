import path from 'path';

const webpackConfig = {
  mode: 'development' as const,
  resolve: {
    extensions: ['.ts', '.tsx', '.js', '.jsx'],
    alias: {
      '~': path.resolve(__dirname, '../../'),
    },
  },
  module: {
    rules: [
      {
        test: /\.(tsx|ts|jsx|js)?$/,
        exclude: [/node_modules/],
        use: [
          {
            loader: 'ts-loader',
            options: {
              transpileOnly: true,
            },
          },
        ],
      },
      {
        test: /\.(svg|ttf|eot|woff|woff2)$/,
        // Handle fonts from PatternFly and other sources
        include: [
          path.resolve(__dirname, '../../../node_modules/patternfly/dist/fonts'),
          path.resolve(
            __dirname,
            '../../../node_modules/@patternfly/react-core/dist/styles/assets/fonts',
          ),
          path.resolve(
            __dirname,
            '../../../node_modules/@patternfly/react-core/dist/styles/assets/pficon',
          ),
          path.resolve(__dirname, '../../../node_modules/@patternfly/patternfly/assets/fonts'),
          path.resolve(__dirname, '../../../node_modules/@patternfly/patternfly/assets/pficon'),
          path.resolve(__dirname, '../../../node_modules/mod-arch-shared'),
        ],
        use: {
          loader: 'file-loader',
          options: {
            limit: 5000,
            outputPath: 'fonts',
            name: '[name].[ext]',
          },
        },
      },
      {
        test: /\.svg$/,
        include: (input: string): boolean => input.indexOf('background-filter.svg') > 1,
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
        // Handle SVG files from mod-arch-shared and other sources
        include: (input: string): boolean =>
          (input.indexOf('mod-arch-shared') > -1 || input.indexOf('images') > -1) &&
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
          path.resolve(__dirname, '../../'),
          path.resolve(__dirname, '../../../node_modules/patternfly'),
          path.resolve(__dirname, '../../../node_modules/@patternfly/patternfly/assets/images'),
          path.resolve(
            __dirname,
            '../../../node_modules/@patternfly/react-styles/css/assets/images',
          ),
          path.resolve(
            __dirname,
            '../../../node_modules/@patternfly/react-core/dist/styles/assets/images',
          ),
          path.resolve(__dirname, '../../../node_modules/mod-arch-shared'),
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
        use: ['style-loader', 'css-loader', 'sass-loader'],
      },
      {
        test: /\.css$/i,
        use: ['style-loader', 'css-loader'],
      },
    ],
  },
};

export default webpackConfig;
