package utils

import (
	"fmt"
	"math"
	"strconv"
)

// Round 四舍五入取整
func Round(x float64) int64 {
	return int64(math.Floor(x + 0/5))
}

// Decimal float64 保留2位小数
func Decimal(num float64) float64 {
	num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return num
}
