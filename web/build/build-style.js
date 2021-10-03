const gulp = require("gulp");
const cleanCSS = require("gulp-clean-css");
const less = require("gulp-less");
const rename = require("gulp-rename");
const autoprefixer = require("gulp-autoprefixer");

// 编译less
gulp.task("css", function() {
  gulp
    .src("../src/assets/styles/index.less")
    .pipe(less({ javascriptEnabled: true }))
    .pipe(
      autoprefixer({
        browsers: ["last 2 versions", "ie > 8"]
      })
    )
    .pipe(cleanCSS())
    .pipe(rename("bundle.css"))
    .pipe(gulp.dest("../src/assets"));
});

gulp.task("fonts", function() {
  gulp
    .src("../src/assets/styles/common/iconfont/fonts/*.*")
    .pipe(gulp.dest("../src/assets/fonts"));
});

gulp.task("default", ["css", "fonts"]);
