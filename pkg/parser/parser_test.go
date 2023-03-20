package parser_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/tiket/TIX-DEPOSIT-BE/pkg/parser"
	"gotest.tools/assert"
)

type ParamPrimitive struct {
	Int     int       `param:"int"`
	Int64   int       `param:"int64"`
	Bool    bool      `param:"bool"`
	Float64 float64   `param:"float64"`
	String  string    `param:"string"`
	Time    time.Time `param:"time"`
}

type ParamSqlNull struct {
	Int64   sql.NullInt64   `param:"int64"`
	Bool    sql.NullBool    `param:"bool"`
	Float64 sql.NullFloat64 `param:"float64"`
	String  sql.NullString  `param:"string"`
	Time    sql.NullTime    `param:"time"`
}

func Test_Decode(t *testing.T) {
	t.Run("Test Decode Primitive Type", func(t *testing.T) {
		testCase := []struct {
			desc  string
			param map[string][]string
			exp   ParamPrimitive
		}{
			{
				desc: "positive case",
				param: map[string][]string{
					"int":     {"1"},
					"int64":   {"1"},
					"bool":    {"true"},
					"float64": {"10.7"},
					"string":  {"active"},
					"time":    {"2022-01-01T00:00:00Z"},
				},
				exp: ParamPrimitive{
					Int:     1,
					Int64:   1,
					Bool:    true,
					Float64: 10.7,
					String:  "active",
					Time:    time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC),
				},
			},
		}

		for _, tc := range testCase {
			result := ParamPrimitive{}
			parser := parser.InitParamParser()
			parser.Decode(&result, tc.param)

			assert.Equal(t, tc.exp.Int, result.Int)
			assert.Equal(t, tc.exp.Int64, result.Int64)
			assert.Equal(t, tc.exp.Bool, result.Bool)
			assert.Equal(t, tc.exp.Float64, result.Float64)
			assert.Equal(t, tc.exp.String, result.String)
			assert.Equal(t, tc.exp.Time, result.Time)
		}
	})

	t.Run("Test Decode Sql Null Type", func(t *testing.T) {
		testCase := []struct {
			desc  string
			param map[string][]string
			exp   ParamSqlNull
		}{
			{
				desc: "positive case",
				param: map[string][]string{
					"int64":   {"99"},
					"bool":    {"true"},
					"float64": {"10.7"},
					"string":  {"active"},
					"time":    {"2022-01-01T00:00:00Z"},
				},
				exp: ParamSqlNull{
					Int64:   sql.NullInt64{Valid: true, Int64: 99},
					Bool:    sql.NullBool{Valid: true, Bool: true},
					Float64: sql.NullFloat64{Valid: true, Float64: 10.7},
					String:  sql.NullString{Valid: true, String: "active"},
					Time:    sql.NullTime{Valid: true, Time: time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC)},
				},
			},
			{
				desc: "negative case",
				param: map[string][]string{
					"int64":   {"hoho"},
					"bool":    {"hoho"},
					"float64": {"hoho"},
					"time":    {"2022-01-01T00:00:00ZABC"},
				},
				exp: ParamSqlNull{
					Int64:   sql.NullInt64{Valid: false},
					Bool:    sql.NullBool{Valid: false},
					Float64: sql.NullFloat64{Valid: false},
					String:  sql.NullString{Valid: false},
					Time:    sql.NullTime{Valid: false},
				},
			},
			{
				desc: "case: empty string",
				param: map[string][]string{
					"string": {""},
				},
				exp: ParamSqlNull{
					Int64:   sql.NullInt64{Valid: false},
					Bool:    sql.NullBool{Valid: false},
					Float64: sql.NullFloat64{Valid: false},
					String:  sql.NullString{Valid: true, String: ""},
					Time:    sql.NullTime{Valid: false},
				},
			},
			{
				desc: "case: empty string",
				param: map[string][]string{
					"string": {""},
				},
				exp: ParamSqlNull{
					Int64:   sql.NullInt64{Valid: false},
					Bool:    sql.NullBool{Valid: false},
					Float64: sql.NullFloat64{Valid: false},
					String:  sql.NullString{Valid: true, String: ""},
					Time:    sql.NullTime{Valid: false},
				},
			},
		}

		for _, tc := range testCase {
			result := ParamSqlNull{}
			parser := parser.InitParamParser()
			parser.Decode(&result, tc.param)

			assert.DeepEqual(t, tc.exp.Int64, result.Int64)
			assert.DeepEqual(t, tc.exp.Bool, result.Bool)
			assert.DeepEqual(t, tc.exp.Float64, result.Float64)
			assert.DeepEqual(t, tc.exp.String, result.String)
			assert.DeepEqual(t, tc.exp.Time, result.Time)
		}
	})
}

func Test_Encode(t *testing.T) {
	t.Run("Test Encode Primitive Type", func(t *testing.T) {
		testCase := []struct {
			desc  string
			input ParamPrimitive
			exp   map[string][]string
		}{
			{
				desc: "positive case",
				input: ParamPrimitive{
					Int:     1,
					Int64:   1,
					Bool:    true,
					Float64: 10.7,
					String:  "active",
					Time:    time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC),
				},
				exp: map[string][]string{
					"int":     {"1"},
					"int64":   {"1"},
					"bool":    {"true"},
					"float64": {"10.700000"},
					"string":  {"active"},
					"time":    {"2022-01-01T00:00:00Z"},
				},
			},
			{
				desc: "positive case",
				input: ParamPrimitive{
					Int:     1,
					Int64:   1,
					Bool:    true,
					Float64: 10.7,
					String:  "active",
					Time:    time.Time{},
				},
				exp: map[string][]string{
					"int":     {"1"},
					"int64":   {"1"},
					"bool":    {"true"},
					"float64": {"10.700000"},
					"string":  {"active"},
					"time":    {""},
				},
			},
		}

		for _, tc := range testCase {
			result := map[string][]string{}
			parser := parser.InitParamParser()
			parser.Encode(tc.input, result)

			assert.DeepEqual(t, tc.exp["int"], result["int"])
			assert.DeepEqual(t, tc.exp["int64"], result["int64"])
			assert.DeepEqual(t, tc.exp["bool"], result["bool"])
			assert.DeepEqual(t, tc.exp["float64"], result["float64"])
			assert.DeepEqual(t, tc.exp["string"], result["string"])
			assert.DeepEqual(t, tc.exp["time"], result["time"])
		}
	})

	t.Run("Test Encode Sql Null Type", func(t *testing.T) {
		testCase := []struct {
			desc  string
			input ParamSqlNull
			exp   map[string][]string
		}{
			{
				desc: "positive case",
				input: ParamSqlNull{
					Int64:   sql.NullInt64{Valid: true, Int64: 99},
					Bool:    sql.NullBool{Valid: true, Bool: true},
					Float64: sql.NullFloat64{Valid: true, Float64: 10.7},
					String:  sql.NullString{Valid: true, String: "active"},
					Time:    sql.NullTime{Valid: true, Time: time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC)},
				},
				exp: map[string][]string{
					"int64":   {"99"},
					"bool":    {"true"},
					"float64": {"10.70"},
					"string":  {"active"},
					"time":    {"2022-01-01T00:00:00Z"},
				},
			},
			{
				desc: "case: sql valid false",
				input: ParamSqlNull{
					Int64:   sql.NullInt64{Valid: false},
					Bool:    sql.NullBool{Valid: false},
					Float64: sql.NullFloat64{Valid: false},
					String:  sql.NullString{Valid: false},
					Time:    sql.NullTime{Valid: false},
				},
				exp: map[string][]string{
					"int64":   {""},
					"bool":    {""},
					"float64": {""},
					"string":  {""},
					"time":    {""},
				},
			},
		}

		for _, tc := range testCase {
			result := map[string][]string{}
			parser := parser.InitParamParser()
			parser.Encode(tc.input, result)

			assert.DeepEqual(t, tc.exp["int64"], result["int64"])
			assert.DeepEqual(t, tc.exp["bool"], result["bool"])
			assert.DeepEqual(t, tc.exp["float64"], result["float64"])
			assert.DeepEqual(t, tc.exp["string"], result["string"])
			assert.DeepEqual(t, tc.exp["time"], result["time"])
		}
	})
}
