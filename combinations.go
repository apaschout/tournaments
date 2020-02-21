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

func (trn *Tournament) MakeMatches() {
	plrs := []PlayerID{}
	for _, par := range trn.Players {
		plrs = append(plrs, par.Player)
	}
	if len(plrs)%2 != 0 {
		// x is dummy player
		plrs = append(plrs, "x")
	}

	numRounds := len(plrs) - 1
	halfSize := len(plrs) / 2

	withoutFirst := plrs[1:]
	plrsLen := len(withoutFirst)
	matches := []Match{}

	for round := 0; round < numRounds; round++ {
		plrIdx := round % plrsLen
		matches = append(matches, Match{Player1: plrs[0], Player2: withoutFirst[plrIdx]})
		for i := 1; i < halfSize; i++ {
			plr1 := (round + plrsLen - i) % plrsLen
			plr2 := (round + i) % plrsLen
			matches = append(matches, Match{Player1: withoutFirst[plr1], Player2: withoutFirst[plr2]})
		}
	}
	trn.Matches = deleteDummyMatches(matches)
}

func deleteDummyMatches(matches []Match) []Match {
	for i := 0; i < len(matches); i++ {
		if matches[i].Player1 == "x" || matches[i].Player2 == "x" {
			matches = append(matches[:i], matches[i+1:]...)
			i--
		}
	}
	return matches
}
