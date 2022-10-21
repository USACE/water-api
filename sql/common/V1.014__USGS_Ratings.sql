-- create usgs ratings table
CREATE TABLE usgs_rating (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR UNIQUE NOT NULL,
    "method" VARCHAR NULL,
    s3key VARCHAR NOT NULL
);

INSERT INTO usgs_rating (id, name, "method", s3key) VALUES
('3935298d-262c-4d29-95e4-adcf5677aa4e', 'browns-fairgrounds-tn-exsa', 'linear', 'water/ratings/usgs/exsa/browns_fairgrounds_tn.csv'),
('102cd653-d1bd-442c-b150-ab89f3cb94c1', 'daddys-hebbertsburg-tn-exsa', 'linear', 'water/ratings/usgs/exsa/daddys_hebbertsburg_tn.csv');
