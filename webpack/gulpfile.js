var del = require("del");
var gulp = require("gulp");
var gulpWebpack = require("gulp-webpack");
var inject = require("gulp-inject");
var util = require("gulp-util");
var webpack = require("webpack");
var WebpackDevServer = require("webpack-dev-server");

var config = {
  src: "./src/",
  dest: "./dist/"
};

var webpackConfig = require("./webpack.config.js");
webpackConfig.entry = config.src + "main.js";

gulp.task("build", function() {
  var webpackProdConfig = Object.create(webpackConfig);
  webpackProdConfig.plugins = webpackProdConfig.plugins || [];
  webpackProdConfig.plugins.push(new webpack.DefinePlugin({
    "process.env.NODE_ENV": JSON.stringify("production"),
    "process.env.API_URL": JSON.stringify("http://dandubois.net/api")
  }));
  webpackProdConfig.plugins.push(new webpack.optimize.UglifyJsPlugin());

  var bundle = gulp.src(webpackConfig.entry)
    .pipe(gulpWebpack(webpackProdConfig, webpack))
    .pipe(gulp.dest(config.dest));

  return gulp.src(config.src + "index.html")
    .pipe(inject(bundle))
    .pipe(gulp.dest(config.dest));
});

gulp.task("serve", function() {
  var webpackDevConfig = Object.create(webpackConfig);
  webpackDevConfig.output = {
    path: __dirname + config.dest,
    filename: "bundle.js"
  };
  webpackDevConfig.devtool = "source-map";
  webpackDevConfig.plugins = webpackDevConfig.plugins || [];
  webpackDevConfig.plugins.push(new webpack.DefinePlugin({
    "process.env.NODE_ENV": JSON.stringify("development"),
    "process.env.API_URL": JSON.stringify("http://localhost:8080/api")
  }));
  webpackDevConfig.jshint.devel = true;

  var compiler = webpack(webpackDevConfig);
  new WebpackDevServer(compiler, {
    stats: { colors: true }
  }).listen(8090, "localhost", function(err) {
    if (err) {
      throw new util.PluginError("webpack-dev-server", err);
    }
    util.log("[webpack-dev-server]",
      util.colors.bgGreen("http://localhost:8090/webpack-dev-server/bundle"));
  });
});

gulp.task("clean", function(done) {
  del(config.dest, done);
});

gulp.task("default", ["serve"]);
