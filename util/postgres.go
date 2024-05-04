package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseQueryOperator(s string) (string, error) {
	parts := strings.Split(s, "=")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid input: %s", s)
	}

	field := parts[0]
	valueStr := parts[1]

	var condition string
	var value int
	if strings.HasPrefix(valueStr, ">") {
		condition = ">"
		valueStr = valueStr[1:]
	} else if strings.HasPrefix(valueStr, "<") {
		condition = "<"
		valueStr = valueStr[1:]
	} else {
		condition = "="
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return "", fmt.Errorf("invalid value: %s", valueStr)
	}

	switch condition {
	case ">":
		return fmt.Sprintf("%s>=%d", field, value), nil
	case "<":
		return fmt.Sprintf("%s<=%d", field, value), nil
	default:
		return fmt.Sprintf("%s=%d", field, value), nil
	}
}
