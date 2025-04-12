CREATE TABLE IF NOT EXISTS geo_entities (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (
        type IN (
            'continent',
            'region',
            'bloc',
            'country',
            'territory',
            'province',
            'subnational-region',
            'state',
            'county',
            'constituency',
            'district'
        )
    )
);

CREATE TABLE IF NOT EXISTS geo_entity_meta (
    geo_entity_id INT NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY (geo_entity_id, key),
    FOREIGN KEY (geo_entity_id) REFERENCES geo_entities(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS geo_relationships (
    parent_id INT NOT NULL,
    child_id INT NOT NULL,
    since DATE NOT NULL,
    until DATE,
    PRIMARY KEY (parent_id, child_id),
    FOREIGN KEY (parent_id) REFERENCES geo_entities(id) ON DELETE RESTRICT,
    FOREIGN KEY (child_id) REFERENCES geo_entities(id) ON DELETE RESTRICT,
    CONSTRAINT no_overlapping_dates CHECK (
        until IS NULL
        OR until >= since
    )
);

CREATE TABLE IF NOT EXISTS indicators (
    id SERIAL PRIMARY KEY,
    code VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    data_type VARCHAR(20) NOT NULL DEFAULT 'numeric' CHECK (
        data_type IN ('numeric', 'categorical', 'boolean')
    )
);

CREATE TABLE IF NOT EXISTS data_sources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    url TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    last_updated DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS data_points (
    id SERIAL PRIMARY KEY,
    geo_entity_id INT NOT NULL,
    indicator_id INT NOT NULL,
    source_id INT NOT NULL,
    date DATE NOT NULL,
    numeric_value DECIMAL(20, 5),
    text_value TEXT,
    FOREIGN KEY (geo_entity_id) REFERENCES geo_entities(id) ON DELETE RESTRICT,
    FOREIGN KEY (indicator_id) REFERENCES indicators(id) ON DELETE RESTRICT,
    FOREIGN KEY (source_id) REFERENCES data_sources(id) ON DELETE RESTRICT,
    CONSTRAINT unique_data_point UNIQUE (geo_entity_id, indicator_id, source_id, date),
    CONSTRAINT check_value_type CHECK (
        (
            numeric_value IS NOT NULL
            AND text_value IS NULL
        )
        OR (
            numeric_value IS NULL
            AND text_value IS NOT NULL
        )
    )
);

CREATE INDEX idx_data_points_date ON data_points(date);

---- create above / drop below ----
DROP INDEX IF EXISTS idx_data_points_date;

DROP TABLE IF EXISTS data_points;

DROP TABLE IF EXISTS data_sources;

DROP TABLE IF EXISTS indicators;

DROP TABLE IF EXISTS geo_relationships;

DROP TABLE IF EXISTS geo_entity_meta;

DROP TABLE IF EXISTS geo_entities;
