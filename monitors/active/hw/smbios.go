package hw

import (
	"fmt"

	"github.com/lijingwei9060/infobeat/monitors"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/lijingwei9060/go-smbios/smbios"
	"github.com/lijingwei9060/infobeat/monitors/jobs"
)

var debugf = logp.MakeDebug("smbios")

func create(name string, cfg *common.Config) ([]jobs.Job, error) {
	return []jobs.Job{jobs.MakeSimpleJob(func(event *beat.Event) error {
		hw := common.MapStr{}
		bios := smbios.GetSMBIOS()

		if bios == nil {
			return nil
		}

		hw.Put("smbios.major", bios.Major)
		hw.Put("smbios.minor", bios.Minor)

		// type 0
		hw.Put("smbios.bios.vendor", bios.BIOSInformation.Vendor)
		hw.Put("smbios.bios.version", bios.BIOSInformation.BIOSVersion)
		hw.Put("smbios.bios.releasedate", bios.BIOSInformation.BIOSReleaseDate)
		// type 1
		hw.Put("smbios.system.manufacturer", bios.SystemInformation.Manufacturer)
		hw.Put("smbios.system.productname", bios.SystemInformation.ProductName)
		hw.Put("smbios.system.serialnumber", bios.SystemInformation.SerialNumber)
		hw.Put("smbios.system.version", bios.SystemInformation.Version)
		hw.Put("smbios.system.uuid", bios.SystemInformation.UUID)
		// type 2
		for n, b := range bios.BaseboardInformations {
			hw.Put(fmt.Sprintf("smbios.baseboard.%d.manufacturer", n), b.Manufacturer)
			hw.Put(fmt.Sprintf("smbios.baseboard.%d.productname", n), b.ProductName)
			hw.Put(fmt.Sprintf("smbios.baseboard.%d.version", n), b.Version)
			hw.Put(fmt.Sprintf("smbios.baseboard.%d.serialnumber", n), b.SerialNumber)
			hw.Put(fmt.Sprintf("smbios.baseboard.%d.boardtype", n), b.BoardType)
		}

		// type 3
		for n, s := range bios.SystemEnclosures {
			hw.Put(fmt.Sprintf("smbios.chassis.%d.manufactor", n), s.Manufacturer)
			hw.Put(fmt.Sprintf("smbios.chassis.%d.type", n), s.Type)
			hw.Put(fmt.Sprintf("smbios.chassis.%d.version", n), s.Version)
			hw.Put(fmt.Sprintf("smbios.chassis.%d.serialnumber", n), s.SerialNumber)
		}

		// type 4
		for n, p := range bios.ProcessorInformations {
			hw.Put(fmt.Sprintf("smbios.processor.%d.version", n), p.ProcessorVersion)
			hw.Put(fmt.Sprintf("smbios.processor.%d.cores", n), p.CoreCount)
			hw.Put(fmt.Sprintf("smbios.processor.%d.threads", n), p.ThreadCount)
		}

		// type 17
		for n, m := range bios.MemoryDevices {
			hw.Put(fmt.Sprintf("smbios.memory.%d.manufacturer", n), m.Manufacturer)
			hw.Put(fmt.Sprintf("smbios.memory.%d.size", n), m.Size)
			hw.Put(fmt.Sprintf("smbios.memory.%d.serialnumber", n), m.SerialNumber)
			hw.Put(fmt.Sprintf("smbios.memory.%d.speed", n), m.Speed)
		}

		if event.Fields == nil {
			event.Fields = hw
		} else {
			event.Fields.DeepUpdate(hw)
		}
		return nil
	})}, nil
}

func init() {
	monitors.RegisterActive("smbios", create)
}
