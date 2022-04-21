package controllers

import (
	"encoding/json"
	"fmt"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/websocket"
	"stream/models"
)

// Controller interface has two methods
type Controller interface {
	// HomeController renders initial home page
	HomeController(e echo.Context) error

	// StreamController responds with live price status over websocket
	StreamController(e echo.Context) error
}

type controller struct {
}

func NewController() Controller {
	return &controller{}
}

var model models.Currency1

// Init Initializes the models
func Init() {
	model = models.NewCurrency1()
}

func (c *controller) HomeController(e echo.Context) error {
	return e.File("views/index.html")
}

func (c *controller) StreamController(e echo.Context) error {

	websocket.Handler(func(ws *websocket.Conn) {
		defer func(ws *websocket.Conn) {
			err := ws.Close()
			if err != nil {
				log.Errorf("Controllers (websocker handler): %v", err)
			}
		}(ws)
		status, err := model.GetLiveCurrency1()
		if err != nil {
			log.Errorf("Controllers (websocker handler status): %v", err)
			return
		}
		for {
			// Write
			newVal := <-status
			jsonResponse, _ := json.Marshal(newVal)
			err := websocket.Message.Send(ws, fmt.Sprintln(string(jsonResponse)))
			if err != nil {
				log.Errorf("Controllers (websocker message): %v", err)
			}
		}
	}).ServeHTTP(e.Response(), e.Request())
	return nil
}
