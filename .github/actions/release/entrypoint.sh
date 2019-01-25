#!/bin/bash
set -e
set -o pipefail

if [[ -z "$GITHUB_TOKEN" ]]; then
	echo "Set the GITHUB_TOKEN env variable."
	exit 1
fi

if [[ -z "$GITHUB_REPOSITORY" ]]; then
	echo "Set the GITHUB_REPOSITORY env variable."
	exit 1
fi

if [[ -z "$FILE_NAME" ]]; then
	echo "Set the FILE_NAME env variable."
	exit 1
fi

URI=https://api.github.com
API_VERSION=v3
API_HEADER="Accept: application/json"
AUTH_HEADER="Authorization: token ${GITHUB_TOKEN}"

main() {
    cat $GITHUB_EVENT_PATH
    
	# validate the GitHub token.
	curl -o /dev/null -sSL -H "${AUTH_HEADER}" -H "${API_HEADER}" "${URI}/repos/${GITHUB_REPOSITORY}" || { echo "Error: Invalid repo, token or network issue!";  exit 1; }

	# get the check run action.
	action=$(jq --raw-output .action "$GITHUB_EVENT_PATH")

	# If it's not synchronize or opened event return early.
	if [[ "$action" != "release" ]]; then
		# Return early we only care about synchronize or opened.
		echo "Check run has action: $action"
		echo "Want: release"
		exit 0
	fi

	# Get the release id.
	ID=$(jq --raw-output .id "$GITHUB_EVENT_PATH")

	# RENAME
	BASENAME="${FILE_NAME%.*}"
	EXTENTION="${FILE_NAME##*.}"
	if [[ EXTENTION != "" ]]; then
	    PUT_NAME="${BASENAME}_${GOOS}_${GOARCH}.${EXTENSION}"
	else
		PUT_NAME="${BASENAME}_${GOOS}_${GOARCH}"
	fi
	mv "./${FILE_NAME}" "./${PUT_NAME}"

	echo "running $GITHUB_ACTION for Release #${ID}, file ${PUT_NAME}"
    GH_ASSET="https://uploads.github.com/repos/${GITHUB_REPOSITORY}/releases/${ID}/assets?name=${PUT_NAME}"
	curl -H "${AUTH_HEADER}" -H "Content-Type: application/zip" --data-binary @"${FILE_NAME}" -X POST "${GH_ASSET}"
}

main