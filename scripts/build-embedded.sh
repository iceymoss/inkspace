#!/bin/sh

set -eu

service="${1:-}"
package_manager="${PACKAGE_MANAGER:-bun}"

case "$service" in
  server)
    frontend="blog"
    binary="server"
    ;;
  admin)
    frontend="admin"
    binary="admin"
    ;;
  *)
    printf 'Usage: %s server|admin\n' "$0" >&2
    exit 1
    ;;
esac

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)

(
  cd "$root/web/$frontend"
  "$package_manager" run build
)

mkdir -p "$root/bin"

(
  cd "$root"
  go build -o "bin/$binary" "./cmd/$binary"
)

printf 'Built %s with embedded %s frontend: %s/bin/%s\n' "$service" "$frontend" "$root" "$binary"
