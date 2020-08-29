package routes

import (
	"database/sql"
	"fmt"
	"os"

	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-oci8"
)

// Repository DBを開
func Repository() *sql.DB {

	// this is oracle way
	var ocistring string
	if ocistring = os.Getenv("OCISTRING"); ocistring == "" {
		fmt.Printf("OCI adapter string not specified in 'OCISTRING' environment variable. Program will quit")
		os.Exit(1)
	}

	db, err := sql.Open("oci8", ocistring)

	// this is mysql way; be sure to import with

	//sqlstring := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true&loc=Asia%%2FTokyo", "oge", "hogehogeA00", "127.0.0.1", "dchild")
	//db, err := sql.Open("mysql", sqlstring)

	if err != nil {
		panic(err)
	}

	return db
}
