package mysqlz

import (
   "database/sql"
   "time"
)

type dbConn struct {
   db *sql.DB
}

func newConn(dns string) *dbConn {
   return &dbConn{
      db: getDb(dns),
   }
}

func getDb(dns string) *sql.DB {
   db, err := sql.Open("mysql", dns)
   if err != nil {
      panic(err)
   }
   // See "Important settings" section.
   db.SetConnMaxLifetime(time.Minute * 3)
   db.SetMaxOpenConns(10)
   db.SetMaxIdleConns(10)

   return db
}
