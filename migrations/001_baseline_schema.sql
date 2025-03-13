-- migrate:up
-- SQL statements to apply the migration (e.g., CREATE TABLE, ALTER TABLE, etc.)
CREATE TABLE data_points (
    id SERIAL PRIMARY KEY,
    location_id INT NOT NULL,
    indicator_id INT NOT NULL,
    source_id INT NOT NULL,
    date DATE NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    is_estimated BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE,
    FOREIGN KEY (indicator_id) REFERENCES indicators(id) ON DELETE CASCADE,
    FOREIGN KEY (source_id) REFERENCES data_sources(id) ON DELETE CASCADE
);

CREATE TABLE data_sources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url TEXT,
    description TEXT
);

CREATE TABLE indicators (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    unit VARCHAR(50),
    description TEXT
);

CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    admin_level INT NOT NULL CHECK (admin_level >= 0),
    parent_id INT,
    code VARCHAR(10) UNIQUE,
    map TEXT,
    FOREIGN KEY (parent_id) REFERENCES locations(id) ON DELETE SET NULL
);

CREATE TABLE regions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL 
);

CREATE TABLE location_regions (
    location_id INT NOT NULL,
    region_id INT NOT NULL, 
    start_date DATE NOT NULL,
    end_date DATE,
    PRIMARY KEY (location_id, region_id),
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE,
    FOREIGN KEY (region_id) REFERENCES regions(id) ON DELETE CASCADE
);

---- create above / drop below ----

-- migrate:down
-- SQL statements to undo the migration (e.g., DROP TABLE, ALTER TABLE DROP COLUMN, etc.)

DROP TABLE location_regions;
DROP TABLE regions;
DROP TABLE locations;
DROP TABLE indicators;
DROP TABLE data_sources;
DROP TABLE data_points;
