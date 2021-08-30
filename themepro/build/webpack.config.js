const path = require('path');

const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const TerserPlugin = require("terser-webpack-plugin");
const ImageMinimizerPlugin = require("image-minimizer-webpack-plugin");

const production = process.env.NODE_ENV === 'production';

module.exports = {
  mode: production ? 'production' : 'development',
  entry: {
    main: './src/main.js',
    'page.package': './src/page.package.js',
  },
  output: {
    path: `${__dirname}/../files`,
    filename: '[name].js',
    assetModuleFilename: '[name][ext]',
  },
  plugins: [
    new MiniCssExtractPlugin({
      filename: '[name].css',
    }),
    new ImageMinimizerPlugin({
      minimizerOptions: {
        plugins: [
          "svgo",
        ],
      },
    })
  ],
  module: {
    rules: [
      {
        test: /\.svg$/,
        type: 'asset/resource',
        include: path.resolve(__dirname, './src/assets'),
      },
      {
        test: /\.svg$/,
        type: 'asset/source',
        include: path.resolve(__dirname, './src/inlineassets'),
      },
      {
        test: /\.css$/,
        use: [
          MiniCssExtractPlugin.loader,
          'css-loader',
          'postcss-loader',
        ],
      },
    ],
  },
  optimization: {
    minimizer: [
      new TerserPlugin({
        extractComments: false,
      }),
      new CssMinimizerPlugin(),
    ],
  },
};
