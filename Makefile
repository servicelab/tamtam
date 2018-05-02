name = tamtam
package = github.com/servicelab/tamtam

PLATFORMS := darwin/amd64 linux/amd64 linux/arm linux/386 windows/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
version = `cat VERSION`
time = `date +%FT%T%z`
hash = `git rev-parse HEAD`
ldflags = "-X $(package)/cmd.Version=$(version) -X $(package)/cmd.BuildTime=$(time) -X $(package)/cmd.GitHash=$(hash)"

build:
	go build -o $(name)

test:
	go test

ensure:
	dep ensure --update

ensure-test:
	dep ensure --update
	go test

prepare: ensure-test

release: $(PLATFORMS)

clean:
	rm -rf dist
	rm $(name)

$(PLATFORMS): prepare
	GOOS=$(os) GOARCH=$(arch) go build -ldflags $(ldflags) -o 'dist/$(name)-$(os)-$(arch)'

.PHONY: build
