package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gustavorohrer/ecl-be-challenge/internal/assets"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssetRepository struct {
	pool *pgxpool.Pool
}

func NewAssetRepository(pool *pgxpool.Pool) *AssetRepository {
	return &AssetRepository{pool: pool}
}

func (r *AssetRepository) ListAssets(ctx context.Context, query assets.ListAssetsQuery) ([]assets.AssetSummary, int, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("begin read transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	whereClause, args := buildFilters(query)

	countSQL := `SELECT COUNT(*) FROM asset a ` + whereClause
	var total int
	if err := tx.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count assets: %w", err)
	}

	pageOrderClause := buildOrder("fa", query.SortBy, query.SortOrder)
	resultOrderClause := buildOrder("pa", query.SortBy, query.SortOrder)

	dataArgs := append([]any{}, args...)
	dataArgs = append(dataArgs, query.PageSize)
	limitPlaceholder := fmt.Sprintf("$%d", len(dataArgs))
	dataArgs = append(dataArgs, (query.Page-1)*query.PageSize)
	offsetPlaceholder := fmt.Sprintf("$%d", len(dataArgs))

	dataSQL := buildListAssetsDataSQL(
		whereClause,
		pageOrderClause,
		resultOrderClause,
		limitPlaceholder,
		offsetPlaceholder,
	)

	rows, err := tx.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query assets: %w", err)
	}
	defer rows.Close()

	result := make([]assets.AssetSummary, 0, query.PageSize)
	for rows.Next() {
		var item assets.AssetSummary
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.CreatedAt,
			&item.LastScan,
			&item.HasVulnerabilities,
			&item.HasThreats,
			&item.VulnerabilityCounts.High,
			&item.VulnerabilityCounts.Medium,
			&item.VulnerabilityCounts.Total,
			&item.ThreatCounts.High,
			&item.ThreatCounts.Medium,
			&item.ThreatCounts.Low,
			&item.ThreatCounts.Total,
		); err != nil {
			return nil, 0, fmt.Errorf("scan asset row: %w", err)
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate asset rows: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, 0, fmt.Errorf("commit read transaction: %w", err)
	}

	return result, total, nil
}

func buildListAssetsDataSQL(
	whereClause,
	pageOrderClause,
	resultOrderClause,
	limitPlaceholder,
	offsetPlaceholder string,
) string {
	return `
WITH filtered_assets AS (
	SELECT
		a.id,
		a.name,
		a.description,
		a.createdat,
		a.lastscan
	FROM asset a
	` + whereClause + `
),
paged_assets AS (
	SELECT
		fa.id,
		fa.name,
		fa.description,
		fa.createdat,
		fa.lastscan
	FROM filtered_assets fa
	ORDER BY ` + pageOrderClause + `
	LIMIT ` + limitPlaceholder + ` OFFSET ` + offsetPlaceholder + `
),
latest_component_scans AS (
	SELECT DISTINCT ON (s.componentid) s.componentid, s.id AS scanid
	FROM scan s
	JOIN component c ON c.id = s.componentid
	JOIN paged_assets pa ON pa.id = c.assetid
	ORDER BY s.componentid, s.performedat DESC, s.id DESC
),
vulnerability_counts_by_asset AS (
	SELECT
		c.assetid AS asset_id,
		COUNT(*) FILTER (WHERE v.severity = 'HIGH') AS vulnerabilities_high,
		COUNT(*) FILTER (WHERE v.severity = 'MEDIUM') AS vulnerabilities_medium,
		COUNT(*) AS vulnerabilities_total
	FROM component c
	JOIN latest_component_scans lcs ON lcs.componentid = c.id
	JOIN vulnerability v ON v.scanid = lcs.scanid
	GROUP BY c.assetid
),
threat_counts_by_asset AS (
	SELECT
		c.assetid AS asset_id,
		COUNT(*) FILTER (WHERE t.risklevel = 'HIGH') AS threats_high,
		COUNT(*) FILTER (WHERE t.risklevel = 'MEDIUM') AS threats_medium,
		COUNT(*) FILTER (WHERE t.risklevel = 'LOW') AS threats_low,
		COUNT(*) AS threats_total
	FROM component c
	JOIN latest_component_scans lcs ON lcs.componentid = c.id
	JOIN threat t ON t.scanid = lcs.scanid
	GROUP BY c.assetid
)
SELECT
	pa.id,
	pa.name,
	pa.description,
	pa.createdat,
	pa.lastscan,
	COALESCE(vca.vulnerabilities_total > 0, FALSE) AS has_vulnerabilities,
	COALESCE(tca.threats_total > 0, FALSE) AS has_threats,
	COALESCE(vca.vulnerabilities_high, 0) AS vulnerabilities_high,
	COALESCE(vca.vulnerabilities_medium, 0) AS vulnerabilities_medium,
	COALESCE(vca.vulnerabilities_total, 0) AS vulnerabilities_total,
	COALESCE(tca.threats_high, 0) AS threats_high,
	COALESCE(tca.threats_medium, 0) AS threats_medium,
	COALESCE(tca.threats_low, 0) AS threats_low,
	COALESCE(tca.threats_total, 0) AS threats_total
FROM paged_assets pa
LEFT JOIN vulnerability_counts_by_asset vca ON vca.asset_id = pa.id
LEFT JOIN threat_counts_by_asset tca ON tca.asset_id = pa.id
ORDER BY ` + resultOrderClause
}

func (r *AssetRepository) GetAssetSummary(ctx context.Context) (assets.AssetRiskSummary, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return assets.AssetRiskSummary{}, fmt.Errorf("begin read transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	summarySQL := `
SELECT
	COUNT(*) AS total,
	COUNT(*) FILTER (WHERE ` + assetHasLatestVulnerabilitiesCondition("a") + `) AS with_vulnerabilities,
	COUNT(*) FILTER (WHERE ` + assetHasLatestThreatsCondition("a") + `) AS with_threats
FROM asset a
`

	var summary assets.AssetRiskSummary
	if err := tx.QueryRow(ctx, summarySQL).Scan(
		&summary.Total,
		&summary.WithVulnerabilities,
		&summary.WithThreats,
	); err != nil {
		return assets.AssetRiskSummary{}, fmt.Errorf("query asset summary: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return assets.AssetRiskSummary{}, fmt.Errorf("commit read transaction: %w", err)
	}

	return summary, nil
}

func (r *AssetRepository) GetAssetDetails(ctx context.Context, assetID string) (assets.AssetDetails, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return assets.AssetDetails{}, fmt.Errorf("begin read transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const detailsSQL = `
WITH filtered_asset AS (
	SELECT
		a.id,
		a.name,
		a.description,
		a.createdat,
		a.lastscan
	FROM asset a
	WHERE a.id = $1
),
latest_component_scans AS (
	SELECT DISTINCT ON (s.componentid) s.componentid, s.id AS scanid
	FROM scan s
	JOIN component c ON c.id = s.componentid
	JOIN filtered_asset fa ON fa.id = c.assetid
	ORDER BY s.componentid, s.performedat DESC, s.id DESC
)
SELECT
	fa.id,
	fa.name,
	fa.description,
	fa.createdat,
	fa.lastscan,
	COALESCE(BOOL_OR(v.id IS NOT NULL), FALSE) AS has_vulnerabilities,
	COALESCE(BOOL_OR(t.id IS NOT NULL), FALSE) AS has_threats
FROM filtered_asset fa
LEFT JOIN component c ON c.assetid = fa.id
LEFT JOIN latest_component_scans lcs ON lcs.componentid = c.id
LEFT JOIN vulnerability v ON v.scanid = lcs.scanid
LEFT JOIN threat t ON t.scanid = lcs.scanid
GROUP BY fa.id, fa.name, fa.description, fa.createdat, fa.lastscan
`

	var details assets.AssetDetails
	if err := tx.QueryRow(ctx, detailsSQL, assetID).Scan(
		&details.ID,
		&details.Name,
		&details.Description,
		&details.CreatedAt,
		&details.LastScan,
		&details.HasVulnerabilities,
		&details.HasThreats,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return assets.AssetDetails{}, assets.ErrAssetNotFound
		}
		return assets.AssetDetails{}, fmt.Errorf("query asset details: %w", err)
	}

	const componentsSQL = `
	SELECT
		c.id,
		c.name,
		c.version,
		c.vendor,
		c.type,
		c.createdat,
		c.lastscan,
		c.assetid
	FROM component c
	WHERE c.assetid = $1
	ORDER BY c.name, c.id
	`

	rows, err := tx.Query(ctx, componentsSQL, assetID)
	if err != nil {
		return assets.AssetDetails{}, fmt.Errorf("query asset components: %w", err)
	}
	defer rows.Close()

	components := make([]assets.AssetComponent, 0)
	for rows.Next() {
		var component assets.AssetComponent
		if err := rows.Scan(
			&component.ID,
			&component.Name,
			&component.Version,
			&component.Vendor,
			&component.Type,
			&component.CreatedAt,
			&component.LastScan,
			&component.AssetID,
		); err != nil {
			return assets.AssetDetails{}, fmt.Errorf("scan asset component row: %w", err)
		}
		components = append(components, component)
	}

	if err := rows.Err(); err != nil {
		return assets.AssetDetails{}, fmt.Errorf("iterate asset component rows: %w", err)
	}

	details.Components = components

	if err := tx.Commit(ctx); err != nil {
		return assets.AssetDetails{}, fmt.Errorf("commit read transaction: %w", err)
	}

	return details, nil
}

func (r *AssetRepository) ListAssetVulnerabilities(ctx context.Context, assetID string, query assets.ListAssetVulnerabilitiesQuery) ([]assets.AssetVulnerability, int, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("begin read transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM asset WHERE id = $1)`, assetID).Scan(&exists); err != nil {
		return nil, 0, fmt.Errorf("check asset existence: %w", err)
	}
	if !exists {
		return nil, 0, assets.ErrAssetNotFound
	}

	severityFilter := ""
	args := []any{assetID}
	if query.Severity != nil {
		args = append(args, string(*query.Severity))
		severityFilter = fmt.Sprintf("WHERE v.severity = $%d", len(args))
	}

	countSQL := `
WITH latest_component_scans AS (
	SELECT DISTINCT ON (s.componentid) s.componentid, s.id AS scanid
	FROM scan s
	JOIN component c ON c.id = s.componentid
	WHERE c.assetid = $1
	ORDER BY s.componentid, s.performedat DESC, s.id DESC
)
SELECT COUNT(*)
FROM vulnerability v
JOIN latest_component_scans lcs ON lcs.scanid = v.scanid
` + severityFilter

	var total int
	if err := tx.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count asset vulnerabilities: %w", err)
	}

	dataArgs := append([]any{}, args...)
	dataArgs = append(dataArgs, query.PageSize)
	limitPlaceholder := fmt.Sprintf("$%d", len(dataArgs))
	dataArgs = append(dataArgs, (query.Page-1)*query.PageSize)
	offsetPlaceholder := fmt.Sprintf("$%d", len(dataArgs))

	dataSQL := `
WITH latest_component_scans AS (
	SELECT DISTINCT ON (s.componentid)
		s.componentid,
		s.id AS scanid,
		s.performedat
	FROM scan s
	JOIN component c ON c.id = s.componentid
	WHERE c.assetid = $1
	ORDER BY s.componentid, s.performedat DESC, s.id DESC
)
SELECT
	v.id,
	v.description,
	v.severity,
	v.scanid,
	c.id AS component_id,
	c.name AS component_name,
	lcs.performedat
FROM vulnerability v
JOIN latest_component_scans lcs ON lcs.scanid = v.scanid
JOIN component c ON c.id = lcs.componentid
` + severityFilter + `
ORDER BY c.name ASC, v.id ASC
LIMIT ` + limitPlaceholder + ` OFFSET ` + offsetPlaceholder

	rows, err := tx.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query asset vulnerabilities: %w", err)
	}
	defer rows.Close()

	result := make([]assets.AssetVulnerability, 0, query.PageSize)
	for rows.Next() {
		var item assets.AssetVulnerability
		if err := rows.Scan(
			&item.ID,
			&item.Description,
			&item.Severity,
			&item.ScanID,
			&item.ComponentID,
			&item.ComponentName,
			&item.PerformedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan asset vulnerability row: %w", err)
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate asset vulnerability rows: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, 0, fmt.Errorf("commit read transaction: %w", err)
	}

	return result, total, nil
}

func (r *AssetRepository) ListAssetThreats(ctx context.Context, assetID string, query assets.ListAssetThreatsQuery) ([]assets.AssetThreat, int, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("begin read transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM asset WHERE id = $1)`, assetID).Scan(&exists); err != nil {
		return nil, 0, fmt.Errorf("check asset existence: %w", err)
	}
	if !exists {
		return nil, 0, assets.ErrAssetNotFound
	}

	riskLevelFilter := ""
	args := []any{assetID}
	if query.RiskLevel != nil {
		args = append(args, string(*query.RiskLevel))
		riskLevelFilter = fmt.Sprintf("WHERE t.risklevel = $%d", len(args))
	}

	countSQL := `
WITH latest_component_scans AS (
	SELECT DISTINCT ON (s.componentid) s.componentid, s.id AS scanid
	FROM scan s
	JOIN component c ON c.id = s.componentid
	WHERE c.assetid = $1
	ORDER BY s.componentid, s.performedat DESC, s.id DESC
)
SELECT COUNT(*)
FROM threat t
JOIN latest_component_scans lcs ON lcs.scanid = t.scanid
` + riskLevelFilter

	var total int
	if err := tx.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count asset threats: %w", err)
	}

	dataArgs := append([]any{}, args...)
	dataArgs = append(dataArgs, query.PageSize)
	limitPlaceholder := fmt.Sprintf("$%d", len(dataArgs))
	dataArgs = append(dataArgs, (query.Page-1)*query.PageSize)
	offsetPlaceholder := fmt.Sprintf("$%d", len(dataArgs))

	dataSQL := `
WITH latest_component_scans AS (
	SELECT DISTINCT ON (s.componentid)
		s.componentid,
		s.id AS scanid,
		s.performedat
	FROM scan s
	JOIN component c ON c.id = s.componentid
	WHERE c.assetid = $1
	ORDER BY s.componentid, s.performedat DESC, s.id DESC
)
SELECT
	t.id,
	t.description,
	t.risklevel,
	t.type,
	t.scanid,
	c.id AS component_id,
	c.name AS component_name,
	lcs.performedat
FROM threat t
JOIN latest_component_scans lcs ON lcs.scanid = t.scanid
JOIN component c ON c.id = lcs.componentid
` + riskLevelFilter + `
ORDER BY c.name ASC, t.id ASC
LIMIT ` + limitPlaceholder + ` OFFSET ` + offsetPlaceholder

	rows, err := tx.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query asset threats: %w", err)
	}
	defer rows.Close()

	result := make([]assets.AssetThreat, 0, query.PageSize)
	for rows.Next() {
		var item assets.AssetThreat
		if err := rows.Scan(
			&item.ID,
			&item.Description,
			&item.RiskLevel,
			&item.Type,
			&item.ScanID,
			&item.ComponentID,
			&item.ComponentName,
			&item.PerformedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan asset threat row: %w", err)
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate asset threat rows: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, 0, fmt.Errorf("commit read transaction: %w", err)
	}

	return result, total, nil
}

func (r *AssetRepository) UpdateAsset(ctx context.Context, assetID string, input assets.UpdateAssetInput) (assets.AssetUpdated, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return assets.AssetUpdated{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	setClauses := make([]string, 0, 3)
	args := make([]any, 0, 4)

	if input.Name != nil {
		args = append(args, *input.Name)
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", len(args)))
	}
	if input.Description != nil {
		args = append(args, *input.Description)
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", len(args)))
	}
	if input.LastScanSet {
		if input.LastScan == nil {
			setClauses = append(setClauses, "lastscan = NULL")
		} else {
			args = append(args, dateForSQL(*input.LastScan))
			setClauses = append(setClauses, fmt.Sprintf("lastscan = $%d::date", len(args)))
		}
	}

	if len(setClauses) == 0 {
		return assets.AssetUpdated{}, fmt.Errorf("no fields to update")
	}

	args = append(args, assetID)
	updateSQL := `
UPDATE asset
SET ` + strings.Join(setClauses, ", ") + `
WHERE id = $` + fmt.Sprintf("%d", len(args)) + `
RETURNING id, name, description, createdat, lastscan
`

	var updated assets.AssetUpdated
	if err := tx.QueryRow(ctx, updateSQL, args...).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.CreatedAt,
		&updated.LastScan,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return assets.AssetUpdated{}, assets.ErrAssetNotFound
		}
		return assets.AssetUpdated{}, fmt.Errorf("update asset: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return assets.AssetUpdated{}, fmt.Errorf("commit transaction: %w", err)
	}

	return updated, nil
}

func (r *AssetRepository) DeleteAsset(ctx context.Context, assetID string) (assets.AssetDeleted, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return assets.AssetDeleted{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var lockedID string
	if err := tx.QueryRow(ctx, `SELECT id FROM asset WHERE id = $1 FOR UPDATE`, assetID).Scan(&lockedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return assets.AssetDeleted{}, assets.ErrAssetNotFound
		}
		return assets.AssetDeleted{}, fmt.Errorf("lock asset: %w", err)
	}

	if _, err := tx.Exec(ctx, `
DELETE FROM threat t
USING scan s, component c
WHERE t.scanid = s.id
  AND s.componentid = c.id
  AND c.assetid = $1
`, assetID); err != nil {
		return assets.AssetDeleted{}, fmt.Errorf("delete asset threats: %w", err)
	}

	if _, err := tx.Exec(ctx, `
DELETE FROM vulnerability v
USING scan s, component c
WHERE v.scanid = s.id
  AND s.componentid = c.id
  AND c.assetid = $1
`, assetID); err != nil {
		return assets.AssetDeleted{}, fmt.Errorf("delete asset vulnerabilities: %w", err)
	}

	if _, err := tx.Exec(ctx, `
DELETE FROM scan s
USING component c
WHERE s.componentid = c.id
  AND c.assetid = $1
`, assetID); err != nil {
		return assets.AssetDeleted{}, fmt.Errorf("delete asset scans: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM component WHERE assetid = $1`, assetID); err != nil {
		return assets.AssetDeleted{}, fmt.Errorf("delete asset components: %w", err)
	}

	var deletedID string
	if err := tx.QueryRow(ctx, `DELETE FROM asset WHERE id = $1 RETURNING id`, assetID).Scan(&deletedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return assets.AssetDeleted{}, assets.ErrAssetNotFound
		}
		return assets.AssetDeleted{}, fmt.Errorf("delete asset: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return assets.AssetDeleted{}, fmt.Errorf("commit transaction: %w", err)
	}

	return assets.AssetDeleted{
		ID:      deletedID,
		Deleted: true,
	}, nil
}

func buildFilters(query assets.ListAssetsQuery) (string, []any) {
	conditions := make([]string, 0, 8)
	args := make([]any, 0, 8)

	add := func(condition string, value any) {
		args = append(args, value)
		placeholder := fmt.Sprintf("$%d", len(args))
		conditions = append(conditions, fmt.Sprintf(condition, placeholder))
	}

	if query.NameContains != "" {
		add("a.name ILIKE %s ESCAPE E'\\\\'", "%"+escapeLikeLiteral(query.NameContains)+"%")
	}
	if query.CreatedFrom != nil {
		add("a.createdat >= %s::date", dateForSQL(*query.CreatedFrom))
	}
	if query.CreatedTo != nil {
		add("a.createdat <= %s::date", dateForSQL(*query.CreatedTo))
	}
	if query.LastScanFrom != nil {
		add("a.lastscan IS NOT NULL AND a.lastscan >= %s::date", dateForSQL(*query.LastScanFrom))
	}
	if query.LastScanTo != nil {
		add("a.lastscan IS NOT NULL AND a.lastscan <= %s::date", dateForSQL(*query.LastScanTo))
	}
	if query.HasVulnerabilities != nil {
		add(assetHasLatestVulnerabilitiesCondition("a")+" = %s", *query.HasVulnerabilities)
	}
	if query.HasThreats != nil {
		add(assetHasLatestThreatsCondition("a")+" = %s", *query.HasThreats)
	}
	if query.HasFindings != nil {
		add(assetHasLatestFindingsCondition("a")+" = %s", *query.HasFindings)
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

func assetHasLatestFindingsCondition(assetAlias string) string {
	return "(" + assetHasLatestVulnerabilitiesCondition(assetAlias) + " OR " + assetHasLatestThreatsCondition(assetAlias) + ")"
}

func assetHasLatestVulnerabilitiesCondition(assetAlias string) string {
	return `
EXISTS (
	SELECT 1
	FROM component c
	JOIN (
		SELECT DISTINCT ON (s.componentid) s.componentid, s.id AS scanid
		FROM scan s
		JOIN component c2 ON c2.id = s.componentid
		WHERE c2.assetid = ` + assetAlias + `.id
		ORDER BY s.componentid, s.performedat DESC, s.id DESC
	) latest_component_scans ON latest_component_scans.componentid = c.id
	JOIN vulnerability v ON v.scanid = latest_component_scans.scanid
	WHERE c.assetid = ` + assetAlias + `.id
)`
}

func assetHasLatestThreatsCondition(assetAlias string) string {
	return `
EXISTS (
	SELECT 1
	FROM component c
	JOIN (
		SELECT DISTINCT ON (s.componentid) s.componentid, s.id AS scanid
		FROM scan s
		JOIN component c2 ON c2.id = s.componentid
		WHERE c2.assetid = ` + assetAlias + `.id
		ORDER BY s.componentid, s.performedat DESC, s.id DESC
	) latest_component_scans ON latest_component_scans.componentid = c.id
	JOIN threat t ON t.scanid = latest_component_scans.scanid
	WHERE c.assetid = ` + assetAlias + `.id
)`
}

func buildOrder(alias string, sortBy assets.SortBy, sortOrder assets.SortOrder) string {
	orderDirection := "DESC"
	if sortOrder == assets.SortOrderAsc {
		orderDirection = "ASC"
	}

	switch sortBy {
	case assets.SortByName:
		return alias + ".name " + orderDirection + ", " + alias + ".id ASC"
	case assets.SortByLastScan:
		return alias + ".lastscan " + orderDirection + " NULLS LAST, " + alias + ".id ASC"
	default:
		return alias + ".createdat " + orderDirection + ", " + alias + ".id ASC"
	}
}

func escapeLikeLiteral(value string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`%`, `\%`,
		`_`, `\_`,
	)
	return replacer.Replace(value)
}

func dateForSQL(value time.Time) string {
	year, month, day := value.Date()
	return fmt.Sprintf("%04d-%02d-%02d", year, int(month), day)
}
