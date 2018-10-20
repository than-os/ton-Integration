package service

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/than-os/sent-dante/dbo"
	"github.com/than-os/sent-dante/models"
	"io/ioutil"
	"net/http"
)

var dao = dbo.TON{}

type msg struct {
	Dsc 	string 			`json:"description"`
	Data 	interface{} 	`json:"data,omitempty"`
}

func GetAllTonNodes(ctx echo.Context) error {
	nodes, err := dao.FindAllTonNodes()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, msg{
			Dsc: "error occurred while getting node details",
		})
	}
		return ctx.JSON(http.StatusOK, nodes)
}

func RootFunc(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, msg{
		Dsc: "your TON node is running in Sentinel dNetwork",
	})
}

func RegisterTonNode(ctx echo.Context) error {
	node := models.TONNode{}

	b, e := ioutil.ReadAll(ctx.Request().Body)
	defer ctx.Request().Body.Close()

	if e != nil {
		panic(e)
	}
	err := json.Unmarshal(b, &node)
	if err != nil {
		panic(err)
	}
	err = dao.RegisterTonNode(node)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusCreated, msg{
		Dsc: "successfully registered with Sentinel Network as TON Node",
		Data: node,
	})

}

func GetTonNodeByIP(ctx echo.Context) error {
	node := models.TONNode{}
	b, e := ioutil.ReadAll(ctx.Request().Body)
	defer ctx.Request().Body.Close()
	if e != nil {
		return ctx.JSON(http.StatusInternalServerError, msg{
			Dsc: "error occurred while getting TON node by ip",
		})
	}

	err := json.Unmarshal(b, &node)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, msg{
			Dsc: "Request body is invalid",
			Data: ctx.Request().Body,
		})
	}

	n, err := dao.FindTonNodeByID(node.IPAddr)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, msg{
			Dsc: "error occured while getting TON node by ip",
		})
	}

	return ctx.JSON(http.StatusOK, &n)
}

func RemoveTonNode(ctx echo.Context) error {
	node := models.TONNode{}
	b, e := ioutil.ReadAll(ctx.Request().Body)
	defer ctx.Request().Body.Close()

	if e != nil {
		return ctx.JSON(http.StatusInternalServerError, msg{
			Dsc: "error occurred while getting TON node by ip",
		})
	}

	err := json.Unmarshal(b, &node)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, msg{
			Dsc: "Request body is invalid",
			Data: ctx.Request().Body,
		})
	}

	err = dao.RemoveTonNode(node.IPAddr)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, msg{
			Dsc: "error occurred while getting TON node by ip",
		})
	}

	return ctx.JSON(http.StatusAccepted, msg{
		Dsc: "your TON node has been removed from Sentinel dNetwork",
	})
}