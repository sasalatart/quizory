CREATE TABLE IF NOT EXISTS choices (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  question_id uuid NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  choice TEXT NOT NULL,
  is_correct BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
