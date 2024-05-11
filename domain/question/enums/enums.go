package enums

import "golang.org/x/exp/rand"

//go:generate go run github.com/dmarkham/enumer -type=Difficulty -trimprefix=Difficulty -transform=whitespace
type Difficulty int

const (
	DifficultyNoviceHistorian Difficulty = iota
	DifficultyAvidHistorian
	DifficultyHistoryScholar
)

//go:generate go run github.com/dmarkham/enumer -type=Topic -trimprefix=Topic -transform=whitespace
type Topic int

const (
	TopicAncientEgypt Topic = iota
	TopicAncientGreece
	TopicAncientChina
	TopicAncientRome
	TopicMiddleAges
	TopicCrusades
	TopicRenaissance
	TopicScientificRevolution
	TopicReformation
	TopicEnlightenment
	TopicAmericanRevolution
	TopicFrenchRevolution
	TopicNapoleonicWars
	TopicIndustrialRevolution
	TopicAmericanCivilWar
	TopicChineseRevolution
	TopicRussianRevolution
	TopicWorldWarI
	TopicWorldWarII
	TopicColdWar
)

// RandomTopic returns a random topic.
func RandomTopic() Topic {
	return Topic(rand.Intn(len(TopicValues())))
}
