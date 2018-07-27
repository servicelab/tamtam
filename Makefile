name = tamtam
namespace = servicelaborg
package = github.com/servicelab/tamtam
image = $(namespace)/$(name)
latest = :latest # make empty when on a maintenance branch
major = 1
minor = 0
patch = 0

PLATFORMS = darwin/amd64/p linux/amd64/p linux/arm/p linux/386/p windows/amd64/p
DOCKER = linux/amd64/d linux/arm/d linux/386/d

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

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
		docker login -u $(DOCKER_USER) -p $(DOCKER_PASSWORD); \
	fi

$(DOCKER): $(PLATFORMS) login
	# build
	docker build --build-arg BIN=dist/$(name)-$(os)-$(arch) -t $(image)/$(os)/$(arch) .

	# tag
	docker tag $(image)/$(os)/$(arch) $(image)/$(os)/$(arch):$(major)
	docker tag $(image)/$(os)/$(arch) $(image)/$(os)/$(arch):$(major).$(minor)
	docker tag $(image)/$(os)/$(arch) $(image)/$(os)/$(arch):$(major).$(minor).$(patch)
	docker tag $(image)/$(os)/$(arch) $(image)/$(os)/$(arch)$(latest)

	# push if user is set
	@if [ "$(DOCKER_USER)" != "" ]; then \
		docker push $(image)/$(os)/$(arch) $(image)/$(os)/$(arch):$(major); \
		docker push $(image)/$(os)/$(arch) $(image)/$(os)/$(arch):$(major).$(minor); \
		docker push $(image)/$(os)/$(arch) $(image)/$(os)/$(arch):$(major).$(minor).$(patch); \
		docker push $(image)/$(os)/$(arch) $(image)/$(os)/$(arch)$(latest); \
	fi

.PHONY: build
