package grammar_test

import (
	"fmt"
	"log"

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

//nolint:errcheck // for example brevity
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
                ${NUMBER}              // start with number
                (?: \s+ ${NUMBER} )+   // followed by one or more numbers, separated by whitespace
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

func ExampleGrammar_AddVerbatim() {
	verbatim := `^\QExactly like this!\E$`
	g := grammar.New("example_raw")

	if err := g.AddVerbatim("RAW_RULE", verbatim); err != nil {
		log.Fatal(err)
	}

	if err := g.Compile(); err != nil {
		log.Fatal(err)
	}

	rx, err := g.Rx("RAW_RULE")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rx)

	// Output:
	// ^\QExactly like this!\E$
}
