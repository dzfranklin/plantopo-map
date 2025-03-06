FROM rust:1.85-bookworm

ARG SPREET_REPO="https://github.com/flother/spreet"
ARG SPREET_VERSION="v0.11.0"

RUN set -eux ;\
    /bin/bash -c 'echo ""; echo ""; echo "##### Build Spreet -- $SPREET_REPO -- version $SPREET_VERSION"' >&2 ;\
    git clone --quiet --depth 1 $SPREET_REPO -b $SPREET_VERSION ;\
    cd spreet ;\
    cargo build --release ;\
    mv target/release/spreet /bin/

COPY entrypoint.sh .

ENTRYPOINT ["./entrypoint.sh"]
