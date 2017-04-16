package main

import (
	"strconv"
	"time"
	"regexp"
	"fmt"
	"flag"
	"os/exec"
)

var gpioPinToStatus = map[int]bool{0: false, 1: false, 2: false, 3: false, 18: false, 19: false}

func main() {
	useShiftRegister()
}

// Experiment 5
func useShiftRegister() {
	flag.Parse()

	input := flag.Args()[0]

	setGPIO(1, false)
	setGPIO(2, false)
	setGPIO(3, false)

	setShiftRegister(1, 2, 3, "00000000")
	setShiftRegister(1, 2, 3, input)
}

// Experiment 4
func switchLED() {
	for {
		value := readGPIO(0)
		setGPIO(1, value)
		time.Sleep(time.Millisecond * 100)
	}
}

// Experiment 3
func dimUpLED() {
	setGPIOOut(0)
	setGPIO(0, true)

	for level := 0; level <= 100; level++ {
		time.Sleep(time.Millisecond * 100)
		pwmGPIO(0, 50, level)
		level++
	}
}

// Experiment 2
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

// Experiment 1
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

func setShiftRegister(serialDataPin, serialClockPin, registerClockPin int, byteString string) {
	reversedBits := reverseBytes([]byte(byteString))

	for _, bit := range reversedBits {
		bitValue := false
		if bit == '1' {
			bitValue = true
		}

		setGPIO(serialDataPin, bitValue)

		setGPIO(serialClockPin, false)
		setGPIO(serialClockPin, true)
		setGPIO(serialClockPin, false)
	}

	setGPIO(registerClockPin, false)
	setGPIO(registerClockPin, true)
	setGPIO(registerClockPin, false)
}

func reverseBytes(bytes []byte) []byte {
	reversedBytes := make([]byte, len(bytes))

	for i, j := len(bytes)-1, 0; i >= 0; i, j = i-1, j+1 {
		reversedBytes[j] = bytes[i]
	}

	return reversedBytes
}

func readGPIO(pin int) bool {
	out := executeGPIOCommand("read", strconv.Itoa(pin))

	regex := regexp.MustCompile(fmt.Sprintf("(> Read GPIO%d: )", pin))

	value := regex.ReplaceAllString(string(out), "")
	if value == "1\n" {
		return true
	}

	return false
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

func executeGPIOCommand(args ...string) []byte {
	//fmt.Println("fast-gpio", args)
	//return []byte{}
	cmd := exec.Command("fast-gpio", args...)

	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	return out
}
