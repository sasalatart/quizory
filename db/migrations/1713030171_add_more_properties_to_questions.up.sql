ALTER TABLE questions
  ADD COLUMN difficulty VARCHAR(255) NOT NULL,
  ADD COLUMN topic VARCHAR(255) NOT NULL,
  ADD COLUMN more_info TEXT NOT NULL;

CREATE INDEX questions_difficulty_idx ON questions (difficulty);
CREATE INDEX questions_topic_idx ON questions (topic);
