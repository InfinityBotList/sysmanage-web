# Only use when building the CI, in other cases, don't
all:
	CGO_ENABLED=0 go build -v 
