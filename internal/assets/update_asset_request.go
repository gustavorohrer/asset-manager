package assets

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"
)

const maxAssetDescriptionLength = 10000

func ParseUpdateAssetRequestBody(body []byte) (UpdateAssetInput, []QueryValidationDetail) {
	var input UpdateAssetInput
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return input, []QueryValidationDetail{
			{
				Field: "body",
				Issue: "must be a non-empty JSON object",
				Value: "",
			},
		}
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(trimmed, &raw); err != nil {
		return input, []QueryValidationDetail{
			{
				Field: "body",
				Issue: "must be a valid JSON object",
				Value: "",
			},
		}
	}

	var details []QueryValidationDetail
	allowedCount := 0
	for field, rawValue := range raw {
		switch field {
		case "name":
			allowedCount++
			if isJSONNull(rawValue) {
				details = append(details, QueryValidationDetail{
					Field: "name",
					Issue: "must be a string, null is not allowed",
					Value: "null",
				})
				continue
			}

			var value string
			if err := json.Unmarshal(rawValue, &value); err != nil {
				details = append(details, QueryValidationDetail{
					Field: "name",
					Issue: "must be a string",
					Value: string(rawValue),
				})
				continue
			}

			value = strings.TrimSpace(value)
			if value == "" {
				details = append(details, QueryValidationDetail{
					Field: "name",
					Issue: "must not be empty",
					Value: "",
				})
				continue
			}
			if len(value) > 255 {
				details = append(details, QueryValidationDetail{
					Field: "name",
					Issue: "must be less than or equal to 255 characters",
					Value: value,
				})
				continue
			}
			input.Name = &value
		case "description":
			allowedCount++
			if isJSONNull(rawValue) {
				details = append(details, QueryValidationDetail{
					Field: "description",
					Issue: "must be a string, null is not allowed",
					Value: "null",
				})
				continue
			}

			var value string
			if err := json.Unmarshal(rawValue, &value); err != nil {
				details = append(details, QueryValidationDetail{
					Field: "description",
					Issue: "must be a string",
					Value: string(rawValue),
				})
				continue
			}
			if len(value) > maxAssetDescriptionLength {
				details = append(details, QueryValidationDetail{
					Field: "description",
					Issue: "must be less than or equal to 10000 characters",
					Value: value,
				})
				continue
			}
			input.Description = &value
		case "lastScan":
			allowedCount++
			input.LastScanSet = true

			if isJSONNull(rawValue) {
				input.LastScan = nil
				continue
			}

			var value string
			if err := json.Unmarshal(rawValue, &value); err != nil {
				details = append(details, QueryValidationDetail{
					Field: "lastScan",
					Issue: "must be an RFC3339 string or null",
					Value: string(rawValue),
				})
				continue
			}

			parsed, err := time.Parse(time.RFC3339, value)
			if err != nil {
				details = append(details, QueryValidationDetail{
					Field: "lastScan",
					Issue: "must be an RFC3339 datetime",
					Value: value,
				})
				continue
			}
			input.LastScan = &parsed
		default:
			details = append(details, QueryValidationDetail{
				Field: field,
				Issue: "is not allowed",
				Value: string(rawValue),
			})
		}
	}

	if allowedCount == 0 {
		details = append(details, QueryValidationDetail{
			Field: "body",
			Issue: "must include at least one updatable field: name, description, lastScan",
			Value: "",
		})
	}

	if input.Name == nil && input.Description == nil && !input.LastScanSet {
		hasFieldError := false
		for _, detail := range details {
			switch detail.Field {
			case "name", "description", "lastScan":
				hasFieldError = true
				break
			}
		}
		if !hasFieldError {
			details = append(details, QueryValidationDetail{
				Field: "body",
				Issue: "must include at least one valid updatable field",
				Value: "",
			})
		}
	}

	return input, details
}

func isJSONNull(value json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(value), []byte("null"))
}
