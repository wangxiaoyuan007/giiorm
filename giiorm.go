package giiorm

import (
	"database/sql"
	"fmt"
	"giiorm/dialect"
	"giiorm/log"
	"giiorm/session"
	"strings"
)

type Engine struct {
	db *sql.DB
	dialect dialect.Dialect
}
type TxFunc func(s *session.Session) (interface{},error)
func NewEngine(driver, source string) (e *Engine,err error)  {
	db,err := sql.Open(driver,source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	dialect ,ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db : db,dialect: dialect}
	log.Info("Connect database success")
	return
}

func (e *Engine) Close() {
	if  err := e.db.Close(); err != nil {
		log.Error("Fail to close databases")
	}
	log.Info("Close database success")
}

func (e * Engine)NewSession()*session.Session  {
	return session.New(e.db,e.dialect)
}

func (e *Engine)Transaction(fun TxFunc) (result interface{},err error) {
	s := e.NewSession()
	if err = s.Begin();err != nil {
		return nil,err
	}
	defer func() {
		if p := recover();p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()
	return fun(s)

}



func GetDiffColum(c1 []string, c2 []string)(diff []string)  {
	flag := make(map[string]bool)
	for _,v := range c2{
		flag[v] = true
	}

	for _,v := range c1{
		if _,ok :=flag[v];!ok {
			diff = append(diff,v)
		}
	}
	return
}

func (e *Engine)Migrate(value interface{}) error  {
	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		table := s.RefTable()
		rows,_ := s.Raw(fmt.Sprintf("SELECT * FROM %s Limit",table.Name)).QueryRows()
		colums,_ := rows.Columns()
		addCols := GetDiffColum(table.FieldName,colums)
		delCols := GetDiffColum(colums,table.FieldName)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)

		for _,col := range addCols{
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUM %s %s",table.Name,f.Name,f.Type)
			if _,err := s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}
		tmpName := "tmp_" + table.Name
		fielStr := strings.Join(delCols,", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmpName, fielStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmpName, table.Name))
		_,err = s.Exec()
		return
	})
	return err
}

