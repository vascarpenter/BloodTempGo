package routes

import (
	"context"
	"database/sql"
	"math"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

/*

 */

type tableBloodTemp struct {
	Date string
	Temp float32
	Memo sql.NullString
}

type indexHTMLtemplate struct {
	Title string
	CSS   string
	Temps []tableBloodTemp
}

// makeGraph Create Graph
func makeGraph(slice []tableBloodTemp) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "Temp"
	labelX := []string{}

	pts := make(plotter.XYs, len(slice))
	for i := range slice {
		labelX = append(labelX, slice[i].Date)
		pts[i].X = float64(i)
		pts[i].Y = float64(slice[i].Temp)
	}
	l, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	p.Add(l)
	p.NominalX(labelX...)

	p.X.Tick.Label.Rotation = math.Pi / 2.5
	p.X.Tick.Label.XAlign = draw.XRight
	p.X.Tick.Label.YAlign = draw.YCenter
	p.Y.Min = 35.0
	p.Y.Max = 38.0
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "static/img/temps.png"); err != nil {
		panic(err)
	}
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

	rows, err = db.QueryContext(ctx, `SELECT TO_CHAR("DATE"),TEMP,MEMO  FROM (SELECT * FROM BLOODTEMP ORDER BY "DATE" DESC) t WHERE ROWNUM<14 ORDER BY t."DATE"`)
	if err != nil {
		panic(err)
	}

	var slice []tableBloodTemp

	for rows.Next() {
		var oneline tableBloodTemp
		err = rows.Scan(&oneline.Date, &oneline.Temp, &oneline.Memo)

		if err != nil {
			panic(err)
		}

		newslice := append(slice, oneline)
		slice = newslice
		//fmt.Printf("%+v\n", oneline)
	}
	rows.Close()

	makeGraph(slice)

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
	memo := c.FormValue("memo")
	tx, _ := db.BeginTx(ctx, nil)
	if _, err := tx.Exec(`INSERT INTO BLOODTEMP VALUES (localtimestamp,:1,:2)`, bloodtemp, memo); err != nil {
		panic(err)
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return c.Redirect(http.StatusFound, "/")
}
