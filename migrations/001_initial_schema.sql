CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) NOT NULL UNIQUE,
    admin_level INT NOT NULL CHECK (admin_level >= 0),
    parent_id INT,
    map TEXT,
    FOREIGN KEY (parent_id) REFERENCES locations(id) ON DELETE RESTRICT
);

CREATE TABLE indicators (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL UNIQUE,
    category VARCHAR(255) NOT NULL,
    description TEXT,
    unit VARCHAR(50),
    data_type VARCHAR(20) NOT NULL DEFAULT 'numeric' CHECK (data_type IN ('numeric', 'categorical', 'boolean'))
);

CREATE TABLE data_sources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url TEXT,
    description TEXT
);

CREATE TABLE data_points (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(10) NOT NULL CHECK (entity_type IN ('location', 'region')),
    entity_id INT NOT NULL,
    indicator_id INT NOT NULL,
    source_id INT NOT NULL,
    date DATE NOT NULL,
    numeric_value DECIMAL(20,5),
    text_value TEXT,
    FOREIGN KEY (indicator_id) REFERENCES indicators(id) ON DELETE RESTRICT,
    FOREIGN KEY (source_id) REFERENCES data_sources(id) ON DELETE RESTRICT,
    CONSTRAINT unique_data_point UNIQUE (entity_type, entity_id, indicator_id, source_id, date),
    CONSTRAINT check_value_type CHECK (
        (numeric_value IS NOT NULL AND text_value IS NULL) OR
        (numeric_value IS NULL AND text_value IS NOT NULL)
    )
);

CREATE TABLE regions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(10) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE location_in_region (
    location_id INT NOT NULL,
    region_id INT NOT NULL,
    join_date DATE NOT NULL,
    exit_date DATE,
    PRIMARY KEY (location_id, region_id),
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE RESTRICT,
    FOREIGN KEY (region_id) REFERENCES regions(id) ON DELETE RESTRICT
);

---- create above / drop below ----

DROP TABLE location_in_region;
DROP TABLE regions;
DROP TABLE data_points;
DROP TABLE data_sources;
DROP TABLE indicators;
DROP TABLE locations;
