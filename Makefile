# Make sure service names match directory names under cmd/
SERVICES := pdf-compressor img-to-pdf pdf-to-img

build-all:
	for service in $(SERVICES); do \
		docker build --build-arg SERVICE=$$service -t $$service:latest . ; \
	done

# Add a build target for a specific service
build:
	docker build --build-arg SERVICE=$(SERVICE) -t $(SERVICE):latest .

run:
	docker run -p 8080:8080 $(SERVICE):latest

push-all:
	for service in $(SERVICES); do \
		docker push $$service:latest ; \
	done

deploy-all:
	kubectl apply -f k8s/

# Add separate targets for each service
pdf-compressor:
	make build SERVICE=pdf-compressor
	make run SERVICE=pdf-compressor

img-to-pdf:
	make build SERVICE=img-to-pdf
	make run SERVICE=img-to-pdf

pdf-to-img:
	make build SERVICE=pdf-to-img
	make run SERVICE=pdf-to-img