package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"strconv"
	"time"
)

type SendMessageRequest struct {
	Message string `json:"message"`
}

type Update struct {
	CreatedAt int64
	Message string
}

func filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func main()  {
	//normal case
	//q:=NewCappedQueue[string](10)
	//e:=echo.New()
	//e.GET("updates", func(c echo.Context) error {
	//	return c.JSON(200,q.Copy())
	//})
	//e.POST("send", func(c echo.Context) error {
	//	var request SendMessageRequest
	//	if err:=c.Bind(&request); err != nil {
	//		return c.String(400,fmt.Sprintf("Bad request: %v", err))
	//	}
	//	q.Append(request.Message)
	//	return c.JSON(201, "Request has sent")
	//})
	//e.Logger.Fatal(e.Start(":8080"))


	q := NewCappedQueue[Update](10)
	e:=echo.New()
	e.GET("updates", func(c echo.Context) error {
		lastUpdate := c.QueryParam("lastUpdate")
		lastUpdateUnix, _ := strconv.ParseInt(lastUpdate, 10, 64)
		var updates []Update
		for{
			updates = filter(q.Copy(), func(update Update) bool {
				return update.CreatedAt > lastUpdateUnix
			})
			if len(updates) != 0 {
				break
			}
			select {
				case <- c.Request().Context().Done():
				case <- time.After(time.Second):
			}
		}
		return c.JSON(200,updates)
	})

	e.POST("send", func(c echo.Context) error {
			var request SendMessageRequest
			if err:=c.Bind(&request); err != nil {
				return c.String(400,fmt.Sprintf("Bad request: %v", err))
			}
			q.Append(Update{
				CreatedAt: time.Now().Unix(),
				Message: request.Message,
			})
			return c.JSON(201, "Request has sent")
	})
	e.Logger.Fatal(e.Start(":8080"))
}
