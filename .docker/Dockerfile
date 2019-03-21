FROM golang:1.12.1 as builder

# Install go-bindata
RUN go get -u github.com/jteeuwen/go-bindata/...

# Copy the source code converter script
COPY .docker /go/src/github.com/egjiri/cliff/.docker

# Specifiy a volume to mount external cli.yml files or other cliff based projects
VOLUME /data

# Specifity the GOOS for the target binary. Ex. darwin, linux, windows, etc.
ENV GOOS_TARGET=""

# Specifity the GOARCH for the target binary. Ex. darwin, linux, windows, etc.
ENV GOARCH_TARGET=""

# Specifity the flags to the build command.
ENV GOBUILD_FLAGS=""

# Specifity the repository name of the project as it is referenced by Go import
ENV REPO=""

# Run the script to build the project
CMD sh /go/src/github.com/egjiri/cliff/.docker/build.sh

LABEL image="egjiri/cliff:0.1.3"
