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

COPY --from=build-stage /apiserver /apiserver
COPY --from=build-stage /controller /controller
COPY --from=build-stage /runner /runner

USER 65532:65532

ENTRYPOINT ["/job-bin"]