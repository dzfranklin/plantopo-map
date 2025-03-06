MAPUTNIK_VERSION := 2.1.1

SPRITE_TARGETS := icons/out/sprite.json icons/out/sprite.png icons/out/sprite@2x.json icons/out/sprite@2x.png

all: \
	example_data/openmaptiles.pmtiles \
	cmd/admin/vendor/maputnik \
	$(SPRITE_TARGETS) \
	fonts

.PHONY: clean
clean: clean-sprites clean-terrain clean-fonts
	rm -f example_data/openmaptiles.pmtiles
	rm -rf cmd/admin/vendor

example_data/openmaptiles.pmtiles:
	curl -L --fail --output example_data/openmaptiles.pmtiles \
 		https://plantopo-storage.b-cdn.net/plantopo-map-examples/openmaptiles_scotland.pmtiles

cmd/admin/vendor/maputnik:
	mkdir -p cmd/admin/vendor
	curl -L --fail --output cmd/admin/vendor/maputnik.zip \
		https://github.com/maplibre/maputnik/releases/download/v${MAPUTNIK_VERSION}/dist.zip
	unzip cmd/admin/vendor/maputnik.zip -d cmd/admin/vendor/maputnik
	rm cmd/admin/vendor/maputnik.zip
	rm cmd/admin/vendor/maputnik/dist/assets/*.map

.PHONY: clean-sprites
clean-sprites:
	rm -rf icons/out

$(SPRITE_TARGETS) &: $(wildcard icons/source/*)
	docker build --load -t icons-builder ./icons
	docker run --rm --mount type=bind,src=$(shell pwd)/icons,dst=/icons icons-builder

.PHONY: clean-fonts
clean-fonts:
	rm -rf fonts

fonts:
	mkdir -p fonts

	curl -L --fail --output fonts/noto_sans.zip \
		https://github.com/openmaptiles/fonts/releases/download/v2.0/noto-sans.zip
	unzip -q fonts/noto_sans.zip -d fonts
	rm fonts/noto_sans.zip
