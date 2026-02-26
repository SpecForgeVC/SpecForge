-- Down Migration

DROP TABLE IF EXISTS variable_dependencies;
DROP TABLE IF EXISTS variable_lineage_events;
DROP TYPE IF EXISTS lineage_event_type;
DROP TYPE IF EXISTS dependency_type;
DROP TABLE IF EXISTS feature_intelligence;
ALTER TABLE roadmap_items DROP COLUMN IF EXISTS readiness_level;
