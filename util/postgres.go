package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type queryOperator struct {
	Operator string
	Value    int
}

func ParseQueryOperator(s string) (*queryOperator, error) {
	pattern := `^([<>]?)([=]?)(\d+)$`

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("Invalid input: %s", s)
	}

	// Extract the result and value from the matches
	condition := matches[1]
	valueStr := matches[3]

	if strings.HasPrefix(condition, ">") {
		condition = ">"
	} else if strings.HasPrefix(condition, "<") {
		condition = "<"
	} else {
		condition = "="
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return nil, fmt.Errorf("invalid value: %s", valueStr)
	}

	return &queryOperator{
		Operator: condition,
		Value:    value,
	}, nil
}
