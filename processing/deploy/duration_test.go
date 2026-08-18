package deploy

import (
	"testing"
	"time"
)

func TestParseISO8601Duration(t *testing.T) {
	cases := []struct {
		in      string
		want    time.Duration
		wantErr bool
	}{
		{in: "PT0S", want: 0},
		{in: "pt1s", want: time.Second},
		{in: "PT0.5S", want: 500 * time.Millisecond},
		{in: "PT1M", want: time.Minute},
		{in: "PT1H", want: time.Hour},
		{in: "PT1H30M", want: time.Hour + 30*time.Minute},
		{in: "PT1H2M3S", want: time.Hour + 2*time.Minute + 3*time.Second},
		{in: "PT", wantErr: true},
		{in: "P1D", wantErr: true},
		{in: "PT1D", wantErr: true},
		{in: "", wantErr: true},
		{in: "1H", wantErr: true},
	}
	for _, tc := range cases {
		got, err := ParseISO8601Duration(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseISO8601Duration(%q) err=nil want error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseISO8601Duration(%q) err=%v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseISO8601Duration(%q)=%v want %v", tc.in, got, tc.want)
		}
	}
}

func TestParseISO8601Date(t *testing.T) {
	cases := []struct {
		in      string
		want    time.Time
		wantErr bool
	}{
		{in: "2000-01-01T00:00:00Z", want: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{in: "2000-01-01T00:00:00.500Z", want: time.Date(2000, 1, 1, 0, 0, 0, 500000000, time.UTC)},
		{in: "2000-01-01T08:00:00+08:00", want: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{in: "2000-01-01T00:00:00", want: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{in: "2000-01-01", want: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{in: "", wantErr: true},
		{in: "PT1H", wantErr: true},
		{in: "not-a-date", wantErr: true},
	}
	for _, tc := range cases {
		got, err := ParseISO8601Date(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseISO8601Date(%q) err=nil want error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseISO8601Date(%q) err=%v", tc.in, err)
			continue
		}
		if !got.Equal(tc.want) {
			t.Errorf("ParseISO8601Date(%q)=%v want %v", tc.in, got, tc.want)
		}
	}
}

func TestParseISO8601Cycle(t *testing.T) {
	cases := []struct {
		in       string
		repeat   int
		period   time.Duration
		hasStart bool
		hasEnd   bool
		wantErr  bool
	}{
		{in: "R/PT1H", repeat: -1, period: time.Hour},
		{in: "R3/PT10S", repeat: 3, period: 10 * time.Second},
		{in: "R/2000-01-01T00:00:00Z/PT1H", repeat: -1, period: time.Hour, hasStart: true},
		{in: "R/PT1H/2099-01-01T00:00:00Z", repeat: -1, period: time.Hour, hasEnd: true},
		{in: "", wantErr: true},
		{in: "PT1H", wantErr: true},
		{in: "R/P1D", wantErr: true},
		{in: "R0/PT1H", wantErr: true},
	}
	for _, tc := range cases {
		got, err := ParseISO8601Cycle(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseISO8601Cycle(%q) err=nil want error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseISO8601Cycle(%q) err=%v", tc.in, err)
			continue
		}
		if got.Repeat != tc.repeat || got.Period != tc.period {
			t.Errorf("ParseISO8601Cycle(%q) repeat=%d period=%v want %d %v", tc.in, got.Repeat, got.Period, tc.repeat, tc.period)
		}
		if tc.hasStart != !got.Start.IsZero() || tc.hasEnd != !got.End.IsZero() {
			t.Errorf("ParseISO8601Cycle(%q) start=%v end=%v", tc.in, got.Start, got.End)
		}
	}
}

func TestCycleFirstDue(t *testing.T) {
	now := time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC)
	c, err := ParseISO8601Cycle("R/PT1H")
	if err != nil {
		t.Fatal(err)
	}
	due, err := c.FirstDue(now)
	if err != nil {
		t.Fatal(err)
	}
	if due != now.Add(time.Hour) {
		t.Fatalf("R/PT1H due=%v want %v", due, now.Add(time.Hour))
	}

	c, err = ParseISO8601Cycle("R/2000-01-01T00:00:00Z/PT1H")
	if err != nil {
		t.Fatal(err)
	}
	due, err = c.FirstDue(now)
	if err != nil {
		t.Fatal(err)
	}
	if due != now {
		t.Fatalf("aligned start due=%v want %v", due, now)
	}

	c, err = ParseISO8601Cycle("R/2099-01-01T00:00:00Z/PT1H")
	if err != nil {
		t.Fatal(err)
	}
	due, err = c.FirstDue(now)
	if err != nil {
		t.Fatal(err)
	}
	if due != time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("future start due=%v", due)
	}

	c, err = ParseISO8601Cycle("R/PT1H/2000-01-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.FirstDue(now); err == nil {
		t.Fatal("expected no occurrence after end")
	}
}
