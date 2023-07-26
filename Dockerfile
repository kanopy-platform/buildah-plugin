FROM golang:1.20 as build
ARG VERSION="0.0.0"
ARG GIT_COMMIT
ARG PLUGIN_TYPE="drone"
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-X 'github.com/kanopy-platform/buildah-plugin/internal/version.version=${VERSION}' -X 'github.com/kanopy-platform/buildah-plugin/internal/version.gitCommit=${GIT_COMMIT}' -X 'github.com/kanopy-platform/buildah-plugin/internal/version.pluginType=${PLUGIN_TYPE}'" -o /go/bin/app ./cmd/

FROM debian:bookworm-slim as ecr-login
ARG ECR_LOGIN_VERSION="0.7.1"
ARG ARCH="amd64"
RUN apt-get update && apt-get install -y wget
RUN wget -O /docker-credential-ecr-login https://amazon-ecr-credential-helper-releases.s3.us-east-2.amazonaws.com/${ECR_LOGIN_VERSION}/linux-${ARCH}/docker-credential-ecr-login
RUN chmod +x /docker-credential-ecr-login

FROM quay.io/buildah/stable:v1.31.0
USER build
RUN mkdir -p /home/build/.docker
COPY --from=ecr-login /docker-credential-ecr-login /usr/bin/
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]
