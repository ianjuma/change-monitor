The official PostgreSQL Docker image https://hub.docker.com/_/postgres/ allows us to place SQL files in the /docker-entrypoint-initb.d folder, and the first time the service starts, it will import and execute those SQL files.
In our Postgres container, we will find this bash script /usr/local/bin/docker-entrypoint.sh where each *.sh, **.sql and *.*sql.gz file will be executed

docker build -t change-monitor -f Dockerfile .
