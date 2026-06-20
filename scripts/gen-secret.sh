#!/usr/bin/env bash
set -euo pipefail

name="${1:-}"
if [[ -z "$name" ]]; then
  echo "Usage: $0 <ENV_VAR_NAME>" >&2
  echo "Example: $0 JWT_SECRET" >&2
  exit 1
fi

password=$(openssl rand -base64 48 | tr -d '\n/+=' | head -c 40)
echo "${name}=${password}"
