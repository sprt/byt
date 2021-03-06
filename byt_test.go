package byt

import (
	"fmt"
	"testing"
)

func ExampleByteSize_Binary() {
	fmt.Printf("%.1f", Size(2560).Binary())
	// Output: 2.5 KiB
}

func TestByteSize_Binary(t *testing.T) {
	tests := []struct {
		in  int64
		out string
	}{
		{1, "1.0 B"},
		{1000, "1000.0 B"},
		{1023, "1023.0 B"},
		{1024, "1.0 KiB"},
		{1536, "1.5 KiB"},
		{1048576, "1.0 MiB"},
		{1 << 30, "1.0 GiB"},
		{1 << 40, "1.0 TiB"},
		{1 << 50, "1.0 PiB"},
		{1 << 60, "1.0 EiB"},
	}

	for _, tt := range tests {
		s := fmt.Sprintf("%.1f", Size(tt.in).Binary())
		if s != tt.out {
			t.Errorf("%d => %q, expected %q", tt.in, s, tt.out)
		}
	}
}

func ExampleByteSize_SI() {
	fmt.Printf("%.2f", Size(2560).Decimal())
	// Output: 2.56 KB
}

func TestFormatSI(t *testing.T) {
	tests := []struct {
		in  int64
		out string
	}{
		{1, "1.0 B"},
		{1e3, "1.0 KB"},
		{1500, "1.5 KB"},
		{1e6, "1.0 MB"},
		{1e9, "1.0 GB"},
		{1e12, "1.0 TB"},
		{1e15, "1.0 PB"},
		{1e18, "1.0 EB"},
	}

	for _, tt := range tests {
		s := fmt.Sprintf("%.1f", Size(tt.in).Decimal())
		if s != tt.out {
			t.Errorf("%d => %q, expected %q", tt.in, s, tt.out)
		}
	}
}

func ExampleParseCLI() {
	fmt.Print(parseCLI("2.5k"))
	// Output: 2560 <nil>
}

func TestParseCLI(t *testing.T) {
	tests := []struct {
		in  string
		out Size
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
		{"1K", 1024, false},
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
		out, err := parseCLI(tt.in)
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
