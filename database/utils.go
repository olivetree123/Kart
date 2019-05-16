package database

import "strings"

func FindOperator(s string) string {
	if strings.Contains(s, ">=") {
		return ">="
	} else if strings.Contains(s, "<=") {
		return "<="
	} else if strings.Contains(s, "=") {
		return "="
	} else if strings.Contains(s, ">") {
		return ">"
	} else if strings.Contains(s, "<") {
		return "<"
	}
	return ""
}

// CompareByOperator 根据 operator 比较两个字符串的大小，满足 operator 返回 True, 否则返回 False
func CompareByOperator(s1 string, s2 string, operator string) bool {
	r := strings.Compare(s1, s2)
	if operator == ">=" {
		if r >= 0 {
			return true
		}
		return false
	} else if operator == "<=" {
		if r <= 0 {
			return true
		}
		return false
	} else if operator == ">" {
		if r > 0 {
			return true
		}
		return false
	} else if operator == "<" {
		if r < 0 {
			return true
		}
		return false
	} else if operator == "=" {
		if r == 0 {
			return true
		}
		return false
	}
	return false
}
