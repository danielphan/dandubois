module.exports = {
  resolve: {
      moduleDirectories: ["node_modules"]
  },
  module: {
    preLoaders: [
      { test: /\.js$/, exclude: /node_modules/, loader: "jshint" },
      { test: /\.jsx$/, exclude: /node_modules/, loader: "jsxhint" }
    ],
    loaders: [
      { test: /\.js$/, loader: "envify-loader" },
      { test: /\.jsx$/, loader: "react-hot-loader!envify-loader!jsx-loader" },
      { test: /\.css$/, loader: "style-loader!css-loader!autoprefixer-loader" },
      { test: /\.less$/, loader: "style-loader!css-loader!autoprefixer-loader!less-loader" },
      { test: /\.(png|jpe?g|gif|svg|ico)$/, loader: "url-loader?limit=100000" }
    ]
  },
  jshint: {
    globalstrict: true,
    browser: true
  }
};
