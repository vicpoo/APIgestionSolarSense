// api/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	alertinfra "github.com/vicpoo/apigestion-solar-go/src/alerts/infrastructure"
	"github.com/vicpoo/apigestion-solar-go/src/core"
	"github.com/vicpoo/apigestion-solar-go/src/email"
	authinfra "github.com/vicpoo/apigestion-solar-go/src/login/infrastructure"
	membershipinfra "github.com/vicpoo/apigestion-solar-go/src/memberships/infrastructure"
	notificationinfra "github.com/vicpoo/apigestion-solar-go/src/notification_settings/infrastructure"
	reportinfra "github.com/vicpoo/apigestion-solar-go/src/reports/infrastructure"
	readinginfra "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/infrastructure"
	thresholdinfra "github.com/vicpoo/apigestion-solar-go/src/sensor_thresholds/infrastructure"
	sensorinfra "github.com/vicpoo/apigestion-solar-go/src/sensors/infrastructure"
	sessioninfra "github.com/vicpoo/apigestion-solar-go/src/sessions/infrastructure"
	systemnewsinfra "github.com/vicpoo/apigestion-solar-go/src/system_news/infrastructure"
    "github.com/vicpoo/apigestion-solar-go/src/alerts/application"
    "github.com/vicpoo/apigestion-solar-go/src/workers"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}
	// Inicializar base de datos
	core.InitDB()

	// Verificar conexión
	db := core.GetBD()
	if err := db.Ping(); err != nil {
		log.Fatal("No se pudo hacer ping a la base de datos:", err)
	}
	fmt.Println("✅ Conexión a la base de datos verificada exitosamente")

	// Configuración del email
	emailService := email.NewEmailService(
		"smtp.gmail.com",
		465,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
		"polarsoftsenss@gmail.com",
	)

// Inicialización del servicio de alertas
alertService := application.NewAlertService(
    alertinfra.NewMySQLAlertRepository(),
    authinfra.NewAuthRepository(db),
    notificationinfra.NewMySQLSettingsRepository(),
    emailService,
)

	// Iniciar worker de alertas
	go func() {
		worker := workers.NewAlertWorker(alertService, 5*time.Minute)
		worker.Start(context.Background())
	}()
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
		AllowOrigins:     []string{"https://solarsense.zapto.org", "http://localhost:4200", "https://frontsolarsense.servepics.com", "http://3.229.144.5/","http://3.229.144.5"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Registrar todas las rutas
	alertinfra.InitAlertRoutes(router, emailService)
	authinfra.InitAuthRoutes(router)
	membershipinfra.InitMembershipRoutes(router)
	notificationinfra.InitSettingsRoutes(router)
	readinginfra.InitReadingRoutes(router)
	reportinfra.InitReportRoutes(router)
	sensorinfra.InitSensorRoutes(router)
	sessioninfra.InitSessionRoutes(router)
	thresholdinfra.InitThresholdRoutes(router)
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