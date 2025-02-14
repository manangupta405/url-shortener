CREATE TABLE IF NOT EXISTS urls (
	short_path VARCHAR(255) PRIMARY KEY NOT NULL,
	original_url TEXT NOT NULL,
	expiry TIMESTAMP,
	created_at TIMESTAMP NOT NULL,
	created_by VARCHAR(255) NOT NULL,
	modified_at TIMESTAMP,
	modified_by VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS urls_archive (
	short_path VARCHAR(255) PRIMARY KEY NOT NULL,
	original_url TEXT NOT NULL,
	expiry TIMESTAMP,
	created_at TIMESTAMP NOT NULL,
	created_by VARCHAR(255) NOT NULL,
	modified_at TIMESTAMP,
	modified_by VARCHAR(255),
	deleted_at TIMESTAMP,
	deleted_by VARCHAR(255)
);