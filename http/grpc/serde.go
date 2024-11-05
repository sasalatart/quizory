package grpc

import (
	"fmt"

	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/domain/question/enums"
	"github.com/sasalatart/quizory/http/grpc/proto"
)

func toLatestQuestions(questions []question.Question) []string {
	result := make([]string, 0, len(questions))
	for _, q := range questions {
		result = append(result, q.Question)
	}
	return result
}

func fromDifficulty(difficulty proto.Difficulty) (*enums.Difficulty, error) {
	var result enums.Difficulty
	switch difficulty {
	case proto.Difficulty_DIFFICULTY_NOVICE_HISTORIAN:
		result = enums.DifficultyNoviceHistorian
	case proto.Difficulty_DIFFICULTY_AVID_HISTORIAN:
		result = enums.DifficultyAvidHistorian
	case proto.Difficulty_DIFFICULTY_HISTORY_SCHOLAR:
		result = enums.DifficultyHistoryScholar
	default:
		return nil, fmt.Errorf("invalid difficulty: %s", difficulty)
	}
	return &result, nil
}
