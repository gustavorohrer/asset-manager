package assets

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ParseListAssetsQuery(values url.Values) (ListAssetsQuery, []QueryValidationDetail) {
	query := ListAssetsQuery{
		Page:      defaultPage,
		PageSize:  defaultPageSize,
		SortBy:    SortByCreatedAt,
		SortOrder: SortOrderDesc,
	}

	var details []QueryValidationDetail

	if rawName := strings.TrimSpace(values.Get("name")); rawName != "" {
		query.NameContains = rawName
	}

	if rawPage := strings.TrimSpace(values.Get("page")); rawPage != "" {
		page, err := strconv.Atoi(rawPage)
		if err != nil || page < 1 {
			details = append(details, QueryValidationDetail{
				Field: "page",
				Issue: "must be a positive integer",
				Value: rawPage,
			})
		} else if page > maxPage {
			details = append(details, QueryValidationDetail{
				Field: "page",
				Issue: "must be less than or equal to 10000",
				Value: rawPage,
			})
		} else {
			query.Page = page
		}
	}

	if rawPageSize := strings.TrimSpace(values.Get("pageSize")); rawPageSize != "" {
		pageSize, err := strconv.Atoi(rawPageSize)
		if err != nil || pageSize < 1 {
			details = append(details, QueryValidationDetail{
				Field: "pageSize",
				Issue: "must be a positive integer",
				Value: rawPageSize,
			})
		} else if pageSize > maxPageSize {
			details = append(details, QueryValidationDetail{
				Field: "pageSize",
				Issue: "must be less than or equal to 100",
				Value: rawPageSize,
			})
		} else {
			query.PageSize = pageSize
		}
	}

	if rawSortBy := strings.TrimSpace(values.Get("sortBy")); rawSortBy != "" {
		switch SortBy(rawSortBy) {
		case SortByCreatedAt, SortByName, SortByLastScan:
			query.SortBy = SortBy(rawSortBy)
		default:
			details = append(details, QueryValidationDetail{
				Field: "sortBy",
				Issue: "must be one of createdAt, name, lastScan",
				Value: rawSortBy,
			})
		}
	}

	if rawSortOrder := strings.TrimSpace(values.Get("sortOrder")); rawSortOrder != "" {
		order := SortOrder(strings.ToLower(rawSortOrder))
		switch order {
		case SortOrderAsc, SortOrderDesc:
			query.SortOrder = order
		default:
			details = append(details, QueryValidationDetail{
				Field: "sortOrder",
				Issue: "must be one of asc, desc",
				Value: rawSortOrder,
			})
		}
	}

	createdFrom, createdFromErr := parseRFC3339(values, "created_from")
	createdTo, createdToErr := parseRFC3339(values, "created_to")
	lastScanFrom, lastScanFromErr := parseRFC3339(values, "last_scan_from")
	lastScanTo, lastScanToErr := parseRFC3339(values, "last_scan_to")

	details = append(details, createdFromErr...)
	details = append(details, createdToErr...)
	details = append(details, lastScanFromErr...)
	details = append(details, lastScanToErr...)

	query.CreatedFrom = createdFrom
	query.CreatedTo = createdTo
	query.LastScanFrom = lastScanFrom
	query.LastScanTo = lastScanTo

	if createdFrom != nil && createdTo != nil && dateOnly(*createdFrom).After(dateOnly(*createdTo)) {
		details = append(details, QueryValidationDetail{
			Field: "created_from",
			Issue: "must be less than or equal to created_to",
			Value: createdFrom.Format(time.RFC3339),
		})
	}

	if lastScanFrom != nil && lastScanTo != nil && dateOnly(*lastScanFrom).After(dateOnly(*lastScanTo)) {
		details = append(details, QueryValidationDetail{
			Field: "last_scan_from",
			Issue: "must be less than or equal to last_scan_to",
			Value: lastScanFrom.Format(time.RFC3339),
		})
	}

	return query, details
}

func parseRFC3339(values url.Values, field string) (*time.Time, []QueryValidationDetail) {
	raw := strings.TrimSpace(values.Get(field))
	if raw == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, []QueryValidationDetail{
			{
				Field: field,
				Issue: "must be a valid RFC3339 datetime",
				Value: raw,
			},
		}
	}

	return &parsed, nil
}

func dateOnly(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
