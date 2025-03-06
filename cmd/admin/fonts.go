package main

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func handleGetFonts(w http.ResponseWriter, r *http.Request) {
	fontstack, err := url.PathUnescape(r.PathValue("fontstack"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fontNames := strings.Split(fontstack, ",")

	rangeAndExt := r.PathValue("rangeAndExt")
	if !strings.HasSuffix(rangeAndExt, ".pbf") {
		http.NotFound(w, r)
		return
	}
	rangeName := strings.TrimSuffix(rangeAndExt, ".pbf")

	var b bytes.Buffer

	for _, fontName := range fontNames {
		filename := fontName + "/" + rangeName + ".pbf"

		file, openErr := fontFS.Open(filename)
		if openErr != nil {
			slog.Error("failed to open font file", "filename", filename, "error", openErr)
			http.NotFound(w, r)
			return
		}

		_, readErr := b.ReadFrom(file)
		_ = file.Close()
		if readErr != nil {
			slog.Error("failed to read from font file", "filename", filename, "error", readErr)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}

	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Header().Set("Content-Length", strconv.Itoa(b.Len()))
	_, _ = w.Write(b.Bytes())
}
