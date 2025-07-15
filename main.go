//main.go
package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"

    "github.com/vicpoo/apigestion-solar-go/src/core"
    authinfra "github.com/vicpoo/apigestion-solar-go/src/login/infrastructure"
    notificationinfra "github.com/vicpoo/apigestion-solar-go/src/notification_settings/infrastructure"
    readinginfra "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/infrastructure"
    sensorinfra "github.com/vicpoo/apigestion-solar-go/src/sensors/infrastructure"
    sessioninfra "github.com/vicpoo/apigestion-solar-go/src/sessions/infrastructure"
    thresholdinfra "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/infrastructure"
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

    // Middleware de headers de seguridad
    router.Use(func(c *gin.Context) {
        c.Header("Cross-Origin-Opener-Policy", "same-origin-allow-popups")
        c.Header("Cross-Origin-Embedder-Policy", "require-corp")
        c.Next()
    })

    // Configuración CORS mejorada
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"https://solarsense.zapto.org", "http://localhost:4200"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Registrar todas las rutas
    authinfra.InitAuthRoutes(router)
    sensorinfra.InitSensorRoutes(router)
    readinginfra.InitReadingRoutes(router)
    thresholdinfra.InitThresholdRoutes(router)
    notificationinfra.InitSettingsRoutes(router)
    sessioninfra.InitSessionRoutes(router)
    systemnewsinfra.NewSystemNewsRouter(router).Run()

    // Iniciar servidor
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }
    if err := router.Run(":" + port); err != nil {
        log.Fatal("Error al iniciar el servidor:", err)
    }
}