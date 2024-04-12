CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS questions (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  question TEXT NOT NULL,
  hint TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
