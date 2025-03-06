#!/usr/bin/env bash
set -euox pipefail

mkdir -p /icons/out
/bin/spreet /icons/source /icons/out/sprite
/bin/spreet --retina /icons/source /icons/out/sprite@2x

echo "All done"
