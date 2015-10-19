package metadata

import (
	"fmt"
	"strings"
)

// simple debugging helper function
func Diff(s1, s2 fmt.Stringer) string {
	var h, w1, w2 int

	l1 := strings.Split(
		strings.Replace(s1.String(), "\t", "  ", -1), "\n",
	)
	l2 := strings.Split(
		strings.Replace(s2.String(), "\t", "  ", -1), "\n",
	)

	if len(l1) > len(l2) {
		h = len(l1)
	} else {
		h = len(l2)
	}
	for _, l := range l1 {
		if len(l) > w1 {
			w1 = len(l)
		}
	}
	for _, l := range l2 {
		if len(l) > w2 {
			w2 = len(l)
		}
	}

	var s []string
	for i := 0; i < h; i++ {
		switch {
		case i < len(l1) && i < len(l2):
			if l1[i] != l2[i] {
				s = append(s, fmt.Sprintf(fmt.Sprintf("\t!!! %%-%ds ! %%-%ds !!!", w1, w2), l1[i], l2[i]))
			}
		case i < len(l1):
			s = append(s, fmt.Sprintf(fmt.Sprintf("\t+++ %%-%ds + %%-%ds +++", w1, w2), l1[i], ""))
		case i < len(l2):
			s = append(s, fmt.Sprintf(fmt.Sprintf("\t--- %%-%ds - %%-%ds ---", w1, w2), "", l2[i]))
		}
	}

	return strings.Join(s, "\n")
}
