/**
 * Helper function to get value of a "single or array" argument. This is a type of argument which expects multiple values
 * in an array but it is acceptable to pass only a single value which will be treated as an array.
 *
 * Examples:
 *     foo("cats"), first arg becomes ["cats"] when run through this function
 *     foo(["cats", "dogs"]), first arg stays as array, function doesn't touch it.
 *
 * @param value Array of values
 */
function singleOrArrayArg (value) {
    // If value is already array
    if (Array.isArray(value) === true) {
        return value;
    }

    // "Convert"
    return [value];
}