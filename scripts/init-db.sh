#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create the database for GlitchTip, if it doesn't exist.
    SELECT 'CREATE DATABASE glitchtip_db'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'glitchtip_db')\gexec

    -- Grant all privileges on the new database to the app_user.
    GRANT ALL PRIVILEGES ON DATABASE glitchtip_db TO app_user;
EOSQL
