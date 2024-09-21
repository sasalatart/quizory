#!/bin/sh

# Function: wait_for_service
# Description: Waits for a service to become healthy by periodically checking a specified URL.
# Parameters:
#   $1 - URL to check
#   $2 - Timeout in seconds (default: 60)
#   $3 - Interval between checks in seconds (default: 2)
wait_for_service() {
  url=$1
  timeout=${2:-60}    # Default timeout is 60 seconds
  interval=${3:-2}    # Default interval is 2 seconds
  elapsed=0

  echo "Waiting for service at '$url' to be healthy..."

  while true; do
    if curl -s --fail "$url" > /dev/null; then
      echo "Service at '$url' is healthy."
      break
    else
      if [ "$elapsed" -ge "$timeout" ]; then
        echo "Service at '$url' did not become healthy after $timeout seconds."
        return 1
      fi
      echo "Service not healthy yet. Waiting for $interval seconds..."
      sleep "$interval"
      elapsed=$((elapsed + interval))
    fi
  done
}
