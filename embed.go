package plantopomap

import (
	_ "embed"
	"strings"
)

//go:embed zero_terrain_tile.webp
var ZeroTerrainTile []byte

//go:embed attribution.html
var attributionHTMLRaw string

var AttributionHTML string

func init() {
	AttributionHTML = strings.ReplaceAll(attributionHTMLRaw, "\n", " ")
}
