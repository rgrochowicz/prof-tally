#!/bin/bash
set -e

usage() {
	base="$(basename "$0")"
	cat <<EOUSAGE

usage: $base command [args]

This script assumes that is is run from the root of your Go package (for
example, "/go/src/app" if your GOPATH is set to "/go").

The main feature of this wrapper over the native "go" tool is that it supports
the addition of a ".godir" file in your package directory which is a plain text
file that is the import path your application expects (ie, it would be something
like "github.com/jsmith/my-cool-app").  If this file is found, the current
package is then symlinked to /go/src/<.gopath> and the commands then use that
full import path to act on it.  This allows for applications to be universally
cloned into a path like "/go/src/app" and then automatically built properly so
that inter-package imports use the existing source instead of downloading it a
second time.  This also makes the final binary have the correct name (ie,
instead of being "/go/bin/app", it can be "/go/bin/my-cool-app"), which is why
the "run" command exists here.

Available Commands:

  $base download
  $base download -u
    (equivalent to "go get -d [args] [godir]")

  $base install
  $base install -race
    (equivalent to "go install [args] [godir]")

  $base run
  $base run -app -specific -arguments
    (assumes "GOPATH/bin" is in "PATH")

EOUSAGE
}

# "shift" so that "$@" becomes the remaining arguments and can be passed along to other "go" subcommands easily
cmd="$1"
if ! shift; then
	usage >&2
	exit 1
fi

dir="$(pwd -P)"
goBin="$(basename "$dir")" # likely "app"

goDir=
if [ -f .godir ]; then
	goDir="$(cat .godir)"
	goPath="${GOPATH%%:*}" # this just grabs the first path listed in GOPATH, if there are multiple (which is the detection logic "go get" itself uses, too)
	goDirPath="$goPath/src/$goDir"
	mkdir -p "$(dirname "$goDirPath")"
	if [ ! -e "$goDirPath" ]; then
		ln -sfv "$dir" "$goDirPath"
	elif [ ! -L "$goDirPath" ]; then
		echo >&2 "error: $goDirPath already exists but is unexpectedly not a symlink!"
		exit 1
	fi
	goBin="$goPath/bin/$(basename "$goDir")"
fi

case "$cmd" in
	download)
		execCommand=( go get -v -d "$@" )
		if [ "$goDir" ]; then execCommand+=( "$goDir" ); fi
		set -x; exec "${execCommand[@]}"
		;;
		
	install)
		execCommand=( go install -v "$@" )
		if [ "$goDir" ]; then execCommand+=( "$goDir" ); fi
		set -x; exec "${execCommand[@]}"
		;;
		
	run)
		set -x; exec "$goBin" "$@"
		;;
		
	*)
		echo >&2 'error: unknown command:' "$cmd"
		usage >&2
		exit 1
		;;
esac