TAG?=latest
IMAGE_NAME=carlosdescalzi/mtjobrunner:$(TAG)
SOURCES=cmd pkg/apis pkg/controller client/v1alpha1

run: go run cmd/* -kubeconfig=$(KUBECONFIG) 

build: 
	docker build -t $(IMAGE_NAME) $(PWD)