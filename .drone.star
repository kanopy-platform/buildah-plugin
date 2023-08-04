def main(ctx):
    platforms = ["amd64", "arm64"]

    volumes = [
        {
            "name": "cache",
            "path": "/go",
        }
    ]

    workspace = {"path": "/go/src/github.com/${DRONE_REPO}"}

    resources = {"requests": {"cpu": 400, "memory": "2Gi"}}

    trigger = {"branch": ["**/*"]}

    test_steps = {
        "test": append_volumes(test_step(), volumes),
    }

    pipelines = [
        {
            "kind": "pipeline",
            "type": "kubernetes",
            "name": "pre-build",
            "resources": resources,
            "steps": [append_volumes(license_step(), volumes)],
            "trigger": trigger,
            "volumes": volumes,
            "workspace": workspace,
        }
    ]
    for plat in platforms:
        pipe = {
            "kind": "pipeline",
            "type": "kubernetes",
            "name": plat,
            "platform": {"arch": plat},
            "resources": resources,
            "steps": [v for v in test_steps.values()],
            "trigger": trigger,
            "volumes": volumes,
            "workspace": workspace,
        }

        pipe = append_depends_on(pipe, ["pre-build"])

        bsnp = build("build", plat, False, False)
        bsnp = set_when(bsnp, {"event": ["pull_request"]})
        bsnp = append_depends_on(bsnp, test_steps.keys())
        bsnp = append_volumes(bsnp, volumes)
        pipe["steps"].append(bsnp)

        bs = build("publish", plat, False, True)
        bs = set_when(bs, {"event": ["push"]})
        bs = append_depends_on(bs, test_steps.keys())
        bs = append_volumes(bs, volumes)
        pipe["steps"].append(bs)

        bstp = build("publish-tag", plat, True, True)
        bstp = set_when(bstp, {"event": ["tag"]})
        bstp = append_volumes(bstp, volumes)
        bstp = append_depends_on(bstp, test_steps.keys())
        pipe["steps"].append(bstp)

        pipelines.append(pipe)

    manifest_pipeline = {
       "kind": "pipeline",
        "type": "kubernetes",
        "name": "multi-arch",
        "platform": {"arch": "arm64"},
        "resources": resources,
        "steps": [],
        "trigger": trigger,
        "volumes": volumes,
        "workspace": workspace,
        "depends_on": platforms,
    }
    
    manifestStepCommit = manifest("manifest-commit", platforms, False)
    manifestStepCommit = set_when(manifestStepCommit, {"event": ["push"]})
    manifestStepCommit = append_volumes(manifestStepCommit, volumes)
    manifest_pipeline["steps"].append(manifestStepCommit)

    manifestStepTag = manifest("manifest-tag", platforms, True)
    manifestStepTag = set_when(manifestStepTag, {"event": ["tag"]})
    manifestStepTag = append_volumes(manifestStepTag, volumes)
    manifest_pipeline["steps"].append(manifestStepTag)

    pipelines.append(manifest_pipeline)

    return pipelines

def build(name, arch, tag, publish):
    step = {
        "name": name,
        "image": "plugins/kaniko-ecr",
        "pull": "always",
        "environment": {
            "GIT_COMMIT": "${DRONE_COMMIT_SHA:0:7}",
            "PLUGIN_TYPE": "drone",
            "ARCH": arch,
        },
        "settings": {
            "repo": "${DRONE_REPO_NAME}",
            "build_args": [
                "GIT_COMMIT",
                "PLUGIN_TYPE",
                "ARCH",
            ],
            "tags": [
                get_commit_name(arch),
            ],
            "create_repository": True,
        },
    }

    if tag:
        step["settings"]["tags"].append(get_tag_name(arch))
        step["environment"]["VERSION"] = get_tag_name(arch)
        step["settings"]["build_args"].append("VERSION")
    else:
        step["settings"]["tags"].append("latest-" + arch)

    if publish:
        step["settings"]["registry"] = {"from_secret": "ecr_registry"}
        step["settings"]["access_key"] = {"from_secret": "ecr_access_key"}
        step["settings"]["secret_key"] = {"from_secret": "ecr_secret_key"}
    else:
        step["settings"]["no_push"] = True

    return step

def manifest(name, platforms, tag):
    imageRepo = "public.ecr.aws/kanopy/buildah-plugin"

    step = {
        "name": name,
        "settings": {
            "registry": {"from_secret": "ecr_registry"},
            "repo": "${DRONE_REPO_NAME}",
            "access_key": {"from_secret": "ecr_access_key"},
            "secret_key": {"from_secret": "ecr_secret_key"},
            "manifest": {
                "sources": [],
                "targets": [],
            },
        },
    }

    if tag:
        # step["image"] = imageRepo + ":" + get_tag_name("arm64")
        step["image"] = imageRepo + ":" + "v0.1.0-arm64"
        for platform in platforms:
            step["settings"]["manifest"]["sources"].append(get_tag_name(platform))
        step["settings"]["manifest"]["targets"].append(get_tag_name())
    else:
        # step["image"] = imageRepo + ":" + get_commit_name("arm64")
        step["image"] = imageRepo + ":" + "v0.1.0-arm64"
        for platform in platforms:
            step["settings"]["manifest"]["sources"].append(get_commit_name(platform))
        step["settings"]["manifest"]["targets"].append(get_commit_name())

    return step

def test_step():
    return {
        "name": "test",
        "image": "golangci/golangci-lint:v1.53.3",
        "pull": "always",
        "commands": ["make test"],
    }


def license_step():
    return {
        "name": "license-check",
        "image": "public.ecr.aws/kanopy/licensed-go:4.0.4-0.1.0",
        "commands": ["licensed cache", "licensed status"],
    }


def set_when(step, when_condition):
    when_cons = getattr(step, "when", {})
    for k, v in when_condition.items():
        when_cons[k] = v

    step["when"] = when_cons
    return step


def append_volumes(step, vols):
    volumes = getattr(step, "volumes", [])
    for i in vols:
        volumes.append(i)

    step["volumes"] = volumes
    return step


def append_depends_on(step, refs):
    deps = getattr(step, "depends_on", [])

    for ref in refs:
        deps.append(ref)

    step["depends_on"] = deps
    return step

def get_commit_name(arch=""):
    if arch:
        return "git-${DRONE_COMMIT_SHA:0:7}-" + arch
    
    return "git-${DRONE_COMMIT_SHA:0:7}"

def get_tag_name(arch=""):
    if arch:
        return "${DRONE_TAG}-" + arch
    
    return "${DRONE_TAG}"
