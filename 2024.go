package main

import (
	"regexp"
	"slices"
	"strings"
)

var y2024 = []func(string) (int, int){
	y2024day1,
	y2024day2,
	y2024day3,
	y2024day4,
}

func y2024day4(input string) (sum1 int, sum2 int) {
	rawLines := strings.Split(input, "\n")
	lines := make([][]rune, len(rawLines))
	for i := range rawLines {
		lines[i] = []rune(rawLines[i])
	}

	target := [...]rune{'X', 'M', 'A', 'S'}
	directions := [...][2]int{
		{0, 1},
		{0, -1},
		{1, 1},
		{1, 0},
		{1, -1},
		{-1, 1},
		{-1, 0},
		{-1, -1},
	}
	for j := 0; j < len(lines); j++ {
		for i := 0; i < len(lines[j]); i++ {
		loopDay1:
			for d := range directions {
				for t := range target {
					x := uint(i + directions[d][0]*t)
					y := uint(j + directions[d][1]*t)
					if x >= uint(len(lines[0])) || y >= uint(len(lines)) || lines[y][x] != target[t] {
						continue loopDay1
					}
				}

				sum1 += 1
			}

			if i != 0 && j != 0 && i != len(lines[j])-1 && j != len(lines)-1 && lines[j][i] == 'A' {
				corners := [...]rune{
					lines[j+1][i+1],
					lines[j-1][i-1],
					lines[j-1][i+1],
					lines[j+1][i-1],
				}
				if (corners[0] == 'M' || corners[0] == 'S') &&
					(corners[1] == 'M' || corners[1] == 'S') &&
					corners[0] != corners[1] &&
					(corners[2] == 'M' || corners[2] == 'S') &&
					(corners[3] == 'M' || corners[3] == 'S') &&
					corners[2] != corners[3] {
					sum2 += 1
				}
			}
		}
	}

	return
}

func y2024day3(input string) (int, int) {
	expDo := regexp.MustCompile(`do\(\)`)
	expDont := regexp.MustCompile(`don't\(\)`)
	expMul := regexp.MustCompile(`mul\(\d+,\d+\)`)
	expNum := regexp.MustCompile(`\d+`)
	sum1 := 0
	for _, v := range expMul.FindAllString(input, -1) {
		found := expNum.FindAllString(v, 2)
		sum1 += quickconv(found[0]) * quickconv(found[1])
	}

	sum2 := 0
	idx := 0
	enabled := true
	for {
		if enabled {
			mulIdx := expMul.FindStringIndex(input[idx:])
			dontIdx := expDont.FindStringIndex(input[idx:])
			if mulIdx == nil && dontIdx == nil {
				break
			}

			if (mulIdx == nil) || (dontIdx != nil && dontIdx[0] < mulIdx[0]) {
				enabled = false
				idx += dontIdx[1]
			} else {
				found := expNum.FindAllString(input[idx+mulIdx[0]:idx+mulIdx[1]], 2)
				sum2 += quickconv(found[0]) * quickconv(found[1])
				idx += mulIdx[1]
			}
		} else {
			doIdx := expDo.FindStringIndex(input[idx:])
			if nil == doIdx {
				break
			}
			enabled = true
			idx += doIdx[1]
		}
	}

	return sum1, sum2
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
