package session

import (
	"fmt"
	"giiorm/log"
	"giiorm/schema"
	"reflect"
	"strings"

)

func (s *Session) Model(value interface{}) *Session  {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value,s.dialect)
	}
	return s
}

func (s *Session)RefTable() *schema.Schema  {
	if s.refTable == nil {
		log.Error("model is not set")
	}
	return s.refTable
}

func (s *Session)CreateTable()error  {
	table := s.refTable
	var colums []string
	for _,filed := range table.Fields {
		colums = append(colums,fmt.Sprintf("%s %s %s", filed.Name,filed.Type,filed.Tag))
	}
	desc := strings.Join(colums,",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);",table.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

