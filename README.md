[![Go Reference](https://pkg.go.dev/badge/github.com/gaissmai/grammar.svg)](https://pkg.go.dev/github.com/gaissmai/grammar)

## usage

```
package grammar // import "github.com/gaissmai/grammar"
```


package grammar allows defining regexp rules with comments, whitespace and
newlines to make them less dense, and easier to read:

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

result: `[+-]?(\d+\.\d+|\d+\.|\.\d+|\d+)([eE][+-]?\d+)?`

Complex rules can be assembled by simpler subrules using string interpolation.

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

The package is in personel production use but the API is still subject to change (with semver in mind).

