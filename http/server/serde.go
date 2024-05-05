package server

import (
	"github.com/sasalatart.com/quizory/answer"
	"github.com/sasalatart.com/quizory/http/oapi"
	"github.com/sasalatart.com/quizory/question"
	"github.com/sasalatart.com/quizory/question/enums"
)

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

func toUnansweredChoice(c question.Choice) oapi.UnansweredChoice {
	return oapi.UnansweredChoice{
		Id:     c.ID,
		Choice: c.Choice,
	}
}

func toSubmitAnswerResult(r answer.SubmissionResponse) oapi.SubmitAnswerResult {
	return oapi.SubmitAnswerResult{
		Id:              r.ID,
		CorrectChoiceId: r.CorrectChoiceID,
		MoreInfo:        r.MoreInfo,
	}
}
