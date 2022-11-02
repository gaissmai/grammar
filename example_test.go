package grammar_test

import (
	"fmt"

	"github.com/gaissmai/grammar"
)

func ExampleTrim() {
	number := `
            [+-]?                  // first, match an optional sign
            (?:                    // then match mantissas:
                \d+ \. \d+         // mantissa of the form a.b
              | \d+ \.             // or mantissa of the form a.
              |     \. \d+         // or mantissa of the form .b
              | \d+                // or integer of the form a
            )
            (?: [eE] [+-]? \d+ )?  // finally, optionally match an exponent
         `
	fmt.Println(grammar.Trim(number))

	// Output:
	// [+-]?(?:\d+\.\d+|\d+\.|\.\d+|\d+)(?:[eE][+-]?\d+)?
}

func ExampleNew() {
	subrule := `                       // NUMBER
              [+-]?                  // first, match an optional sign
              (?:                    // then match mantissas:
                  \d+ \. \d+         // mantissa of the form a.b
                | \d+ \.             // or mantissa of the form a.
                |     \. \d+         // or mantissa of the form .b
                | \d+                // or integer of the form a
              )
              (?: [eE] [+-]? \d+ )?  // finally, optionally match an exponent
            `

	rule := `^ \s*                       // MANY NUMBERS
                ${NUMBER}              // star with number
                (?: \s+ ${NUMBER} )+   // foolowed by one or more numbers, separated by whitespace
              $
             `

	g := grammar.New("example_interpolation")

	// error handling neglected in this example for better clarity
	g.Add("MANY", rule)
	g.Add("NUMBER", subrule)
	g.Compile()
	rx, _ := g.Rx("MANY")

	fmt.Println(rx)

	// Output:
	// ^\s*[+-]?(?:\d+\.\d+|\d+\.|\.\d+|\d+)(?:[eE][+-]?\d+)?(?:\s+[+-]?(?:\d+\.\d+|\d+\.|\.\d+|\d+)(?:[eE][+-]?\d+)?)+$
}

func ExampleGrammar_AddRaw() {
	verbatim := `^\QExactly like this!\E$`
	g := grammar.New("example_raw")

	// error handling neglected in this example for better clarity
	g.AddRaw("RAW_RULE", verbatim)
	g.Compile()
	rx, _ := g.Rx("RAW_RULE")

	fmt.Println(rx)

	// Output:
	// ^\QExactly like this!\E$
}
