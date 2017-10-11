// var webpack = require('webpack');

// module.exports = {
//   entry: './src/App.js',
//   output: {
//     path: './public',
//     filename: 'bundle.js'
//   },
//   plugins: [
//     new webpack.DefinePlugin({
//       'process.env': {
//         // This has effect on the react lib size
//         'NODE_ENV': JSON.stringify('development'),
//       }
//     }),
//     //new ExtractTextPlugin("bundle.css", {allChunks: false}),
//     // new webpack.optimize.AggressiveMergingPlugin(),
//     // new webpack.optimize.OccurrenceOrderPlugin(),
//     // new webpack.optimize.DedupePlugin(),
//     // new webpack.optimize.UglifyJsPlugin({
//     //   mangle: true,
//     //   compress: {
//     //     warnings: false, // Suppress uglification warnings
//     //     pure_getters: true,
//     //     unsafe: true,
//     //     unsafe_comps: true,
//     //     screw_ie8: true
//     //   },
//     //   output: {
//     //     comments: false,
//     //   },
//     //   exclude: [/\.min\.js$/gi] // skip pre-minified libs
//     // }),
//     new webpack.IgnorePlugin(/^\.\/locale$/, [/moment$/])
//     // new webpack.CompressionPlugin({
//     //   asset: "[path].gz[query]",
//     //   algorithm: "gzip",
//     //   test: /\.js$|\.css$|\.html$/,
//     //   threshold: 10240,
//     //   minRatio: 0
//     // })
//   ],
//   devtool: 'inline-source-map',
//   module: {
//     loaders: [
//       {
//         test: /\.jsx?$/,
//         exclude: /(node_modules|bower_components)/,
//         loader: 'babel',
//         query: {
//           presets: ['es2015', 'react']
//         }
//       },
//       {
//         test: /\.css$/,
//         loader: 'style-loader!css-loader'
//       },
//       {
//         test: /\.(png|jpg|gif|svg|eot|ttf|woff|woff2)$/,
//         loader: 'url-loader',
//           options: {
//             limit: 10000
//           }
//       }
//     ]
//   }
// }

var path = require('path');
var webpack = require('webpack');

var javascriptEntryPath = path.resolve(__dirname, 'src', 'App.js');
var htmlEntryPath = path.resolve(__dirname, 'src', 'index.html');
var buildPath = path.resolve(__dirname, 'public');

module.exports = {
  entry: javascriptEntryPath,
  output: {
    path: buildPath,
    filename: 'bundle.js',
  },
  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: /(node_modules|bower_components)/,
        loader: 'babel',
        query: {
          presets: ['es2015', 'react']
        }
      },
      {
        test: /\.css$/,
        loader: 'style-loader!css-loader'
      },
      {
        test: /\.(png|jpg|gif|svg|eot|ttf|woff|woff2)$/,
        loader: 'url-loader',
          options: {
            limit: 10000
          }
      }
    ]
  },
  plugins: [
    // new webpack.optimize.OccurenceOrderPlugin(),
    // new webpack.HotModuleReplacementPlugin(),
    // new webpack.NoErrorsPlugin()
  ]
}
