var
	autoprefixer = require('gulp-autoprefixer'),
	browserify = require('gulp-browserify'),
	concat = require('gulp-concat'),
	eventStream = require('event-stream'),
	gulp = require('gulp'),
	gulpif = require('gulp-if'),
	inject = require('gulp-inject'),
	less = require('gulp-less'),
	livereload = require('gulp-livereload'),
	minifyCss = require('gulp-minify-css'),
	minifyHtml = require('gulp-minify-html'),
	minimist = require('minimist'),
	plumber = require('gulp-plumber'),
	process = require('process'),
	recess = require('gulp-recess'),
	uglify = require('gulp-uglify'),
	watch = require('gulp-watch')
;

var argv = minimist(process.argv);
var config = {
	src: 'src',
	dist: 'dist',
	minify: argv.minify
};

var css = function() {
	return gulp.src(config.src + '/**/*.less')
		.pipe(plumber())
		// .pipe(recess())
		.pipe(less())
		.pipe(autoprefixer())
		.pipe(gulpif(config.minify, concat('main.css')))
		.pipe(gulpif(config.minify, minifyCss()))
		.pipe(gulp.dest(config.dist));
};
gulp.task('css', css);

var js = function() {
	return gulp.src(config.src + '/main.js')
		.pipe(plumber())
		.pipe(browserify({
			transform: ['reactify']
		}))
		.pipe(gulpif(config.minify, concat('main.js')))
		.pipe(gulpif(config.minify, uglify()))
		.pipe(gulp.dest(config.dist));
};
gulp.task('js', js);

gulp.task('html', ['css', 'js'], function() {
	gulp.src(config.src + '/**/*.html')
		.pipe(plumber())
		.pipe(inject(eventStream.merge(js(), css())))
		.pipe(gulpif(config.minify, minifyHtml()))
		.pipe(gulp.dest(config.dist));
});

gulp.task('build', ['html']);

gulp.task('watch', ['build'], function() {
	gulp.watch(config.src + '/**/*', ['build']);
	gulp.src(config.dist + '/**/*')
		.pipe(watch())
		.pipe(livereload());
});

gulp.task('default', ['watch']);
