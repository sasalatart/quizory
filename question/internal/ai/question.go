package ai

// Question represents a question generated via AI.
type Question struct {
	Question   string   `json:"question"`
	Hint       string   `json:"hint"`
	Choices    []Choice `json:"choices"`
	MoreInfo   string   `json:"moreInfo"`
	Difficulty string   `json:"difficulty"`
}

// Choice represents a possible answer to a question generated via AI.
type Choice struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}
