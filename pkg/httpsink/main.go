package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net/http"
	"strings"
)

type key int

const (
	requestIDKey key = 0
)

var (
	listenAddr string
	healthy    int32
	maxRequestsPerSecond int
	queues     map[string]chan string
)


type Payload struct {
	Source string
	Data  []byte

}

func main() {
	flag.StringVar(&listenAddr, "listen", ":2110", "server listen address")
	flag.IntVar(&maxRequestsPerSecond, "mrps", 10, "Max-Requests-Per-Seconds: define the throttle limit in requests per seconds")
	flag.Parse()

	// setup server
	e := echo.New()
	e.Use(middleware.Logger())
	e.HideBanner = true
	e.StdLogger.Println("starting cosmos-cash-resolver rest server")
	e.StdLogger.Println("target node is ", listenAddr)

	// start the rest server
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(middleware.RateLimiter(
		middleware.NewRateLimiterMemoryStore(rate.Limit(maxRequestsPerSecond)),
	))

	// initialize the queues
	queues = make(map[string]chan string)
	dispatcherQueue := make(chan Payload)

	go func(dq chan Payload) {
		for {
			p := <-dq
			e.StdLogger.Println("sender: ", p.Source )

			q, exists := queues[p.Source]
			if !exists {
				q = make(chan string, 2048)
				queues[p.Source] = q
			}

			q <- string(p.Data)


			e.StdLogger.Println("enqueuing event from", p.Source, " - new queue size", len(q))
		}
	}(dispatcherQueue)

	e.POST("/wh", func(c echo.Context) error {

		defer c.Request().Body.Close()
		bodyBytes, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			e.StdLogger.Println(err)
		}
		dispatcherQueue <- Payload{c.Request().Header.Get("sender"), bodyBytes}
		// track the resolution
		// atomic.AddUint64(&rt.resolves, 1)
		return c.JSON(http.StatusOK, map[string]string{})
	})

	e.GET("/messages/:agent_sender", func(c echo.Context) error {
		sender := c.Param("agent_sender")
		q, found := queues[sender]
		if !found || len(q) == 0{
			return c.JSON(http.StatusOK, []string{})
		}


		var sb strings.Builder
		prefix := "["
		for  len(q)>0 {
			sb.WriteString(prefix)
			sb.WriteString(<- q)
			prefix = ","
		}
		sb.WriteString("]")

		// track the resolution
		// atomic.AddUint64(&rt.resolves, 1)
		return c.Blob(http.StatusOK, "application/json", []byte(sb.String()))

	})
	// start the server
	e.StdLogger.Fatal(e.Start(listenAddr))
}

