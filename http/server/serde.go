package server

import (
	"github.com/sasalatart.com/quizory/http/oapi"
	"github.com/sasalatart.com/quizory/question"
	"github.com/sasalatart.com/quizory/question/enums"
)

// toUnansweredQuestion converts a question.Question to an oapi.UnansweredQuestion.
func toUnansweredQuestion(q question.Question) oapi.UnansweredQuestion {
	choices := make([]oapi.UnansweredChoice, len(q.Choices))
	for _, c := range q.Choices {
		choices = append(choices, toUnansweredChoice(c))
	}
	return oapi.UnansweredQuestion{
		Id:         q.ID,
		Topic:      q.Topic.String(),
		Question:   q.Question,
		Hint:       q.Hint,
		Difficulty: toDifficulty(q.Difficulty),
		Choices:    choices,
	}
}

// toDifficulty converts a enums.Difficulty to an oapi.Difficulty.
func toDifficulty(d enums.Difficulty) oapi.Difficulty {
	switch d {
	case enums.DifficultyNoviceHistorian:
		return oapi.DifficultyNoviceHistorian
	case enums.DifficultyAvidHistorian:
		return oapi.DifficultyAvidHistorian
	case enums.DifficultyHistoryScholar:
		return oapi.DifficultyHistoryScholar
	default:
		return oapi.DifficultyNoviceHistorian
	}
}

// toUnansweredChoice converts a question.Choice to an oapi.UnansweredChoice.
func toUnansweredChoice(c question.Choice) oapi.UnansweredChoice {
	return oapi.UnansweredChoice{
		Id:     c.ID,
		Choice: c.Choice,
	}
}
