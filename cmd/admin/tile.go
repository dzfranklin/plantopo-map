package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	plantopomap "github.com/dzfranklin/plantopo-map"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strconv"
)

type tile struct {
	z int
	x int
	y int
}

var sources = map[string]func(w http.ResponseWriter, r *http.Request, t tile){
	"plantopo":                      plantopoSource,
	"natural_earth_2_shaded_relief": naturalEarth2ShadedReliefSource,
	"terrain":                       terrainSource,
}

func handleGetTile(w http.ResponseWriter, r *http.Request) {
	source := r.PathValue("source")

	handler, ok := sources[source]
	if !ok {
		http.Error(w, fmt.Sprintf(`source "%s" not found`, source), http.StatusNotFound)
		return
	}

	z, err := strconv.Atoi(r.PathValue("z"))
	if err != nil {
		http.Error(w, "tile not found", http.StatusNotFound)
		return
	}
	x, err := strconv.Atoi(r.PathValue("x"))
	if err != nil {
		http.Error(w, "tile not found", http.StatusNotFound)
		return
	}
	y, err := strconv.Atoi(r.PathValue("y"))
	if err != nil {
		http.Error(w, "tile not found", http.StatusNotFound)
		return
	}

	maxXY := int(math.Pow(2, float64(z))) - 1
	if x < 0 || x > maxXY || y < 0 || y > maxXY {
		http.Error(w, "tile not found", http.StatusNotFound)
		return
	}

	t := tile{z: z, x: x, y: y}

	handler(w, r, t)
}

func plantopoSource(w http.ResponseWriter, r *http.Request, t tile) {
	var buf bytes.Buffer
	out := gzip.NewWriter(&buf)

	omtStatus, _, omtTile := openmaptilesServer.Get(r.Context(), fmt.Sprintf("/openmaptiles/%d/%d/%d.mvt", t.z, t.x, t.y))
	if omtStatus != 200 && omtStatus != 204 {
		slog.Error("failed to get openmaptiles tile", "t", t, "status", omtStatus)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if len(omtTile) != 0 {
		omtR, err := gzip.NewReader(bytes.NewBuffer(omtTile))
		if err != nil {
			slog.Error("omtTile not valid gzip", "t", t, "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(out, omtR)
		if err != nil {
			slog.Error("omtTile not valid gzip", "t", t, "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	if t.z >= 9 && t.z <= 14 && maptilerKey != "" {
		contourTile, err := fetchTile(fmt.Sprintf("https://api.maptiler.com/tiles/contours-v2/%d/%d/%d.pbf?key=%s", t.z, t.x, t.y, maptilerKey))
		if err == nil {
			_, _ = out.Write(contourTile)
		} else {
			slog.Error("failed to fetch contour tile", "t", t, "error", err)
		}
	}

	_ = out.Close()
	if buf.Len() > 0 {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
		w.Header().Set("Content-Type", "application/vnd.mapbox-vector-tile")
		_, _ = w.Write(buf.Bytes())
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func naturalEarth2ShadedReliefSource(w http.ResponseWriter, _ *http.Request, t tile) {
	f, err := naturalEarth2ShadedReliefFS.Open(fmt.Sprintf("%d/%d/%d.webp", t.z, t.x, t.y))
	if err != nil {
		http.Error(w, "tile not found", http.StatusNotFound)
	}
	defer f.Close()

	w.Header().Set("Content-Type", "image/webp")

	_, _ = io.Copy(w, f)
}

func terrainSource(w http.ResponseWriter, r *http.Request, t tile) {
	var data []byte
	if maptilerKey != "" {
		var err error
		data, err = fetchTile(fmt.Sprintf("https://api.maptiler.com/tiles/terrain-rgb-v2/%d/%d/%d.webp?key=%s", t.z, t.x, t.y, maptilerKey))
		if err != nil {
			slog.Error("failed to fetch terrain tile", "t", t, "error", err)
			data = nil
		}
	}

	if len(data) == 0 {
		data = plantopomap.ZeroTerrainTile
	}

	w.Header().Set("Content-Type", "image/webp")
	_, _ = w.Write(data)
}

func fetchTile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		return nil, nil
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get %s: status %d", url, resp.StatusCode)
	}

	var reader io.ReadCloser
	defer func() {
		if reader != nil {
			_ = reader.Close()
		}
	}()

	contentEncoding := resp.Header.Get("Content-Encoding")
	switch contentEncoding {
	case "":
		reader = resp.Body
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read gzip: %s", err)
		}
	default:
		return nil, fmt.Errorf("unknown content encoding %q", contentEncoding)
	}

	return io.ReadAll(reader)
}
