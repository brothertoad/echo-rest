package main

import (
  "database/sql"
  "fmt"
  "math/rand"
  "net/http"
  "sort"
  "strings"
  "time"
  _ "github.com/jackc/pgx/stdlib"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "github.com/urfave/cli/v2"
)

var serveCommand = cli.Command {
  Name: "serve",
  Usage: "run as a REST service",
  Flags: []cli.Flag {
    &cli.IntFlag {Name: "port", Aliases: []string{"p"}, Value: 9903},
  },
  Action: doServe,
}

func doServe(c *cli.Context) error {
  port := c.Int("port")
  db := openDB()
	defer db.Close()

	e := echo.New()
  e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Format: "${time_rfc3339} ${method} uri=${uri} status=${status} error=${error}\n",
  }))
  e.Use(middleware.CORS()) // allow all requests

  // routes for blocks
  e.GET("/list/:name", func(c echo.Context) error {
		return getListForREST(c, db, false)
	})
  e.GET("/randomlist/:name", func(c echo.Context) error {
		return getListForREST(c, db, true)
	})
  e.GET("/asclist/:name", func(c echo.Context) error {
		return getAscListForREST(c, db)
	})
  e.GET("/desclist/:name", func(c echo.Context) error {
		return getDescListForREST(c, db)
	})
	e.GET("/block/:name", func(c echo.Context) error {
		return getBlockForREST(c, db)
	})
	e.POST("/block", func(c echo.Context) error {
		return putBlockForREST(c, db)
	})

  // routes for weight
  e.POST("/weight/daily/add", func(c echo.Context) error {
    return addDailyWeight(c, db)
  })
  e.GET("/weight/latest", func(c echo.Context) error {
    return getLatest(c, db)
  })
  e.GET("/weight/lowest", func(c echo.Context) error {
    return getLowest(c, db)
  })
  e.GET("/weight/latest-months", func(c echo.Context) error {
    return getLatestMonths(c, db)
  })
  e.GET("/weight/month-averages/low", func(c echo.Context) error {
    return getMonthAverages(c, db, true)
  })
  e.GET("/weight/month-averages/high", func(c echo.Context) error {
    return getMonthAverages(c, db, false)
  })
  e.GET("/weight/year-averages", func(c echo.Context) error {
    return getYearAverages(c, db)
  })
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
  return nil
}

func getListForREST(c echo.Context, db *sql.DB, randomize bool) error {
  response, err := getListResponse(c, db)
  if err != nil {
    return err
  }
  if randomize {
    // logic taken from https://www.calhoun.io/how-to-shuffle-arrays-and-slices-in-go/
    r := rand.New(rand.NewSource(time.Now().Unix()))
    for n := len(response.Items); n > 0; n-- {
      randIndex := r.Intn(n)
      response.Items[n-1], response.Items[randIndex] = response.Items[randIndex], response.Items[n-1]
    }
  }
	return c.JSON(http.StatusOK, response)
}

func getAscListForREST(c echo.Context, db *sql.DB) error {
  response, err := getListResponse(c, db)
  if err != nil {
    return err
  }
  sort.Slice(response.Items, func(i, j int) bool {
    return response.Items[i] < response.Items[j]
  })
  return c.JSON(http.StatusOK, response)
}

func getDescListForREST(c echo.Context, db *sql.DB) error {
  response, err := getListResponse(c, db)
  if err != nil {
    return err
  }
  sort.Slice(response.Items, func(i, j int) bool {
    return response.Items[j] < response.Items[i]
  })
  return c.JSON(http.StatusOK, response)
}

func getListResponse(c echo.Context, db *sql.DB) (*ListResponse, error) {
  block := new(BlockRequest)
	if err := c.Bind(block); err != nil {
		return nil, err
	}
  getBlock(db, block)
  response := new(ListResponse)
  response.Name = block.Name
  response.ModTime = block.ModTime
  rawList := strings.Split(block.Contents, "\n")
  response.Items = make([]string, 0, len(rawList))
  for _, rawItem := range(rawList) {
    item := strings.TrimSpace(rawItem)
    if len(item) > 0 && !strings.HasPrefix(item, "#") {
      response.Items = append(response.Items, item)
    }
  }
  return response, nil
}

// Should return a struct
func parseMarkdown(s string) string {
  if !strings.HasPrefix(s, "[") {
    return s
  }
  j := strings.Index(s, "]")
  if j < 0 {
    fmt.Printf("Can't find closing bracket in markdown '%s'\n", s)
    return s // something is wrong here
  }
  text := s[1:j]
  remaining := s[j:]
  j = strings.Index(remaining, "(") // probably right after closing brace
  k := strings.Index(remaining, ")")
  url := remaining[(j+1):k]
  // fmt.Printf("text is '%s', and url is '%s'\n", text, url)
  return fmt.Sprintf("<a href=\"%s\">%s</a>", url, text)
}

func getBlockForREST(c echo.Context, db *sql.DB) error {
	block := new(BlockRequest)
	if err := c.Bind(block); err != nil {
		return err
	}
  getBlock(db, block)
	return c.JSON(http.StatusOK, block)
}

func putBlockForREST(c echo.Context, db *sql.DB) error {
	block := new(BlockRequest)
	if err := c.Bind(block); err != nil {
		return err
	}
  err := putBlock(db, block)
  if err != nil {
    return c.String(http.StatusInternalServerError, err.Error())
  }
	return c.String(http.StatusOK, "")
}
