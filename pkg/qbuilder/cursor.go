package qbuilder

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type cursor struct {
	field   reflect.Value // struct field
	param   string        // tag:"param"
	db      string        // tag:"db"
	jsonKey string        // tag:"json_key"
}

func newCursor(field reflect.Value, param, db, jsonKey string) cursor {
	return cursor{
		field:   field,
		param:   param,
		db:      db,
		jsonKey: jsonKey,
	}
}

func (c *cursor) IsPage() bool {
	return c.param == "page"
}

func (c *cursor) IsLimit() bool {
	return c.param == "limit"
}

func (c *cursor) IsSortBy() bool {
	return c.param == "sortBy"
}

func (c *cursor) IsEmpty() bool {
	if c.param == "-" ||
		c.param == "" ||
		c.db == "" ||
		c.db == "-" {
		return true
	}
	return false
}

func (c *cursor) GetOperand() string {
	operand := "="

	// update operand
	switch param := c.param; {
	case strings.HasSuffix(param, "__gt"):
		operand = ">"
	case strings.HasSuffix(param, "__gte"):
		operand = ">="
	case strings.HasSuffix(param, "__lt"):
		operand = "<"
	case strings.HasSuffix(param, "__lte"):
		operand = "<="
	case strings.HasSuffix(param, "__neq"):
		operand = "!="
	default:
		operand = "="
	}

	return operand
}

func (c *cursor) GetOperandMulti() string {
	operand := "IN"

	// update operand
	if strings.HasSuffix(c.param, "__nin") {
		operand = "NOT IN"
	} else {
		// skip
	}

	return operand
}

func (c *cursor) Make() (clause string, args []interface{}, skip bool) {
	switch c.field.Interface().(type) {
	case string, int, int32, int64, float32, float64:
		return c.makeClausePrimitiveType()
	case time.Time, sql.NullTime:
		return c.makeClauseTimeType()
	case []string, []int, []int32, []int64, []float32, []float64:
		return c.makeClauseArrayType()
	case sql.NullString, sql.NullInt32, sql.NullInt64, sql.NullFloat64, sql.NullBool:
		return c.makeClauseSqlNullType()
	default:
		skip = true
	}

	return
}

func (c *cursor) makeClausePrimitiveType() (clause string, args []interface{}, skip bool) {
	operand := c.GetOperand()

	switch val := c.field.Interface().(type) {
	case string:
		clause, args = c.makeClauseStringType(operand, val)
	case int, int32, int64, float32, float64:
		clause, args = c.makeClause(whereClauseFmt, operand, val)
	default:
		skip = true
	}

	return
}

var regexMysqlJsonKey = regexp.MustCompile(`^(\$\[([\d]+)\]).*`)

func (c *cursor) makeClauseStringType(operand, val string) (clause string, args []interface{}) {
	if c.jsonKey != "" {
		clause += fmt.Sprintf(" AND "+whereClauseJsonMemberFmt, val, c.db, c.jsonKey)
		fmt.Println("clause = ", clause)
		return
	}

	clause, args = c.makeClause(whereClauseFmt, operand, val)
	return
}

func (c *cursor) makeClauseTimeType() (clause string, args []interface{}, skip bool) {
	operand := c.GetOperand()

	switch val := c.field.Interface().(type) {
	case time.Time:
		if val.IsZero() {
			skip = true
			return
		}
		clause, args = c.makeClause(whereClauseFmt, operand, val)
	case sql.NullTime:
		if !val.Valid {
			skip = true
			return
		}
		clause, args = c.makeClause(whereClauseFmt, operand, val.Time)
	default:
		skip = true
	}

	return
}

func (c *cursor) makeClauseArrayType() (clause string, args []interface{}, skip bool) {
	switch val := c.field.Interface().(type) {
	case []string, []int, []int32, []int64, []float32, []float64:
		return c.makeClauseMulti(val)
	default:
		skip = true
	}

	return
}

func (c *cursor) makeClauseSqlNullType() (clause string, args []interface{}, skip bool) {
	operand := c.GetOperand()

	switch val := c.field.Interface().(type) {
	case sql.NullString:
		clause, args, skip = c.makeClauseNullString(whereClauseFmt, operand, val)
	case sql.NullInt32:
		clause, args, skip = c.makeClauseNullInt32(whereClauseFmt, operand, val)
	case sql.NullInt64:
		clause, args, skip = c.makeClauseNullInt64(whereClauseFmt, operand, val)
	case sql.NullFloat64:
		clause, args, skip = c.makeClauseNullFloat64(whereClauseFmt, operand, val)
	case sql.NullBool:
		clause, args, skip = c.makeClauseNullBool(whereClauseFmt, operand, val)
	default:
		skip = true
	}

	return
}

func (c *cursor) makeClause(layout, operand string, val interface{}) (clause string, args []interface{}) {
	clause = fmt.Sprintf(layout, c.db, operand)
	args = append(args, val)

	return
}

func (c *cursor) makeClauseMulti(val interface{}) (clause string, args []interface{}, skip bool) {
	operandMulti := c.GetOperandMulti()
	tempQuery := fmt.Sprintf(whereClauseMultiFmt, c.db, operandMulti)
	tempQuery, tempArgs, _ := sqlx.In(tempQuery, val)
	clause = tempQuery
	if len(tempArgs) < 1 {
		skip = true
	}
	args = append(args, tempArgs...)

	return
}

func (c *cursor) makeClauseNullString(layout, operand string, val sql.NullString) (clause string, args []interface{}, skip bool) {
	if !val.Valid {
		skip = true
		return
	}

	clause, args = c.makeClauseStringType(operand, val.String)
	return
}

func (c *cursor) makeClauseNullInt32(layout, operand string, val sql.NullInt32) (clause string, args []interface{}, skip bool) {
	if !val.Valid {
		skip = true
		return
	}
	clause = fmt.Sprintf(layout, c.db, operand)
	args = append(args, val.Int32)

	return
}

func (c *cursor) makeClauseNullInt64(layout, operand string, val sql.NullInt64) (clause string, args []interface{}, skip bool) {
	if !val.Valid {
		skip = true
		return
	}
	clause = fmt.Sprintf(layout, c.db, operand)
	args = append(args, val.Int64)

	return
}

func (c *cursor) makeClauseNullFloat64(layout, operand string, val sql.NullFloat64) (clause string, args []interface{}, skip bool) {
	if !val.Valid {
		skip = true
		return
	}
	clause = fmt.Sprintf(layout, c.db, operand)
	args = append(args, val.Float64)

	return
}

func (c *cursor) makeClauseNullBool(layout, operand string, val sql.NullBool) (clause string, args []interface{}, skip bool) {
	if !val.Valid {
		skip = true
		return
	}
	clause = fmt.Sprintf(layout, c.db, operand)
	args = append(args, val.Bool)

	return
}
