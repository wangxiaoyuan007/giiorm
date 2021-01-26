package main

import (
	"giiorm"
	"giiorm/log"
	"giiorm/session"
	_ "github.com/go-sql-driver/mysql"
)
type User struct {
	Name string `giiorm:"PRIMARY KEY"`
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

type Account struct {
	ID       int `geeorm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *session.Session) error {
	log.Info("before inert", account.ID)
	account.ID += 1000
	return nil
}

func (account *Account) AfterQuery(s *session.Session) error {
	log.Info("after query", account.ID)
	account.Password = "******"
	return nil
}
func main()  {
	engine, _ := giiorm.NewEngine("mysql","root:admin@tcp(127.0.0.1:3306)/school")

	s := engine.NewSession().Model(&Account{})
	s.DropTable()
	s.CreateTable()
	_, _ = s.Insert(&Account{1, "123456"})
	_, _ = s.Insert(&Account{2, "qwerty"})
	u := &Account{}
	err := s.First(u)
	if err != nil || u.ID != 1001 || u.Password != "******" {
		log.Error("Failed to call hooks after query, got", u)
	}
}


