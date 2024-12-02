package main

import (
	"slices"
	"strings"
)

var y2024 = []func(string) (int, int){
	y2024day1,
	y2024day2,
}

func y2024day2(input string) (int, int) {
	lines := strings.Split(input, "\n")
	reports := make([][]int, len(lines))
	for i, v := range lines {
		for _, n := range strings.Split(v, " ") {
			reports[i] = append(reports[i], quickconv(n))
		}
	}

	isValid := func(report []int) bool {
		prev := report[0]
		direction := 0
		for _, entry := range report[1:] {
			diff := entry - prev
			prev = entry
			if direction == 0 {
				direction = sign(diff)
			}
			if abs(diff) <= 0 || abs(diff) > 3 || sign(diff) != direction {
				return false
			}
		}

		return true
	}

	sum1 := 0
	sum2 := 0
	for _, report := range reports {
		if isValid(report) {
			sum1++
			sum2++
		} else {
			// Brute-force
			newSlice := make([]int, len(report)-1)
			for j := range report {
				newSlice = newSlice[:0]
				for k, v := range report {
					if k != j {
						newSlice = append(newSlice, v)
					}
				}

				if isValid(newSlice) {
					sum2++
					break
				}
			}
		}
	}

	return sum1, sum2
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
