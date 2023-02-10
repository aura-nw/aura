#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE DATABASE "$POSTGRES_DB_TEST";
	GRANT ALL PRIVILEGES ON DATABASE "$POSTGRES_DB_TEST" TO "$POSTGRES_USER";
	CREATE DATABASE hasura;
	GRANT ALL PRIVILEGES ON DATABASE hasura TO "$POSTGRES_USER";
EOSQL

# psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -a -f ./docker-entrypoint-initdb.d/00-cosmos.sql
# psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -a -f ./docker-entrypoint-initdb.d/01-auth.sql
# psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -a -f ./docker-entrypoint-initdb.d/02-bank.sql
# psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -a -f ./docker-entrypoint-initdb.d/04-consensus.sql
# psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -a -f ./docker-entrypoint-initdb.d/05-mint.sql
# psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -a -f ./docker-entrypoint-initdb.d/06-distribution.sql
# psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -a -f ./docker-entrypoint-initdb.d/07-pricefeed.sql
# psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -a -f ./docker-entrypoint-initdb.d/08-gov.sql
# psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -a -f ./docker-entrypoint-initdb.d/09-modules.sql
# psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -a -f ./docker-entrypoint-initdb.d/10-slashing.sql
# psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -a -f ./docker-entrypoint-initdb.d/11-feegrant.sql
# psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -a -f ./docker-entrypoint-initdb.d/12-upgrade.sql