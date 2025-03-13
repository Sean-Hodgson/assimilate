 
# Define variables for image and container names
IMAGE_NAME := assimilateproxy
CONTAINER_NAME_TEST := assimilateproxyTEST
CONTAINER_NAME_DEPLOY := assimilateproxyDEPLOY
PORTS := -p 9051:9051 -p 9001:9001

build:
	docker build -t $(IMAGE_NAME) .

#interactive
test:
	docker run -it --rm --name $(CONTAINER_NAME_TEST) $(PORTS) $(IMAGE_NAME) /bin/bash

# Run the container in detached mode
deploy:
	docker run -d --rm --name $(CONTAINER_NAME_DEPLOY) $(PORTS) $(IMAGE_NAME)

# Stop the container
stop:
	docker stop $(CONTAINER_NAME_DEPLOY) || true
	docker stop $(CONTAINER_NAME_TEST) ||

# Reload containers
reload:
	$(MAKE) stop
	$(MAKE) rm

# Remove the container
rm:
	docker rm $(CONTAINER_NAME_DEPLOY) || true
	docker rm $(CONTAINER_NAME_TEST) || true

# Remove the image
rmi:
	docker rmi $(IMAGE_NAME)

# List all containers (running and stopped)
ps:
	docker ps -a

# Show logs of the container
logtest:
	docker logs $(CONTAINER_NAME_TEST)

logdeploy:
	docker logs $(CONTAINER_NAME_DEPLOY)

# Remove all stopped containers
cleancontainers:
	docker container prune -f

cleannetwork:
	docker network prune -f

# Remove all unused images
cleanimages:
	docker image prune -a -f
