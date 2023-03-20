package parser

import (
	"database/sql"
	"reflect"
	"time"
)

func (p *paramparser) InitDecoder() {
	nullString, nullBool, nullInt64, nullFloat64, nullTime := sql.NullString{}, sql.NullBool{}, sql.NullInt64{}, sql.NullFloat64{}, sql.NullTime{}
	p.decoder.RegisterConverter(nullString, convertsqlNullString)
	p.decoder.RegisterConverter(nullBool, convertsqlNullBool)
	p.decoder.RegisterConverter(nullInt64, convertsqlNullInt64)
	p.decoder.RegisterConverter(nullFloat64, convertsqlNullFloat64)
	p.decoder.RegisterConverter(nullTime, p.convertsqlNullTime)
}

func convertsqlNullString(value string) reflect.Value {
	v := sql.NullString{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func convertsqlNullBool(value string) reflect.Value {
	v := sql.NullBool{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func convertsqlNullInt64(value string) reflect.Value {
	v := sql.NullInt64{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func convertsqlNullFloat64(value string) reflect.Value {
	v := sql.NullFloat64{}
	if err := v.Scan(value); err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func (p *paramparser) convertsqlNullTime(value string) reflect.Value {
	// handle multi time format
	v := sql.NullTime{}
	formatToCheck := []string{time.RFC3339, `2006-01-02`, `2006-01-02 15:04:05`, `2006-01-02T15:04:05`, `2006-01-02T15:04:05.000Z`,
		`2006-01-02 15:04:05.000Z`, `2006-01-02 15:04:05-07:00`, `2006-01-02T15:04:05-07:00`, `2006-01-02 15:04:05 -07:00 MST`,
		`2006-01-02T15:04:05 -07:00 MST`, `2006-01-02T15:04:05 -07:00MST`, `2006-01-02 15:04:05 -07:00MST`}

	for _, format := range formatToCheck {
		if t0, err := time.Parse(format, value); err == nil {
			return p.generateTime(v, t0)
		}
	}

	return reflect.Value{}
}

func (p *paramparser) generateTime(result sql.NullTime, value time.Time) reflect.Value {
	result.Valid = true
	result.Time = value
	return reflect.ValueOf(result)
}
