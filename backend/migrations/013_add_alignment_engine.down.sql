-- Down Migration
DROP TABLE IF EXISTS alignment_reports;
DROP TABLE IF EXISTS roadmap_dependencies;
ALTER TABLE projects DROP COLUMN IF EXISTS alignment_score;
