FROM golang:latest as build-env

WORKDIR /go/src/Oceanus
ADD . /go/src/Oceanus

RUN go get -d -v ./...

RUN go build -o /go/bin/Oceanus

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/Oceanus /
CMD ["/Oceanus"]