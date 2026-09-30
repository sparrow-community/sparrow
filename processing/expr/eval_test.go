package expr

import "testing"

func TestEval(t *testing.T) {
	vars := map[string]string{
		"approved":          "true",
		"flag":              "false",
		"name":              `"alice"`,
		"n":                 "2",
		"Vacation Approval": `"Approved"`,
		"clarified":         `"no"`,
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
		{"n > 1", true},
		{"approved && !flag", true},
		{"!approved == true", false},
		{"bpmn:getDataObject('approved')", true},
		{"not(bpmn:getDataObject('approved'))", false},
		{"bpmn:getDataObject('name') = 'alice'", true},
		{"bpmn:getDataObject('clarified') = 'yes'", false},
		{"= approved", true},
		{"= not(approved)", false},
		{"Vacation Approval = 'Approved'", true},
		{"Vacation Approval = 'Rejected'", false},
		{`= some risk in riskLevels satisfies risk = "red"`, true},
		{`= some risk in riskLevels satisfies risk = "blue"`, false},
		{`= every risk in riskLevels satisfies risk = "yellow"`, false},
		{`= every risk in allYellow satisfies risk = "yellow"`, true},
	}
	vars["riskLevels"] = `["red","green"]`
	vars["allYellow"] = `["yellow","yellow"]`
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
	if _, err := Eval("approved +", map[string]string{"approved": "true"}); err == nil {
		t.Fatal("expected error")
	}
}

func TestEvalJSON(t *testing.T) {
	vars := map[string]string{
		"orderId": `"o1"`,
		"total":   "9",
	}
	got, err := EvalJSON(`${orderId + "-x"}`, vars)
	if err != nil {
		t.Fatal(err)
	}
	if got != `"o1-x"` {
		t.Fatalf("got %q", got)
	}
	got, err = EvalJSON("${total * 2}", vars)
	if err != nil {
		t.Fatal(err)
	}
	if got != "18" {
		t.Fatalf("got %q", got)
	}
}
