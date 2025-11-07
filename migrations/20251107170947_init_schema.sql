-- +goose Up
-- +goose StatementBegin

-- Drop triggers if they exist
DROP TRIGGER IF EXISTS pages_ai;
DROP TRIGGER IF EXISTS pages_ad;
DROP TRIGGER IF EXISTS pages_au;

-- Create tables (safe if they exist)
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS pages (
  title TEXT PRIMARY KEY UNIQUE,
  url TEXT NOT NULL UNIQUE,
  content TEXT NOT NULL,
  language TEXT NOT NULL CHECK(language IN ('en', 'da')) DEFAULT 'en',
  createdAt DATETIME NOT NULL,
  updatedAt DATETIME NOT NULL
);

-- FTS tables
CREATE VIRTUAL TABLE IF NOT EXISTS pages_fts USING fts5(
  title, url, content, language, createdAt, updatedAt,
  content='pages', content_rowid='rowid'
);

CREATE TABLE IF NOT EXISTS pages_fts_data(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS pages_fts_idx(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS pages_fts_docsize(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS pages_fts_config(k PRIMARY KEY, v) WITHOUT ROWID;

-- Create triggers
CREATE TRIGGER pages_ai AFTER INSERT ON pages BEGIN
  INSERT INTO pages_fts(rowid, title, url, content, language, createdAt, updatedAt)
  VALUES (new.rowid, new.title, new.url, new.content, new.language, new.createdAt, new.updatedAt);
END;

CREATE TRIGGER pages_ad AFTER DELETE ON pages BEGIN
  INSERT INTO pages_fts(pages_fts, rowid, title, url, content, language, createdAt, updatedAt)
  VALUES('delete', old.rowid, old.title, old.url, old.content, old.language, old.createdAt, old.updatedAt);
END;

CREATE TRIGGER pages_au AFTER UPDATE ON pages BEGIN
  INSERT INTO pages_fts(pages_fts, rowid, title, url, content, language, createdAt, updatedAt)
  VALUES('delete', old.rowid, old.title, old.url, old.content, old.language, old.createdAt, old.updatedAt);
  INSERT INTO pages_fts(rowid, title, url, content, language, createdAt, updatedAt)
  VALUES (new.rowid, new.title, new.url, new.content, new.language, new.createdAt, new.updatedAt);
END;

-- +goose StatementEnd
