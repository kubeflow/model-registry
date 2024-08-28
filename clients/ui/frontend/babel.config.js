module.exports = {
    presets: [
      [
        '@babel/preset-env',
        {
          targets: {
            chrome: 110,
          },
          useBuiltIns: 'usage',
          corejs: '3',
        },
      ],
      '@babel/preset-react',
      '@babel/preset-typescript',
    ],
  };
  