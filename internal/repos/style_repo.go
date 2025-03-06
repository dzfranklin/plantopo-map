package repos

import (
	"context"
	"encoding/json"
	"errors"
	plantopomap "github.com/dzfranklin/plantopo-map"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tidwall/sjson"
	"reflect"
)

var emptyStyle = []byte(`{"version": 8, "sources": {}, "layers": []}`)

type StyleRepo struct {
	db *pgxpool.Pool
}

func NewStyleRepo(db *pgxpool.Pool) *StyleRepo {
	return &StyleRepo{db: db}
}

func (r *StyleRepo) Get() ([]byte, error) {
	var value []byte
	err := r.db.QueryRow(context.Background(), "SELECT value FROM style_versions ORDER BY id DESC LIMIT 1").Scan(&value)
	if errors.Is(err, pgx.ErrNoRows) {
		value = emptyStyle
	} else if err != nil {
		return nil, err
	}

	value, _ = sjson.SetBytes(value, "sources.plantopo.attribution", plantopomap.AttributionHTML)

	return value, nil
}

func (r *StyleRepo) Set(value []byte) error {
	existing, err := r.Get()
	if err != nil {
		return err
	}

	var parsedExisting interface{}
	var parsedNew interface{}
	if err := json.Unmarshal(existing, &parsedExisting); err != nil {
		return err
	}
	if err := json.Unmarshal(value, &parsedNew); err != nil {
		return err
	}

	if reflect.DeepEqual(parsedExisting, parsedNew) {
		return nil
	}

	_, err = r.db.Exec(context.Background(), "INSERT INTO style_versions (value) VALUES ($1)", value)
	return err
}
