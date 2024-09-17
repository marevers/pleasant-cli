FROM golang:1.23-alpine AS builder

ENV GO111MODULE=on

# Add our code
COPY ./ /src

# build
WORKDIR /src
RUN GOGC=off go build -v -o /pleasant-cli .

# multistage
# FROM gcr.io/distroless/static-debian11
FROM alpine:latest

COPY --from=builder /pleasant-cli /pleasant-cli

ENTRYPOINT  [ "/pleasant-cli" ]
