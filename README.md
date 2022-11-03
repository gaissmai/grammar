[![Go Reference](https://pkg.go.dev/badge/github.com/gaissmai/grammar.svg)](https://pkg.go.dev/github.com/gaissmai/grammar)


`package grammar` is designed to make long and convoluted regular expressions easier to handle. 

## Features
* Enable the use of whitespaces within regexes to make them less dense.
* Allow spreading regexes across multiple lines to allow easier grouping. 
* Allow inline comments in regexes.
* Interpolation of regexes: Allows a regex to be composed of multiple sub-rules. 

## Usage

```
package grammar // import "github.com/gaissmai/grammar"
```
### Whitespaces / Newlines
Example regex containing white-spaces and newlines to make it more readable. 
This regex matches number of the form `1.2e+42, .3, 11, 42., 3.1415,...`. 
```
   `                       // NUMBER: 1.2e+42, .3, 11, 42., 3.1415, ...
    [+-]?                  // first, match an optional sign
    (?:                    // then match f.p mantissas or integer:
         \d+ \. \d+        // mantissa of the form a.b
       | \d+ \.            // or mantissa of the form a.
       |     \. \d+        // or mantissa of the form .b
       | \d+               // or integer of the form a
    )
    (?: [eE] [+-]? \d+ )?  // finally, optionally match an exponent
    `
```
This yields `[+-]?(\d+\.\d+|\d+\.|\.\d+|\d+)([eE][+-]?\d+)?` but is much more readable. 

### Regex Interpolation
Complex rules can be comprised of simpler subrules using string interpolation.
For example: Using the above regex as `${NUMBER}`, you can easily assemble a rule
that matches many numbers using the following snippet:

```
     ^                  // 1.23   3.1415 0.5E3 ...
       ${NUMBER}        // start with number
       (?:
         \s+ ${NUMBER}  // followed by one ore more numbers, separated by whitespace(s)
       )+
    $
```

Any number of rules can be added to a grammar, dependent or independent,
as long as there are no cyclic dependencies.

The package has a very thin API, please see the examples in the documentation.

## ATTENTION

The package is already in production use, but the API may still change.
