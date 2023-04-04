//go:build linux
// +build linux

package cpu

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/cpu"
)

func GetInfo() (cpuObj CpuObj) {
	usageCpus := getUsage()
	cpuObj.Usage = usageCpus

	temperatureCpus := getTemperature()
	cpuObj.Temperature = temperatureCpus

	return cpuObj
}

type CpuInfo struct {
	CoreID int `json:"cpu"`
	CpuId  int `json:"cores"`
}

func getUsage() (cpuAttrs []CpuAttr) {
	is, err := cpu.Info()
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	var mapCore = make(map[int]string)
	for _, i := range is {
		_id := i.PhysicalID + "-" + i.CoreID
		mapCore[int(i.CPU)] = _id
	}
	ts, err := cpu.Times(true)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	for index, t := range ts {
		total := t.Total()
		usage := float64(total-t.Idle) / float64(total) * 100.0
		_value := strconv.Itoa(int(usage))

		_id := mapCore[index]

		cpuAttr := CpuAttr{ID: _id, Value: _value}
		fmt.Printf("Core %s Usage: %s\n", _id, _value)

		cpuAttrs = append(cpuAttrs, cpuAttr)
	}

	return cpuAttrs
}

func getTemperature() (cpuAttrs []CpuAttr) {
	// 获取 CPU 温度
	tempPath := ""
	tempPath_prefix := ""
	if _, err := os.Stat("/sys/class/thermal/thermal_zone0/temp"); err == nil {
		// Intel CPU 温度传感器路径
		tempPath_prefix = "/sys/class/thermal"
		tempPath = tempPath_prefix + "/thermal_zone0/temp"
	} else if _, err := os.Stat("/sys/class/hwmon/hwmon0/temp1_input"); err == nil {
		// AMD CPU 温度传感器路径
		tempPath_prefix = "/sys/class/hwmon"
		tempPath = tempPath_prefix + "/hwmon0/temp1_input"
	} else {
		panic("Unable to find CPU temperature sensor")
	}
	tempData, err := ioutil.ReadFile(tempPath)
	if err != nil {
		panic(err)
	}
	temp, err := strconv.ParseFloat(strings.TrimSpace(string(tempData)), 64)
	if err != nil {
		panic(err)
	}
	temp = temp / 1000.0 // 转换为摄氏度
	fmt.Println(temp)

	// 获取 CPU 核心数量
	cpuInfoData, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		panic(err)
	}
	coreCount := 0
	threadRegexp := regexp.MustCompile(`^cpu cores\s*:.*$`)
	reader := strings.NewReader(string(cpuInfoData))
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if threadRegexp.MatchString(line) {
			// fmt.Println("line:", line)
			fields := strings.Fields(line)
			// fmt.Println("field-len:", len(fields))
			if len(fields) > 1 {
				coreCount, _ = strconv.Atoi(fields[3])
				break
			}
		}
	}
	fmt.Println("coreCount:", coreCount)
	if coreCount == 0 {
		panic("Unable to determine CPU core count")
	}

	// 获取每个 CPU 核心的温度
	cpuIndex := 0
coreOuter:
	for num := 0; num <= coreCount; num++ {
		for i := 1; i <= coreCount; i++ {
			coreNamePath := fmt.Sprintf("/sys/class/hwmon/hwmon%d/name", num)
			coreName, err := ioutil.ReadFile(coreNamePath)
			trimmedCoreName := strings.TrimSpace(string(coreName))
			if err != nil || !isCpuName(trimmedCoreName) {
				continue coreOuter
			}
			corePath := fmt.Sprintf("/sys/class/hwmon/hwmon%d/temp%d_input", num, i)
			coreData, err := ioutil.ReadFile(corePath)
			if err != nil {
				continue
			}

			coreTemp, err := strconv.ParseFloat(strings.TrimSpace(string(coreData)), 64)
			if err != nil {
				continue
			}
			coreTemp = coreTemp / 1000.0 // 转换为摄氏度

			_id := strconv.Itoa(cpuIndex) + "-" + strconv.Itoa(i-1)
			_value := strconv.Itoa(int(coreTemp))

			cpuAttr := CpuAttr{ID: _id, Value: _value}
			fmt.Printf("%s Temperature: %s C\n", _id, _value)

			cpuAttrs = append(cpuAttrs, cpuAttr)
		}
		cpuIndex++
	}
	return cpuAttrs
}

func isCpuName(name string) bool {
	return name == "coretemp" || // Intel CPUs
		name == "k10temp" // AMD Ryzen CPUs
}
