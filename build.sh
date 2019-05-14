#!/usr/bin/env bash

build_for_platform() {
    local platform="${1}"
    if [[ -z "${platform}" ]]; then
        echo "platform required"
        return 1
    fi
    local platform_split=(${platform//\// })

    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    local output_name="${PROJECT_NAME}-${GOOS}-${GOARCH}"
    if [[ ${GOOS} = "windows" ]]; then
        output_name+='.exe'
    fi

    GOOS=${GOOS} GOARCH=${GOARCH} go build -o "${BUILD_DIR}/${output_name}" main.go
}

platforms=("linux/amd64" "darwin/amd64")
for platform in "$@"; do
    build_for_platform "${platform}"
done
