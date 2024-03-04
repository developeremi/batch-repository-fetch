package main

import (
	"github.com/labstack/echo/v4"
	"go-find-by-id-conncurrency/batch"
	"go-find-by-id-conncurrency/database"
	"math/rand"
)

func main() {

	db, err := database.ConnectDB()
	if err != nil {
		panic("failed to connect database")
	}

	newBatch := batch.NewBatch(batch.NewByChannel(db))

	e := echo.New()
	e.GET("/read-by-channel", func(c echo.Context) error {
		id := rand.Intn(100_000) + 1
		campaign := newBatch.ReadByID(c.Request().Context(), uint(id))
		return c.JSON(200, campaign)
	})

	e.Logger.Error(e.Start(":8080"))

}
