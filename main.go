package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/vicpoo/apigestion-solar-go/src/core"
	authinfra "github.com/vicpoo/apigestion-solar-go/src/login/infrastructure"
	systemnewsinfra "github.com/vicpoo/apigestion-solar-go/src/system_news/infrastructure"
)

func main() {
	// Inicializar base de datos
	core.InitDB()

	// Verificar conexión
	db := core.GetBD()
	if err := db.Ping(); err != nil {
		log.Fatal("No se pudo hacer ping a la base de datos:", err)
	}
	fmt.Println("✅ Conexión a la base de datos verificada exitosamente")

	// Crear motor Gin
	router := gin.Default()

	// Middleware CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Cambiar en producción
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Registrar rutas
	authinfra.InitAuthRoutes(router)
	systemnewsinfra.NewSystemNewsRouter(router).Run()

	// Iniciar servidor en puerto 8000
	if err := router.Run(":8000"); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}