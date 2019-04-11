protoc-setup:
	cd meshes
	wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto

proto:	
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

docker:
	docker build -t layer5/meshery-consul .

docker-run:
	(docker rm -f meshery-consul) || true
	docker run --name meshery-consul -d \
	-p 10002:10002 \
	-e DEBUG=true \
	layer5/meshery-consul

run:
	DEBUG=true go run main.go