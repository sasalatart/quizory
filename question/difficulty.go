package question

import "golang.org/x/exp/rand"

//go:generate go run github.com/dmarkham/enumer -type=Difficulty -trimprefix=Difficulty -transform=whitespace
type Difficulty int

const (
	DifficultyNoviceHistorian Difficulty = iota
	DifficultyAvidHistorian
	DifficultyHistoryScholar
)

// RandomDifficulty returns a random difficulty.
func RandomDifficulty() Difficulty {
	return Difficulty(rand.Intn(len(DifficultyValues())))
}
