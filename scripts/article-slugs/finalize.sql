-- PHASE 3 / MySQL 8.0: validate repaired data, then constraints, marker LAST.
-- Phase 1 is schema.sql; phase 2 is cmd/article-slugs (local pinyin and strict
-- textual calendar-date repair). This file NEVER repairs/deletes article data.
-- STOP ALL WRITERS; back up; explicitly set @inkspace_article_seo_offline = 1
-- in this same mysql-client connection. Do not use mysql --force.
DELIMITER $$
DROP PROCEDURE IF EXISTS inkspace_finalize_article_seo$$
CREATE PROCEDURE inkspace_finalize_article_seo()
BEGIN
  DECLARE index_rows INT DEFAULT 0;
  DECLARE valid_rows INT DEFAULT 0;
  IF COALESCE(@inkspace_article_seo_offline, 0) <> 1 OR DATABASE() IS NULL THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Select the reviewed database and explicitly confirm offline writers';
  END IF;
  IF (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='articles' AND column_name IN ('slug','publish_at','created_at','updated_at','status','deleted_at','id')) <> 7 THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Run reviewed schema preparation first';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='articles' AND column_name='slug' AND data_type='varchar' AND character_maximum_length=120)
      OR NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='articles' AND column_name='publish_at' AND data_type='datetime') THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Unexpected prepared column types; manual review required';
  END IF;
  IF EXISTS (SELECT 1 FROM articles WHERE slug IS NULL OR slug='' OR OCTET_LENGTH(slug)>120
      OR REGEXP_LIKE(slug, '[^a-z0-9-]', 'c')
      OR NOT REGEXP_LIKE(slug, '^[a-z0-9]+(-[a-z0-9]+)*$', 'c') OR REGEXP_LIKE(slug, '^[0-9]+$', 'c')) THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Missing/invalid slugs; run application repair, never rewrite populated slugs';
  END IF;
  IF EXISTS (SELECT slug FROM articles GROUP BY slug HAVING COUNT(*)>1) THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Duplicate populated slugs; manual resolution required';
  END IF;
  -- Date parts are checked without driver normalization or SQL-mode changes.
  -- LAST_DAY receives a valid first-of-month, not a legacy impossible date.
  IF EXISTS (SELECT 1 FROM articles WHERE publish_at IS NULL
      OR YEAR(publish_at)<1000 OR MONTH(publish_at) NOT BETWEEN 1 AND 12
      OR DAYOFMONTH(publish_at)<1
      OR DAYOFMONTH(publish_at)>DAYOFMONTH(LAST_DAY(CONCAT(YEAR(publish_at),'-',LPAD(MONTH(publish_at),2,'0'),'-01')))) THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Publication dates not repaired; run cmd/article-slugs first';
  END IF;
  -- Validate BOTH existing constraint shapes before installing either one.
  SELECT COUNT(*), COALESCE(SUM(non_unique=0 AND seq_in_index=1 AND column_name='slug' AND sub_part IS NULL AND index_type='BTREE' AND is_visible='YES'),0)
    INTO index_rows, valid_rows FROM information_schema.statistics
    WHERE table_schema=DATABASE() AND table_name='articles' AND index_name='ux_articles_slug';
  IF index_rows<>0 AND (index_rows<>1 OR valid_rows<>1) THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid existing ux_articles_slug; manual review required';
  END IF;
  SELECT COUNT(*), COALESCE(SUM(non_unique=1 AND sub_part IS NULL AND index_type='BTREE' AND is_visible='YES' AND
      ((seq_in_index=1 AND column_name='status' AND collation='A') OR
       (seq_in_index=2 AND column_name='deleted_at' AND collation='A') OR
       (seq_in_index=3 AND column_name='publish_at' AND collation='D') OR
       (seq_in_index=4 AND column_name='id' AND collation='D'))),0)
    INTO index_rows, valid_rows FROM information_schema.statistics
    WHERE table_schema=DATABASE() AND table_name='articles' AND index_name='idx_articles_first_publication_v1';
  IF index_rows<>0 AND (index_rows<>4 OR valid_rows<>4) THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid existing publication marker; manual review required';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.statistics WHERE table_schema=DATABASE() AND table_name='articles' AND index_name='ux_articles_slug') THEN
    CREATE UNIQUE INDEX ux_articles_slug ON articles(slug);
  END IF;
  -- Completion marker is LAST, after all row/constraint validations above.
  IF NOT EXISTS (SELECT 1 FROM information_schema.statistics WHERE table_schema=DATABASE() AND table_name='articles' AND index_name='idx_articles_first_publication_v1') THEN
    CREATE INDEX idx_articles_first_publication_v1 ON articles(status, deleted_at, publish_at DESC, id DESC);
  END IF;
END$$
CALL inkspace_finalize_article_seo()$$
DROP PROCEDURE inkspace_finalize_article_seo$$
DELIMITER ;
