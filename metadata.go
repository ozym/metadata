package metadata

import (
	"fmt"
	"strings"
)

const DateTimeFormat = "2006-01-02 15:04:05"

type Keys []string

func (k Keys) Len() int           { return len(k) }
func (k Keys) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
func (k Keys) Less(i, j int) bool { return k[i] < k[j] }
func (k Keys) get(i int) string   { return k[i] }

// simple debugging helper function
func SimpleDiff(s1, s2 string) string {

	l1 := strings.Split(
		strings.Replace(s1, "\t", "  ", -1), "\n",
	)
	l2 := strings.Split(
		strings.Replace(s2, "\t", "  ", -1), "\n",
	)

	var w [2]int
	for _, l := range l1 {
		if len(l) > w[0] {
			w[0] = len(l)
		}
	}
	for _, l := range l2 {
		if len(l) > w[1] {
			w[1] = len(l)
		}
	}

	var n int
	if len(l1) > len(l2) {
		n = len(fmt.Sprintf("%d", len(l1)))
	} else {
		n = len(fmt.Sprintf("%d", len(l2)))
	}

	var s []string
	for i := 0; i < len(l1) && i < len(l2); i++ {
		if l1[i] == l2[i] {
			continue
		}
		s = append(s, fmt.Sprintf(fmt.Sprintf("\t[%%%dd]!!! %%-%ds ! %%-%ds !!!", n, w[0], w[1]), i, l1[i], l2[i]))
	}
	for i := len(l2); i < len(l1); i++ {
		s = append(s, fmt.Sprintf(fmt.Sprintf("\t[%%%dd]+++ %%-%ds + %%-%ds +++", n, w[0], w[1]), i, l1[i], ""))
	}
	for i := len(l1); i < len(l2); i++ {
		s = append(s, fmt.Sprintf(fmt.Sprintf("\t[%%%dd]--- %%-%ds - %%-%ds ---", n, w[0], w[1]), i, "", l2[i]))
	}

	return strings.Join(s, "\n")
}
