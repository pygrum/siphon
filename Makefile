.PHONY: build build-agent build-gen
build:
	go build cmd/siphon/siphon.go

build-agent:
	go build cmd/agent/siphon_agent.go

build-gen:
	go build cmd/generator/siphon_gen.go