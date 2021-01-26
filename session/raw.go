package session

import (
	"database/sql"
	"errors"
	"giiorm/clause"
	"giiorm/dialect"
	"giiorm/log"
	"giiorm/schema"
	"reflect"
	"strings"
)

type Session struct {
	db *sql.DB
	tx *sql.Tx
	sql strings.Builder
	clause clause.Clause
	dialect dialect.Dialect
	refTable *schema.Schema
	sqlVars []interface{}
}
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func New(db *sql.DB,dialect dialect.Dialect) *Session  {
	return &Session{db: db,dialect: dialect}
}

func (s *Session)DB()CommonDB  {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}
func (s *Session)Clear()  {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

func (s *Session)Raw (sql string, values ...interface{}) *Session  {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars,values...)
	return s
}

func (s *Session)Exec() (result sql.Result, err error)  {
	defer s.Clear()
	log.Info(s.sql.String(),s.sqlVars)
	if result,err = s.db.Exec(s.sql.String(),s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session)QueryRow() *sql.Row  {
	defer s.Clear()
	log.Info(s.sql.String(),s.sqlVars)
	return s.db.QueryRow(s.sql.String(),s.sqlVars...)
}

func (s *Session)QueryRows() (rows *sql.Rows,err error)  {
	defer s.Clear()
	log.Info(s.sql.String(),s.sqlVars)
	if rows, err = s.db.Query(s.sql.String(),s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session)Limit(num int) *Session  {
	s.clause.Set(clause.LIMIT,num)
	return s
}

func (s *Session)Where(desc string,args ...interface{}) *Session  {
	var vars []interface{}
	s.clause.Set(clause.WHERE,append(append(vars,desc),args...)...)
	return s
}

func (s *Session)OrderBy(desc string) *Session  {
	s.clause.Set(clause.ORDERBY,desc)
	return s
}

func (s *Session)First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface());err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil

}