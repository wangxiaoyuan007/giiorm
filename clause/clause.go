package clause

import "strings"

type Clause struct {
	sql map[Type]string
	sqlVars map[Type][]interface{}
}

func (c *Clause)Set(name Type,vars ...interface{})  {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql,vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

func (c *Clause)Build(types ...Type) (string,[]interface{})  {

	var sqls []string
	var vars []interface{}
	for _,tpe := range types{
		if sql,ok := c.sql[tpe];ok {
			sqls = append(sqls,sql)
			vars = append(vars,c.sqlVars[tpe]...)
		}
	}
	return strings.Join(sqls," "),vars
}
