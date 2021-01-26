package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct {}
func init() {
	RegisterDialect("mysql", &mysql{})
}
func (m *mysql)DataTypeOf(typ reflect.Value) string  {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int8:
		return "tinyint"
	case reflect.Int16:
		return "smallint"
	case reflect.Int,reflect.Int32:
		return "integer"
	case reflect.Int64:
		return "bigint"
	case reflect.Uint,reflect.Uint32:
		return "integer unsigned"
	case reflect.String:
		return "varchar(255)"
	case reflect.Float32,reflect.Float64:
		return "double precision"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}

	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func  (m *mysql)TableExistSql(tableName string) (string,[]interface{})  {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}