#!/bin/bash
set -x
set -eo pipefail

SKIP_DOCKER=$1

DATABASE_USER="${DB_USER}"
DATABASE_PASSWORD="${DB_PASSWORD}"
DATABASE_NAME="${DB_NAME}"
DATABASE_PORT="${DB_PORT}"


if [[ -z "${SKIP_DOCKER}" ]]
then
  # Start PostgreSQL container
  docker run --rm -d --name postgresql \
    -p ${DATABASE_PORT}:5432 \
    -e POSTGRES_PASSWORD=${DATABASE_PASSWORD} \
    -e POSTGRES_DB=${DATABASE_NAME} \
    -e POSTGRES_USER=${DATABASE_USER} \
    postgres:15.2
fi

# Wait for PostgreSQL to be available
export PGPASSWORD="${DATABASE_PASSWORD}"
until docker exec -it postgresql psql -U "${DATABASE_USER}" "${DATABASE_NAME}" -c '\q'; do
  >&2 echo "Postgres is still unavailable - sleeping"
  sleep 3
done

>&2 echo "Postgres is up and running on port ${DATABASE_PORT}"

