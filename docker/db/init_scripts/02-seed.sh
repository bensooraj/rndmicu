#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL

    BEGIN;
        INSERT INTO creators (name, email) VALUES ('ben_0', 'ben_0@gmail.com');
        INSERT INTO creators (name, email) VALUES ('ben_1', 'ben_1@gmail.com');
        INSERT INTO creators (name, email) VALUES ('ben_2', 'ben_2@gmail.com');
        INSERT INTO creators (name, email) VALUES ('ben_3', 'ben_3@gmail.com');
    COMMIT;

EOSQL