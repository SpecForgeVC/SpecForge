-- name: CreateFeatureIntelligence :one
INSERT INTO feature_intelligence (
    feature_id,
    completeness_score,
    contract_integrity_score,
    variable_coverage_score,
    dependency_stability_score,
    drift_risk_score,
    test_coverage_score,
    llm_confidence_score,
    overall_score
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetFeatureIntelligence :one
SELECT * FROM feature_intelligence
WHERE feature_id = $1;

-- name: UpdateFeatureIntelligence :one
UPDATE feature_intelligence
SET
    completeness_score = $2,
    contract_integrity_score = $3,
    variable_coverage_score = $4,
    dependency_stability_score = $5,
    drift_risk_score = $6,
    test_coverage_score = $7,
    llm_confidence_score = $8,
    overall_score = $9,
    last_calculated_at = NOW()
WHERE feature_id = $1
RETURNING *;

-- name: DeleteFeatureIntelligence :exec
DELETE FROM feature_intelligence
WHERE feature_id = $1;
