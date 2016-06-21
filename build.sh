#!/bin/bash

set -euo pipefail

LIB="$GOPATH/src/hkjn.me/lib/"

load() {
	source "$LIB/$1" 2>/dev/null || {
		echo "[$0] FATAL: Couldn't find '$LIB/$1'." >&2
		exit 1
	}
}
load "logging.sh"

info "Building rr using Dockerfile.build.."
docker build -t rr-build -f Dockerfile.build .
docker run --rm -it -v $(pwd):/build rr-build
info "Building hkjn/rr image.."
docker build -t hkjn/rr .
#info "Cleaning up.."
#docker rmi rr-build
