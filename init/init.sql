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

-- Index for fast lookups by original URL
CREATE INDEX IF NOT EXISTS idx_urls_original_url ON urls(original_url);

-- Index for fast deletion and expiry lookups
CREATE INDEX IF NOT EXISTS idx_urls_expiry ON urls(expiry);

-- Composite index for expiry and short_path (bulk cleanups)
CREATE INDEX IF NOT EXISTS idx_urls_expiry_short_path ON urls(expiry, short_path);


CREATE OR REPLACE FUNCTION move_expired_urls_to_archive() RETURNS void AS $$
BEGIN
    -- Insert expired URLs into the archive table
    INSERT INTO urls_archive (short_path, original_url, expiry, created_at, created_by, modified_at, modified_by, deleted_at, deleted_by)
    SELECT short_path, original_url, expiry, created_at, created_by, modified_at, modified_by, NOW(), 'system'
    FROM urls
    WHERE expiry < NOW();

    -- Delete expired URLs from the main table
    DELETE FROM urls WHERE expiry < NOW();
END;
$$ LANGUAGE plpgsql;