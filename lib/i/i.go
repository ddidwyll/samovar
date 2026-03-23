package i

import "slices"

func N(s string, ss ...string) bool {
	return slices.Contains(ss, s)
}
