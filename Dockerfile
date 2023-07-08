FROM golang:1.20 AS build-stage

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /apiserver cmd/apiserver/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /controller cmd/controller/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /runner cmd/runner/main.go

####################################################################################################
# Final image
####################################################################################################

FROM alpine AS build-release-stage

WORKDIR /

RUN apk add --no-cache tini

COPY --from=build-stage /apiserver /usr/local/bin/apiserver
COPY --from=build-stage /controller /usr/local/bin/controller
COPY --from=build-stage /runner /usr/local/bin/runner

USER 999
ENTRYPOINT ["/sbin/tini", "--"]