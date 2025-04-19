ALTER TABLE
    data_points
ALTER COLUMN
    numeric_value TYPE DECIMAL(35, 5);

---- create above / drop below ----
ALTER TABLE
    data_points
ALTER COLUMN
    numeric_value TYPE DECIMAL(20, 5);
