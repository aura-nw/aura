FROM postgres:13-alpine
# COPY ./00-cosmos.sql /docker-entrypoint-initdb.d/
# COPY ./01-auth.sql /docker-entrypoint-initdb.d/
# COPY ./02-bank.sql /docker-entrypoint-initdb.d/
# COPY ./03-staking.sql /docker-entrypoint-initdb.d/
# COPY ./04-consensus.sql /docker-entrypoint-initdb.d/
# COPY ./05-mint.sql /docker-entrypoint-initdb.d/
# COPY ./06-distribution.sql /docker-entrypoint-initdb.d/
# COPY ./07-pricefeed.sql /docker-entrypoint-initdb.d/
# COPY ./08-gov.sql /docker-entrypoint-initdb.d/
# COPY ./09-modules.sql /docker-entrypoint-initdb.d/
# COPY ./10-slashing.sql /docker-entrypoint-initdb.d/
# COPY ./11-feegrant.sql /docker-entrypoint-initdb.d/
# COPY ./12-upgrade.sql /docker-entrypoint-initdb.d/
COPY ./schema /docker-entrypoint-initdb.d/
COPY ./test_hasura_init.sh /docker-entrypoint-initdb.d/