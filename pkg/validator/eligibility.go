package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// EligibilityResult represents the result of eligibility validation
type EligibilityResult struct {
	IsEligible bool                   `json:"is_eligible"`
	Details    map[string]interface{} `json:"details"`
}

// ValidateEligibility checks if applicant credentials meet scheme eligibility criteria
func ValidateEligibility(criteria json.RawMessage, credentials json.RawMessage) (*EligibilityResult, error) {
	var criteriaMap map[string]interface{}
	if err := json.Unmarshal(criteria, &criteriaMap); err != nil {
		return nil, fmt.Errorf("invalid eligibility criteria format: %v", err)
	}

	var credentialsMap map[string]interface{}
	if err := json.Unmarshal(credentials, &credentialsMap); err != nil {
		return nil, fmt.Errorf("invalid applicant credentials format: %v", err)
	}

	result := &EligibilityResult{
		IsEligible: true,
		Details:    make(map[string]interface{}),
	}

	// Check each criteria against provided credentials
	for key, criteriaValue := range criteriaMap {
		credValue, exists := credentialsMap[key]

		// If credential not provided, mark as ineligible
		if !exists {
			result.IsEligible = false
			result.Details[key] = map[string]interface{}{
				"reason":   "credential_missing",
				"expected": criteriaValue,
				"provided": nil,
			}
			continue
		}

		// Check if the credential meets the criteria
		eligible, reason := compareValues(criteriaValue, credValue)
		if !eligible {
			result.IsEligible = false
			result.Details[key] = map[string]interface{}{
				"reason":   reason,
				"expected": criteriaValue,
				"provided": credValue,
			}
		} else {
			result.Details[key] = map[string]interface{}{
				"reason":   "match",
				"expected": criteriaValue,
				"provided": credValue,
			}
		}
	}

	return result, nil
}

// compareValues compares two values to see if they match the eligibility criteria
func compareValues(criteria, credential interface{}) (bool, string) {
	// Handle different criteria types
	switch c := criteria.(type) {
	case map[string]interface{}:
		// Handle operators like $gt, $lt, $eq, etc.
		return evaluateOperators(c, credential)
	case []interface{}:
		// Handle arrays (one of values must match)
		return evaluateArray(c, credential)
	default:
		// Direct comparison
		return reflect.DeepEqual(criteria, credential), "direct_comparison"
	}
}

// evaluateOperators checks if credential meets operator-based criteria
func evaluateOperators(operators map[string]interface{}, credential interface{}) (bool, string) {
	for op, value := range operators {
		switch op {
		case "$eq":
			if !reflect.DeepEqual(value, credential) {
				return false, "not_equal"
			}
		case "$ne":
			if reflect.DeepEqual(value, credential) {
				return false, "should_not_equal"
			}
		case "$gt":
			if !isGreaterThan(credential, value) {
				return false, "not_greater_than"
			}
		case "$gte":
			if !isGreaterThanOrEqual(credential, value) {
				return false, "not_greater_than_or_equal"
			}
		case "$lt":
			if !isLessThan(credential, value) {
				return false, "not_less_than"
			}
		case "$lte":
			if !isLessThanOrEqual(credential, value) {
				return false, "not_less_than_or_equal"
			}
		case "$in":
			if !isInArray(credential, value) {
				return false, "not_in_array"
			}
		case "$nin":
			if isInArray(credential, value) {
				return false, "should_not_be_in_array"
			}
		}
	}
	return true, "match"
}

// evaluateArray checks if credential matches any value in array criteria
func evaluateArray(array []interface{}, credential interface{}) (bool, string) {
	for _, item := range array {
		if reflect.DeepEqual(item, credential) {
			return true, "match"
		}
	}
	return false, "not_in_options"
}

// isGreaterThan checks if a is greater than b
func isGreaterThan(a, b interface{}) bool {
	aVal, bVal, ok := convertToComparable(a, b)
	if !ok {
		return false
	}
	return aVal > bVal
}

// isGreaterThanOrEqual checks if a is greater than or equal to b
func isGreaterThanOrEqual(a, b interface{}) bool {
	aVal, bVal, ok := convertToComparable(a, b)
	if !ok {
		return false
	}
	return aVal >= bVal
}

// isLessThan checks if a is less than b
func isLessThan(a, b interface{}) bool {
	aVal, bVal, ok := convertToComparable(a, b)
	if !ok {
		return false
	}
	return aVal < bVal
}

// isLessThanOrEqual checks if a is less than or equal to b
func isLessThanOrEqual(a, b interface{}) bool {
	aVal, bVal, ok := convertToComparable(a, b)
	if !ok {
		return false
	}
	return aVal <= bVal
}

// isInArray checks if a is in an array
func isInArray(a, b interface{}) bool {
	bArray, ok := b.([]interface{})
	if !ok {
		return false
	}

	for _, item := range bArray {
		if reflect.DeepEqual(a, item) {
			return true
		}
	}
	return false
}

// convertToComparable converts two values to float64 for comparison
func convertToComparable(a, b interface{}) (float64, float64, bool) {
	var aVal, bVal float64
	var ok bool

	aVal, ok = toFloat64(a)
	if !ok {
		return 0, 0, false
	}

	bVal, ok = toFloat64(b)
	if !ok {
		return 0, 0, false
	}

	return aVal, bVal, true
}

// toFloat64 converts an interface{} to float64 if possible
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	case string:
		// Try to parse the string as a number
		if num, err := parseStringNumber(val); err == nil {
			return num, true
		}
	}
	return 0, false
}

// parseStringNumber attempts to parse a string as a number
func parseStringNumber(s string) (float64, error) {
	// Try to parse as float
	var result float64
	var err error
	result, err = parseFloat(s)
	if err == nil {
		return result, nil
	}

	// Handle percentage format
	if strings.HasSuffix(s, "%") {
		percentValue, err := parseFloat(strings.TrimSuffix(s, "%"))
		if err == nil {
			return percentValue / 100.0, nil
		}
	}

	return 0, fmt.Errorf("could not parse as number: %s", s)
}

// parseFloat tries to parse a string to float64
func parseFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	return result, err
}
