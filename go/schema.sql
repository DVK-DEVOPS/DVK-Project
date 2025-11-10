-- Drop tables in correct dependency order
-- DROP TABLE IF EXISTS pages_fts; -- (Obsolete: was used for FTS in old schema. Remove if not needed for migration.)
DROP TABLE IF EXISTS pages;
DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username TEXT NOT NULL UNIQUE,
                                     email TEXT NOT NULL UNIQUE,
                                     password TEXT NOT NULL
);

INSERT INTO users (username, email, password)
VALUES ('admin', 'keamonk1@stud.kea.dk', '5f4dcc3b5aa765d61d8327deb882cf99')
    ON CONFLICT DO NOTHING;


CREATE TABLE IF NOT EXISTS pages (
                                     id SERIAL PRIMARY KEY,
                                     title TEXT NOT NULL UNIQUE,
                                     url TEXT NOT NULL UNIQUE,
                                     content TEXT NOT NULL,
                                     language TEXT NOT NULL CHECK (language IN ('en', 'da')),
    createdAt TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP NOT NULL
    );


INSERT INTO pages (title, url, content, language, createdAt, updatedAt) VALUES
                                                                            ('Python Asyncio Guide v2', 'progurl1.dev-v2', 'Overview of asyncio event loop, tasks, and coroutines with examples.', 'en', '2025-09-12 14:30:00', '2025-09-12 14:30:00'),
                                                                            ('Java Streams Examples v2', 'progurl2.dev-v2', 'Common stream operations, collectors, and parallel processing patterns.', 'en', '2025-09-12 14:30:00', '2025-09-12 14:30:00'),
                                                                            ('C# LINQ Recipes v2', 'progurl3.dev-v2', 'LINQ queries, deferred execution, and performance tips.', 'en', '2025-09-12 14:30:00', '2025-09-12 14:30:00'),
                                                                            ('JavaScript Promises Patterns v2', 'progurl4.dev-v2', 'Chaining, error handling, and converting callbacks to promises.', 'en', '2025-09-12 14:30:00', '2025-09-12 14:30:00'),
                                                                            ('Rust Ownership Cheatsheet v2', 'progurl5.dev-v2', 'Ownership, borrowing, lifetimes, and common borrow checker fixes.', 'en', '2025-09-12 14:30:00', '2025-09-12 14:30:00'),
                                                                            ('Go Concurrency Patterns v2', 'progurl6.dev-v2', 'Goroutines, channels, worker pools, and select usage.', 'en', '2025-09-12 14:30:00', '2025-09-12 14:30:00')
    ON CONFLICT DO NOTHING;


-- ===========================
-- FULL-TEXT SEARCH SETUP
-- ===========================
-- Add a generated column for tsvector
ALTER TABLE pages
    ADD COLUMN search_vector tsvector
        GENERATED ALWAYS AS (
            setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
            setweight(to_tsvector('english', coalesce(content, '')), 'B')
            ) STORED;

-- Create index for FTS
CREATE INDEX IF NOT EXISTS pages_search_idx ON pages USING GIN (search_vector);


-- ===========================
-- EXAMPLE FTS QUERY
-- ===========================
-- SELECT title, ts_rank_cd(search_vector, query) AS rank
-- FROM pages, to_tsquery('english', 'promises | goroutines') query
-- WHERE search_vector @@ query
-- ORDER BY rank DESC;
