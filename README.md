## go-unit-test-example
Example showing common unit testing patterns for Go. The example used was derived from Andrei Avram's blog about unit testing interfaces: https://blog.andreiavram.ro/golang-unit-testing-interfaces/

I took his example and applied the patterns he mentioned in the blog to achieve 100% test coverage.

### Run Tests

run container using golang alpine image:
```
docker container run \
--env HTTP_PROXY \
--env HTTPS_PROXY \
--rm \
-it \
-v $PWD:/go/src/sandbox \
-w /go/src/sandbox \
golang:1.11-alpine \
/bin/sh
```

in the container execute:
```
export CGO_ENABLED=0
export GO111MODULE=on
apk add --no-cache git
go test -v
```