package assets

import (
	"net/url"
	"strconv"
	"strings"
)

func ParseListAssetVulnerabilitiesQuery(values url.Values) (ListAssetVulnerabilitiesQuery, []QueryValidationDetail) {
	query := ListAssetVulnerabilitiesQuery{
		Page:     defaultPage,
		PageSize: defaultPageSize,
	}

	var details []QueryValidationDetail

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

	if rawSeverity := strings.TrimSpace(values.Get("severity")); rawSeverity != "" {
		severity := Severity(strings.ToUpper(rawSeverity))
		if !isValidSeverity(severity) {
			details = append(details, QueryValidationDetail{
				Field: "severity",
				Issue: "must be one of LOW, MEDIUM, HIGH, CRITICAL",
				Value: rawSeverity,
			})
		} else {
			query.Severity = &severity
		}
	}

	return query, details
}

func isValidSeverity(severity Severity) bool {
	switch severity {
	case SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical:
		return true
	default:
		return false
	}
}
