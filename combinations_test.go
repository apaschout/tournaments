package tournaments

import "testing"

func TestBinomialCoefficient(t *testing.T) {
	bin, _ := binomialCoefficient(3, 2)
	if bin != 28 {
		t.Errorf("want: %d, got %d", 28, bin)
	}
}

func TestCombinations(t *testing.T) {
	res := combinations(8, 2)
	if len(res) != 1 {
		t.Errorf("got: %v", res)
	}
}

func TestMakeMatches(t *testing.T) {
	plrs := []Participant{
		{
			Player: "1",
		},
		{
			Player: "2",
		},
		{
			Player: "3",
		},
		{
			Player: "4",
		},
		{
			Player: "5",
		},
		{
			Player: "6",
		},
	}
	trn := Tournament{
		Players: plrs,
	}

	trn.MakeMatches()

	if len(trn.Matches) != 0 {
		t.Errorf("got %v", trn.Matches)
	}
}
