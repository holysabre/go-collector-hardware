//go:build windows
// +build windows

package ipmi

import (
	"bytes"
	"collector/bin"
	"fmt"
	"regexp"
	"strings"
)

type IPMI_SensorData struct {
	Name             string
	SensorType       string
	CurrentReading   float64
	Units            string
	LowerThreshold   float64
	UpperThreshold   float64
	Enabled          bool
	TimeOfLastRead   string
	OperationalState uint16
}

func GetInfo() (sensors []Sensor) {

	args := []string{"sensor"}
	output := bin.RunCommandAndReturnBytes("ipmiutil\\ipmiutil.exe", args...)

	lines := bytes.Split(output.Bytes(), []byte{'\n'})

	pattren := regexp.MustCompile(`[0-9A-Fa-f]+\s+.*\s+snum [0-9A-Fa-f]+ (.*)\s+\= [0-9A-Fa-f]+ (?:OK)?(.*)?`)

	for _, line := range lines {
		matches := pattren.FindStringSubmatch(string(line))
		// fmt.Printf("num %d : %s - %d \n", num, string(line), len(matches))

		if len(matches) > 1 {
			sensor := Sensor{Key: strings.TrimSpace(matches[1]), Value: strings.TrimSpace(matches[2])}
			fmt.Printf("Sensor: %s, Value: %s\n", sensor.Key, sensor.Value)
			sensors = append(sensors, sensor)
		}
	}

	return sensors
}
