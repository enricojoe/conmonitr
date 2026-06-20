#!/usr/bin/env bash
# Start the ConMonitr backend and frontend together for local development.
set -euo pipefail

cd "$(dirname "$0")"

cleanup() {
  echo
  echo "shutting down..."
  kill 0
}
trap cleanup EXIT INT TERM

echo "starting backend on :8081..."
(cd backend && go run .) &

if [ ! -d frontend/node_modules ]; then
  echo "installing frontend deps..."
  (cd frontend && npm install)
fi

echo "starting frontend on :5174..."
(cd frontend && npm run dev) &

wait
