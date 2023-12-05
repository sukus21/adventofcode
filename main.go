package main

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var days = []func(){
	day1,
	day2,
	day3,
	nil,
	day5,
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
		if dayChosen >= len(days) {
			fmt.Println("day not implemented yet")
			os.Exit(1)
		}
	default:
		fmt.Printf("usage: %s <day #>\n", os.Args[1])
		os.Exit(1)
	}

	//Call day function
	days[dayChosen-1]()
}

//go:embed day5.txt
var inputDay5 string

func day5() {
	cats := strings.Split(inputDay5, "\r\n\r\n")
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

	seedToSoil := readMap(cats[1])
	soilToFert := readMap(cats[2])
	fertToWater := readMap(cats[3])
	waterToLight := readMap(cats[4])
	lightToTemp := readMap(cats[5])
	tempToHumid := readMap(cats[6])
	humidToLoc := readMap(cats[7])

	findInMap := func(val int, xtoy [][3]int) int {
		for _, v := range xtoy {
			if val >= v[1] && val < v[1]+v[2] {
				diff := val - v[1]
				return v[0] + diff
			}
		}
		return val
	}

	sum1 := uint(0)
	sum1 -= 1
	for _, seed := range seeds {
		soil := findInMap(seed, seedToSoil)
		fert := findInMap(soil, soilToFert)
		water := findInMap(fert, fertToWater)
		light := findInMap(water, waterToLight)
		temp := findInMap(light, lightToTemp)
		humid := findInMap(temp, tempToHumid)
		location := findInMap(humid, humidToLoc)
		if uint(location) < sum1 {
			sum1 = uint(location)
		}
	}
	fmt.Println("part 1:", sum1)

	sumChan := make(chan uint)
	for i := 0; i < len(seeds); i += 2 {
		go func(seedInit int, seedRange int) {
			chanSum := uint(0)
			chanSum -= 1
			for seed := seedInit; seed < seedInit+seedRange; seed++ {
				soil := findInMap(seed, seedToSoil)
				fert := findInMap(soil, soilToFert)
				water := findInMap(fert, fertToWater)
				light := findInMap(water, waterToLight)
				temp := findInMap(light, lightToTemp)
				humid := findInMap(temp, tempToHumid)
				location := findInMap(humid, humidToLoc)
				if uint(location) < chanSum {
					chanSum = uint(location)
				}
			}
			sumChan <- chanSum
		}(seeds[i], seeds[i+1])
	}

	sum2 := uint(0)
	sum2 -= 1
	for i := 0; i < len(seeds)/2; i++ {
		newSum := int(<-sumChan)
		if uint(newSum) < sum2 {
			sum2 = uint(newSum)
		}
		fmt.Println("finished set", i+1, "/", len(seeds)/2)
	}

	fmt.Println("part 2:", sum2)
}

//go:embed day3.txt
var inputDay3 string

func day3() {
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

	lines := strings.Split(inputDay3, "\r\n")
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

	fmt.Println("part 1:", sum1)
	fmt.Println("part 2:", sum2)
}

//go:embed day2.txt
var inputDay2 string

func day2() {
	lines := strings.Split(inputDay2, "\r\n")
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
	fmt.Printf("Part 1: %d\nPart 2: %d\n", sum1, sum2)
}

//go:embed day1.txt
var inputDay1 string

func day1() {
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

	lines := strings.Split(inputDay1, "\n")
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

	fmt.Printf("Part 1: %d\nPart 2: %d\n", sum1, sum2)
}

func quickconv(str string) int {
	n, _ := strconv.ParseInt(str, 10, 64)
	return int(n)
}
