NAME   := kevinmichaelchen/go-sqlboiler-user-api
TAG    := $(shell git log -1 --pretty=%h)
IMG    := ${NAME}:${TAG}
LATEST := ${NAME}:latest

dc-build:
	@docker build -t ${IMG} .
	@docker tag ${IMG} ${LATEST}

dc-push:
	@docker push ${NAME}

dc-login:
	@docker login -u ${DOCKER_USER} --password-stdin