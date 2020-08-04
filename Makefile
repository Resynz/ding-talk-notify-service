EXE = ding-talk-notify-service
IMAGE_NAME = $(EXE):0.0.1
build:
	go build -ldflags "-linkmode external -extldflags -static"
docker:
	mv $(EXE) docker-configs
	@echo "building docker image $(IMAGE_NAME) ..."
	docker build -t $(IMAGE_NAME) docker-configs
	rm docker-configs/$(EXE)
run:
	@echo "runing $(EXE) ..."
	docker-compose up -d
publish:
	@echo "publishing images to server ..."
	docker save $(IMAGE_NAME) | gzip > $(EXE).image.tar.gz
	scp $(EXE).image.tar.gz jianxin@47.92.254.232:~
	rm ./$(EXE).image.tar.gz
	@echo "publish Done!"
stop:
	@echo "stopping $(EXE) ..."
	docker-compose down
clean:
	@echo "cleaning docker image $(IMAGE_NAME) ..."
	docker rmi $(IMAGE_NAME)
log:
	docker-compose logs -f
help:
	@echo "make			-- build $(EXE) application"
	@echo "make docker		-- build docker image"
	@echo "make run 		-- run $(EXE) docker server"
	@echo "make stop		-- stop $(EXE) docker server"
	@echo "make clean 		-- clean docker image"
	@echo "make log 		-- show logs"
	@echo "make publish		-- publish docker image to develop server"