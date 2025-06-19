package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"regexp"
)

// VerifyCronExpression 验证和解析 cron 表达式。
func VerifyCronExpression(expr string) (cron.Schedule, error) {
	// 先处理自定义格式
	standardExpr, err := parseCustomCron(expr)
	if err != nil {
		return nil, fmt.Errorf("无效的 cron 表达式: %w", err)
	}

	// 尝试解析 cron 表达式
	schedule, err := cron.ParseStandard(standardExpr)
	if err != nil {
		return nil, fmt.Errorf("cron 表达式解析失败: %w", err)
	}
	return schedule, nil
}

// parseCustomCron 将自定义的格式转换为标准 cron 表达式
func parseCustomCron(expr string) (string, error) {
	switch expr {
	case "每小时执行一次":
		return "0 * * * *", nil
	case "每天午夜执行":
		return "0 0 * * *", nil
	case "每分钟执行一次":
		return "* * * * *", nil
	// 添加更多自定义的解析规则
	default:
		// 如果没有匹配的自定义规则，尝试标准表达式
		return validateStandardCron(expr)
	}
}

// validateStandardCron 使用正则表达式验证标准 cron 表达式格式
func validateStandardCron(expr string) (string, error) {
	cronRegex := `^(\*|[0-5]?\d) (\*|[01]?\d|2[0-3]) (\*|[01]?\d|2[0-3]) (\*|[0-2]?\d|3[01]) (\*|[1-9]|1[0-2])(?: (\*|[0-6]))?$`
	if matched, _ := regexp.MatchString(cronRegex, expr); matched {
		return expr, nil
	}
	return "", fmt.Errorf("无效的cron 表达式")
}
