package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("mysql", "3C4SuSETtbVjxEt.root:T22V6spwkiL9R6B@tcp(gateway01.us-west-2.prod.aws.tidbcloud.com:4000)/test")
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	var dbName string
	err = db.QueryRow("SELECT DATABASE();").Scan(&dbName)
	if err != nil {
		log.Fatal("failed to execute query", err)
	}
	fmt.Println(dbName)
	fmt.Println("db init done, max connection", db.Stats().MaxOpenConnections)

	// db.
}

func NewTable() error {
	res, err := db.Exec(`CREATE TABLE IF NOT EXISTS tasks (task_id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		priority TINYINT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)  ENGINE=INNODB;`)
	if err != nil {
		fmt.Println("create table error:", err)
		return err
	}
	fmt.Println("create table result:", res)
	return nil
}

func Close() {
	db.Close()
}

var (
	wg    sync.WaitGroup
	sqlCh = make(chan string)
)

func dbWorker() {
	defer wg.Done()
	for {
		sqlStr, ok := <-sqlCh
		if !ok {
			fmt.Println("chan closed? quit worker")
			return
		}
		res, err := db.Exec(sqlStr)
		if err != nil {
			fmt.Println("insert fail: ", err)
		} else {
			id, _ := res.LastInsertId()
			fmt.Printf("sql: %s, inserted id: %d\n", sqlStr, id)
		}
	}
}

func NewTask() {
	con := 200
	wg.Add(con)
	for i := 0; i < con; i++ {
		go dbWorker()
	}

	go func() {
		for i := 0; i < 1000000; i++ {
			now := time.Now()
			sqlStr := fmt.Sprintf("INSERT INTO tasks(title,priority) VALUES ('%s', '%d')", "task-"+now.Format(time.RFC1123Z), now.UnixNano()%10)
			sqlCh <- sqlStr
		}
		close(sqlCh)
	}()
	wg.Wait()
}
