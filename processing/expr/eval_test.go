package expr

import "testing"

func TestEval(t *testing.T) {
	vars := map[string]string{
		"approved": "true",
		"flag":     "false",
		"name":     `"alice"`,
		"n":        "2",
	}
	cases := []struct {
		in   string
		want bool
	}{
		{"approved", true},
		{"${approved}", true},
		{"!approved", false},
		{"flag", false},
		{"!flag", true},
		{"missing", false},
		{"!missing", true},
		{"approved == true", true},
		{"${approved == false}", false},
		{"flag == false", true},
		{"name == 'alice'", true},
		{`name == "alice"`, true},
		{"name != 'bob'", true},
		{"n == 2", true},
		{"n != 1", true},
	}
	for _, tc := range cases {
		got, err := Eval(tc.in, vars)
		if err != nil {
			t.Fatalf("%s: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("%s: got %v want %v", tc.in, got, tc.want)
		}
	}
}

func TestEvalRejects(t *testing.T) {
	if _, err := Eval("", nil); err == nil {
		t.Fatal("expected error")
	}
	if _, err := Eval("!approved == true", map[string]string{"approved": "true"}); err == nil {
		t.Fatal("expected error")
	}
	if _, err := Eval("approved > 1", map[string]string{"approved": "true"}); err == nil {
		t.Fatal("expected error")
	}
}
