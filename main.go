package main

import (
	"context"
	"fmt"
	"htmxgo/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	trate "golang.org/x/time/rate"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Menu struct {
	Name   string `json:"name"`
	Link   string `json:"link"`
	Active bool   `json:"active"`
}

type User struct {
	Name   string `json:"name" form:"name"`
	Phone  string `phone:"phone" form:"phone"`
	Email  string `json:"email" form:"email"`
	Active bool   `json:"active" form:"active"`
}

func generateUsers(startAt int) []*User {
	users := []*User{}
	for i := startAt; i < (startAt + 20); i++ {
		users = append(users, &User{
			Name:   fmt.Sprintf("Oliver Sykes %d", i+1),
			Phone:  "08123123123",
			Email:  fmt.Sprintf("oliver%d@sykes.com", i+1),
			Active: i%2 == 0,
		})
	}
	return users
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	user := &User{
		Name:   "John Doe",
		Phone:  "081234567890",
		Email:  "john@doe.com",
		Active: true,
	}
	users := generateUsers(0)

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		trate.Limit(20),
	)))

	// Initiate our template renderer
	template.NewTemplateRenderer(e, "public/*.html")

	// client route-ish

	e.GET("/", func(c echo.Context) error {
		res := map[string]interface{}{
			"Greeting": "Hola!",
			"Title":    Menu{Name: "HTMX GO!", Link: "/"},
			"Menus": []Menu{
				{Name: "Users", Link: "/users"},
				{Name: "SSE", Link: "/sse"},
			},
		}
		return c.Render(http.StatusOK, "index", res)
	})

	e.GET("/users", func(c echo.Context) error {
		return c.Render(http.StatusOK, "users", nil)
	})

	e.GET("/sse", func(c echo.Context) error {
		return c.Render(http.StatusOK, "sse", nil)
	})

	// rest api-ish

	e.GET("/api/navbar", func(c echo.Context) error {
		activeLink := c.QueryParam("active-link")

		res := map[string]interface{}{
			"Title": Menu{Name: "HTMX GO!", Link: "/"},
			"Menus": []Menu{
				{
					Name:   "Users",
					Link:   "/users",
					Active: activeLink == "users",
				},
				{
					Name:   "SSE",
					Link:   "/sse",
					Active: activeLink == "sse",
				},
			},
		}
		return c.Render(http.StatusOK, "navbar", res)
	})

	e.GET("/api/username", func(c echo.Context) error {
		return c.String(http.StatusOK, user.Name)
	})

	e.GET("/api/user", func(c echo.Context) error {
		return c.Render(http.StatusOK, "user_card", user)
	})

	e.GET("/api/user/edit", func(c echo.Context) error {
		return c.Render(http.StatusOK, "user_form", user)
	})

	e.PUT("/api/user", func(c echo.Context) error {
		// showing delay purposes
		time.Sleep(2 * time.Second)

		u := &User{}
		if err := c.Bind(u); err != nil {
			return err
		}

		user.Name = u.Name
		user.Phone = u.Phone
		user.Email = u.Email

		return c.Render(http.StatusOK, "user_card", user)
	})

	e.GET("/api/users", func(c echo.Context) error {
		// showing delay purposes
		time.Sleep(1 * time.Second)

		res := map[string]interface{}{
			"Users":    []*User{},
			"NextPage": 1,
		}

		qPage := c.QueryParam("page")
		if qPage != "" {
			page, err := strconv.Atoi(qPage)
			if err != nil {
				return err
			}

			if page > 0 {
				newUsers := generateUsers(page * 20)
				res["Users"] = newUsers
				res["NextPage"] = page + 1
				return c.Render(http.StatusOK, "row_user", res)
			}
		}

		res["Users"] = users
		return c.Render(http.StatusOK, "list_user", res)
	})

	e.GET("/api/pool-random", func(c echo.Context) error {
		res := RandStringRunes(8)
		return c.String(http.StatusOK, res)
	})

	e.GET("/sse/events", func(c echo.Context) error {
		// showing delay purposes
		time.Sleep(5 * time.Second)

		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

		// Create a channel to send data
		dataCh := make(chan string)

		// Create a context for handling client disconnection
		_, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		// Send data to the client
		go func() {
			for data := range dataCh {
				fmt.Fprintf(c.Response(), "data: %s\n\n", data)
				c.Response().Flush()
			}
		}()

		// Simulate sending data periodically
		for {
			dataCh <- time.Now().Format(time.RFC3339)
			time.Sleep(1 * time.Second)
		}
	})

	e.Logger.Fatal(e.Start(":4040"))
}
