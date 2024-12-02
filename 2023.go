package main

import (
	"fmt"
	"math/bits"
	"slices"
	"strconv"
	"strings"
)

var y2023 = []func(string) (int, int){
	y2023day1,
	y2023day2,
	y2023day3,
	y2023day4,
	y2023day5,
	y2023day6,
	y2023day7,
	y2023day8,
	y2023day9,
	y2023day10,
	y2023day11,
	y2023day12,
	y2023day13,
}

func y2023day13(input string) (int, int) {
	mirrorLines := strings.Split(input, "\n\n")
	mirrors := make([][]uint, len(mirrorLines))
	rotated := make([][]uint, len(mirrorLines))
	for i, mline := range mirrorLines {
		lines := strings.Split(mline, "\n")
		mirror := make([]uint, len(lines))
		for j, line := range lines {
			for k, char := range line {
				mirror[j] |= uint(ternary(char == '#', 1, 0)) << k
			}
		}
		mirrors[i] = mirror
		rotate := make([]uint, len(lines[0]))
		for j := range rotate {
			for k := range lines {
				rotate[j] |= ((mirror[k] >> j) & 1) << k
			}
		}
		rotated[i] = rotate
	}

	getMirror := func(mirror []uint, off int) int {
		for i := range mirror {
			if i == 0 {
				continue
			}
			diff := 0
			low := i - 1
			high := i
			fail := false
			for low >= 0 && high < len(mirror) {
				diff += bits.OnesCount(mirror[low] ^ mirror[high])
				if diff > off {
					fail = true
					break
				}
				low--
				high++
			}
			if fail {
				continue
			} else if diff == off {
				return i
			}
		}

		return 0
	}

	sum1 := 0
	sum2 := 0
	for i := range mirrors {
		sum1 += getMirror(mirrors[i], 0) * 100
		sum1 += getMirror(rotated[i], 0)
		sum2 += getMirror(mirrors[i], 1) * 100
		sum2 += getMirror(rotated[i], 1)
	}

	return sum1, sum2
}

func y2023day12(input string) (int, int) {
	type springline struct {
		springs     []rune
		brokenCount []int
		missing     int
		unknowns    int
	}

	lines := strings.Split(input, "\n")
	springLines := make([]springline, len(lines))
	for i, v := range lines {
		spring := &springLines[i]
		split1 := strings.Split(v, " ")
		spring.springs = []rune(split1[0])

		knownBroken := 0
		for _, c := range spring.springs {
			if c == '#' {
				knownBroken++
			} else if c == '?' {
				spring.unknowns++
			}
		}

		totalBroken := 0
		brokenCounts := strings.Split(split1[1], ",")
		spring.brokenCount = make([]int, len(brokenCounts))
		for j, n := range brokenCounts {
			num := quickconv(n)
			spring.brokenCount[j] = num
			totalBroken += num
		}
		spring.missing = totalBroken - knownBroken
	}

	valid := func(s *springline, attempt uint) bool {
		currentBroken := 0
		brokenRow := 0

		validSequence := func() bool {
			if currentBroken != 0 {
				if brokenRow >= len(s.brokenCount) || s.brokenCount[brokenRow] != currentBroken {
					return false
				}
				currentBroken = 0
				brokenRow++
			}
			return true
		}

		for _, v := range s.springs {
			if v == '?' {
				v = ternary((attempt&1) != 0, '#', '.')
				attempt >>= 1
			}
			if v == '#' {
				currentBroken++
			} else if !validSequence() {
				return false
			}
		}
		return validSequence()
	}

	debugPrint := func(unknowns int, attempt uint) string {
		return fmt.Sprintf(fmt.Sprintf("%%0%db", unknowns), attempt)
	}
	_ = debugPrint

	var moveBit func(s *springline, attempt, upperbits uint, mask uint, bit int) int
	moveBit = func(s *springline, attempt, upperbits, mask uint, bit int) int {
		numValid := ternary(valid(s, attempt|upperbits), 1, 0)

		for bit >= 0 {
			attempt &= mask ^ (1 << bit)
			bit++
			attempt |= 1 << bit
			if attempt&mask != attempt {
				break
			}

			//Move other bits to the start
			naked := ((1 << bit) - 1) & attempt
			otherBits := bits.OnesCount(uint(naked))
			numValid += moveBit(s, naked, upperbits|(1<<bit), (1<<bit)-1, otherBits-1)
		}
		return numValid
	}

	sum1 := 0
	for _, line := range springLines {
		attempt := uint(0)
		for i := 0; i < line.missing; i++ {
			attempt <<= 1
			attempt |= 1
		}
		mask := uint(0)
		for i := 0; i < line.unknowns; i++ {
			mask <<= 1
			mask |= 1
		}
		sum1 += moveBit(&line, attempt, 0, mask, line.missing-1)
	}

	return sum1, 0
}

func y2023day11(input string) (int, int) {
	type vec2 struct{ x, y int }
	lines := strings.Split(input, "\n")
	galaxies := make([]vec2, 0, 256)
	spaceRow := make([]int, len(lines))
	spaceCol := make([]int, len(lines))
	for i := range spaceCol {
		spaceCol[i] = 1
	}
	for i, line := range lines {
		if !strings.ContainsRune(line, '#') {
			spaceRow[i] = 1
			continue
		}
		for j, char := range line {
			if char == '#' {
				galaxies = append(galaxies, vec2{j, i})
				spaceCol[j] = 0
			}
		}
	}

	//Expand empty space
	offsetX := make([]int, len(spaceRow))
	offsetY := make([]int, len(spaceCol))
	offsetTemp := 0
	for i, v := range spaceRow {
		offsetTemp += v
		offsetY[i] = offsetTemp
	}
	offsetTemp = 0
	for i, v := range spaceCol {
		offsetTemp += v
		offsetX[i] = offsetTemp
	}

	sum1 := 0
	sum2 := 0
	for i, v := range galaxies {
		for j := i; j < len(galaxies); j++ {
			g := galaxies[j]
			distanceGalaxies := abs(v.x-g.x) + abs(v.y-g.y)
			spaceExpansion := abs(offsetX[v.x]-offsetX[g.x]) + abs(offsetY[v.y]-offsetY[g.y])
			sum1 += distanceGalaxies + spaceExpansion
			sum2 += distanceGalaxies + spaceExpansion*999_999
		}
	}

	return sum1, sum2
}

func y2023day10(input string) (int, int) {
	type vec2 struct {
		x, y int
	}

	lines := strings.Split(input, "\n")
	chars := make([][]rune, len(lines))
	distances := make([][]int, len(lines))
	startPos := vec2{}
	for i := range lines {
		chars[i] = []rune(lines[i])
		distances[i] = make([]int, len(chars[i]))
		if pos := strings.IndexRune(lines[i], 'S'); pos != -1 {
			startPos.x = pos
			startPos.y = i
		}
	}

	move := func(pos vec2, dir vec2) (vec2, vec2, bool) {
		if pos.x >= len(chars[0]) || pos.x < 0 || pos.y >= len(chars) || pos.y < 0 {
			return vec2{}, vec2{}, false
		} else if chars[pos.y][pos.x] == 'S' {
			return pos, vec2{}, true
		}
		switch {
		case dir.x == 1 && dir.y == 0:
			switch chars[pos.y][pos.x] {
			case 'J':
				return vec2{pos.x, pos.y - 1}, vec2{0, -1}, true
			case '7':
				return vec2{pos.x, pos.y + 1}, vec2{0, 1}, true
			case '-':
				return vec2{pos.x + 1, pos.y}, vec2{1, 0}, true
			}

		case dir.x == -1 && dir.y == 0:
			switch chars[pos.y][pos.x] {
			case 'L':
				return vec2{pos.x, pos.y - 1}, vec2{0, -1}, true
			case 'F':
				return vec2{pos.x, pos.y + 1}, vec2{0, 1}, true
			case '-':
				return vec2{pos.x - 1, pos.y}, vec2{-1, 0}, true
			}

		case dir.x == 0 && dir.y == 1:
			switch chars[pos.y][pos.x] {
			case 'L':
				return vec2{pos.x + 1, pos.y}, vec2{1, 0}, true
			case 'J':
				return vec2{pos.x - 1, pos.y}, vec2{-1, 0}, true
			case '|':
				return vec2{pos.x, pos.y + 1}, vec2{0, 1}, true
			}

		case dir.x == 0 && dir.y == -1:
			switch chars[pos.y][pos.x] {
			case 'F':
				return vec2{pos.x + 1, pos.y}, vec2{1, 0}, true
			case '7':
				return vec2{pos.x - 1, pos.y}, vec2{-1, 0}, true
			case '|':
				return vec2{pos.x, pos.y - 1}, vec2{0, -1}, true
			}
		}
		return pos, vec2{}, false
	}

	painted := make([]vec2, 0, 256)
	sum1 := 0
	for _, diff := range []vec2{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
		pos := vec2{startPos.x + diff.x, startPos.y + diff.y}
		distance := 1
		var ok bool
		painted = painted[:0]
		for {
			var npos vec2
			npos, diff, ok = move(pos, diff)
			if !ok || chars[pos.y][pos.x] == 'S' {
				if !ok {
					for _, v := range painted {
						distances[v.y][v.x] = 0
					}
				}
				break
			}
			painted = append(painted, vec2{pos.x, pos.y})
			if distances[pos.y][pos.x] == 0 || distance < distances[pos.y][pos.x] {
				distances[pos.y][pos.x] = distance
			} else if distances[pos.y][pos.x] >= distance {
				sum1 = ternary(sum1 > distance, sum1, distance)
			}
			pos = npos
			distance++
		}
	}

	checkstart := func(x, y int) bool {
		x += startPos.x
		y += startPos.y
		if x > len(chars[0]) || x < 0 || y > len(chars) || y < 0 {
			return false
		}
		return distances[y][x] != 0
	}

	distances[startPos.y][startPos.x] = -1
	switch {
	case checkstart(1, 0) && checkstart(-1, 0):
		chars[startPos.y][startPos.x] = '-'
	case checkstart(0, 1) && checkstart(0, -1):
		chars[startPos.y][startPos.x] = '|'
	case checkstart(1, 0) && checkstart(0, 1):
		chars[startPos.y][startPos.x] = 'F'
	case checkstart(-1, 0) && checkstart(0, 1):
		chars[startPos.y][startPos.x] = '7'
	case checkstart(1, 0) && checkstart(0, -1):
		chars[startPos.y][startPos.x] = 'J'
	case checkstart(-1, 0) && checkstart(0, -1):
		chars[startPos.y][startPos.x] = 'L'
	}

	sum2 := 0
	inside := false
	switchOn := ' '
	for i := range distances {
		for j := range distances[i] {
			char := chars[i][j]
			distance := distances[i][j]
			if distance != 0 {
				if char == 'F' {
					switchOn = 'J'
				} else if char == 'L' {
					switchOn = '7'
				}
				if char == '|' || char == 'J' || char == '7' {
					if char == '|' || char == switchOn {
						inside = !inside
					}
				}
			} else if inside {
				sum2++
			}
		}
	}

	return sum1, sum2
}

func y2023day9(input string) (int, int) {
	lines := strings.Split(input, "\n")
	sequences := make([][][]int, len(lines))
	for i, line := range lines {
		nums := strings.Split(line, " ")
		sequences[i] = make([][]int, 1, 8)
		sequences[i][0] = make([]int, len(nums))
		for j, num := range nums {
			sequences[i][0][j] = quickconv(num)
		}
	}

	all0s := func(sequence []int) bool {
		for _, v := range sequence {
			if v != 0 {
				return false
			}
		}
		return true
	}

	for i := range sequences {
		previous := sequences[i][0]
		for !all0s(previous) {
			next := make([]int, len(previous)-1)
			sequences[i] = append(sequences[i], next)
			for j := range next {
				next[j] = previous[j+1] - previous[j]
			}
			previous = next
		}
	}

	predictNext := func(sequence [][]int) int {
		diff := 0
		for i := len(sequence) - 2; i >= 0; i-- {
			diff = sequence[i][len(sequence[i])-1] + diff
		}
		return diff
	}

	predictPrevious := func(sequence [][]int) int {
		diff := 0
		for i := len(sequence) - 2; i >= 0; i-- {
			diff = sequence[i][0] - diff
		}
		return diff
	}

	sum1 := 0
	sum2 := 0
	for _, v := range sequences {
		sum1 += predictNext(v)
		sum2 += predictPrevious(v)
	}

	return sum1, sum2
}

func y2023day8(input string) (int, int) {
	lines := strings.Split(input, "\n")
	turns := []rune(lines[0])
	paths := make(map[string][2]string)
	nodes := make([]string, 0)
	for _, v := range lines[2:] {
		pathName := v[0:3]
		paths[pathName] = [2]string{v[7:10], v[12:15]}
		if pathName[2] == 'A' {
			nodes = append(nodes, pathName)
		}
	}

	//Find exit to path AAA
	node := "AAA"
	instruction := 0
	sum1 := 0
	for node != "ZZZ" {
		sum1++
		node = paths[node][ternary(turns[instruction] == 'R', 1, 0)]
		instruction = (instruction + 1) % len(turns)
	}

	//Find iterations until things repeat (which they do, apparently???)
	instruction = 0
	iterations := 0
	found := len(nodes)
	repeats := make([]int, found)
	for found != 0 {
		iterations++
		direction := ternary(turns[instruction] == 'R', 1, 0)
		instruction = (instruction + 1) % len(turns)
		for i, v := range nodes {
			if repeats[i] != 0 {
				continue
			}
			nodes[i] = paths[v][direction]
			if nodes[i][2] == 'Z' {
				found--
				repeats[i] = iterations
			}
		}
	}

	//Find lowest common multiple of repeat iterations
	//Modified version of this: https://siongui.github.io/2017/06/03/go-find-lcm-by-gcd/
	var LCM func(v ...int) int
	LCM = func(v ...int) int {
		a := v[0]
		b := v[1]
		for b != 0 {
			t := b
			b = a % b
			a = t
		}
		result := v[0] * v[1] / a
		for i := 2; i < len(v); i++ {
			result = LCM(result, v[i])
		}
		return result
	}
	sum2 := LCM(repeats...)

	return sum1, sum2
}

func y2023day7(input string) (int, int) {
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
	lines := strings.Split(input, "\n")
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

func y2023day6(input string) (int, int) {
	lines := strings.Split(input, "\n")
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

func y2023day5(input string) (int, int) {
	cats := strings.Split(input, "\n\n")
	readMap := func(str string) [][3]int {
		lines := strings.Split(str, "\n")[1:]
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

func y2023day4(input string) (int, int) {
	lines := strings.Split(input, "\n")
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

func y2023day3(input string) (int, int) {
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

	lines := strings.Split(input, "\n")
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

func y2023day2(input string) (int, int) {
	lines := strings.Split(input, "\n")
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

func y2023day1(input string) (int, int) {
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
