FROM golang:1.20 as build
ARG VERSION="0.0.0"
ARG GIT_COMMIT
ARG PLUGIN_TYPE="drone"
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-X 'github.com/kanopy-platform/buildah-plugin/internal/version.version=${VERSION}' -X 'github.com/kanopy-platform/buildah-plugin/internal/version.gitCommit=${GIT_COMMIT}' -X 'github.com/kanopy-platform/buildah-plugin/internal/version.pluginType=${PLUGIN_TYPE}'" -o /go/bin/app ./cmd/
RUN go install github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cli/docker-credential-ecr-login@latest

FROM quay.io/buildah/stable:v1.31.0
RUN groupadd -r app && useradd --no-log-init -r -g app app
# Add ranges needed for buildah commands
RUN echo app:10000:65536 >> /etc/subuid
RUN echo app:10000:65536 >> /etc/subgid
# Create directory needed to store credentials
ENV HOME=/buildah
RUN mkdir -m 777 -p $HOME/.docker
RUN chown -R app $HOME
USER app
COPY --from=build /go/bin/app /
COPY --from=build /go/bin/docker-credential-ecr-login /usr/bin/
ENTRYPOINT ["/app"]
