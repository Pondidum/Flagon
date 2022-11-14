#!/bin/sh

set -eu

AUTH="Bearer ${GITHUB_TOKEN}"
URL="${GITHUB_API_URL}/repos/${GITHUB_REPOSITORY}/releases"

version=$(./flagon version --short)
body=$(./flagon version --changelog | awk '{printf "%s\\n", $0}' )

json=$(cat <<EOF
{
  "tag_name": "${version}",
  "name": "${version}",
  "body": "${body}",
  "draft": false
}
EOF
)


json=$(curl -X POST \
  --url "${URL}" \
  --header "Authorization: ${AUTH}" \
  --header "Accept: application/vnd.github+json" \
  -d "${json}")

upload_url=$(echo "$json" | sed -n 's,.*upload_url.*https://\(.*\){.*,\1,p')

curl -X POST \
  --url "https://${upload_url}?name=flagon" \
  --header "Authorization: ${AUTH}" \
  --header "Accept: application/vnd.github+json" \
  --header "Content-Type: application/octet-stream" \
  -d @flagon
