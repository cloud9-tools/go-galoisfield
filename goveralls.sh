#!/bin/bash
set -eu -o pipefail
export PATH="${HOME}/gopath/bin:${PATH}"
go list ./... | ( while read pkg; do
  gocov test "$@" "$pkg" || exit 1
done ) | gocov-merge >gocov.json
goveralls -service=travis-ci -gocovdata=gocov.json
