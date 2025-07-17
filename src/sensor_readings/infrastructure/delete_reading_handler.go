// src/sensor_readings/infrastructure/delete_reading_handler.go
package infrastructure

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/vicpoo/apigestion-solar-go/src/sensor_readings/application"
)

type DeleteReadingHandler struct {
    useCase *application.DeleteReadingUseCase
}

func NewDeleteReadingHandler(useCase *application.DeleteReadingUseCase) *DeleteReadingHandler {
    return &DeleteReadingHandler{useCase: useCase}
}

func (h *DeleteReadingHandler) DeleteReading(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reading ID"})
        return
    }

    if err := h.useCase.DeleteReading(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}