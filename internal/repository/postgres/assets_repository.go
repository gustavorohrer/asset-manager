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

	orderClause := buildOrder("fa", query.SortBy, query.SortOrder)

	dataArgs := append([]any{}, args...)
	dataArgs = append(dataArgs, query.PageSize)
	limitPlaceholder := fmt.Sprintf("$%d", len(dataArgs))
	dataArgs = append(dataArgs, (query.Page-1)*query.PageSize)
	offsetPlaceholder := fmt.Sprintf("$%d", len(dataArgs))

	dataSQL := `
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
latest_component_scans AS (
	SELECT DISTINCT ON (s.componentid) s.componentid, s.id AS scanid
	FROM scan s
	JOIN component c ON c.id = s.componentid
	JOIN filtered_assets fa ON fa.id = c.assetid
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
FROM filtered_assets fa
LEFT JOIN component c ON c.assetid = fa.id
LEFT JOIN latest_component_scans lcs ON lcs.componentid = c.id
LEFT JOIN vulnerability v ON v.scanid = lcs.scanid
LEFT JOIN threat t ON t.scanid = lcs.scanid
GROUP BY fa.id, fa.name, fa.description, fa.createdat, fa.lastscan
ORDER BY ` + orderClause + `
LIMIT ` + limitPlaceholder + ` OFFSET ` + offsetPlaceholder

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

func buildFilters(query assets.ListAssetsQuery) (string, []any) {
	conditions := make([]string, 0, 5)
	args := make([]any, 0, 5)

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

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
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
