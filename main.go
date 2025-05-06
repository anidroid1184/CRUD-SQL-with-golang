package main

import (
	_ "fmt"
	"proyectoMySQL/connect"
	"proyectoMySQL/handlers"
	_ "proyectoMySQL/models"
)

func main(){
	connect.CloseConnection()
	defer connect.CloseConnection()

	// funcion para comenzar conexiones
	handler := handlers.NuevoClienteHandler()
	handler.Menu()

}
