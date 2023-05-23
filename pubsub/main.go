package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"longpolling/model"
	"net/http"
	"strconv"
	"time"
)

func main() {
	q := model.NewCappedQueue[model.Update](10)
	ps := NewPubSub()
	e := echo.New()
	e.GET("updates", func(c echo.Context) error {
		lastUpdate := c.QueryParam("lastUpdate")
		lastUpdateUnix, _ := strconv.ParseInt(lastUpdate, 10, 64)
		getUpdates := func() []model.Update {
			return model.Filter(q.Copy(), func(update model.Update) bool {
				return update.CreatedAt > lastUpdateUnix
			})
		}
		if updates := getUpdates(); len(updates) > 0 {
			return c.JSON(200, updates)
		}
		ch, close := ps.Subcribe()
		defer close()

		select {
		case <-ch:
			return c.JSON(200, getUpdates())
		case <-c.Request().Context().Done():
			return c.String(http.StatusRequestTimeout, "timeout")
		}
	})

	e.POST("send", func(c echo.Context) error {
		var request model.SendMessageRequest
		if err := c.Bind(&request); err != nil {
			return c.String(400, fmt.Sprintf("Bad request: %v", err))
		}
		q.Append(model.Update{
			CreatedAt: time.Now().Unix(),
			Message:   request.Message,
		})
		ps.Publish()
		return c.JSON(201, "Request has sent")
	})
	e.Logger.Fatal(e.Start(":8080"))

}
