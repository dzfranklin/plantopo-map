package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/dzfranklin/plantopo-map/internal/repos"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protomaps/go-pmtiles/pmtiles"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
)

//go:embed vendor/maputnik/dist
var MaputnikFS embed.FS

//go:embed index.html
var IndexHTML []byte

var fontFS http.FileSystem
var spritesFS http.FileSystem
var naturalEarth2ShadedReliefFS http.FileSystem

var maputnikIndexHTML []byte

var styleRepo *repos.StyleRepo

var openmaptilesServer *pmtiles.Server

var maptilerKey string

func main() {
	adminPort := "8201"

	maptilerKey = os.Getenv("MAPTILER_KEY")
	if maptilerKey == "" {
		slog.Warn("Missing MAPTILER_KEY")
	}

	maputnikAssets, err := fs.Sub(MaputnikFS, "vendor/maputnik/dist/assets")
	if err != nil {
		panic(err)
	}
	maputnikAssetFS := http.FS(maputnikAssets)
	maputnikIndexHTML, err = MaputnikFS.ReadFile("vendor/maputnik/dist/index.html")
	if err != nil {
		panic(err)
	}

	fontFS = http.Dir("fonts")
	spritesFS = http.Dir("icons/out")
	naturalEarth2ShadedReliefFS = http.Dir("natural_earth_2_shaded_relief")

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	openmaptilesServer, err = pmtiles.NewServer("", "example_data", log.New(io.Discard, "", 0), 64, "")
	if err != nil {
		panic(err)
	}
	openmaptilesServer.Start()

	styleRepo = repos.NewStyleRepo(db)

	adminSrv := http.NewServeMux()

	adminSrv.HandleFunc("GET /{$}", handleGetRoot)

	adminSrv.HandleFunc("GET /style.json", handleGetStyle)

	// maputnik static routes
	adminSrv.HandleFunc("GET /mapnik", handleGetMapnikIndex)
	adminSrv.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(maputnikAssetFS)))

	// maputnik desktop api routes
	adminSrv.HandleFunc("GET /styles", handleMapnikGetStyles)
	adminSrv.HandleFunc("GET /styles/style", handleMapnikGetStyle)
	adminSrv.HandleFunc("PUT /styles/style", handleMapnikPutStyle)

	adminSrv.HandleFunc("GET /tiles/{source}/{z}/{x}/{y}", handleGetTile)

	adminSrv.HandleFunc("GET /fonts/{fontstack}/{rangeAndExt}", handleGetFonts)

	adminSrv.Handle("GET /sprites/", http.StripPrefix("/sprites/", http.FileServer(spritesFS)))

	addr := ":" + adminPort
	fmt.Println("Admin server listening on", addr)
	if err := http.ListenAndServe(addr, adminSrv); err != nil {
		panic(err)
	}
}

func handleGetRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(IndexHTML)
}

func handleGetMapnikIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(maputnikIndexHTML)
}
