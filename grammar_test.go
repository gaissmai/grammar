package grammar_test

import (
	"regexp"
	"testing"

	"github.com/gaissmai/grammar"
)

func TestTrim(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "whitespace",
			input: " ",
			want:  "",
		},
		{
			name: "complex",
			input: `
				foo bar  // baz
				taz asdf$
				42       // result
			`,
			want: "foobartazasdf$42",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := grammar.Trim(tc.input)
			if got != tc.want {
				t.Errorf("strip: %q, want: %q, got: %q", tc.input, tc.want, got)
			}
		})
	}
}

func TestAddOneRule(t *testing.T) {
	t.Parallel()
	g := grammar.New("TEST")

	raw := `^ .* // test`
	checkErr(t, g.Add("ONE", raw))
	checkErr(t, g.Compile())

	rx, err := g.Rx("ONE")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	want := regexp.MustCompile(`^.*`)
	if rx.String() != want.String() {
		t.Errorf("Rx(): %q, want: %q, got: %q", raw, want, rx)
	}
}

func TestAddRuleTwice(t *testing.T) {
	t.Parallel()
	g := grammar.New("TEST")

	raw := `^ .* // test`
	checkErr(t, g.Add("ONE", raw))
	err := g.Add("ONE", raw)

	// expect error
	if err == nil {
		t.Error("expected error adding duplicate rules")
	}
}

func TestRegexpCompile(t *testing.T) {
	t.Parallel()
	raw := `^ ( ` // <-- missing closing )

	g := grammar.New("TEST")
	checkErr(t, g.Add("ONE", raw))

	// expect error
	err := g.Compile()
	if err == nil {
		t.Fatal(err)
	}
}

func TestTemplateError(t *testing.T) {
	t.Parallel()
	one := `world`
	two := `hello ${{ONE}}` // ${foo} and not ${{foo}}

	g := grammar.New("TEST")
	checkErr(t, g.Add("ONE", one))
	checkErr(t, g.Add("TWO", two))

	// expect error
	err := g.Compile()
	if err == nil {
		t.Fatal(err)
	}
}

func TestMissingRule(t *testing.T) {
	t.Parallel()
	g := grammar.New("TEST")

	raw := `interpolate ${TWO} here`
	checkErr(t, g.Add("ONE", raw))

	// expect error
	err := g.Compile()
	if err == nil {
		t.Error("expected error, missing rule")
	}

	// expect error
	_, err = g.Rx("TWO")
	if err == nil {
		t.Error("expected error, missing rule")
	}

	// expect error
	_, err = g.Rx("ONE")
	if err == nil {
		t.Error("expected error, not compiled")
	}
}

func TestSelfReference(t *testing.T) {
	t.Parallel()
	g := grammar.New("TEST")

	raw := `interpolate ${ONE} here`

	// expect error
	err := g.Add("ONE", raw)
	if err == nil {
		t.Error("expected error, already compiled")
	}
}

func TestCyclicReference(t *testing.T) {
	t.Parallel()
	g := grammar.New("TEST")

	one := `interpolate ${TWO} here`
	two := `interpolate ${TRE} here`
	tre := `interpolate ${ONE} here`

	checkErr(t, g.Add("ONE", one))
	checkErr(t, g.Add("TWO", two))
	checkErr(t, g.Add("TRE", tre))

	err := g.Compile()
	if err == nil {
		t.Error("expected error, cyclic reference")
	}
}

func TestCompile(t *testing.T) {
	t.Parallel()
	g := grammar.New("TEST")

	raw := ``
	checkErr(t, g.Add("ONE", raw))
	checkErr(t, g.Compile())

	// expect error, already compiled
	err := g.Add("TWO", raw)
	if err == nil {
		t.Error("expected error, already compiled")
	}
	// expect error, already compiled
	err = g.Compile()
	if err == nil {
		t.Error("expected error, already compiled")
	}
}

func TestVariableName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input      string
		expectFail bool
	}{
		{" ", false},                // no variable at all
		{" ${Ident} ", false},       // identifier
		{" ${_underscore} ", false}, // identifier starting with underscore
		{" ${1digit} ", true},       // not allowed, starts with digit
		{" ${Iden)t}} ", true},      // typo
	}

	for _, tc := range tests {
		g := grammar.New("TESTS")
		err := g.Add("VAR", tc.input)
		if tc.expectFail && err == nil {
			t.Fatalf("AddRule(): %q, want: error, got: %q", tc.input, err)
		}

		if !tc.expectFail && err != nil {
			t.Fatalf("AddRule(): %q, want: success, got: %q", tc.input, err)
		}
	}
}

func TestIP(t *testing.T) {
	t.Parallel()
	rawIP := ` ${LikeIPv4} | ${LikeIPv6} `

	rawIPv4 := ` \d{1,3} \. \d{1,3} \. \d{1,3} \. \d{1,3} ` // just to feed netip.ParseAddr

	rawIPv6 := `(?:                                   // very minimalistic, just to feed netip.ParseAddr
                	[[:xdigit:] :]+ : [[:xdigit:] :]+ // hexdigits and colon, NO dot, no IP4in6
							  |
						    	::                                // unspecified IPv6 against all regular rules
					    )`

	g := grammar.New("IP")
	checkErr(t, g.Add("LikeIP", rawIP))
	checkErr(t, g.Add("LikeIPv4", rawIPv4))
	checkErr(t, g.Add("LikeIPv6", rawIPv6))
	checkErr(t, g.Compile())

	rx, err := g.Rx("LikeIP")
	if err != nil {
		t.Fatal(err)
	}

	got := rx.String()
	want := `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|(?:[[:xdigit:]:]+:[[:xdigit:]:]+|::)`
	if want != got {
		t.Errorf("minimalistic IP rules\nwant: %s\ngot: %s\n", want, got)
	}
}

func checkErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
