package hw

import (
	"fmt"

	"github.com/lijingwei9060/go-smbios/smbios"
	"github.com/lijingwei9060/scout/lib/common"
)

// FetchHWInfo 通过smbios获取服务器硬件信息
func FetchHWInfo() common.MapStr {
	hw := common.MapStr{}
	bios := smbios.GetSMBIOS()

	if bios == nil {
		return hw
	}

	hw.Put("smbios.major", bios.Major)
	hw.Put("smbios.minor", bios.Minor)

	// type 0
	hw.Put("bios.vendor", bios.BIOSInformation.Vendor)
	hw.Put("bios.version", bios.BIOSInformation.BIOSVersion)
	hw.Put("bios.releasedate", bios.BIOSInformation.BIOSReleaseDate)
	// type 1
	hw.Put("system.manufacturer", bios.SystemInformation.Manufacturer)
	hw.Put("system.productname", bios.SystemInformation.ProductName)
	hw.Put("system.serialnumber", bios.SystemInformation.SerialNumber)
	hw.Put("system.version", bios.SystemInformation.Version)
	hw.Put("system.uuid", bios.SystemInformation.UUID)
	// type 2
	for n, b := range bios.BaseboardInformations {
		hw.Put(fmt.Sprintf("baseboard.%d.manufacturer", n), b.Manufacturer)
		hw.Put(fmt.Sprintf("baseboard.%d.productname", n), b.ProductName)
		hw.Put(fmt.Sprintf("baseboard.%d.version", n), b.Version)
		hw.Put(fmt.Sprintf("baseboard.%d.serialnumber", n), b.SerialNumber)
		hw.Put(fmt.Sprintf("baseboard.%d.boardtype", n), b.BoardType)
	}

	// type 3
	for n, s := range bios.SystemEnclosures {
		hw.Put(fmt.Sprintf("chassis.%d.manufactor", n), s.Manufacturer)
		hw.Put(fmt.Sprintf("chassis.%d.type", n), s.Type)
		hw.Put(fmt.Sprintf("chassis.%d.version", n), s.Version)
		hw.Put(fmt.Sprintf("chassis.%d.serialnumber", n), s.SerialNumber)
	}

	// type 4
	for n, p := range bios.ProcessorInformations {
		hw.Put(fmt.Sprintf("processor.%d.version", n), p.ProcessorVersion)
		hw.Put(fmt.Sprintf("processor.%d.cores", n), p.CoreCount)
		hw.Put(fmt.Sprintf("processor.%d.threads", n), p.ThreadCount)
	}

	// type 17
	for n, m := range bios.MemoryDevices {
		hw.Put(fmt.Sprintf("memory.%d.manufacturer", n), m.Manufacturer)
		hw.Put(fmt.Sprintf("memory.%d.size", n), m.Size)
		hw.Put(fmt.Sprintf("memory.%d.serialnumber", n), m.SerialNumber)
		hw.Put(fmt.Sprintf("memory.%d.speed", n), m.Speed)
	}
	return hw
}
