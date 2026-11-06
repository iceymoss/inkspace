-- Read-only post-migration checks. Run with a protected MySQL client defaults
-- file against the explicitly selected deployment database, not a test/dev DB.
-- Includes soft-deleted/private/draft reservations. Outputs counts, not content.
SELECT COUNT(*) AS article_rows,
       COALESCE(SUM(slug IS NULL OR slug = ''), 0) AS missing_slugs,
       COALESCE(SUM(slug IS NOT NULL AND
           (OCTET_LENGTH(slug) > 120 OR slug NOT REGEXP '^[a-z0-9]+(-[a-z0-9]+)*$'
            OR slug REGEXP '^[0-9]+$')), 0) AS invalid_slugs
FROM articles;

SELECT COUNT(*) AS duplicate_slug_groups FROM
  (SELECT slug FROM articles WHERE slug IS NOT NULL AND slug <> ''
   GROUP BY slug HAVING COUNT(*) > 1) AS duplicates;

SELECT INDEX_NAME, NON_UNIQUE, SEQ_IN_INDEX, COLUMN_NAME, COLLATION
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'articles'
  AND INDEX_NAME IN ('ux_articles_slug', 'idx_articles_first_publication_v1')
ORDER BY INDEX_NAME, SEQ_IN_INDEX;

-- The extended Go command does full textual calendar validation (including
-- ALLOW_INVALID_DATES history). These cheap counts are NOT a substitute for it.
SELECT COUNT(*) AS missing_or_zero_publication_dates FROM articles
WHERE publish_at IS NULL OR CAST(publish_at AS CHAR) LIKE '0000-%';

SELECT COUNT(*) AS public_articles FROM articles
WHERE status = 1 AND deleted_at IS NULL;
