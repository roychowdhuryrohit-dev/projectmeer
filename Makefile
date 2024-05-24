.PHONY: all build_local run_local local build_docker clean_docker run_docker docker

all: local

build_local:
	@cd assets && npm run build
	@go build -o bin/meer
run_local:
	@. ./.env && ./bin/meer
local: build_local run_local

build_docker:
	docker build -t meer . && \
	docker images meer
run_docker:
	. ./.env && \
	docker run --env REPLICA_ID=$$REPLICA_ID --env NODE_LIST=$$NODE_LIST --env PORT=$$PORT -p 80:80 meer
clean_docker:
	docker container kill $$(docker container ls -aq); \
	docker rm $$(docker ps -a -q) \
	docker system prune --all --force --volumes; \
	docker image rm $$(docker image ls -a -q)
docker : build_docker run_docker
