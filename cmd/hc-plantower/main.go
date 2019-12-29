package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tidwall/gjson"

	"github.com/brutella/hc/service"
)

var (
	pm10 int64 = -1
	pm25 int64 = -1
	done       = make(chan bool)
)

func connect(clientID string, uri *url.URL) (mqtt.Client, error) {
	var opts = mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientID)

	var client = mqtt.NewClient(opts)
	var token = client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	return client, token.Error()
}

func main() {
	var port = flag.String("port", "", "port for HC")
	var pin = flag.String("pin", "00102003", "pairing PIN for the accessory")
	var storagePath = flag.String("storagePath", "hc-plantower-storage", "storage path")

	var name = flag.String("name", "PMS5003", "name for the accessory")
	var manufacturer = flag.String("manufacturer", "Plantower", "manufacturer for the accessory")
	var model = flag.String("model", "pms5003", "model for the accessory")

	var brokerURI = flag.String("brokerURI", "mqtt://127.0.0.1:1883", "URI of the MQTT broker")
	var clientID = flag.String("clientID", "hc-mqtt-temperature", "client ID for MQTT")

	var topic = flag.String("topic", "air", "topic to subscribe to in MQTT")
	var pm25JSONPath = flag.String("pm2.5JSONPath", "pm2\\.5", "JSON path to pm2.5 data")
	var pm10JSONPath = flag.String("pm10JSONPath", "pm10", "JSON path to pm10 data")

	flag.Parse()

	mqttURI, err := url.Parse(*brokerURI)
	if err != nil {
		log.Fatal(err)
	}

	mqttClient, err := connect(*clientID, mqttURI)
	if err != nil {
		log.Fatal(err)
	}

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

	mqttClient.Subscribe(*topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		pm25 = gjson.Get(string(msg.Payload()), *pm25JSONPath).Int()
		pm10 = gjson.Get(string(msg.Payload()), *pm10JSONPath).Int()

		pm25Val.UpdateValue(pm25)
		pm10Val.UpdateValue(pm10)
		aqSensor.AirQuality.UpdateValue(getRating(pm25))
	})

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

func getRating(pm25 int64) int {
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
