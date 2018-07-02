name = tamtam
package = github.com/servicelab/tamtam

PLATFORMS := darwin/amd64 linux/amd64 linux/arm linux/386 windows/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
version = `cat VERSION`
time = `date +%FT%T%z`
hash = `git rev-parse HEAD`
ldflags = "-s -w -X $(package)/cmd.Version=$(version) -X $(package)/cmd.BuildTime=$(time) -X $(package)/cmd.GitHash=$(hash)"

build:
	go build -o $(name)

test:
	go test

docker: test
	GOOS=linux GOARCH=amd64 go build -ldflags $(ldflags) -o 'dist/$(name)'

release: $(PLATFORMS)

clean:
	rm -rf dist
	rm $(name)

$(PLATFORMS): test
	GOOS=$(os) GOARCH=$(arch) go build -ldflags $(ldflags) -o 'dist/$(name)-$(os)-$(arch)'

.PHONY: build
