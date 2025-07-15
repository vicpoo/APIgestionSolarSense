//api/src/sensors/application/sensor_service.go

package domain

type SensorThreshold struct {
    ID             int     `json:"id"`
    SensorID       int     `json:"sensor_id"`
    MinTemperature *float64 `json:"min_temperature,omitempty"`
    MaxTemperature *float64 `json:"max_temperature,omitempty"`
    MinHumidity    *float64 `json:"min_humidity,omitempty"`
    MaxHumidity    *float64 `json:"max_humidity,omitempty"`
    MinPressure    *float64 `json:"min_pressure,omitempty"`
    MaxPressure    *float64 `json:"max_pressure,omitempty"`
    MinVoltage     *float64 `json:"min_voltage,omitempty"`
    MaxVoltage     *float64 `json:"max_voltage,omitempty"`
    MinCurrent     *float64 `json:"min_current,omitempty"`
    MaxCurrent     *float64 `json:"max_current,omitempty"`
    ConfiguredBy   *int     `json:"configured_by,omitempty"`
}