name = tamtam
namespace = servicelaborg
package = github.com/servicelab/tamtam
image = $(namespace)/$(name)

# make $(latest) empty when on a maintenance branch
latest = :latest

major = 1
minor = 0
patch = 0

# Targets for [b]uilding binaries for various platforms
PLATFORMS = b/darwin/amd64 b/linux/amd64 b/linux/arm b/linux/arm64 b/linux/386 b/windows/amd64

# Targets for [d]ockerizing containers for various platforms
DOCKER = d/linux/amd64 d/linux/arm d/linux/arm64 d/linux/386

temp = $(subst /, ,$@)
os = $(word 2, $(temp))
arch = $(word 3, $(temp))

time = `date +%FT%T%z`
hash = `git rev-parse HEAD`
ldflags = "-s -w -X $(package)/cmd.Version=$(major).$(minor).$(patch) -X $(package)/cmd.BuildTime=$(time) -X $(package)/cmd.GitHash=$(hash)"

build:
	go build -o $(name)

test:
	go test

docker:
	GOOS=linux GOARCH=amd64 go build -ldflags $(ldflags) -o 'dist/$(name)-$(os)-$(arch)'
	docker build --build-arg BIN=dist/$(name)-$(os)-$(arch) -t $(image) .

images: $(DOCKER)

release: $(PLATFORMS)

clean:
	rm -rf dist
	rm $(name)

$(PLATFORMS): test
	GOOS=$(os) GOARCH=$(arch) go build -ldflags $(ldflags) -o 'dist/$(name)-$(os)-$(arch)'

login:
	@if [ "$(DOCKER_USER)" != "" ]; then \
		docker login -u $(DOCKER_USER) -p $(DOCKER_PASSWORD) ; \
	fi

$(DOCKER): login
	# build
	docker build --build-arg BIN=dist/$(name)-$(os)-$(arch) -t $(image) .

	# tag
	docker tag $(image) $(image):$(major)-$(os)-$(arch)
	docker tag $(image) $(image):$(major).$(minor)-$(os)-$(arch)
	docker tag $(image) $(image):$(major).$(minor).$(patch)-$(os)-$(arch)
	docker tag $(image) $(image)$(latest)-$(os)-$(arch)

	# push if user is set
	@if [ "$(DOCKER_USER)" != "" ]; then \
		docker push $(image):$(major)-$(os)-$(arch) ; \
		docker push $(image):$(major).$(minor)-$(os)-$(arch) ; \
		docker push $(image):$(major).$(minor).$(patch)-$(os)-$(arch) ; \
		docker push $(image)$(latest)-$(os)-$(arch) ; \
	fi

.PHONY: build
