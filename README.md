# grammar

package grammar allows defining regexp rules with comments, whitespace and
newlines to make them less dense, and easier to read:

```
   `                       // number: 1.2e+42, .3, 11, 42., 3.1415, ...
    [+-]?                  // first, match an optional sign
    (?:                    // then match integers or f.p. mantissas:
         \d+ \. \d+        // mantissa of the form a.b
       | \d+ \.            // mantissa of the form a.
       |     \. \d+        // mantissa of the form .b
       | \d+               // integer of the form a
    )
    (?: [eE] [+-]? \d+ )?  // finally, optionally match an exponent
    `
```

result: `[+-]?(\d+\.\d+|\d+\.|\.\d+|\d+)([eE][+-]?\d+)?`

Complex rules can be assembled by simpler rules using string interpolation.

```
     ^                  // 1.23   3.1415 0.5E3 ...
       ${number}        // start with number
       (?:
         \s+ ${number}  // followed by one ore more numbers, separated by whitespace(s)
       )+
    $
```

Any number of rules can be added to a grammar, dependent or independent,
as long as there are no cyclic dependencies.

## ATTENTION: it's in personel production use but the API is still subject to change

```
package grammar // import "github.com/gaissmai/grammar"

FUNCTIONS

func Trim(s string) string
    Trim removes all comments and whitespace from string.

        input: `
           foo bar // baz
            taz
        `

        result: `foobartaz`

    Trim is a helper function and would normally not be public, but it is also
    helpful if you don't want to build whole grammars, but just want to remove
    whitespace and comments from patterns.


TYPES

type Grammar struct {
	// Has unexported fields.
}
    Grammar is a container for related and maybe dependent rules. Subrules are
    string interpolated in other rules before compiling to regexp.

func New(name string) *Grammar
    New initializes a new grammar.

func (g *Grammar) Add(name string, pattern string) error
    Add rule to grammar, returns error if rule with same name already exists or
    grammar is already compiled. The pattern string gets trimmed.

func (g *Grammar) AddRaw(name string, pattern string) error
    AddRaw is similar to Add, but no trimming takes place. Use this method if
    whitespace is significant.

func (g *Grammar) Compile() error
    Compile all rules in grammar. Resolve dependencies, interpolate strings and
    compile all rules to regexp.

func (g *Grammar) Rx(name string) (*regexp.Regexp, error)
    Rx returns the compiled regexp for named rule or error if rule is not added
    or not compiled.

```
