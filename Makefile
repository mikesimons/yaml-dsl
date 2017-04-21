export CGO_LDFLAGS=-lyaml

all: go-mruby
	go build -ldflags -s

test:
	ginkgo -r -cover
	gover
	mkdir -p _reports
	go tool cover -html gover.coverprofile -o _reports/coverage.html

go-mruby:
	cd vendor/github.com/mitchellh/go-mruby && MRUBY_CONFIG=$(shell pwd)/mruby_config.rb make libmruby.a
	cp vendor/github.com/mitchellh/go-mruby/libmruby.a .

.PHONY: all go-mruby
