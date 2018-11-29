package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/labstack/echo"
	"github.com/than-os/sent-dante/dbo"
	"github.com/than-os/sent-dante/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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
		Dsc: "Sentinel TON Interface",
	})
}

func RegisterTonNode(ctx echo.Context) error {
	node := models.TONNode{}

	b, e := ioutil.ReadAll(ctx.Request().Body)
	defer ctx.Request().Body.Close()

	log.Printf("request body: %s", b)

	if e != nil {
		//panic("error 1")
		log.Fatalf("error while reading body: \n%v", e)
		return nil
	}
	err := json.Unmarshal(b, &node)
	if err != nil {
		//panic("error 2")
		log.Println("error 1: ", err.Error())
		return nil
		//panic(err)
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

func GetActiveConnections() {

	cmd := exec.Command("netstat", "-ant")
	c2 := exec.Command("grep -E ':30001.*ESTABLISHED'")
	//c3 := exec.Command("awk", "'{printf $4}'")

	pr, pw := io.Pipe()
	cmd.Stdout = pw
	c2.Stdin = pr

	var b2 bytes.Buffer
	c2.Stdout = &b2

	cmd.Run()
	c2.Run()
	cmd.Wait()
	defer pw.Close()
	c2.Wait()
	io.Copy(os.Stdout, &b2)
	fmt.Printf("here's final output: %s", os.Stdout)
}

func KeepAlive()  {

	nodes , err := dao.FindAllTonNodes()
	if err != nil {
		log.Println("error while getting nodes: ", err)
		return
	}
	for _, node := range nodes {
		color.Red("IP Address: %s", node.IPAddr)
		b, err := MakeGetRequest("https://"+ node.IPAddr + "/live")
		if err != nil {
			log.Println("error while getting node status: ", err)
			return
		}

		if b.Message != "up" {
			err := dao.RemoveTonNode(node.IPAddr)
			if err != nil {
				log.Println("error while removing node from db: ", err)
				return
			}
		}

	}
	return
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func MakeGetRequest(url string) (models.KeepAliveResponse, error) {
	var body models.KeepAliveResponse
	resp, err := http.Get(url)
	if err != nil {
		log.Println("error while checking node status: ", err)
		return body, err
	}
	err = json.NewDecoder(resp.Body).Decode(&body)

	return body, err
}