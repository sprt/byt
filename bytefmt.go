// Package bytefmt provides utility functions to format and parse byte quantities.
package bytefmt

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type unit int64

const (
	byt unit = 1 << (10 * iota)

	kibibyte
	mebibyte
	gibibyte
	tebibyte
	pebibyte
	exbibyte

	kilobyte = 1e3
	megabyte = 1e6
	gigabyte = 1e9
	terabyte = 1e12
	petabyte = 1e15
	exabyte  = 1e18
)

var unitSymbols = map[unit]string{
	byt: "B",

	kibibyte: "KiB",
	mebibyte: "MiB",
	gibibyte: "GiB",
	tebibyte: "TiB",
	pebibyte: "PiB",
	exbibyte: "EiB",

	kilobyte: "KB",
	megabyte: "MB",
	gigabyte: "GB",
	terabyte: "TB",
	petabyte: "PB",
	exabyte:  "EB",
}

var cliSuffixes = map[string]unit{
	"k": kibibyte,
	"m": mebibyte,
	"g": gibibyte,
	"t": tebibyte,
	"p": pebibyte,
	"e": exbibyte,

	"kb": kilobyte,
	"mb": megabyte,
	"gb": gigabyte,
	"tb": terabyte,
	"pb": petabyte,
	"eb": exabyte,
}

// ByteSize represents a quantity in bytes.
type ByteSize int64

// Format returns a human-friendly string with a binary prefix.
func (s ByteSize) Format() string {
	return s.format(kibibyte)
}

// FormatSI returns a human-friendly string with an SI prefix.
func (s ByteSize) FormatSI() string {
	return s.format(kilobyte)
}

func (s ByteSize) format(un unit) string {
	ss := unit(s)
	u := byt
	for ss >= un {
		ss /= un
		u *= un
	}
	return fmt.Sprintf("%.1f %s", float64(s)/float64(u), unitSymbols[u])
}

// ParseCLI parses s and returns the corresponding size in bytes.
// s is a number followed by a unit (optional), no whitespace allowed.
// Units are K,M,G,T,P,E,Z,Y (powers of 1024) or KB,MB,... (powers of 1000).
// ParseCLI is not case sensitive.
func ParseCLI(s string) (ByteSize, error) {
	s = strings.ToLower(s)

	for suffix, size := range cliSuffixes {
		if strings.HasSuffix(s, suffix) {
			x, err := parseFloat(strings.TrimSuffix(s, suffix))
			if err != nil {
				return 0, err
			}
			return ByteSize(x * float64(size)), nil
		}
	}

	x, err := parseFloat(s)
	if err != nil {
		return 0, err
	}

	return ByteSize(x), nil
}

func parseFloat(s string) (float64, error) {
	x, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsInf(x, 0) || math.IsNaN(x) {
		return 0, fmt.Errorf("cannot parse %q", s)
	}
	return x, nil
}
