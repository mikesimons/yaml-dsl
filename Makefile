all:
	cd vendor/github.com/mitchellh/go-mruby && MRUBY_CONFIG=$(shell pwd)/mruby_config.rb make libmruby.a
	cp vendor/github.com/mitchellh/go-mruby/libmruby.a .
	go build

.PHONY: all
