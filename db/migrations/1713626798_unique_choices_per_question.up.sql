ALTER TABLE choices ADD CONSTRAINT unique_choices_per_question UNIQUE (question_id, choice);
