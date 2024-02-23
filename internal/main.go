package main

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/razullya/restAPI/pkg/swagger/server/restapi"
	"github.com/razullya/restAPI/pkg/swagger/server/restapi/operations"
)

func main() {
	// настроили сваггер
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	// создали апи
	api := operations.NewHelloAPIAPI(swaggerSpec)

	// иницируем сервер
	server := restapi.NewServer(api)

	server.Port = 8080

	api.CheckHealthHandler = operations.CheckHealthHandlerFunc(Healths)
	api.GetHelloUserHandler = operations.GetHelloUserHandlerFunc(GetHelloUser)
	api.GetGopherNameHandler = operations.GetGopherNameHandlerFunc(GetGopherByName)

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func Healths(operations.CheckHealthParams) middleware.Responder {
	return operations.NewCheckHealthOK().WithPayload("OK")
}

func GetHelloUser(user operations.GetHelloUserParams) middleware.Responder {
	return operations.NewGetHelloUserOK().WithPayload("Hello, " + user.User)
}

func GetGopherByName(gopher operations.GetGopherNameParams) middleware.Responder {

	var URL string

	if gopher.Name != "" {
		URL = "https://github.com/scraly/gophers/raw/main/" + gopher.Name + ".png"
	} else {
		URL = "https://github.com/scraly/gophers/raw/main/dr-who.png"
	}

	response, err := http.Get(URL)
	if err != nil {
		log.Println(err)
	}

	return operations.NewGetGopherNameOK().WithPayload(response.Body)
}
