package generator

import "io/ioutil"

func (gen *caGen) GenDockerfile(dirName string) error {
	docker := []byte(`FROM golang:alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .

# Fetch dependencies.
RUN go get -d -v

# Do test
RUN CGO_ENABLED=0 go test -v ./...

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
			-o /go/bin/yourappname .


FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable
COPY --from=builder /go/bin/yourappname /go/bin/yourappname

# Use an unprivileged user.
USER appuser:appuser

ENV GIN_MODE=release
ENV SERVER_ECHO_PORT=9090
ENV SERVER_GIN_PORT=8090
ENV SERVER_GORILLA_MUX_PORT=7090
ENV SERVER_NET_HTTP_SERVER_MUX_PORT=6090
ENV SERVER_GRAPHQL_SERVER_MUX_PORT=5090
ENV SERVER_GRPC_PORT=50051

ENV DATABASE_HOST=
ENV DATABASE_PORT=
ENV DATABASE_USER=
ENV DATABASE_PASSWORD=
ENV DATABASE_NAME=

ENTRYPOINT ["/go/bin/yourappname"]
`)

	err := ioutil.WriteFile("./"+dirName+"/Dockerfile", docker, 0644)
	if err != nil {
		return err
	}

	return nil
}
