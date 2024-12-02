package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var years = map[int][]func(string) (int, int){
	2023: y2023,
}

func main() {
	if len(os.Args) > 3 {
		fmt.Printf("usage: %s <day #>\n", os.Args[1])
		os.Exit(1)
	}

	yearChosen := -1
	dayChosen := -1

	// Get chosen year
	if len(os.Args) < 2 {
		for y := range years {
			if y > yearChosen {
				yearChosen = y
			}
		}
		fmt.Printf("no year specified, choosing year %d\n", yearChosen)
	} else {
		yearChosen = quickconv(os.Args[1])
		if years[yearChosen] == nil {
			fmt.Printf("invalid year specified")
			os.Exit(1)
		}
	}
	yearDays := years[yearChosen]

	// Get chosen day
	if len(os.Args) < 3 {
		dayChosen = len(yearDays)
		fmt.Printf("no day specified, choosing day %d\n", dayChosen)
	} else {
		dayChosen = quickconv(os.Args[2])
		if dayChosen < 1 || dayChosen > 25 {
			fmt.Println("invalid day specified")
			os.Exit(1)
		}
		if dayChosen > len(yearDays) {
			fmt.Printf("day %d not implemented yet\n", dayChosen)
			os.Exit(1)
		}
	}

	// Read input file
	fname := fmt.Sprintf("input/year%dday%d.txt", yearChosen, dayChosen)
	raw, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println("could not open input file:", err)
		os.Exit(1)
	}

	// Call day function
	start := time.Now()
	content := strings.ReplaceAll(string(raw), "\r\n", "\n")
	part1, part2 := yearDays[dayChosen-1](content)
	end := time.Now()
	fmt.Println("part 1:", part1)
	fmt.Println("part 2:", part2)
	fmt.Println("time taken:", end.UnixMicro()-start.UnixMicro(), "μs")
}
