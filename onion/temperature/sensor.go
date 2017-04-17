package temperature

import (
	"github.com/DanShu93/golang-playground/onion/onewire"
	"strconv"
	"strings"
)

type TemperatureSensor struct {
	driver onewire.OneWire
	Ready  bool
}

func NewTemperatureSensor(address string, gpio int) TemperatureSensor {
	driver := onewire.NewOneWire(address, gpio)
	temperatureSensor := TemperatureSensor{
		driver: driver,
		Ready:  driver.SetupComplete,
	}

	return temperatureSensor
}

func (t *TemperatureSensor) ReadValue() float64 {
	rawValue := t.driver.ReadDevice()

	temperatureData := strings.Split(rawValue[1], " ")

	valueString := strings.Split(temperatureData[len(temperatureData)-1], "=")[1]

	value, err := strconv.ParseFloat(valueString, 64)
	if err != nil {
		panic(err)
	}

	value /= 1000

	return value
}
