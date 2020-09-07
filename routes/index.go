package routes

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

/*

 */

type tableBloodTemp struct {
	Date string
	Temp float32
}

type indexHTMLtemplate struct {
	Title string
	CSS   string
	Temps []tableBloodTemp
}

// IndexRouter  GET "/" を処理
func IndexRouter(c echo.Context) error {

	db := Repository()
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)

	err := db.PingContext(ctx)
	cancel()
	if err != nil {
		panic(err)
	}

	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 55*time.Second)
	defer cancel()

	// 潜伏期間　14日
	// 新しいほうから14日分取り出す SQL; 逆順ソートしたものを、14個分順ソートしなおして取り出す
	// select * from (select * from bloodtemp order by "DATE" desc) t where rownum<14 order by t."DATE" ;

	rows, err = db.QueryContext(ctx, `SELECT TO_CHAR("DATE"),TEMP FROM (SELECT * FROM BLOODTEMP ORDER BY "DATE" DESC) t WHERE ROWNUM<14 ORDER BY t."DATE"`)
	if err != nil {
		panic(err)
	}

	var slice []tableBloodTemp

	for rows.Next() {
		var oneline tableBloodTemp
		err = rows.Scan(&oneline.Date, &oneline.Temp)

		if err != nil {
			panic(err)
		}

		newslice := append(slice, oneline)
		slice = newslice
		//fmt.Printf("%+v\n", oneline)
	}
	rows.Close()

	htmlvariable := indexHTMLtemplate{
		Title: "体温一覧",
		CSS:   "/css/index.css",
		Temps: slice,
	}

	return c.Render(http.StatusOK, "index", htmlvariable)
}

// IndexRouterPost  POST "/" を処理
func IndexRouterPost(c echo.Context) error {

	db := Repository()
	defer db.Close()
	ctx := context.Background()

	bloodtemp := c.FormValue("bloodtemp")
	tx, _ := db.BeginTx(ctx, nil)
	if _, err := tx.Exec(`INSERT INTO BLOODTEMP VALUES (localtimestamp,:1)`, bloodtemp); err != nil {
		panic(err)
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return c.Redirect(http.StatusFound, "/")
}
