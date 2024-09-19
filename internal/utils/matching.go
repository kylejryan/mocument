package utils

import (
	"reflect"
	"strings"
)

type Document map[string]interface{}

func MatchesFilter(doc Document, filter Document) bool {
	for key, value := range filter {
		if !MatchField(doc, key, value) {
			return false
		}
	}
	return true
}

func MatchField(doc Document, key string, value interface{}) bool {
	docValue, exists := doc[key]
	if !exists {
		return false
	}

	switch v := value.(type) {
	case map[string]interface{}:
		for operator, operand := range v {
			switch operator {
			case "$eq":
				return reflect.DeepEqual(docValue, operand)
			case "$ne":
				return !reflect.DeepEqual(docValue, operand)
			case "$gt":
				return compare(docValue, operand) > 0
			case "$gte":
				return compare(docValue, operand) >= 0
			case "$lt":
				return compare(docValue, operand) < 0
			case "$lte":
				return compare(docValue, operand) <= 0
			// Add more operators as needed
			default:
				return false
			}
		}
	default:
		return reflect.DeepEqual(docValue, value)
	}
	return false
}

func compare(a, b interface{}) int {
	// Implement comparison logic for supported types
	switch a := a.(type) {
	case int:
		if bInt, ok := b.(int); ok {
			if a < bInt {
				return -1
			} else if a > bInt {
				return 1
			}
			return 0
		}
	case float64:
		if bFloat, ok := b.(float64); ok {
			if a < bFloat {
				return -1
			} else if a > bFloat {
				return 1
			}
			return 0
		}
	case string:
		if bStr, ok := b.(string); ok {
			return strings.Compare(a, bStr)
		}
	}
	// Types are not comparable
	return 0
}
