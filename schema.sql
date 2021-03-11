CREATE TABLE IF NOT EXISTS topics (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
  id integer PRIMARY KEY,
  topic_id integer REFERENCES topics (id) NOT NULL,
  content text NOT NULL,
  author_initials character(2) NOT NULL CHECK(author_initials GLOB '[A-Z][A-Z]'),
	author_theme integer NOT NULL,
  posted TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	FOREIGN KEY (author_theme) REFERENCES colors (id)
);