package hw

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/lijingwei9060/go-smbios/smbios"
)

func getBIOSInformation() {
	rc, ep, err := smbios.Stream()
	if err != nil {
		log.Fatalf("failed to open stream: %v", err)
	}
	// Be sure to close the stream!
	defer rc.Close()

	// Decode SMBIOS structures from the stream.
	d := smbios.NewDecoder(rc)
	ss, err := d.Decode()
	if err != nil {
		log.Fatalf("failed to decode structures: %v", err)
	}

	major, minor, rev := ep.Version()
	fmt.Printf("SMBIOS %d.%d.%d\n", major, minor, rev)

	for _, s := range ss {
		// Only look at memory devices.
		// if s.Header.Type == 17 {
		// 	out, err := smbios.ParseMemoryDevice(s)
		// 	if err != nil {
		// 		fmt.Print(err)
		// 	}
		// 	str, _ := json.Marshal(out)
		// 	fmt.Print(string(str))
		// }

		// Code based on: https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.1.1.pdf.

		if s.Header.Type == 0 {
			out, err := smbios.ParseBIOSInformation(s)
			if err != nil {
				fmt.Print(err)
			}
			str, _ := json.Marshal(out)
			fmt.Print(string(str))
		}
		if s.Header.Type == 1 {
			out, err := smbios.ParseSystemInformation(s)
			if err != nil {
				fmt.Print(err)
			}
			str, _ := json.Marshal(out)
			fmt.Print(string(str))
		}
		if s.Header.Type == 2 {
			out, err := smbios.ParseBaseboardInformation(s)
			if err != nil {
				fmt.Print(err)
			}
			str, _ := json.Marshal(out)
			fmt.Print(string(str))
		}
		if s.Header.Type == 3 {
			out, err := smbios.ParseSystemEnclosure(s)
			if err != nil {
				fmt.Print(err)
			}
			str, _ := json.Marshal(out)
			fmt.Print(string(str))
		}
		if s.Header.Type == 4 {
			out, err := smbios.ParseProcessorInformation(s)
			if err != nil {
				fmt.Print(err)
			}
			str, _ := json.Marshal(out)
			fmt.Print(string(str))
		}

	}
}
