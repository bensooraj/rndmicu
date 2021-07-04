-- Description: Create table creators
CREATE TABLE IF NOT EXISTS creators (
	id       SERIAL PRIMARY KEY,
	name     TEXT NOT NULL,
	email    TEXT NOT NULL UNIQUE
);

-- Description: Create table audio_shorts
CREATE TABLE IF NOT EXISTS audio_shorts (
	id               UUID,
	creator_id       INTEGER NOT NULL,
	title            TEXT NOT NULL,
	description      TEXT NOT NULL,
	category         TEXT NOT NULL,
	audio_file_url   TEXT NOT NULL,
	date_created     TIMESTAMP NOT NULL,
	date_updated     TIMESTAMP NOT NULL,

	PRIMARY KEY (id)
);