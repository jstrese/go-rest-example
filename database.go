package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

// Dummy data
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Password1!"
	dbname   = "postgres"
)

var DbConn *sql.DB

/*
	Note on pointers

	& creates a pointer to a space in memory
	* deferences the pointer

	foo := "bar"

	baz(&foo, "foo")

	func baz(strValue *string, addValue string) {
		// Dereference with * to get value in memory for pointer (&)
		*strValue = fmt.Sprintf("%s%s", addValue, *strValue)
	}

	// Prints foobar
	fmt.Println(foo)
*/
func connect() {
	var err error
	println("Connecting to database..")

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DbConn, err = sql.Open("postgres", connString)

	if err != nil {
		panic(err)
	}

	if status, _ := connected(); status == true {
		DbConn.SetMaxIdleConns(2)
		DbConn.SetMaxOpenConns(10)
	} else {
		panic("Could not connect to database")
	}

	fmt.Println("Connected to Postgres!")
}

// Helper function to make sure there is a database connection
// prior to executing a query
func connected() (status bool, err error) {
	// Default to false for connected status
	status = false
	// Check to see if the DB has been connected to
	if DbConn != nil {
		// Check to see if the connection is alive
		conn := *DbConn
		if err = conn.Ping(); err == nil {
			status = true
			return
		}

		return
	}

	err = errors.New("Not connected to a database")
	return
}

func GetPostById(id int) (post Post, err error) {
	if _, err = connected(); err == nil {
		var row *sql.Row
		var stmt string = `SELECT "id", "title", "body", "user" FROM "public"."posts" WHERE "id" = $1`
		db := *DbConn
		row = db.QueryRow(stmt, id)
		ret := new(Post)
		if err = row.Scan(&ret.Id, &ret.Title, &ret.Body, &ret.User); err == nil {
			post = *ret
			return
		} else {
			post = Post{
				Id: -1,
			}
		}
	}
	return
}

// TODO: Finish this
func CreatePost(title string, body string, user string) (id int, err error) {
	id = -1

	if _, err = connected(); err == nil {
		var row *sql.Row
		var stmt string = `insert into posts (title, body, "user") values($1, $2, $3) returning id;`
		db := *DbConn
		row = db.QueryRow(stmt, id)
		//var newId int
		if err = row.Scan(&id); err == nil {
			return
		} else {
			return
		}
	}
	return
}

func GetPostList(limit ...int) (posts []Post, err error) {
	if _, err = connected(); err == nil {
		var rows *sql.Rows
		var stmt string = `SELECT "id", "title", "body", "user" FROM "public"."posts" ORDER BY "id" DESC`

		if limit != nil && limit[0] > 0 {
			stmt = fmt.Sprintf(`%s LIMIT %d`, stmt, limit[0])
		}

		db := *DbConn
		if rows, err = db.Query(stmt); err == nil {
			defer rows.Close()

			var ret []Post
			for rows.Next() {
				newPost := new(Post)
				if err = rows.Scan(&newPost.Id, &newPost.Title, &newPost.Body, &newPost.User); err == nil {
					ret = append(ret, *newPost)
				}
			}

			posts = ret
		}
	}

	return
}
