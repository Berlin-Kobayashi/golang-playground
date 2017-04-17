package onewire

import (
	"strconv"
	"os/exec"
	"os"
	"time"
	"io/ioutil"
	"strings"
	"errors"
	"fmt"
)

var setupDelay = 3 * time.Second

var oneWireDir = "/sys/devices/w1_bus_master1"

var paths = map[string]string{
	"slaveCount": oneWireDir + "/w1_master_slave_count",
	"slaves":     oneWireDir + "/w1_master_slaves",
}

type OneWire struct {
	gpio                   int
	address, slaveFilePath string
	SetupComplete          bool
}

func NewOneWire(address string, gpio int) OneWire {
	oneWire := OneWire{
		gpio:          gpio,
		address:       address,
		slaveFilePath: oneWireDir + "/" + address + "/" + "w1_slave",
	}

	oneWire.prepare()

	return oneWire
}

func (o *OneWire) prepare() {
	if !SetupOneWire(o.gpio) {
		fmt.Println("Could not set up 1-Wire on GPIO " + strconv.Itoa(o.gpio))
		o.SetupComplete = false

		return
	}

	if !checkSlaves() {
		fmt.Println("Kernel is not recognizing slaves.")
		o.SetupComplete = false

		return
	}

	if !checkRegistered(o.address) {
		fmt.Println("Device is not registered on the bus.")
		o.SetupComplete = false

		return
	}

	o.SetupComplete = true
}

func (o *OneWire) ReadDevice() []string {
	contents, err := ioutil.ReadFile(o.slaveFilePath)
	if err != nil {
		panic(err)
	}

	message := strings.Split(string(contents), "\n")

	return message
}

func SetupOneWire(gpio int) bool {
	for i := 0; i < 2; i++ {
		if checkFilesystem() {
			return true
		}

		insertKernelModule(gpio)

		time.Sleep(setupDelay)
	}

	return false
}

func insertKernelModule(gpio int) {
	argBus := "bus0=0," + strconv.Itoa(gpio) + ",0"

	cmd := exec.Command("insmod", "w1-gpio-custom", argBus)

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func checkFilesystem() bool {
	return isDir(oneWireDir)
}

func isDir(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func checkSlaves() bool {
	contents, err := ioutil.ReadFile(paths["slaveCount"])
	if err != nil {
		panic(err)
	}

	slaveCount := strings.Split(string(contents), "\n")[0]

	if slaveCount == "0" {
		return false
	}

	return true
}

func checkRegistered(address string) bool {
	slaveList, err := scanAddresses()
	if err != nil {
		panic(err)
	}

	for _, line := range slaveList {
		if strings.Contains(line, address) {
			return true
		}
	}

	return false
}

func scanAddresses() ([]string, error) {
	if ! checkFilesystem() {
		return nil, errors.New("Directory oneWireDir does not exist!")
	}

	contents, err := ioutil.ReadFile(paths["slaves"])
	if err != nil {
		panic(err)
	}

	slaveList := strings.Split(string(contents), "\n")
	slaveList = slaveList[0:len(slaveList)-1]

	return slaveList, nil
}

func ScanOneAddress() string {
	addresses, err := scanAddresses()
	if err != nil {
		panic(err)
	}

	return addresses[0]
}
