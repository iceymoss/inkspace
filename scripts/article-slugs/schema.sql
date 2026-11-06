-- PHASE 1 / MySQL 8.0: guarded, repeatable preparation ONLY.
-- Back up and STOP ALL WRITERS. In this same mysql-client connection set:
-- SET @inkspace_article_seo_offline = 1;
-- Then source this file. Run cmd/article-slugs for pinyin/date repair, and
-- source finalize.sql only after reviewing its results. DDL commits implicitly.
DELIMITER $$
DROP PROCEDURE IF EXISTS inkspace_prepare_article_seo$$
CREATE PROCEDURE inkspace_prepare_article_seo()
BEGIN
  IF COALESCE(@inkspace_article_seo_offline, 0) <> 1 OR DATABASE() IS NULL THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Select the reviewed database and explicitly confirm offline writers';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema=DATABASE() AND table_name='articles') THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'articles table is missing';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='articles' AND column_name='slug') THEN
    ALTER TABLE articles ADD COLUMN slug varchar(120) CHARACTER SET ascii COLLATE ascii_bin NULL DEFAULT NULL;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='articles' AND column_name='slug' AND data_type='varchar' AND character_maximum_length=120) THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Unexpected slug column; review manually, never rewrite populated slugs';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='articles' AND column_name='publish_at') THEN
    ALTER TABLE articles ADD COLUMN publish_at datetime(3) NULL DEFAULT NULL;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name='articles' AND column_name='publish_at' AND data_type='datetime') THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Unexpected publish_at column; review manually';
  END IF;
END$$
CALL inkspace_prepare_article_seo()$$
DROP PROCEDURE inkspace_prepare_article_seo$$
DELIMITER ;
