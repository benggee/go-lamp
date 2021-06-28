package mysqlz

import (
   "database/sql"
   "log"
   "time"
)

type Option func(db *sql.DB)

func newMysql(dns string, opt ...Option) *sql.DB {
   db, err := sql.Open("mysql", dns)
   if err != nil {
      log.Fatal("mysql open fail")
   }

   for _, op := range opt {
      op(db)
   }

   return db
}

func setConnMaxLifetime(d time.Duration) Option {
   return func(db *sql.DB) {
      db.SetConnMaxLifetime(d)
   }
}

func setMaxOpenConns(n int) Option {
   return func(db *sql.DB) {
      db.SetMaxOpenConns(n)
   }
}

func setMaxIdleConns(n int) Option {
   return func(db *sql.DB) {
      db.SetMaxIdleConns(n)
   }
}
