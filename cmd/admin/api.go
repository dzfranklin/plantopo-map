package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/tidwall/sjson"
)

func handleMapnikGetStyles(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`["style"]`))
}

func handleGetStyle(w http.ResponseWriter, _ *http.Request) {
	value, err := styleRepo.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(value)
}

func handleMapnikGetStyle(w http.ResponseWriter, r *http.Request) {
	value, err := styleRepo.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	value, _ = sjson.SetBytes(value, "id", "style")
	value = bytes.ReplaceAll(value, []byte("https://mapping.plantopo.com/"), []byte("http://"+r.Host+"/"))

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(value)
}

func handleMapnikPutStyle(w http.ResponseWriter, r *http.Request) {
	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	value, _ = sjson.DeleteBytes(value, "id")
	value, _ = sjson.DeleteBytes(value, "metadata.maputnik:renderer")
	value = bytes.ReplaceAll(value, []byte("http://"+r.Host+"/"), []byte("https://mapping.plantopo.com/"))

	err = styleRepo.Set(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
