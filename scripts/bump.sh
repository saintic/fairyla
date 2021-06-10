#!/bin/bash
#
# Version: 1.0
# Author:  Hiroshi.tao <me@tcw.im>
# Create:  2020-05-26
# Desc:    bump version
#

SHELL_DIR=$(cd $(dirname $0);pwd)
BASE_DIR=$(dirname $SHELL_DIR)
PKG_DIR=${BASE_DIR}/dist

BINARY=fairyla
LINUX=${BINARY}.linux-amd64
MACOS=${BINARY}.darwin-amd64
WIN=${BINARY}.windows-amd64.exe
UI=ui

SERVER_DIR=${BASE_DIR}/server
VERSION=$(grep "const version" ${SERVER_DIR}/main.go | tr -d '"' | awk '{print $NF}')
BUMP_DIR=${BASE_DIR}/.tmp-${VERSION}

usage() {
    echo $"Usage: $0 docker|pack|release"
    exit 1
}

tarpkg() {
    os=$1
    if [ "${os}" != "windows" ]; then
        tar zcvf ${PKG_DIR}/${BINARY}.${VERSION}-${os}-amd64.tar.gz ${BINARY} ${UI} NOTICE LICENSE && rm -f ${BINARY}
    else
        zip -r ${PKG_DIR}/${BINARY}.${VERSION}-${os}-amd64.zip ${BINARY}.exe ${UI} NOTICE LICENSE && rm -f ${BINARY}.exe
    fi
}

verpack() {
    mkdir ${PKG_DIR}
    set -e
    # compile frontend
    cd ${BASE_DIR}/client && yarn build --outDir ${BUMP_DIR}/${UI}
    # copy license
    cd ${BASE_DIR} && cp -f LICENSE NOTICE ${BUMP_DIR}/
    # compile backend for popular arch
    cd ${SERVER_DIR}
    make gotool && make build-linux && make build-macos && make build-windows
    # tar unpack
    cd ${BUMP_DIR}
    mv ${SERVER_DIR}/bin/${LINUX} ${BINARY}
    tarpkg linux
    mv ${SERVER_DIR}/bin/${MACOS} ${BINARY}
    tarpkg darwin
    mv ${SERVER_DIR}/bin/${WIN} ${BINARY}.exe
    tarpkg windows
    # tar assets
    tar zcvf ${PKG_DIR}/${BINARY}-assets.${VERSION}.tar.gz ${UI}
    # clean tmp
    cd ${PKG_DIR}
    rm -rf ${BUMP_DIR} ${SERVER_DIR}/bin/
    set +e
    echo "Packaged successfully"
    return 0
}

main() {
    case $1 in
    docker)
        cd ${BASE_DIR}
        docker build -t staugur/fairyla .
        ;;
    pack)
        verpack
        ;;
    release)
        cd ${BASE_DIR}
        git add . && git ci -m "bump version ${VERSION}" && git tag ${VERSION}
        ;;
    *)
        usage
        ;;
    esac
}

main $1
