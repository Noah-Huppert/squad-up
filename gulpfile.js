var gulp = require("gulp");
var babel = require("gulp-babel");
var watch = require("gulp-watch");

var targets = ["client/components/**/*.js", "client/js/**/*.js"];

var runI = 0;

gulp.task("default", () => {
    return watch(targets, () => {
        console.log(`Transpiling (run #${runI})`);
        runI++;

        gulp.src(targets)
            .pipe(babel())
            .pipe(gulp.dest("dist"));
        // TODO: Figure out how to make it so dist has the sub dirs of {components,js}
    });
});