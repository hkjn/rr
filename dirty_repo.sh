#!/bin/bash
#
# Checks that all files under REPOBASE are:
# 1. Directories
# 2. With git repos inside
# 3. With a clean working tree
#
# TODO(hkjn): Reimplement git-wtf.rb calls in lib/ bash scripts, to
# avoid wrapping it in shell call here (and having dependency on
# ruby).
# TODO(hkjn): Also look for LICENSE, README.md?
#
set -euo pipefail

GOPATH=${GOPATH:-""}
[ "$GOPATH" ] || { echo "[$0] FATAL: GOPATH is not set." >&2; exit 1; }
REPOBASE="$GOPATH/src/hkjn.me"

load() {
	local p="$GOPATH/src/hkjn.me/lib/$1"
	source "$p" 2>/dev/null || { echo "[$0] FATAL: Couldn't find $p." >&2; exit 1; }	
}
load "logging.sh"

which ruby 1>/dev/null || fatal "No 'ruby' found on PATH."
which git-wtf.rb 1>/dev/null || fatal "No 'git-wtf.rb' found on PATH."

check() {
	cd "$1"
	local dirty=0
	for d in $(ls); do
	if [ ! -d "$d" ]; then
			fatal "Not a directory: '$d'"
		fi
		cd "$d"
		if [ ! -d ".git" ]; then
			fatal "Not a git repo: '$d/.git' doesn't exist"
		fi

		if ! msg="$(git-wtf.rb 2>&1)"; then
			error "Dirty tree in '$d' repo:\n'$msg'"
			dirty=$(($dirty + 1))
		fi
		cd ..
	done
	[ $dirty -eq 0 ] || error "There were $dirty dirty repos."
	return $dirty
}

check "$REPOBASE"
