SHELL = bash
COMPOSE_FILE_RUNTIME = docker-compose.yml

buildImage:
	 docker build -t mygithubapp --build-arg GITHUB_CLIENT_ID --build-arg GITHUB_CLIENT_SECRET --build-arg SESSION_KEY  .

up:
	docker-compose -f ${COMPOSE_FILE_RUNTIME} up -d --remove-orphans

down:
	docker-compose -f ${COMPOSE_FILE_RUNTIME} down --remove-orphans