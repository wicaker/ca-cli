############################
# STEP 1 build executable binary
############################
# golang alpine 1.13.5
FROM golang@sha256:0991060a1447cf648bab7f6bb60335d1243930e38420bee8fec3db1267b84cfa as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
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

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /go/bin/todolist .

############################
# STEP 2 build a small image
############################
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable
COPY --from=builder /go/bin/todolist /go/bin/todolist

# Use an unprivileged user.
USER appuser:appuser

ENV GIN_MODE=release
ENV SERVER_ECHO_PORT=9090
ENV SERVER_GIN_PORT=8090
ENV SERVER_GORILLA_MUX_PORT=7090
ENV SERVER_NET_HTTP_SERVER_MUX_PORT=6090
ENV SERVER_GRAPHQL_SERVER_MUX_PORT=5090
ENV SERVER_GRPC_PORT=50051

ENV DATABASE_HOST_PG=docker.for.mac.localhost
ENV DATABASE_PORT_PG=5432
ENV DATABASE_USER_PG=root
ENV DATABASE_PASSWORD_PG=root
ENV DATABASE_NAME_PG=ca_cli

ENV DATABASE_HOST_MYSQL=docker.for.mac.localhost
ENV DATABASE_PORT_MYSQL=3306
ENV DATABASE_USER_MYSQL=root
ENV DATABASE_PASSWORD_MYSQL=
ENV DATABASE_NAME_MYSQL=ca_cli

# Run the todolist binary.
ENTRYPOINT ["/go/bin/todolist"]