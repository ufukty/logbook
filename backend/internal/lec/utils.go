package lec

import "slices"

func pow(number, power int) int {
	if power == 0 {
		return 1
	}
	result := 1
	for i := 0; i < power; i++ {
		result *= number
	}
	return result
}

// returns the min n for t <= 2^n
func minpow(t int) int {
	if t == 1 {
		return 1
	}
	r := 0
	for n := 1; n < t; n *= 2 {
		r++
	}
	return r
}

func twodim(rows, columns int, def int) [][]int {
	row := []int{}
	for range columns {
		row = append(row, def)
	}
	dst := [][]int{}
	for i := 0; i < rows; i++ {
		dst = append(dst, slices.Clone(row))
	}
	return dst
}

func longestline(src [][]int) int {
	chars := -1
	for i := 0; i < len(src); i++ {
		if chars < len(src[i]) {
			chars = len(src[i])
		}
	}
	return chars
}

func transpose(src [][]int) [][]int {
	dst := twodim(longestline(src), len(src), -1)
	for i := 0; i < len(src); i++ {
		for j := 0; j < len(src[i]); j++ {
			dst[j][i] = src[i][j]
		}
	}
	return dst
}
