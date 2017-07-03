PROJECT_NAME = mqti

.PHONY: test

run:
	docker-compose up

test:
	cd $(PROJECT_NAME) ; make test

docker_image:
	docker build -t ashmckenzie/mqti .
