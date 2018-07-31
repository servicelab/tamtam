name = tamtam
group = servicelab
package = github.com/$(group)/$(name)

# docker name space
namespace = servicelaborg
image = $(namespace)/$(name)

# make $(branchtag) empty when on a maintenance branch
branchtag = :dev

# Targets for [b]uilding binaries for various platforms
PLATFORMS = b/darwin/amd64 b/linux/amd64 b/linux/arm b/linux/arm64 b/linux/386 b/windows/amd64

# Targets for [d]ockerizing containers for various platforms
DOCKER = d/linux/amd64 d/linux/arm d/linux/arm64 d/linux/386

temp = $(subst /, ,$@)
os = $(word 2, $(temp))
arch = $(word 3, $(temp))

time = `date +%FT%T%z`
hash = `git rev-parse HEAD`
version = `git describe --tags`
ldflags = "-s -w -X $(package)/cmd.Version=$(version) -X $(package)/cmd.BuildTime=$(time) -X $(package)/cmd.GitHash=$(hash)"

build:
	go build -o $(name)

test:
	go test ./...

docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(ldflags) -o 'dist/$(name)-$(os)-$(arch)'
	docker build --build-arg BIN=dist/$(name)-$(os)-$(arch) -t $(image) .

images: $(DOCKER)

release: $(PLATFORMS)

buildall: $(PLATFORMS)

clean:
	rm -rf dist
	rm $(name)

$(PLATFORMS): test
	@mkdir -p dist/$(name)-$(os)-$(arch)
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build -ldflags $(ldflags) -o 'dist/$(name)-$(os)-$(arch)/$(name)'
	@if [ "$(os)" == "windows" ]; then \
		mv dist/$(name)-$(os)-$(arch)/$(name) dist/$(name)-$(os)-$(arch)/$(name).exe ; \
	fi
	zip -j dist/$(name)-$(version)-$(os)-$(arch).zip dist/$(name)-$(os)-$(arch)/*
	@if [ "$(GITHUB_TOKEN)" != "" ]; then \
		curl --data-binary @"dist/$(name)-$(version)-$(os)-$(arch).zip" \
			-H "Authorization: token $(GITHUB_TOKEN)" \
			-H "Content-Type: application/octet-stream" \
			https://uploads.github.com/repos/$(group)/$(name)/releases/$(version)/assets?name=$(name)-$(version)-$(os)-$(arch).zip ; \
	fi

login:
	@if [ "$(DOCKER_USER)" != "" ]; then \
		docker login -u $(DOCKER_USER) -p $(DOCKER_PASSWORD) ; \
	fi

$(DOCKER): login
	# build
	docker build --build-arg BIN=dist/$(name)-$(os)-$(arch) -t $(image)$(branchtag)-$(os)-$(arch) .

	# tag
	docker tag $(image)$(branchtag)-$(os)-$(arch) $(image):$(version)-$(os)-$(arch)

	# push if user is set
	@if [ "$(DOCKER_USER)" != "" ]; then \
		docker push $(image):$(version)-$(os)-$(arch) ; \
		if ["$(branchtag)" != "" ]; then \
			docker push $(image)$(branchtag)-$(os)-$(arch) ; \
		fi \
	fi

.PHONY: build
