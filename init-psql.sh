#!/bin/bash
psql -v ON_ERROR_STOP=1 --username postgres <<-EOSQL
    CREATE USER registry_stat;
    CREATE DATABASE registry;
    GRANT ALL PRIVILEGES ON DATABASE registry TO registry_stat;
EOSQL

psql -v ON_ERROR_STOP=1 --username postgres <<-EOSQL
	CREATE TABLE REQUEST(
		requestID VARCHAR,
        timestamp TIMESTAMP,
		remoteAddr VARCHAR(40),
		imageName VARCHAR(128),
		imageTag VARCHAR(128),
		PRIMARY KEY(requestID)
	);
EOSQL
