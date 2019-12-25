package main

import (
	"flag"
	"log"
	"time"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	"github.com/duncanleo/go-plantower/devices"
)

var (
	pm10 = -1
	pm25 = -1
	done = make(chan bool)
)

func main() {
	var port = flag.String("port", "", "port for HC")
	var pin = flag.String("pin", "00102003", "pairing PIN for the accessory")
	var storagePath = flag.String("storagePath", "hc-plantower-storage", "storage path")

	var name = flag.String("name", "PMS5003", "name for the accessory")
	var manufacturer = flag.String("manufacturer", "Plantower", "manufacturer for the accessory")
	var model = flag.String("model", "pms5003", "model for the accessory")

	var serialDevice = flag.String("device", "/dev/ttyAMA0", "name of the serial device. e.g. COM1 on Windows, /dev/ttyAMA0 on Linux")
	var waitTime = flag.Int("wait", 2, "time to wait before getting reading from sensor device")
	var fetchInterval = flag.Int("fetchInterval", 120, "time interval between fetching data (in seconds)")

	flag.Parse()

	var ticker = time.NewTicker(time.Duration(*fetchInterval) * time.Second)

	hcConfig := hc.Config{
		Pin:         *pin,
		StoragePath: *storagePath,
		Port:        *port,
	}

	accInfo := accessory.Info{
		Name:         *name,
		Manufacturer: *manufacturer,
		Model:        *model,
	}

	acc := accessory.New(accInfo, accessory.TypeSensor)

	aqSensor := service.NewAirQualitySensor()
	aqSensor.AirQuality.OnValueGet(func() interface{} {
		return getRating(pm25)
	})

	pm25Val := characteristic.NewPM2_5Density()
	aqSensor.AddCharacteristic(pm25Val.Characteristic)
	pm25Val.OnValueGet(func() interface{} {
		return pm25
	})

	pm10Val := characteristic.NewPM10Density()
	aqSensor.AddCharacteristic(pm10Val.Characteristic)

	pm10Val.OnValueGet(func() interface{} {
		return pm10
	})

	statusActive := characteristic.NewStatusActive()
	statusActive.SetValue(true)
	aqSensor.AddCharacteristic(statusActive.Characteristic)

	statusFault := characteristic.NewStatusFault()
	statusFault.SetValue(characteristic.StatusFaultNoFault)
	aqSensor.AddCharacteristic(statusFault.Characteristic)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				data, err := devices.DeviceFuncs[*model](*serialDevice, map[string]interface{}{
					"waitTime": *waitTime,
				})
				if err != nil {
					log.Println(err)
					statusFault.UpdateValue(characteristic.StatusFaultGeneralFault)
					continue
				}
				statusFault.UpdateValue(characteristic.StatusFaultNoFault)
				pm10 = data.Atmospheric.PM10
				pm10Val.UpdateValue(pm10)
				pm25 = data.Atmospheric.PM25
				pm25Val.UpdateValue(pm25)
				break
			}
		}

	}()

	acc.AddService(aqSensor.Service)

	t, err := hc.NewIPTransport(hcConfig, acc)
	if err != nil {
		log.Fatal(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
		done <- true
	})

	t.Start()
}

func getRating(pm25 int) int {
	if pm25 >= 201 {
		return characteristic.AirQualityPoor
	} else if pm25 >= 151 {
		return characteristic.AirQualityInferior
	} else if pm25 >= 101 {
		return characteristic.AirQualityFair
	} else if pm25 >= 51 {
		return characteristic.AirQualityGood
	} else if pm25 >= 0 {
		return characteristic.AirQualityExcellent
	}
	return characteristic.AirQualityUnknown
}
