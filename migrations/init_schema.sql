CREATE TABLE short_urls ( 
	id serial, url VARCHAR ( 512 ) UNIQUE NOT NULL  , 
	hash VARCHAR ( 512 ) NOT NULL PRIMARY KEY, 
	short_base32 VARCHAR ( 64 ) UNIQUE NOT NULL, 
	short_base32_inc INTEGER, 
	url_http_status INTEGER, 
	last_check_time TIMESTAMP, 
	redirect_count INTEGER DEFAULT 0 NOT NULL 
	);

CREATE INDEX idx_short_base32 ON short_urls USING HASH (short_base32);

CREATE INDEX idx_hash ON short_urls USING HASH (hash);
