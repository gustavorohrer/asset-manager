package assets

import (
	"net/url"
	"strconv"
	"strings"
)

func ParseListAssetThreatsQuery(values url.Values) (ListAssetThreatsQuery, []QueryValidationDetail) {
	query := ListAssetThreatsQuery{
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

	if rawRiskLevel := strings.TrimSpace(values.Get("riskLevel")); rawRiskLevel != "" {
		riskLevel := RiskLevel(strings.ToUpper(rawRiskLevel))
		if !isValidRiskLevel(riskLevel) {
			details = append(details, QueryValidationDetail{
				Field: "riskLevel",
				Issue: "must be one of LOW, MEDIUM, HIGH",
				Value: rawRiskLevel,
			})
		} else {
			query.RiskLevel = &riskLevel
		}
	}

	return query, details
}

func isValidRiskLevel(riskLevel RiskLevel) bool {
	switch riskLevel {
	case RiskLevelLow, RiskLevelMedium, RiskLevelHigh:
		return true
	default:
		return false
	}
}
