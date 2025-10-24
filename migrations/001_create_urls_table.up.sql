 CREATE TABLE urls (
     id SERIAL PRIMARY KEY,
     original_url TEXT NOT NULL,
     short_code VARCHAR(20)  NOT NULL UNIQUE,
     created_at TIMESTAMP NOT NULL,
     expires_at TIMESTAMP NOT NULL,
     click_count INTEGER DEFAULT 0

 );
 CREATE INDEX idx_short_code ON urls(short_code);
 CREATE INDEX idx_created_at ON urls(created_at);