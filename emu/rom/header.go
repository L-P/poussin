package rom

import "fmt"

type Header struct {
	Title         string
	DMG           bool // original DMG game before CGB/SGB was a thing
	CGBOnly       bool
	SGBSupport    bool
	CartridgeType byte
	ROMSize       byte
	RAMSize       byte
	Japanese      bool
	Version       byte
}

const (
	HeaderOldTitleStart = 0x0134
	HeaderOldTitleEnd   = 0x0143
	HeaderCGBFlag       = 0x0143
	HeaderSGBFlag       = 0x0146
	HeaderCartridgeType = 0x0147
	HeaderROMSize       = 0x0148
	HeaderRAMSize       = 0x0148
	HeaderDestination   = 0x014A
	HeaderVersion       = 0x014C
)

var CartridgeTypes = map[byte]string{
	0x00: "ROM ONLY",
	0x01: "MBC1",
	0x02: "MBC1+RAM",
	0x03: "MBC1+RAM+BATTERY",
	0x05: "MBC2",
	0x06: "MBC2+BATTERY",
	0x08: "ROM+RAM",
	0x09: "ROM+RAM+BATTERY",
	0x0B: "MMM01",
	0x0C: "MMM01+RAM",
	0x0D: "MMM01+RAM+BATTERY",
	0x0F: "MBC3+TIMER+BATTERY",
	0x10: "MBC3+TIMER+RAM+BATTERY",
	0x11: "MBC3",
	0x12: "MBC3+RAM",
	0x13: "MBC3+RAM+BATTERY",
	0x19: "MBC5",
	0x1A: "MBC5+RAM",
	0x1B: "MBC5+RAM+BATTERY",
	0x1C: "MBC5+RUMBLE",
	0x1D: "MBC5+RUMBLE+RAM",
	0x1E: "MBC5+RUMBLE+RAM+BATTERY",
	0x20: "MBC6",
	0x22: "MBC7+SENSOR+RUMBLE+RAM+BATTERY",
	0xFC: "POCKET CAMERA",
	0xFD: "BANDAI TAMA5",
	0xFE: "HuC3",
	0xFF: "HuC1+RAM+BATTERY",
}

func NewHeader(rom []byte) Header {
	h := Header{}

	for i := HeaderOldTitleStart; i < HeaderOldTitleEnd; i++ {
		if rom[i] == 0 {
			break
		}
		h.Title += string(rom[i])
	}

	h.DMG = true
	if rom[HeaderCGBFlag] == 0xC0 || rom[HeaderCGBFlag] == 0x80 {
		h.DMG = false
	}

	h.CGBOnly = rom[HeaderCGBFlag] == 0xC0
	h.SGBSupport = rom[HeaderSGBFlag] == 0x03
	h.CartridgeType = rom[HeaderCartridgeType]
	h.ROMSize = rom[HeaderROMSize]
	h.RAMSize = rom[HeaderRAMSize]
	h.Japanese = rom[HeaderDestination] == 0x00

	return h
}

func (h *Header) String() string {
	str := `"` + h.Title + `"`

	if h.DMG {
		str += ", DMG"
	}

	if h.Japanese {
		str += ", Japan"
	} else {
		str += ", World"
	}

	if !h.DMG && !h.CGBOnly {
		str += ", DMG+CGB"
	} else if h.CGBOnly {
		str += ", CGB"
	}

	if h.SGBSupport {
		str += "+SGB"
	}

	str += ", " + CartridgeTypes[h.CartridgeType]
	str += fmt.Sprintf(", RAM: %02X", h.ROMSize)
	str += fmt.Sprintf(", ROM: %02X", h.RAMSize)
	str += fmt.Sprintf(", Version: %02X", h.Version)

	return str
}
