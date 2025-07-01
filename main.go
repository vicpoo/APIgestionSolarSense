// main.go
package main

import (
	"fmt"
	"log"

	"github.com/vicpoo/apigestion-solar-go/src/core" 
)

func main() {

	// Inicializar la base de datos
	core.InitDB()

	// Verifica si la conexión sigue viva
	db := core.GetBD()
	err := db.Ping()
	if err != nil {
		log.Fatal("No se pudo hacer ping a la base de datos:", err)
	}

	fmt.Println("¡Conexión a la base de datos verificada exitosamente!")
}
