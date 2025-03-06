CREATE TABLE style_versions (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    inserted_at TIMESTAMPTZ DEFAULT now(),
    value bytea
);
