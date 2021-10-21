############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
ARG DOCKER_TAG=0.0.0
# checkout the project 
WORKDIR /builder
COPY . .
# Fetch dependencies.
# Using go get.
RUN go mod download
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /webhook-relayer -ldflags="-s -w -extldflags \"-static\" -X main.Version=$DOCKER_TAG" cmd/webhook-relayer/main.go
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /webhook-relayer /
# Run the hello binary.
EXPOSE 2110/tcp
ENTRYPOINT [ "/webhook-relayer" ]
CMD []
