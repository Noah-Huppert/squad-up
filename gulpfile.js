var gulp = require("gulp");
var babel = require("gulp-babel");
var watch = require("gulp-watch");

var targets = ["client/components/**/*.js", "client/js/**/*.js"];

var runI = 0;

func babelWatch (targets, out) {
    return watch(targets, () => {
        console.log(`Transpiling ${out} (run #${runI})`);
        runI++;

        gulp.src(targets)
            .pipe(babel())
            .pipe(gulp.dest(`dist/${out}`));
    });
}

gulp.task("default", () => {
    // TODO: Call babelWatch for components and js in parrallel
});