/**
 * Custom error type which holds specific information about an argument check failure.
 */
class CheckError extends Error {
    constructor (arg, option, expected, was) {
        super(`Argument check "${option}" failed, expected: "${expected}", was: "${was}"`);// Set message

        /**
         * Name of argument, prefixed with `check:` if the Check Error occurred in an argument for check function
         * @type {string}
         */
        this.arg = arg;
        /**
         * Check type, see check()#checks
         * @type {string}
         */
        this.option = option;
        /**
         * Description of values which would meet check
         * @type {string}
         */
        this.expected = expected;
        /**
         * Description of how argument didn't meet the check
         * @type {string}
          */
        this.was = was;
    }
}

class CheckFailedError extends Error {
    constructor(checkErrs) {
        // String output which will show more info on CheckErrors
        var checkErrsDetail = "";
        // Object which organizes CheckErrors by argument id
        var argCheckErrs = {};

        // Organize CheckErrors by argument id for string output later
        for (var i = 0; i < checkErrs.length; i++) {
            // Current err in loop
            var err = checkErrs[i];

            // If first error for argument init array
            if (argCheckErrs[err.arg] === undefined) {
                argCheckErrs[err.arg] = [];
            }

            // Add error to argument section
            argCheckErrs[err.arg].push(err);
        }

        // Generate string output
        for (var argId in argCheckErrs) {
            var errs = argCheckErrs[argId];

            // Output: argument title
            checkErrsDetail += `    ${argId}:\n`;

            // Display CheckErrors for argument
            for (var i = 0; i < errs.length; i++) {
                var err = errs[i];
                checkErrsDetail += `        Failed "${err.option}" check, expected: "${err.expected}", was: "${err.was}"\n`;
            }
        }

        super(`Argument checks failed, see:\n${checkErrsDetail}`);// Set message
        this.checkErrs = checkErrs;
    }
}

/**
 * Behaves the same as check but throws an error when Check Errors occur
 * @param cargs See check()#cargs
 * @throws {CheckFailedError} If any checks fail
 */
function checkAndThrow (cargs /* checks... */) {
    var errs = check.apply(this, arguments);

    if (errs.length > 0) {
        throw new CheckFailedError(errs);
    }
}

/**
 * Runs a series of checks on an arguments object to check that the correct parameters where given to a function.
 *
 * @param cargs - Short for "check arguments", the arguments object to check
 * @param checks - Each argument passed to this function after cargs is considered a check. A check is an object defined
 *                 by the following syntax:
 *
 *                      {
 *                          // Expected type
 *                          //
 *                          // If type of argument is "object" and "typ" option is not "object" check will attempt to call
 *                          // the class() method to check against application class types (Classes defined in src with
 *                          // "class" keyword).
 *                          //
 *                          // For convenience you can pass in an object of the type you want to check for, however there
 *                          // is a slight conflict of intentions if you pass in a string as an "Object of type" because
 *                          // normally a string would contain the type name. Therefore a string is only considered an
 *                          // "Object of type" if it is an empty string (length == 0).
 *                          //
 *                          // The "typ" parameter also supports the special value of "array". This checks to make sure
 *                          // the argument is an array. This is a special argument because usually JS views arrays as
 *                          // objects so extra steps have to be taken to check if the argument is an array.
 *                          typ: String with type name | Object of type,
 *
 *                          // If the argument is optional, default false
 *                          opt: boolean,
 *
 *                          // Values argument can't be
 *                          not: Value | Array<Value>,
 *
 *                          // Fields to check for in argument, must exist but doesn't have to be defined
 *                          has: String | Array<String>
 *                      }
 *
 *                 Each check given is applied to its corresponding argument in cargs. So the first check applies to cargs[0]
 *                 the second cargs[1] and so on. If you would like to skip an argument just pass in an empty object. If
 *                 you provide checks for the first few arguments but then do not want to check any further arguments you
 *                 can simply not pass any more checks, ex:
 *                     You have the function signature: function foo(a, b, c, d)
 *                     If you only want to check arguments a and b you can call check like so:
 *                     check({typ: "number", not: -1}, {type: {}, has: "friends"}), you do not have to pass empty objects
 *                     in for arguments c and d.
 *
 * @returns Array of `CheckError`s  that where found, empty if all checks passed.
 */
function check(cargs /*,checks...*/) {
    // Check cargs is provided and an object
    var cargsFailOption = undefined;
    var cargsFailExpected = undefined;
    var cargsFailWas = undefined;

    if (cargs === undefined) {
        cargsFailOption = "opt";
        cargsFailExpected = "not undefined"
        cargsFailWas = "undefined";
    } else if (typeof cargs !== "object") {
        cargsFailOption = "typ";
        cargsFailExpected = "object";
        cargsFailWas = typeof cargs;
    }

    // Manual checks on `cargs` argument failed, return error
    if (cargsFailOption !== undefined) {
        return [new CheckError("check:cargs", cargsFailOption, cargsFailExpected, cargsFailWas)];
    }

    // Convert args objects to arrays b/c args objects are useless and pretend to be arrays
    var args = Array.prototype.slice.call(arguments, 1);// Start args array at second argument so we can leave out `cargs`
    cargs = Array.prototype.slice.call(cargs, 0);

    // Errors to return at the end
    var errors = [];

    // Get checks
    var checks = [];
    for (var i = 0; i < args.length; i++) {
        var check = args[i];

        // Make sure check is object
        if (typeof check !== "object") {
            errors.push(new CheckError(`check:checks:${i}`, "typ", "object", typeof check));
        }

        // Add to array of checks
        checks.push(check);
    }

    // Check that we have enough arguments for our checks
    if (cargs.length < checks.length) {
        errors.push(new CheckError("check:cargs", "opt", "cargs.length at least checks.length", "cargs.length less than checks.length"));
    }

    // If errors occurred while getting checks, return
    if (errors.length > 0) {
        return errors;
    }

    // Check
    for (var i = 0; i < cargs.length; i++) {
        var arg = cargs[i];

        // Ensure there is a check for this argument, if not:
        if (checks.length < i) {
            // No more checks for any arguments, stop checking
            break;
        }

        var check = checks[i];

        // Do the checking and add errors
        errors.push(..._checkArg(i, arg, check));
    }

    // Return errors
    return errors;
}

/**
 * Checks provided argument against check object. Syntax of check object is defined in check method.
 *
 * This is a private function used by the check function, do not call alone, use check function (This function does no
 * "checking" of its own arguments, it relies on the caller (check) to do that).
 *
 * @param arg Argument to check
 * @param check Check object
 * @private
 * @returns An array of `CheckError`s that occured
 */
function _checkArg (argI, arg, check) {
    var errors = [];

    // Type check
    if (check["typ"] !== undefined) {
        // Get expected type
        var typ = undefined;

        if (typeof check["typ"] === "string" && check["typ"].length > 0) {// If "typ" is string of type name
            typ = check["typ"];
        } else if (typeof check["typ"] === "string" && check["typ"].length === 0) {// If "typ" is empty string, aka string being passed as "object of type"
            typ = "string";
        } else {// Otherwise use "typ" as "Object of type"
            typ = typeof check["typ"];
        }

        // Check
        if (Array.isArray(arg) === true) {// Handle special "array" type
            // If wasn't expecting array, fail
            if (typ !== "array") {
                errors.push(new CheckError(argI, "typ", typ, "array"));
            }
        }
        // Handle "application classes" if arg is application class (if arg as check method)
        else if (typeof arg === "object" && typ !== "object" && arg["class"] !== undefined && typeof arg["class"] === "function") {
            var cOk = true;// Set to false if error occurs while reading c
            var c = arg.class();

            // Check that class() function is behaving correctly
            if (typeof c !== "string") {
                cOk = false;
                errors.push(new CheckError(`${argI}:class()`, "typ", "string", typeof  c));
            }

            // Check if c is ok
            if (cOk === true && c !== typ) {
                // We can consider this a typ mismatch because to even enter this if statement branch there has to be
                // a type mismatch (expected: object, was: != object)
                errors.push(new CheckError(argI, "typ", typ, c));
            }
        } else if (typeof arg !== typ) { // Everything else
            errors.push(new CheckError(argI, "typ", typ, typeof arg));
        }
    }

    // Optional check
    var optOk  = true;// Set to false if an error occurs while reading opt
    var opt = false;
    if (check["opt"] !== undefined) {
        // Get
        if (typeof check["opt"] !== "boolean") {
            optOk = false;
            errors.push(new CheckError(`check:checks:${argI}:opt`, "typ", "boolean", typeof check["opt"]));
        }

        opt = check["opt"];
    }

    // -- -- Check
    if (optOk === true && opt === false && (arg === undefined || arg === null)) {
        errors.push(new CheckError(argI, "opt", "not undefined", arg));
    }

    // Not check
    if (check["not"] !== undefined) {
        var not = singleOrArrayArg(check["not"]);

        // Check
        for (var i = 0; i < not.length; i++) {
            if (arg === not[i]) {// Is equal, NOT OK
                errors.push(new CheckError(argI, "not", `not ${not[i]}`, arg));
            }
        }
    }

    // Has check
    if (check["has"] !== undefined) {
        var has = singleOrArrayArg(check["has"]);

        // Check
        for (var i = 0; i < has.length; i++) {
            var key = has[i];
            var value = arg[key];

            if (value === undefined || value === null) {
                errors.push(new CheckError(argI, "has", `has ${key}`, `didn't have ${key}`));
            }
        }
    }

    // All checks complete, valid
    return errors;
}