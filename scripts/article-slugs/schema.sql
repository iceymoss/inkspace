-- REVIEW REFERENCE ONLY: do not run this file wholesale.
-- The Go migration checks information_schema before each DDL, validates existing
-- slugs and performs deterministic local-pinyin backfill BETWEEN these statements.
-- MySQL DDL commits implicitly. Stop all writers; take a backup first.
ALTER TABLE articles ADD COLUMN slug varchar(120)
  CHARACTER SET ascii COLLATE ascii_bin NULL DEFAULT NULL;

-- Only after complete nonempty/format/duplicate validation and backfill:
CREATE UNIQUE INDEX ux_articles_slug ON articles(slug);

-- Extended offline command only: add if missing, then preflight all textual
-- publish_at/created_at values and repair invalid dates in ID batches. Do not
-- replace this with a blanket COALESCE (legacy invalid dates are not only NULL).
ALTER TABLE articles ADD COLUMN publish_at datetime(3) NULL DEFAULT NULL;
-- Conceptual per-invalid-row update (the command binds ID and old value):
-- UPDATE articles SET publish_at = created_at, updated_at = updated_at WHERE id = ?;
-- Only AFTER complete publication repair; also the startup readiness marker:
CREATE INDEX idx_articles_first_publication_v1
  ON articles(status, deleted_at, publish_at DESC, id DESC);
