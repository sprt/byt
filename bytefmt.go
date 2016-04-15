// Package bytefmt provides utility functions to format and parse byte quantities.
package bytefmt

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ByteSize represents a quantity in bytes.
type ByteSize int64

const (
	_ = 1 << (10 * iota)

	Kibibyte
	Mebibyte
	Gibibyte
	Tebibyte
	Pebibyte
	Exbibyte

	Kilobyte = 1e3
	Megabyte = 1e6
	Gigabyte = 1e9
	Terabyte = 1e12
	Petabyte = 1e15
	Exabyte  = 1e18
)

var cliSuffixes = map[string]int64{
	"k": Kibibyte,
	"m": Mebibyte,
	"g": Gibibyte,
	"t": Tebibyte,
	"p": Pebibyte,
	"e": Exbibyte,

	"kb": Kilobyte,
	"mb": Megabyte,
	"gb": Gigabyte,
	"tb": Terabyte,
	"pb": Petabyte,
	"eb": Exabyte,
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
