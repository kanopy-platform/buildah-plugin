FROM golang:1.20 as build
ARG VERSION="0.0.0"
ARG GIT_COMMIT
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-X 'github.com/kanopy-platform/multi-arch-images/internal/version.version=${VERSION}' -X 'github.com/kanopy-platform/multi-arch-images/internal/version.gitCommit=${GIT_COMMIT}'" -o /go/bin/app ./cmd/

FROM debian:buster-slim
RUN groupadd -r app && useradd --no-log-init -r -g app app
USER app
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]
