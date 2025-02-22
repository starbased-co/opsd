FROM golang:1.24.0-alpine3.21 AS build
ARG VERSION="dev"

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v

RUN GOAMD64=v2 GCO_ENABLED=0 \
  go build -o /go/bin/opsd \
  -ldflags="-w -X main.version=${VERSION}" \
  -trimpath \
  -tags netgo,osusergo \
  -buildmode=pie \
  -buildvcs=false
RUN mkdir -p /run/docker/plugins

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/opsd /
COPY --from=build /run/docker/plugins /run/docker/plugins
CMD ["/opsd"]
