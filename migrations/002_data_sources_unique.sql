ALTER TABLE data_sources
ADD CONSTRAINT unique_name UNIQUE (name),
ADD CONSTRAINT unique_url UNIQUE (url);

---- create above / drop below ----

ALTER TABLE data_sources
DROP CONSTRAINT unique_name,
DROP CONSTRAINT unique_url;
