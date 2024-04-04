package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
	"strconv"
	"time"
)

var patlites_from_socket = []string{
	"192.168.10.1",
}

var patlites = []string{
	"192.168.10.2",
}

func handleWebSocket(c echo.Context) error {
	uuid, _ := uuid.NewUUID()
	Clients[uuid] = &Client{
		OutCh: make(chan string),
	}

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		go func() {
			for {
				select {
				case output := <-Clients[uuid].OutCh:
					err := websocket.Message.Send(ws, string(output))
					if err != nil {
						c.Logger().Error(err)
						break
					}
				}
			}
		}()

		for {
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
				break
			}
			switch msg {
			case "pause":
				globalStatus.IsPause = true
			case "resume":
				globalStatus.IsPause = false
			case "stop":
				globalStatus.IsStart = false
			case "0":
				globalStatus.StopwatchMode = true
			default:
				globalStatus.Second, _ = strconv.Atoi(msg)
				globalStatus.InputSecond, _ = strconv.Atoi(msg)
				globalStatus.IsPause = false
				globalStatus.IsStart = true
				globalStatus.StopwatchMode = false
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if !globalStatus.IsStart {
				send_patlites("000000", 0x00)
				continue
			}
			if globalStatus.IsPause {
				continue
			}
			globalStatus.Second--
			for _, client := range Clients {
				client.OutCh <- strconv.Itoa(globalStatus.Second)
			}

			if globalStatus.StopwatchMode {
				continue
			}

			if 0 < globalStatus.Second && globalStatus.Second <= 10 {
				send_patlites("020000", 0x40)
			} else if 10 < globalStatus.Second && globalStatus.Second <= 30 {
				send_patlites("010000", 0x02)
			} else if -5 < globalStatus.Second && globalStatus.Second <= 0 {
				send_patlites("100001", 0x09)
			} else if globalStatus.Second <= -5 {
				send_patlites("200000", 0x20)
			} else {
				send_patlites("001000", 0x04)
			}
		}
	}()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/", "public")
	e.GET("/ws", handleWebSocket)
	e.Logger.Fatal(e.Start(":8080"))
}
