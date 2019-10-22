package sort

import (
	"sort"
	"strconv"
	"strings"
)

type SectionSlice []string

func (s SectionSlice) Len() int {
	return len(s)
}

func (s SectionSlice) Less(i, j int) bool {

	var a, b []int

	for _, e := range strings.Split(s[i], ".") {
		n, _ := strconv.Atoi(e)
		a = append(a, n)
	}

	for _, e := range strings.Split(s[j], ".") {
		n, _ := strconv.Atoi(e)
		b = append(b, n)
	}

	for idx := 0; ; idx++ {

		if a[idx] < b[idx] {
			return true
		}

		if a[idx] > b[idx] {
			return false
		}

		var (
			la = len(a)-1 == idx
			lb = len(b)-1 == idx
		)

		// we cannot continue, left or right are the last element
		if la || lb {
			return len(a) < len(b)
		}

	}

}

func (s SectionSlice) Sort() {
	sort.Sort(s)
}

func (s SectionSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
