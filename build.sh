#!/bin/bash

CURRENT_DIR='./bin/current'
OS_AND_ARCH=('aix-ppc' 'darwin-amd64' 'darwin-arm64' 'dragonfly-amd64' 'freebsd-386' 'freebsd-amd64' 'freebsd-arm64' 'freebsd-arm' 'freebsd-riscv64' 'illumos-amd64' 'linux-386' 'linux-amd64' 'linux-arm64' 'linux-arm' 'linux-arm6l' 'linux-long64' 'linux-mips' 'linux-mips64' 'linux-mips64le' 'linux-mipsle' 'linux-ppc64' 'linux-ppc64le' 'linux-riscv64' 'linux-s390x' 'netbsd-386' 'netbsd-amd64' 'netbsd-arm64' 'netbsd-arm' 'openbsd-386' 'openbsd-amd64' 'openbsd-arm64' 'openbsd-arm' 'plan9-386' 'plan9-amd64' 'plan9-arm' 'solaris-amd64' 'windows-386' 'windows-amd64' 'windows-arm64' 'windows-arm')
BINARY_FILE_NAME="archive.today"

if [ ! -d "CURRENT_DIR" ]; then
    mkdir -p "${CURRENT_DIR}"
fi
for os_arch in "${OS_AND_ARCH[@]}"
do :
    CURRENT_OS_ARCH_DIR="${CURRENT_DIR}/${os_arch}"
    if [ ! -d "CURRENT_OS_ARCH_DIR" ]; then
        mkdir -p "${CURRENT_OS_ARCH_DIR}"
    fi
    #echo "${os_arch}"
    # shellcheck disable=SC2162
    IFS=- read GOOS_1 GOARCH_1 <<< "${os_arch}"
    echo "${os_arch} → GOOS_1: '${GOOS_1}', GOARCH_1: '${GOARCH_1}'"
    GOOS="${GOOS_1}" GOARCH="${GOARCH_1}" build -v "${BINARY_FILE_NAME}" && \
    mv "./${BINARY_FILE_NAME}" "${CURRENT_OS_ARCH_DIR}/${BINARY_FILE_NAME}"
    if [ "$GOOS_1" = "windows" ]; then
        mv "${CURRENT_OS_ARCH_DIR}/${BINARY_FILE_NAME}" "${CURRENT_OS_ARCH_DIR}/${BINARY_FILE_NAME}.exe"
    fi

done