package clause

import (
	"fmt"
	"strings"
)

type Type int
type generator func(values ...interface{}) (string,[]interface{})
var generators map[Type]generator

const (
	INSERT = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)



func init()  {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[COUNT] = _count
	generators[DELETE] = _detete
}
func genBindVars(num int) string  {
	var vars []string
	for i := 0;i < num;i++ {
		vars = append(vars,"?")
	}
	return strings.Join(vars,",")
}

func _insert(values ...interface{}) (string,[]interface{}){
	// INSERT INTO $tableName ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string),",")
	value := genBindVars(len(values[2].([]interface{})))
	return fmt.Sprintf("INSERT INTO %s (%v) VALUES (%v)",tableName,fields,value),values[2].([]interface{})
}

func _select(values ...interface{}) (string,[]interface{}){
	tableName := values[0]
	fieldsSlice := values[1].([]string)
	fields := strings.Join(fieldsSlice,",")
	return fmt.Sprintf("SELECT %v FROM %s",fields,tableName),[]interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

func _update(values ...interface{}) (string, []interface{})  {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k,v := range m {
		keys = append(keys,k + "= ?")
		vars = append(vars,v)
	}
	return fmt.Sprintf("UPDATE %s SET %s",tableName, strings.Join(keys,", ")), vars
}

func _detete(values ...interface{}) (string, []interface{})  {
	return fmt.Sprintf("DELETE FROM %s ",values[0]), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{})  {
	return _select(values[0],[]string{"COUNT(*)"})
}






