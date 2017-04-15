package main

import (
	"os/exec"
	"strconv"
	"time"
)

var gpioPinToStatus = map[int]bool{0: false, 1: false, 2: false, 3: false, 18: false, 19: false}

func main() {
	dimUpLED()
}

func dimUpLED() {
	setGPIOOut(0)
	setGPIO(0, true)

	for level := 0; level <= 100; level++ {
		time.Sleep(time.Millisecond * 100)
		pwmGPIO(0, 50, level)
		level++
	}
}

func blinkLEDs() {
	for gpioPin := range gpioPinToStatus {
		setGPIOOut(gpioPin)
		setGPIO(gpioPin, false)
	}

	for {
		for gpioPin, status := range gpioPinToStatus {
			setGPIO(gpioPin, !status)
			gpioPinToStatus[gpioPin] = !status
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func turnOffAllGPIOs() {
	for gpioPin := range gpioPinToStatus {
		setGPIO(gpioPin, false)
	}
}

func blinkLED() {
	setGPIOOut(0)
	setGPIO(0, true)

	for {
		time.Sleep(time.Second)
		setGPIO(0, false)
		time.Sleep(time.Second)
		setGPIO(0, true)
	}
}

func pwmGPIO(pin, frequency, dutyCyclePercentage int) {
	executeGPIOCommand("pwm", strconv.Itoa(pin), strconv.Itoa(frequency), strconv.Itoa(dutyCyclePercentage))
}

func setGPIOIn(id int) {
	executeGPIOCommand("set-input", strconv.Itoa(id))
}

func setGPIOOut(id int) {
	executeGPIOCommand("set-output", strconv.Itoa(id))
}

func setGPIO(id int, status bool) {
	statusString := "0"
	if status {
		statusString = "1"
	}
	executeGPIOCommand("set", strconv.Itoa(id), statusString)
}

func executeGPIOCommand(args ...string) {
	cmd := exec.Command("fast-gpio", args...)

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
