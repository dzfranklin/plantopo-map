FROM protomaps/go-pmtiles:v1.25.3 AS pmtiles

FROM ghcr.io/osgeo/gdal:ubuntu-small-3.10.2

RUN apt update && apt install -y --no-install-recommends g++ gdb make git python3-dev python3-pip

COPY --from=pmtiles /go-pmtiles /bin/pmtiles

RUN pip install --break-system-packages git+https://github.com/mapbox/rio-rgbify#18f196e2c7a99e08ea4af318fa3bde19b09c67cd

COPY entrypoint.sh /

WORKDIR /out
ENTRYPOINT ["/entrypoint.sh"]
