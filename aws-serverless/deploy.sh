#!/bin/bash

set -eou pipefail

SAM_DIR="$(dirname "$0")"
FUNC_DIR_RELATIVE="../function"

cd "${SAM_DIR}"

SANITY="${FUNC_DIR_RELATIVE}/main.go"
if [ ! -f "${SANITY}" ] ; then
  echo "ERROR: File not found: ${SANITY}" 1>&2
  exit 1
fi

cd "${FUNC_DIR_RELATIVE}"
GOOS=linux GOARCH=arm64 go build -o bootstrap
cd -

mv "${FUNC_DIR_RELATIVE}/bootstrap" ./
sam deploy --tags "project=s3_new_file_email_lambda" $@
rm -f ./bootstrap



