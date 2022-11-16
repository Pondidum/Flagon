#!/bin/sh

set -eu

AUTH="Bearer ${GITHUB_TOKEN}"
URL="${GITHUB_API_URL}/repos/${GITHUB_REPOSITORY}/releases"

version=$(./flagon version --short)
body=$(./flagon version --changelog | awk '{printf "%s\\n", $0}' )


release_code=$(curl -sSL \
  --url "${URL}/tags/${version}" \
  --header "Authorization: ${AUTH}" \
  --header "Accept: application/vnd.github+json" \
  -o /dev/null -w "%{http_code}")

if [ "${release_code}" = "200" ]; then
  echo "Release ${version} already exists, exiting"
  exit 0
fi

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
  --header "Content-Type: $(file -b --mime-type flagon)" \
  --data-binary @flagon
