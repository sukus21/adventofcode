package main

import "strconv"

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

func abs(x int) int {
	return ternary(x < 0, -x, x)
}

func sign(x int) int {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}
