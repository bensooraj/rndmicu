#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL

    BEGIN;

        CREATE USER healthcheck_user;

        CREATE TABLE IF NOT EXISTS creators (
            id       SERIAL PRIMARY KEY,
            name     TEXT NOT NULL,
            email    TEXT NOT NULL UNIQUE
        );

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
        CREATE INDEX idx_audio_shorts_id ON audio_shorts (id);

    COMMIT;

EOSQL