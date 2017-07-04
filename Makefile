PROJECT_NAME = mqti

.PHONY: test docker_image

test:
	cd $(PROJECT_NAME) ; make test

docker_image:
	docker build -t ashmckenzie/mqti .
