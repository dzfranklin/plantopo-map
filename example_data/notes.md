# Notes on example data

## Generating openmaptiles_scotland.pmtiles:

1. Run openmaptiles
    - Checkout github.com/openmaptiles/openmaptiles
    - In .env edit MAX_ZOOM to 14
    - Run ./quickstart.sh scotland
2. Generate tilestats: `tile-join --no-tile-size-limit --force -o ./data/tiles_with_stats.mbtiles ./data/tiles.mbtiles` (tile-join is part of tippecanoe)
3. Convert to pmtiles: `pmtiles convert data/tiles_with_stats.mbtiles data/openmaptiles_scotland.pmtiles`
