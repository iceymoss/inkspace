-- REVIEW REFERENCE ONLY: do not run this file wholesale.
-- The Go migration checks information_schema before each DDL, validates existing
-- slugs and performs deterministic local-pinyin backfill BETWEEN these statements.
-- MySQL DDL commits implicitly. Stop all writers; take a backup first.
ALTER TABLE articles ADD COLUMN slug varchar(120)
  CHARACTER SET ascii COLLATE ascii_bin NULL DEFAULT NULL;

-- Only after complete nonempty/format/duplicate validation and backfill:
CREATE UNIQUE INDEX ux_articles_slug ON articles(slug);
