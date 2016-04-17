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
	Byte unit = 1 << (10 * iota)

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

var unitSymbols = map[unit]string{
	Byte: "B",

	Kibibyte: "KiB",
	Mebibyte: "MiB",
	Gibibyte: "GiB",
	Tebibyte: "TiB",
	Pebibyte: "PiB",
	Exbibyte: "EiB",

	Kilobyte: "KB",
	Megabyte: "MB",
	Gigabyte: "GB",
	Terabyte: "TB",
	Petabyte: "PB",
	Exabyte:  "EB",
}

var cliSuffixes = map[string]unit{
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

// ByteSize represents a quantity in bytes.
type ByteSize int64

// Format returns a human-friendly string with a binary prefix.
func (s ByteSize) Format() string {
	return s.format(Kibibyte)
}

// FormatSI returns a human-friendly string with an SI prefix.
func (s ByteSize) FormatSI() string {
	return s.format(Kilobyte)
}

func (s ByteSize) format(un unit) string {
	ss := unit(s)
	u := Byte
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
