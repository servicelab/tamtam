name = tamtam
group = servicelab
package = github.com/$(group)/$(name)

# docker name space
namespace = servicelaborg
image = $(namespace)/$(name)
version = `git describe --tags`

# make $(branchtag) empty when on a maintenance branch
branchtag = :latest

# Targets for [b]uilding binaries for various platforms
PLATFORMS = b/darwin/amd64 b/linux/amd64 b/linux/arm b/linux/arm64 b/linux/386 b/windows/amd64

# Targets for [d]ockerizing containers for various platforms
DOCKER = d/linux/amd64 d/linux/arm_6 d/linux/arm64 d/linux/386

temp = $(subst /, ,$@)
os = $(word 2, $(temp))
arch = $(word 3, $(temp))

time = `date +%FT%T%z`
hash = `git rev-parse HEAD`
branch = `git rev-parse --abbrev-ref HEAD`
ldflags = "-s -w -X $(package)/cmd.Version=$(branch) -X $(package)/cmd.BuildTime=$(time) -X $(package)/cmd.GitHash=$(hash)"

build:
	go build -ldflags $(ldflags) -o $(name)

test:
	go test ./...

docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(ldflags) -o 'dist/$(name)-$(os)-$(arch)'
	docker build --build-arg BIN=dist/$(name)-$(os)-$(arch) -t $(image) .

images: $(DOCKER)

cross: $(PLATFORMS)

buildall: $(PLATFORMS)

clean:
	rm -rf dist
	rm -f $(name)

$(PLATFORMS): test
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -ldflags $(ldflags) -o 'dist/$(name)-$(os)-$(arch)'

login:
	@if [ "$(DOCKER_USER)" != "" ]; then \
		docker login -u $(DOCKER_USER) -p $(DOCKER_PASSWORD) ; \
	fi

$(DOCKER): login
	# build
	docker build --build-arg BIN=dist/$(os)_$(arch)/$(name) -t $(image)$(branchtag)-$(os)-$(arch) .

	# tag
	docker tag $(image)$(branchtag)-$(os)-$(arch) $(image):$(version)-$(os)-$(arch)

	# push if user is set
	@if [ "$(DOCKER_USER)" != "" ]; then \
		docker push $(image):$(version)-$(os)-$(arch) ; \
	fi

.PHONY: build
