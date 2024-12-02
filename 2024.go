package main

import (
	"slices"
	"strings"
)

var y2024 = []func(string) (int, int){
	y2024day1,
}

func y2024day1(input string) (int, int) {
	lines := strings.Split(input, "\n")
	row1 := make([]int, len(lines))
	row2 := make([]int, len(lines))
	for i, v := range lines {
		strs := strings.Split(v, "   ")
		row1[i] = quickconv(strs[0])
		row2[i] = quickconv(strs[1])
	}

	slices.Sort(row1)
	slices.Sort(row2)
	sum1 := 0
	sum2 := 0
	for i := range row1 {
		diff := abs(row1[i] - row2[i])
		sum1 += diff

		idxOf := slices.Index(row2, row1[i])
		if idxOf != -1 {
			for j := idxOf; j < len(row2) && row2[j] == row1[i]; j++ {
				sum2 += row1[i]
			}
		}
	}

	return sum1, sum2
}
