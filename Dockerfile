FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Create appuser.
ENV USER=appuser
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"


WORKDIR $GOPATH/src/calenvite_svc
COPY . .

# Fetch and install dependencies.
RUN go get -d -v

# Using go mod.
RUN go mod download
RUN go mod verify

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -ldflags="-w -s" -o /go/bin/calenvite_svc

# STEP 2 build a small image
FROM alpine:latest

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable.
COPY --from=builder /go/bin/calenvite_svc /go/bin/calenvite_svc

# Use an unprivileged user.
USER appuser:appuser

# Run the service
EXPOSE 8000
ENTRYPOINT ["/go/bin/calenvite_svc"]
