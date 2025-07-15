//api/src/sensors/infrastructure/mysql_sensor_repository.go

package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/application"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/domain"
)

type SensorHandlers struct {
    service *application.SensorService
}

func NewSensorHandlers(service *application.SensorService) *SensorHandlers {
    return &SensorHandlers{service: service}
}

func (h *SensorHandlers) CreateSensor(c *gin.Context) {
    var sensor domain.Sensor
    if err := c.ShouldBindJSON(&sensor); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.CreateSensor(c.Request.Context(), &sensor); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, sensor)
}

func (h *SensorHandlers) GetSensor(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    sensor, err := h.service.GetSensor(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
        return
    }
    c.JSON(http.StatusOK, sensor)
}

func (h *SensorHandlers) GetUserSensors(c *gin.Context) {
    userID, _ := strconv.Atoi(c.Param("user_id"))
    sensors, err := h.service.GetUserSensors(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, sensors)
}

func (h *SensorHandlers) UpdateSensor(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var sensor domain.Sensor
    if err := c.ShouldBindJSON(&sensor); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    sensor.ID = id
    if err := h.service.UpdateSensor(c.Request.Context(), &sensor); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, sensor)
}

func (h *SensorHandlers) DeleteSensor(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.service.DeleteSensor(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}