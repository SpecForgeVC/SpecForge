package db

import (
	"database/sql"
	"encoding/json"

	"github.com/sqlc-dev/pqtype"
)

func TextToSql(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func SqlToText(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func BoolToSql(b bool) sql.NullBool {
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}

func JSONToSql(m map[string]interface{}) pqtype.NullRawMessage {
	if m == nil {
		return pqtype.NullRawMessage{Valid: false}
	}
	b, _ := json.Marshal(m)
	return pqtype.NullRawMessage{
		RawMessage: b,
		Valid:      true,
	}
}

func SqlToJSON(nrm pqtype.NullRawMessage) map[string]interface{} {
	if !nrm.Valid {
		return make(map[string]interface{})
	}
	var m map[string]interface{}
	json.Unmarshal(nrm.RawMessage, &m)
	return m
}

func BytesToPQRawMessage(b []byte) pqtype.NullRawMessage {
	return pqtype.NullRawMessage{
		RawMessage: b,
		Valid:      true,
	}
}

func RawMessageToJSON(rm json.RawMessage) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(rm, &m)
	return m
}

func JSONToRawMessage(m map[string]interface{}) json.RawMessage {
	b, _ := json.Marshal(m)
	return b
}

func FloatToSql(f float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}
}

func ToSqlString(s string) sql.NullString {
	return TextToSql(s)
}

func ToSqlBool(b bool) sql.NullBool {
	return BoolToSql(b)
}

func ToSqlInt32(i int32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: i,
		Valid: true,
	}
}
