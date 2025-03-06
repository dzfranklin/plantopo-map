#!/usr/bin/env bash
set -euox pipefail

BBOX="-7.74 58.69 -1.64 54.61" # Scotland

if [ ! -f dem.tif ]; then
  gdal_translate -projwin $BBOX /vsicurl/https://plantopo-storage.b-cdn.net/dem/copernicus-dem-30m/copernicus-dem-30m.vrt dem.tif
fi

rio rgbify \
  --base-val -10000 --interval 0.1 \
  --min-z 0 --max-z 14 --format webp \
  --workers 8 \
  dem.tif terrain.mbtiles

pmtiles convert terrain.mbtiles terrain.pmtiles
rm terrain.mbtiles
