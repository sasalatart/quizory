package rest

import (
	"slices"

	"github.com/sasalatart/quizory/domain/answer"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/domain/question/enums"
	"github.com/sasalatart/quizory/http/rest/oapi"
)

func toRemainingTopics(m map[enums.Topic]uint) []oapi.RemainingTopic {
	result := make([]oapi.RemainingTopic, 0, len(m))
	for topic, amount := range m {
		result = append(result, oapi.RemainingTopic{
			Topic:             topic.String(),
			AmountOfQuestions: int(amount),
		})
	}
	slices.SortStableFunc(result, func(rt1, rt2 oapi.RemainingTopic) int {
		if rt1.Topic < rt2.Topic {
			return -1
		}
		if rt1.Topic > rt2.Topic {
			return 1
		}
		return 0
	})
	return result
}

func toQuestion(q question.Question) oapi.Question {
	choices := make([]oapi.Choice, 0, len(q.Choices))
	for _, c := range q.Choices {
		choices = append(choices, toChoice(c))
	}
	return oapi.Question{
		Id:         q.ID,
		Topic:      q.Topic.String(),
		Question:   q.Question,
		Hint:       q.Hint,
		MoreInfo:   q.MoreInfo,
		Difficulty: toDifficulty(q.Difficulty),
		Choices:    choices,
	}
}

func toChoice(c question.Choice) oapi.Choice {
	return oapi.Choice{
		Id:        c.ID,
		Choice:    c.Choice,
		IsCorrect: c.IsCorrect,
	}
}

func toUnansweredQuestion(q question.Question) oapi.UnansweredQuestion {
	choices := make([]oapi.UnansweredChoice, 0, len(q.Choices))
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

func toAnswersLogItem(l answer.LogItem) oapi.AnswersLogItem {
	return oapi.AnswersLogItem{
		Id:       l.ID,
		Question: toQuestion(l.Question),
		ChoiceId: l.ChoiceID,
	}
}
