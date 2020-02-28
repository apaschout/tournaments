package tournaments

import "fmt"

func binomialCoefficient(n, k int) (int, error) {
	if n < 0 || k < 0 {
		return 0, fmt.Errorf("Input negative")
	}
	if n < k {
		return 0, fmt.Errorf("Set too small for k")
	}
	b := 1
	for i := 1; i <= k; i++ {
		b = (n - k + i) * b / i
	}
	return b, nil
}

func combinations(n, k int) [][]int {
	combins, _ := binomialCoefficient(n, k)
	res := make([][]int, combins)
	if len(res) == 0 {
		return res
	}
	res[0] = make([]int, k)
	for i := range res[0] {
		res[0][i] = i
	}
	for i := 1; i < combins; i++ {
		next := make([]int, k)
		copy(next, res[i-1])
		nextCombination(next, n, k)
		res[i] = next
	}
	return res
}

func nextCombination(s []int, n, k int) {
	for j := k - 1; j >= 0; j-- {
		if s[j] == n+j-k {
			continue
		}
		s[j]++
		for l := j + 1; l < k; l++ {
			s[l] = s[j] + 1 - j
		}
		break
	}
}
