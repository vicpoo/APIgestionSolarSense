// src/sensors/infrastructure/get_sensor_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/application"
)

type GetSensorHandler struct {
    useCase *application.GetSensorUseCase
}

func NewGetSensorHandler(useCase *application.GetSensorUseCase) *GetSensorHandler {
    return &GetSensorHandler{useCase: useCase}
}

func (h *GetSensorHandler) GetSensor(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }

    sensor, err := h.useCase.GetSensor(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
        return
    }

    c.JSON(http.StatusOK, sensor)
}

func (h *GetSensorHandler) GetUserSensors(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("user_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    sensors, err := h.useCase.GetUserSensors(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, sensors)
}