// src/sensors/infrastructure/delete_sensor_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensors/application"
)

type DeleteSensorHandler struct {
    useCase *application.DeleteSensorUseCase
}

func NewDeleteSensorHandler(useCase *application.DeleteSensorUseCase) *DeleteSensorHandler {
    return &DeleteSensorHandler{useCase: useCase}
}

func (h *DeleteSensorHandler) DeleteSensor(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
        return
    }

    if err := h.useCase.DeleteSensor(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}