proto: ## Generate protobuf files
	docker run --volume "$(shell pwd):/workspace" --workdir /workspace bufbuild/buf generate
