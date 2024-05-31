#!/bin/bash\

#################################################
#       DEVELOPMENT ENVIRONMENT SETUP           #
#################################################

# SET ENVIRONMENT VARIABLES

echo "[Setting up env variables]"

export APP_ENV=development
export SERVER_ADDRESS=""
export SERVER_PORT=4000
export DATABASE_URL="postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"

# RUN DOCKER CONTAINERS

echo "[Running docker containers]"
echo "Running database container"

sudo docker remove --force /ar-db-dev 
sudo docker run --name ar-db-dev -e POSTGRES_PASSWORD=mysecretpassword -p "5432:5432" -d postgres

