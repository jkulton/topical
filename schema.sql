CREATE TABLE IF NOT EXISTS topics (
  id serial PRIMARY KEY,
  title text NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
  id serial PRIMARY KEY,
  topic_id integer REFERENCES topics (id) NOT NULL,
  content text NOT NULL,
  author_initials char(2) NOT NULL CHECK (author_initials ~ '^[A-Z]{2}$'),
  author_theme integer NOT NULL,
  posted timestamp NOT NULL DEFAULT NOW()
);