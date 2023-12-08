package main

import (
	_ "embed"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var days = []func(string) (int, int){
	day1,
	day2,
	day3,
	day4,
	day5,
	day6,
	day7,
}

func main() {
	dayChosen := -1
	switch len(os.Args) {
	case 1:
		dayChosen = len(days)
		fmt.Printf("no day specified, choosing day %d\n", dayChosen)
	case 2:
		n, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil || n < 1 || n > 25 {
			fmt.Println("invalid day specified")
			os.Exit(1)
		}
		dayChosen = int(n)
		if dayChosen > len(days) {
			fmt.Printf("day %d not implemented yet\n", dayChosen)
			os.Exit(1)
		}
	default:
		fmt.Printf("usage: %s <day #>\n", os.Args[1])
		os.Exit(1)
	}

	//Read input file
	fname := fmt.Sprintf("input/day%d.txt", dayChosen)
	raw, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println("could not open input file:", err)
		os.Exit(1)
	}

	//Call day function
	start := time.Now()
	part1, part2 := days[dayChosen-1](string(raw))
	end := time.Now()
	fmt.Println("part 1:", part1)
	fmt.Println("part 2:", part2)
	fmt.Println("time taken:", end.UnixMicro()-start.UnixMicro(), "μs")
}

func day7(input string) (int, int) {
	type Hand struct {
		bet       int
		strength1 int
		strength2 int
	}
	labels := [2][]int{
		{'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'J': 11, 'Q': 12, 'K': 13, 'A': 14},
		{'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10, 'J': 01, 'Q': 12, 'K': 13, 'A': 14},
	}

	findKind := func(matches [16]int) int {
		m3 := false
		m2 := false
		for _, v := range matches {
			if v == 5 {
				return 7
			} else if v == 4 {
				return 6
			} else if v == 3 {
				m3 = true
			} else if v == 2 {
				if m2 {
					return 3
				}
				m2 = true
			}
		}
		if m2 && m3 {
			return 5
		} else if m3 {
			return 4
		} else if m2 {
			return 2
		} else {
			return 1
		}
	}

	//Construct cards
	lines := strings.Split(input, "\r\n")
	hands := make([]*Hand, len(lines))
	for i, line := range lines {
		hands[i] = new(Hand)
		matches := [16]int{}
		parts := strings.Split(line, " ")
		for j, char := range parts[0] {
			matches[labels[0][char]]++
			hands[i].strength1 += labels[0][char] << ((4 - j) * 4)
			hands[i].strength2 += labels[1][char] << ((4 - j) * 4)
		}
		hands[i].bet = quickconv(parts[1])
		hands[i].strength1 += findKind(matches) << 20
		maxCards := 11
		for i, v := range matches {
			if (v >= matches[maxCards] && i != 11) || (maxCards == 11 && v > 0) {
				maxCards = i
			}
		}
		jokers := matches[11]
		matches[11] = 0
		matches[maxCards] += jokers
		hands[i].strength2 += findKind(matches) << 20
	}

	//Sort cards
	handNum := len(hands)
	sort1 := make([]*Hand, handNum)
	sort2 := make([]*Hand, handNum)
	copy(sort1, hands)
	copy(sort2, hands)
	for i := 0; i < handNum-1; i++ {
		for j := 0; j < handNum-1-i; j++ {
			if sort1[j].strength1 < sort1[j+1].strength1 {
				temp := sort1[j]
				sort1[j] = sort1[j+1]
				sort1[j+1] = temp
			}
			if sort2[j].strength2 < sort2[j+1].strength2 {
				temp := sort2[j]
				sort2[j] = sort2[j+1]
				sort2[j+1] = temp
			}
		}
	}
	sum1 := 0
	sum2 := 0
	for i := range hands {
		sum1 += (handNum - i) * sort1[i].bet
		sum2 += (handNum - i) * sort2[i].bet
	}

	return sum1, sum2
}

func day6(input string) (int, int) {
	lines := strings.Split(input, "\r\n")
	for i := range lines {
		lines[i] = strings.Split(lines[i], ":")[1]
	}
	times := make([]int, 0)
	distances := make([]int, 0)
	for _, v := range strings.Split(lines[0], " ") {
		if v != "" {
			times = append(times, quickconv(v))
		}
	}
	for _, v := range strings.Split(lines[1], " ") {
		if v != "" {
			distances = append(distances, quickconv(v))
		}
	}
	sum1 := 1
	for i, record := range distances {
		recordBreaks := 0
		timeLimit := times[i]
		for time := 1; time < timeLimit; time++ {
			distance := time * (timeLimit - time)
			if distance > record {
				recordBreaks++
			}
		}
		sum1 *= recordBreaks
	}

	timeLimit := quickconv(strings.ReplaceAll(lines[0], " ", ""))
	record := quickconv(strings.ReplaceAll(lines[1], " ", ""))
	sum2 := 0
	for time := 1; time < timeLimit; time++ {
		distance := time * (timeLimit - time)
		if distance > record {
			sum2++
		}
	}

	return sum1, sum2
}

func day5(input string) (int, int) {
	cats := strings.Split(input, "\r\n\r\n")
	readMap := func(str string) [][3]int {
		lines := strings.Split(str, "\r\n")[1:]
		output := make([][3]int, len(lines))
		for i, line := range lines {
			for j, num := range strings.Split(line, " ") {
				output[i][j] = quickconv(num)
			}
		}
		return output
	}
	seedsStr := strings.Split(cats[0], " ")[1:]
	seeds := make([]int, len(seedsStr))
	for i, v := range seedsStr {
		seeds[i] = quickconv(v)
	}
	maps := [7][][3]int{}
	for i := range maps {
		maps[i] = readMap(cats[i+1])
	}

	findInMap := func(val int, xtoy [][3]int) (int, int) {
		next := uint(0)
		next -= 1
		for _, v := range xtoy {
			if val >= v[1] && val < v[1]+v[2] {
				diff := uint((v[1] + v[2]) - val)
				return v[0] + (val - v[1]), int(diff)
			} else if v[1] > val {
				diff := uint(v[1] - val)
				if diff < next {
					next = diff
				}
			}
		}
		return val, int(next)
	}

	findLocation := func(seed int) (int, int) {
		diffs := [7]int{}
		for i, v := range maps {
			seed, diffs[i] = findInMap(seed, v)
		}
		diff := -1
		for _, v := range diffs {
			if (v < diff || diff == -1) && v > 0 {
				diff = v
			}
		}
		return seed, diff
	}

	sum1 := uint(0)
	sum1 -= 1
	for _, seed := range seeds {
		location, _ := findLocation(seed)
		if uint(location) < sum1 {
			sum1 = uint(location)
		}
	}

	sum2 := uint(0)
	sum2 -= 1
	for i := 0; i < len(seeds); i += 2 {
		seedInit := seeds[i]
		seedRange := seeds[i+1]
		seedSkip := 0
		for seed := seedInit; seed < seedInit+seedRange; seed += seedSkip {
			var location int
			location, seedSkip = findLocation(seed)
			if uint(location) < sum2 {
				sum2 = uint(location)
			}
			if seedSkip <= 0 {
				break
			}
		}
	}

	return int(sum1), int(sum2)
}

func day4(input string) (int, int) {
	lines := strings.Split(input, "\r\n")
	repeats := make([]int, len(lines)+5)
	sum1 := 0
	sum2 := 0
	for i, line := range lines {
		repeats[i+1]++
		points := 0
		nums := strings.Split(strings.Split(line, ":")[1], "|")
		winners := make([]int, 0, 10)
		for _, v := range strings.Split(nums[0], " ") {
			if v != "" {
				winners = append(winners, quickconv(strings.Trim(v, " ")))
			}
		}
		winCount := 0
		for _, v := range strings.Split(nums[1], " ") {
			if v != "" {
				n := quickconv(strings.Trim(v, " "))
				if slices.Contains(winners, n) {
					winCount++
					if points == 0 {
						points = 1
					} else {
						points *= 2
					}
				}
			}
		}
		sum1 += points
		sum2 += repeats[i+1]
		for j := 0; j < winCount; j++ {
			index := i + j + 2
			if index >= len(repeats) {
				break
			}
			repeats[index] += repeats[i+1]
		}
	}

	return sum1, sum2
}

func day3(input string) (int, int) {
	type Symbol struct {
		val rune
		row int
		col int
	}
	symbols := make([]Symbol, 0)
	type PartNumber struct {
		number int
		row    int
		col    int
		length int
	}
	numbers := make([]PartNumber, 0)
	addNumber := func(row int, col int, str string) {
		n, _ := strconv.ParseInt(str, 10, 64)
		numbers = append(numbers, PartNumber{
			number: int(n),
			row:    row,
			col:    col,
			length: len(str),
		})
	}

	lines := strings.Split(input, "\r\n")
	for i, line := range lines {
		strBuf := ""
		numStartCol := 0
		for j, char := range line {
			if char < '0' || char > '9' {
				if strBuf != "" {
					addNumber(i+1, numStartCol, strBuf)
					strBuf = ""
				}
				if char != '.' {
					symbols = append(symbols, Symbol{
						row: i + 1,
						col: j + 1,
						val: char,
					})
				}
			} else {
				if strBuf == "" {
					numStartCol = j + 1
				}
				strBuf += string(char)
				if j == len(line)-1 {
					addNumber(i+1, numStartCol, strBuf)
					strBuf = ""
				}
			}
		}
	}

	sum1 := 0
	for _, n := range numbers {
		for _, s := range symbols {
			if s.row >= n.row-1 &&
				s.row <= n.row+1 &&
				s.col >= n.col-1 &&
				s.col <= n.col+n.length {
				sum1 += n.number
				break
			}
		}
	}

	sum2 := 0
	for _, s := range symbols {
		if s.val != '*' {
			continue
		}
		gnums := make([]int, 0, 4)
		for _, n := range numbers {
			if s.row >= n.row-1 &&
				s.row <= n.row+1 &&
				s.col >= n.col-1 &&
				s.col <= n.col+n.length {
				gnums = append(gnums, n.number)
			}
		}
		if len(gnums) == 2 {
			sum2 += gnums[0] * gnums[1]
		}
	}

	return sum1, sum2
}

func day2(input string) (int, int) {
	lines := strings.Split(input, "\r\n")
	sum1 := 0
	sum2 := 0
	for i, line := range lines {
		maxRed := 0
		maxBlue := 0
		maxGreen := 0
		gameString, _ := strings.CutPrefix(line, fmt.Sprintf("Game %d: ", i+1))
		rounds := strings.Split(gameString, ";")
		for _, r := range rounds {
			balls := strings.Split(strings.Trim(r, " "), ",")
			for _, ball := range balls {
				ball = strings.Trim(ball, " ")
				if n, ok := strings.CutSuffix(ball, " red"); ok {
					num, _ := strconv.ParseInt(n, 10, 64)
					if int(num) > maxRed {
						maxRed = int(num)
					}
				} else if n, ok := strings.CutSuffix(ball, " green"); ok {
					num, _ := strconv.ParseInt(n, 10, 64)
					if int(num) > maxGreen {
						maxGreen = int(num)
					}
				} else if n, ok := strings.CutSuffix(ball, " blue"); ok {
					num, _ := strconv.ParseInt(n, 10, 64)
					if int(num) > maxBlue {
						maxBlue = int(num)
					}
				} else {
					fmt.Println("didnt work :U")
				}
			}
		}

		if maxRed <= 12 && maxGreen <= 13 && maxBlue <= 14 {
			sum1 += i + 1
		}
		sum2 += maxRed * maxGreen * maxBlue
	}
	return sum1, sum2
}

func day1(input string) (int, int) {
	numLookup := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
	}

	lines := strings.Split(input, "\n")
	sum1 := 0
	sum2 := 0
	for _, v := range lines {
		sorted1 := ""
		sorted2 := ""
		for i, r := range v {
			if r >= '0' && r <= '9' {
				sorted1 += string(r)
				sorted2 += string(r)
			} else {
				for j, n := range numLookup {
					if len(v)-i < len(n) {
						continue
					}
					if v[i:i+len(n)] == n {
						sorted2 += fmt.Sprint(j + 1)
					}
				}
			}
		}
		shorten := func(str string) string {
			switch len(str) {
			case 0:
				return "0"
			case 1:
				return str + str
			case 2:
				return str
			default:
				return string(str[0]) + string(str[len(str)-1])
			}
		}

		num, _ := strconv.ParseInt(shorten(sorted1), 10, 64)
		sum1 += int(num)
		num, _ = strconv.ParseInt(shorten(sorted2), 10, 64)
		sum2 += int(num)
	}

	return sum1, sum2
}

func quickconv(str string) int {
	n, _ := strconv.ParseInt(str, 10, 64)
	return int(n)
}

func ternary[T any](condition bool, truthy T, falsy T) T {
	if condition {
		return truthy
	} else {
		return falsy
	}
}
