CREATE TABLE IF NOT EXISTS answers (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid NOT NULL,
  choice_id uuid NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  FOREIGN KEY (choice_id) REFERENCES choices(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_answers_user_id ON answers (user_id);
