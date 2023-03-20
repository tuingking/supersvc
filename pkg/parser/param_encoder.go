package parser

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func (p *paramparser) InitEncoder() {
	p.encoder.RegisterEncoder(sql.NullString{}, encodesqlNullString)
	p.encoder.RegisterEncoder(sql.NullBool{}, encodesqlNullBool)
	p.encoder.RegisterEncoder(sql.NullInt64{}, encodesqlNullInt64)
	p.encoder.RegisterEncoder(sql.NullFloat64{}, encodesqlNullFloat64)
	p.encoder.RegisterEncoder(sql.NullTime{}, encodesqlNullTime)
	p.encoder.RegisterEncoder(time.Time{}, encodeTime)
}

func encodesqlNullString(v reflect.Value) string {
	nullString, _ := v.Interface().(sql.NullString)

	if !nullString.Valid {
		return ""
	}

	return nullString.String
}

func encodesqlNullBool(v reflect.Value) string {
	nullBool, _ := v.Interface().(sql.NullBool)

	if !nullBool.Valid {
		return ""
	}

	return strconv.FormatBool(nullBool.Bool)
}

func encodesqlNullInt64(v reflect.Value) string {
	nullInt, _ := v.Interface().(sql.NullInt64)

	if !nullInt.Valid {
		return ""
	}

	return strconv.FormatInt(nullInt.Int64, 10)
}

func encodesqlNullFloat64(v reflect.Value) string {
	nullFloat, _ := v.Interface().(sql.NullFloat64)

	if !nullFloat.Valid {
		return ""
	}

	return fmt.Sprintf("%.2f", nullFloat.Float64)
}

func encodesqlNullTime(v reflect.Value) string {
	nullTime, _ := v.Interface().(sql.NullTime)

	if !nullTime.Valid {
		return ""
	}

	return nullTime.Time.Format(time.RFC3339)
}

func encodeTime(v reflect.Value) string {
	currTime, _ := v.Interface().(time.Time)

	if currTime.IsZero() {
		return ""
	}

	return currTime.Format(time.RFC3339)
}
