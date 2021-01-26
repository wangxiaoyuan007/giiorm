package session

import(
	"giiorm/clause"
	"reflect"
)
func (s *Session)Insert(value interface{})(int64,error)  {
	s.CallMethod(BeforeInsert,value)
	table := s.Model(value).refTable
	s.clause.Set(clause.INSERT,table.Name,table.FieldName,table.RecordValues(value))
	sql, vars := s.clause.Build(clause.INSERT)
	result ,err := s.Raw(sql,vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert,nil)
	return result.RowsAffected()

}

func (s *Session) Find(values interface{}) error {
	s.CallMethod(BeforeQuery,nil)
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldName)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldName {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))

	}

	return rows.Close()
}

func (s *Session)Update(kv ...interface{})(int64,error)  {
	s.CallMethod(BeforeUpdate,nil)
	m,ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE,s.refTable.Name,m)
	sql ,vars := s.clause.Build(clause.UPDATE,clause.WHERE)
	result, err := s.Raw(sql,vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate,nil)
	return result.RowsAffected()
}

func (s *Session)Delete() (int64,error) {
	s.CallMethod(BeforeDelete,nil)
	s.clause.Set(clause.DELETE,s.refTable.Name)
	sql ,vars := s.clause.Build(clause.DELETE,clause.WHERE)
	result, err := s.Raw(sql,vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete,nil)
	return result.RowsAffected()
}

func (s *Session)Count()  (int64,error)  {
	s.clause.Set(clause.COUNT,s.refTable.Name)
	sql ,vars := s.clause.Build(clause.COUNT,clause.WHERE)
	row := s.Raw(sql,vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp,nil
}