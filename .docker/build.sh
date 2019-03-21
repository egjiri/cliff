# Create a new directory for the repository in the go path
mkdir -p /go/src/$REPO

# Navigate to that directory
cd /go/src/$REPO

# Copy all the mounted data to the repo directory
cp -R /data/* .

# Generate the data package with all the config data
go-bindata -pkg data -o /go/src/github.com/egjiri/cliff/data/go-bindata.go cli.yml

# Manipulate the cliff package source code to use the previously generated data package
go run /go/src/github.com/egjiri/cliff/.docker/main.go `pwd`/vendor/

# Run the standard go build command to generate the binary
echo $GOBUILD_FLAGS | xargs -I {} sh -c 'GOOS=$GOOS_TARGET GOARCH=$GOARCH_TARGET go build {}'
