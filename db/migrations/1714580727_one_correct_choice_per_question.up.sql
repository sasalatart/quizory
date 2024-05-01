CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_correct_choice_on_question
ON choices (question_id)
WHERE is_correct = true;
