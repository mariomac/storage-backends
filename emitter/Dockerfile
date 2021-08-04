FROM golang:1.16 as build

WORKDIR /code
COPY . .
RUN go build -o emitter cmd/emitter.go

FROM ubuntu:focal

WORKDIR /
COPY --from=build /code/emitter /emitter
CMD ["/emitter"]
