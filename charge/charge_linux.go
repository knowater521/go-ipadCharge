package charge

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/gousb"
)

const appleVid = 0x05ac
const ipadAdditionValue = 1600

func FoundAndSet() {
	ctx := gousb.NewContext()
	ctx.Debug(0)
	defer ctx.Close()
	dev, err := getDevice(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	defer dev.Close()
	setCharge(dev)
}

func setCharge(dev *gousb.Device) {
	log.Println("set charge", dev.Desc)
	dev.Control(gousb.ControlOut|gousb.ControlVendor, 0x40, 500, ipadAdditionValue, nil)
}

func getDevice(ctx *gousb.Context) (*gousb.Device, error) {

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		if appleVid == desc.Vendor {
			log.Println("Found", desc)
			return true
		}
		return false
	})

	if err != nil {
		log.Printf("Warning: OpenDevices: %s.", err)
		return nil, err
	}

	// All Devices returned from OpenDevices must be closed.

	if len(devs) != 1 {
		for _, d := range devs {
			d.Close()
		}
		return nil, errors.New(fmt.Sprint("set ipad Charge:Devices Num incorecct", len(devs)))
	}
	dev := devs[0]
	dev.SetAutoDetach(true)

	return dev, nil
}
