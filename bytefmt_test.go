package bytefmt

import (
	"fmt"
	"testing"
)

func ExampleParseCLI() {
	fmt.Println(ParseCLI("2.5k"))
	// Output: 2560 <nil>
}

func TestParseCLI(t *testing.T) {
	tests := []struct {
		in  string
		out ByteSize
		err bool
	}{
		{"", 0, true},
		{" ", 0, true},
		{"inf", 0, true},
		{"NaN", 0, true},
		{"infk", 0, true},
		{"nanMB", 0, true},
		{"foo", 0, true},
		{"BARGB", 0, true},
		{" 0 ", 0, true},

		{"-1.5", -1, false},
		{"-1", -1, false},
		{"0", 0, false},
		{"1", 1, false},
		{"1.5", 1, false},
		{"1500", 1500, false},

		{"1Kb", 1e3, false},
		{"1.5kb", 1.5e3, false},
		{"1mb", 1e6, false},
		{"1GB", 1e9, false},
		{"1tB", 1e12, false},
		{"1pb", 1e15, false},
		{"1eb", 1e18, false},
		{"5eb", 5e18, false},

		{"1k", 1024, false},
		{"1k", 1024, false},
		{"1.5k", 1536, false},
		{"8k", 8192, false},
		{"1m", 1048576, false},
		{"1M", 1048576, false},
		{"1g", 1 << 30, false},
		{"1G", 1 << 30, false},
		{"1t", 1 << 40, false},
		{"1T", 1 << 40, false},
		{"1p", 1 << 50, false},
		{"1P", 1 << 50, false},
		{"1e", 1 << 60, false},
		{"1E", 1 << 60, false},
	}

	for _, tt := range tests {
		out, err := ParseCLI(tt.in)
		if out != tt.out {
			t.Errorf("ParseCLI(%q) => %g, expected %g", tt.in, float64(out), float64(tt.out))
		}
		if err != nil && !tt.err {
			t.Errorf("ParseCLI(%q) => err, expected nil", tt.in)
			t.Log("Error:", err)
		} else if err == nil && tt.err {
			t.Errorf("ParseCLI(%q) => nil, expected err", tt.in)
		}
	}
}
