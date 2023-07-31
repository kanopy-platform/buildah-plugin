# buildah-plugin
This is a Drone compatible plugin that executes [buildah](https://github.com/containers/buildah) commands.

## Use cases
### Multi-architecture image
To build a multi-architecture image, see the following .drone.yml pipeline snippet.
- `sources` are the image tags that support a single architecture.
- `targets` are the resulting image tags supporting multiple architectures.
```
steps:
- name: manifest
  image: public.ecr.aws/kanopy/buildah-plugin:latest-amd64
  pull: always
  settings:
    registry:
      from_secret: registry
    repo: devops/${DRONE_REPO_NAME}
    access_key:
      from_secret: ecr_access_key
    secret_key:
      from_secret: ecr_secret_key
    manifest:
      sources:
      - git-${DRONE_COMMIT_SHA:0:7}-amd64
      - git-${DRONE_COMMIT_SHA:0:7}-arm64
      targets:
      - git-${DRONE_COMMIT_SHA:0:7}
```
