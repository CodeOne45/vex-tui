package app

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/CodeOne45/vex-tui/pkg/models"
)

type FormulaEngine struct {
	sheet *models.Sheet
}

// evaluateFormula evaluates a formula and returns the result
func (m *Model) evaluateFormula(formula string) string {
	engine := &FormulaEngine{sheet: &m.sheets[m.currentSheet]}
	result, err := engine.Evaluate(formula)
	if err != nil {
		return "#ERROR!"
	}
	return result
}

// Evaluate evaluates a formula string
func (fe *FormulaEngine) Evaluate(formula string) (string, error) {
	formula = strings.TrimSpace(formula)
	if formula == "" {
		return "", nil
	}

	formula = strings.ToUpper(formula)

	if strings.HasPrefix(formula, "SUM(") {
		return fe.evaluateSum(formula)
	} else if strings.HasPrefix(formula, "AVERAGE(") || strings.HasPrefix(formula, "AVG(") {
		return fe.evaluateAverage(formula)
	} else if strings.HasPrefix(formula, "COUNT(") {
		return fe.evaluateCount(formula)
	} else if strings.HasPrefix(formula, "MAX(") {
		return fe.evaluateMax(formula)
	} else if strings.HasPrefix(formula, "MIN(") {
		return fe.evaluateMin(formula)
	} else if strings.HasPrefix(formula, "IF(") {
		return fe.evaluateIf(formula)
	} else if strings.HasPrefix(formula, "CONCATENATE(") || strings.HasPrefix(formula, "CONCAT(") {
		return fe.evaluateConcatenate(formula)
	} else if strings.HasPrefix(formula, "UPPER(") {
		return fe.evaluateUpper(formula)
	} else if strings.HasPrefix(formula, "LOWER(") {
		return fe.evaluateLower(formula)
	} else if strings.HasPrefix(formula, "LEN(") {
		return fe.evaluateLen(formula)
	} else if strings.HasPrefix(formula, "ROUND(") {
		return fe.evaluateRound(formula)
	} else if strings.HasPrefix(formula, "ABS(") {
		return fe.evaluateAbs(formula)
	} else if strings.HasPrefix(formula, "SQRT(") {
		return fe.evaluateSqrt(formula)
	} else if strings.HasPrefix(formula, "POWER(") || strings.HasPrefix(formula, "POW(") {
		return fe.evaluatePower(formula)
	} else if strings.Contains(formula, ":") {
		return fe.evaluateCellReference(formula)
	}

	result, err := fe.evaluateExpression(formula)
	if err != nil {
		return "", err
	}
	return result, nil
}

// evaluateSum evaluates SUM function
func (fe *FormulaEngine) evaluateSum(formula string) (string, error) {
	rangeStr := fe.extractFunctionArg(formula, "SUM")
	values, err := fe.getRangeValues(rangeStr)
	if err != nil {
		return "", err
	}

	sum := 0.0
	for _, val := range values {
		if num, err := strconv.ParseFloat(val, 64); err == nil {
			sum += num
		}
	}

	return formatNumber(sum), nil
}

// evaluateAverage evaluates AVERAGE/AVG function
func (fe *FormulaEngine) evaluateAverage(formula string) (string, error) {
	funcName := "AVERAGE"
	if strings.HasPrefix(formula, "AVG(") {
		funcName = "AVG"
	}

	rangeStr := fe.extractFunctionArg(formula, funcName)
	values, err := fe.getRangeValues(rangeStr)
	if err != nil {
		return "", err
	}

	sum := 0.0
	count := 0
	for _, val := range values {
		if num, err := strconv.ParseFloat(val, 64); err == nil {
			sum += num
			count++
		}
	}

	if count == 0 {
		return "0", nil
	}

	return formatNumber(sum / float64(count)), nil
}

// evaluateCount evaluates COUNT function
func (fe *FormulaEngine) evaluateCount(formula string) (string, error) {
	rangeStr := fe.extractFunctionArg(formula, "COUNT")
	values, err := fe.getRangeValues(rangeStr)
	if err != nil {
		return "", err
	}

	count := 0
	for _, val := range values {
		if _, err := strconv.ParseFloat(val, 64); err == nil {
			count++
		}
	}

	return strconv.Itoa(count), nil
}

// evaluateMax evaluates MAX function
func (fe *FormulaEngine) evaluateMax(formula string) (string, error) {
	rangeStr := fe.extractFunctionArg(formula, "MAX")
	values, err := fe.getRangeValues(rangeStr)
	if err != nil {
		return "", err
	}

	max := math.Inf(-1)
	for _, val := range values {
		if num, err := strconv.ParseFloat(val, 64); err == nil {
			if num > max {
				max = num
			}
		}
	}

	if math.IsInf(max, -1) {
		return "0", nil
	}

	return formatNumber(max), nil
}

// evaluateMin evaluates MIN function
func (fe *FormulaEngine) evaluateMin(formula string) (string, error) {
	rangeStr := fe.extractFunctionArg(formula, "MIN")
	values, err := fe.getRangeValues(rangeStr)
	if err != nil {
		return "", err
	}

	min := math.Inf(1)
	for _, val := range values {
		if num, err := strconv.ParseFloat(val, 64); err == nil {
			if num < min {
				min = num
			}
		}
	}

	if math.IsInf(min, 1) {
		return "0", nil
	}

	return formatNumber(min), nil
}

// evaluateIf evaluates IF function
func (fe *FormulaEngine) evaluateIf(formula string) (string, error) {
	args := fe.extractFunctionArgs(formula, "IF")
	if len(args) < 3 {
		return "", fmt.Errorf("IF requires 3 arguments")
	}

	condition := strings.TrimSpace(args[0])
	trueVal := strings.TrimSpace(args[1])
	falseVal := strings.TrimSpace(args[2])

	result, err := fe.evaluateCondition(condition)
	if err != nil {
		return "", err
	}

	if result {
		if strings.HasPrefix(trueVal, "\"") && strings.HasSuffix(trueVal, "\"") {
			return trueVal[1 : len(trueVal)-1], nil
		}
		return fe.evaluateExpression(trueVal)
	} else {
		if strings.HasPrefix(falseVal, "\"") && strings.HasSuffix(falseVal, "\"") {
			return falseVal[1 : len(falseVal)-1], nil
		}
		return fe.evaluateExpression(falseVal)
	}
}

// evaluateConcatenate evaluates CONCATENATE/CONCAT function
func (fe *FormulaEngine) evaluateConcatenate(formula string) (string, error) {
	funcName := "CONCATENATE"
	if strings.HasPrefix(formula, "CONCAT(") {
		funcName = "CONCAT"
	}

	args := fe.extractFunctionArgs(formula, funcName)
	var result strings.Builder

	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"") {
			result.WriteString(arg[1 : len(arg)-1])
		} else {
			val, err := fe.evaluateExpression(arg)
			if err == nil {
				result.WriteString(val)
			}
		}
	}

	return result.String(), nil
}

// evaluateUpper evaluates UPPER function
func (fe *FormulaEngine) evaluateUpper(formula string) (string, error) {
	arg := fe.extractFunctionArg(formula, "UPPER")
	val, err := fe.evaluateExpression(arg)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(val), nil
}

// evaluateLower evaluates LOWER function
func (fe *FormulaEngine) evaluateLower(formula string) (string, error) {
	arg := fe.extractFunctionArg(formula, "LOWER")
	val, err := fe.evaluateExpression(arg)
	if err != nil {
		return "", err
	}
	return strings.ToLower(val), nil
}

// evaluateLen evaluates LEN function
func (fe *FormulaEngine) evaluateLen(formula string) (string, error) {
	arg := fe.extractFunctionArg(formula, "LEN")
	val, err := fe.evaluateExpression(arg)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(len(val)), nil
}

// evaluateRound evaluates ROUND function
func (fe *FormulaEngine) evaluateRound(formula string) (string, error) {
	args := fe.extractFunctionArgs(formula, "ROUND")
	if len(args) < 2 {
		return "", fmt.Errorf("ROUND requires 2 arguments")
	}

	numStr, err := fe.evaluateExpression(strings.TrimSpace(args[0]))
	if err != nil {
		return "", err
	}
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return "", err
	}

	digitsStr, err := fe.evaluateExpression(strings.TrimSpace(args[1]))
	if err != nil {
		return "", err
	}
	digits, err := strconv.Atoi(digitsStr)
	if err != nil {
		return "", err
	}

	multiplier := math.Pow(10, float64(digits))
	rounded := math.Round(num*multiplier) / multiplier

	return formatNumber(rounded), nil
}

// evaluateAbs evaluates ABS function
func (fe *FormulaEngine) evaluateAbs(formula string) (string, error) {
	arg := fe.extractFunctionArg(formula, "ABS")
	val, err := fe.evaluateExpression(arg)
	if err != nil {
		return "", err
	}
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return "", err
	}
	return formatNumber(math.Abs(num)), nil
}

// evaluateSqrt evaluates SQRT function
func (fe *FormulaEngine) evaluateSqrt(formula string) (string, error) {
	arg := fe.extractFunctionArg(formula, "SQRT")
	val, err := fe.evaluateExpression(arg)
	if err != nil {
		return "", err
	}
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return "", err
	}
	if num < 0 {
		return "", fmt.Errorf("SQRT of negative number")
	}
	return formatNumber(math.Sqrt(num)), nil
}

// evaluatePower evaluates POWER/POW function
func (fe *FormulaEngine) evaluatePower(formula string) (string, error) {
	funcName := "POWER"
	if strings.HasPrefix(formula, "POW(") {
		funcName = "POW"
	}

	args := fe.extractFunctionArgs(formula, funcName)
	if len(args) < 2 {
		return "", fmt.Errorf("POWER requires 2 arguments")
	}

	baseStr, err := fe.evaluateExpression(strings.TrimSpace(args[0]))
	if err != nil {
		return "", err
	}
	base, err := strconv.ParseFloat(baseStr, 64)
	if err != nil {
		return "", err
	}

	expStr, err := fe.evaluateExpression(strings.TrimSpace(args[1]))
	if err != nil {
		return "", err
	}
	exp, err := strconv.ParseFloat(expStr, 64)
	if err != nil {
		return "", err
	}

	return formatNumber(math.Pow(base, exp)), nil
}

// evaluateCellReference gets value from cell reference
func (fe *FormulaEngine) evaluateCellReference(ref string) (string, error) {
	ref = strings.TrimSpace(ref)

	if strings.Contains(ref, ":") {
		values, err := fe.getRangeValues(ref)
		if err != nil {
			return "", err
		}
		if len(values) > 0 {
			return values[0], nil
		}
		return "", nil
	}

	row, col, err := fe.parseCellReference(ref)
	if err != nil {
		return "", err
	}

	if row < 0 || row >= len(fe.sheet.Rows) {
		return "", nil
	}
	if col < 0 || col >= len(fe.sheet.Rows[row]) {
		return "", nil
	}

	cell := fe.sheet.Rows[row][col]
	return cell.Value, nil
}

// evaluateExpression evaluates a simple expression
func (fe *FormulaEngine) evaluateExpression(expr string) (string, error) {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"") {
		return expr[1 : len(expr)-1], nil
	}

	if fe.isCellReference(expr) {
		return fe.evaluateCellReference(expr)
	}

	if _, err := strconv.ParseFloat(expr, 64); err == nil {
		return expr, nil
	}

	result, err := fe.evaluateArithmetic(expr)
	if err == nil {
		return formatNumber(result), nil
	}

	return expr, nil
}

// evaluateArithmetic evaluates arithmetic expressions
func (fe *FormulaEngine) evaluateArithmetic(expr string) (float64, error) {
	expr = strings.ReplaceAll(expr, " ", "")

	for i := len(expr) - 1; i >= 0; i-- {
		if expr[i] == '+' || expr[i] == '-' {
			left, err := fe.evaluateArithmetic(expr[:i])
			if err != nil {
				continue
			}
			right, err := fe.evaluateArithmetic(expr[i+1:])
			if err != nil {
				continue
			}
			if expr[i] == '+' {
				return left + right, nil
			}
			return left - right, nil
		}
	}

	for i := len(expr) - 1; i >= 0; i-- {
		if expr[i] == '*' || expr[i] == '/' {
			left, err := fe.evaluateArithmetic(expr[:i])
			if err != nil {
				continue
			}
			right, err := fe.evaluateArithmetic(expr[i+1:])
			if err != nil {
				continue
			}
			if expr[i] == '*' {
				return left * right, nil
			}
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		}
	}

	if fe.isCellReference(expr) {
		val, err := fe.evaluateCellReference(expr)
		if err != nil {
			return 0, err
		}
		return strconv.ParseFloat(val, 64)
	}

	return strconv.ParseFloat(expr, 64)
}

// evaluateCondition evaluates a conditional expression
func (fe *FormulaEngine) evaluateCondition(condition string) (bool, error) {
	operators := []string{">=", "<=", "<>", "=", ">", "<"}

	for _, op := range operators {
		if strings.Contains(condition, op) {
			parts := strings.SplitN(condition, op, 2)
			if len(parts) != 2 {
				continue
			}

			left, err := fe.evaluateExpression(strings.TrimSpace(parts[0]))
			if err != nil {
				continue
			}
			right, err := fe.evaluateExpression(strings.TrimSpace(parts[1]))
			if err != nil {
				continue
			}

			leftNum, leftErr := strconv.ParseFloat(left, 64)
			rightNum, rightErr := strconv.ParseFloat(right, 64)

			if leftErr == nil && rightErr == nil {
				switch op {
				case "=":
					return leftNum == rightNum, nil
				case ">":
					return leftNum > rightNum, nil
				case "<":
					return leftNum < rightNum, nil
				case ">=":
					return leftNum >= rightNum, nil
				case "<=":
					return leftNum <= rightNum, nil
				case "<>":
					return leftNum != rightNum, nil
				}
			} else {
				switch op {
				case "=":
					return left == right, nil
				case "<>":
					return left != right, nil
				}
			}
		}
	}

	return false, fmt.Errorf("invalid condition")
}

// extractFunctionArg extracts single argument from function
func (fe *FormulaEngine) extractFunctionArg(formula, funcName string) string {
	start := strings.Index(formula, funcName+"(")
	if start == -1 {
		return ""
	}
	start += len(funcName) + 1

	depth := 1
	for i := start; i < len(formula); i++ {
		if formula[i] == '(' {
			depth++
		} else if formula[i] == ')' {
			depth--
			if depth == 0 {
				return formula[start:i]
			}
		}
	}
	return ""
}

// extractFunctionArgs extracts multiple arguments from function
func (fe *FormulaEngine) extractFunctionArgs(formula, funcName string) []string {
	arg := fe.extractFunctionArg(formula, funcName)
	if arg == "" {
		return nil
	}

	var args []string
	var current strings.Builder
	depth := 0
	inQuote := false

	for i := 0; i < len(arg); i++ {
		ch := arg[i]

		if ch == '"' {
			inQuote = !inQuote
			current.WriteByte(ch)
		} else if ch == '(' && !inQuote {
			depth++
			current.WriteByte(ch)
		} else if ch == ')' && !inQuote {
			depth--
			current.WriteByte(ch)
		} else if ch == ',' && depth == 0 && !inQuote {
			args = append(args, current.String())
			current.Reset()
		} else {
			current.WriteByte(ch)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

// getRangeValues gets all values in a range
func (fe *FormulaEngine) getRangeValues(rangeStr string) ([]string, error) {
	rangeStr = strings.TrimSpace(rangeStr)

	if !strings.Contains(rangeStr, ":") {
		val, err := fe.evaluateCellReference(rangeStr)
		return []string{val}, err
	}

	parts := strings.Split(rangeStr, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid range")
	}

	startRow, startCol, err := fe.parseCellReference(parts[0])
	if err != nil {
		return nil, err
	}
	endRow, endCol, err := fe.parseCellReference(parts[1])
	if err != nil {
		return nil, err
	}

	if startRow > endRow {
		startRow, endRow = endRow, startRow
	}
	if startCol > endCol {
		startCol, endCol = endCol, startCol
	}

	var values []string
	for row := startRow; row <= endRow && row < len(fe.sheet.Rows); row++ {
		for col := startCol; col <= endCol && col < len(fe.sheet.Rows[row]); col++ {
			values = append(values, fe.sheet.Rows[row][col].Value)
		}
	}

	return values, nil
}

// parseCellReference parses a cell reference like "A1" to row, col
func (fe *FormulaEngine) parseCellReference(ref string) (int, int, error) {
	ref = strings.ToUpper(strings.TrimSpace(ref))
	if ref == "" {
		return 0, 0, fmt.Errorf("empty reference")
	}

	col := 0
	i := 0
	for i < len(ref) && ref[i] >= 'A' && ref[i] <= 'Z' {
		col = col*26 + int(ref[i]-'A') + 1
		i++
	}
	col--

	if i >= len(ref) {
		return 0, 0, fmt.Errorf("no row number")
	}

	row, err := strconv.Atoi(ref[i:])
	if err != nil {
		return 0, 0, err
	}
	row--

	return row, col, nil
}

// isCellReference checks if string is a cell reference
func (fe *FormulaEngine) isCellReference(s string) bool {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "" {
		return false
	}

	i := 0
	for i < len(s) && s[i] >= 'A' && s[i] <= 'Z' {
		i++
	}
	if i == 0 || i >= len(s) {
		return false
	}

	_, err := strconv.Atoi(s[i:])
	return err == nil
}

// formatNumber formats a number for display
func formatNumber(num float64) string {
	if num == float64(int64(num)) {
		return strconv.FormatInt(int64(num), 10)
	}
	return strconv.FormatFloat(num, 'f', -1, 64)
}

// recalculateFormulas recalculates all formulas in the sheet
func (m *Model) recalculateFormulas() {
	sheet := &m.sheets[m.currentSheet]

	for i := 0; i < 3; i++ {
		for _, row := range sheet.Rows {
			for j := range row {
				cell := &row[j]
				if cell.Formula != "" {
					cell.Value = m.evaluateFormula(cell.Formula)
				}
			}
		}
	}
}
