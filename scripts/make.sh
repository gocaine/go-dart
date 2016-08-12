#!/bin/bash -e

# List of bundles to create when no argument is passed
DEFAULT_BUNDLES=(
	validate-gofmt
	validate-govet
	generate
	binary

	test-unit
)

bundle() {
    local bundle="$1"; shift
    echo "---> Making bundle: $(basename "$bundle") (in $DEST)"
    source "scripts/$bundle" "$@" || exit 1
}

if [ $# -lt 1 ]; then
    bundles=(${DEFAULT_BUNDLES[@]})
else
    bundles=($@)
fi
for bundle in ${bundles[@]}; do
    export DEST=.
    ABS_DEST="$(cd "$DEST" && pwd -P)"
    bundle "$bundle"
    echo
done
