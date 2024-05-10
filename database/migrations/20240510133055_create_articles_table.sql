-- migrate:up
CREATE TABLE articles (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    contents TEXT NOT NULL,
    slug VARCHAR(126) NOT NULL UNIQUE,
    author_id BIGINT REFERENCES users (id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- migrate:down
DROP TABLE IF EXISTS articles;