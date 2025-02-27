CREATE TABLE IF NOT EXISTS urls (
	short_path VARCHAR(255) PRIMARY KEY NOT NULL,
	original_url TEXT NOT NULL,
	expiry TIMESTAMP WITHOUT TIME ZONE,
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	created_by VARCHAR(255) NOT NULL,
	modified_at TIMESTAMP WITHOUT TIME ZONE,
	modified_by VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS urls_archive (
	short_path VARCHAR(255) PRIMARY KEY NOT NULL,
	original_url TEXT NOT NULL,
	expiry TIMESTAMP WITHOUT TIME ZONE,
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	created_by VARCHAR(255) NOT NULL,
	modified_at TIMESTAMP WITHOUT TIME ZONE,
	modified_by VARCHAR(255),
	deleted_at TIMESTAMP WITHOUT TIME ZONE,
	deleted_by VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS url_access_logs (
    id SERIAL PRIMARY KEY,
    short_path VARCHAR(255) NOT NULL,
    accessed_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_accessed_at ON url_access_logs(short_path, accessed_at);

-- Index for fast lookups by original URL
CREATE INDEX IF NOT EXISTS idx_urls_original_url ON urls(original_url);

-- Index for fast deletion and expiry lookups
CREATE INDEX idx_urls_expiry_filtered
ON urls(expiry)
WHERE expiry IS NOT NULL;

-- Composite index for expiry and short_path (bulk cleanups)
CREATE INDEX IF NOT EXISTS idx_urls_expiry_short_path ON urls(expiry, short_path) WHERE expiry IS NOT NULL;;


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