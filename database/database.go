package database

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
)

// struct representation of the XML Data.

// The toplevel database. One per file.
type Database struct {
	Name     string       `xml:"name,attr"`
	Revision string       `xml:"revision,attr"`
	Mem      []MemorySlot `xml:"mem"`
}

// Get a memory slot by its number (zero-padded two-char string)
func (d Database) GetMemorySlotByNumber(num string) *MemorySlot {
	for _, mem := range d.Mem {
		if mem.Number() == num {
			return &mem
		}
	}
	return nil
}

// Memory slot, one per <mem> tag.
type MemorySlot struct {
	XmlId  string   `xml:"id,attr"`
	Name   SlotName `xml:"NAME"`
	Master Master   `xml:"MASTER"`
}

// The Memory Slot number, as a two-char padded string
func (m MemorySlot) Number() string {
	ret, err := strconv.Atoi(m.XmlId)
	if err != nil {
		return "00"
	}
	ret += 1
	return fmt.Sprintf("%02d", ret)
}

// A text descroption of a memory slot.
func (m MemorySlot) Describe() string {
	var ret string
	ret += fmt.Sprintf("Mem: %v\n", m.Number())
	ret += fmt.Sprintf("Name: %v\n", m.Name.String())

	return ret
}

// Slot names are represented in this funny way.
type SlotName struct {
	C01 int `xml:"C01"`
	C02 int `xml:"C02"`
	C03 int `xml:"C03"`
	C04 int `xml:"C04"`
	C05 int `xml:"C05"`
	C06 int `xml:"C06"`
	C07 int `xml:"C07"`
	C08 int `xml:"C08"`
	C09 int `xml:"C09"`
	C10 int `xml:"C10"`
	C11 int `xml:"C11"`
	C12 int `xml:"C12"`
}

// The string representation of the slot name.
func (s SlotName) String() string {
	var ret string
	// Ugh.
	ret += string(rune(s.C01))
	ret += string(rune(s.C02))
	ret += string(rune(s.C03))
	ret += string(rune(s.C04))
	ret += string(rune(s.C05))
	ret += string(rune(s.C06))
	ret += string(rune(s.C07))
	ret += string(rune(s.C08))
	ret += string(rune(s.C09))
	ret += string(rune(s.C10))
	ret += string(rune(s.C11))
	ret += string(rune(s.C12))
	return ret
}

func SlotNameFromString(str string) (*SlotName, error) {
	ret := SlotName{}
	if len(str) != 12 {
		return nil, fmt.Errorf("memory slot name must be 12 characters long")
	}
	// Ugh.
	ret.C01 = int(str[0])
	ret.C02 = int(str[1])
	ret.C03 = int(str[2])
	ret.C04 = int(str[3])
	ret.C05 = int(str[4])
	ret.C06 = int(str[5])
	ret.C07 = int(str[6])
	ret.C08 = int(str[7])
	ret.C09 = int(str[8])
	ret.C10 = int(str[9])
	ret.C11 = int(str[10])
	ret.C12 = int(str[11])
	return &ret, nil
}

type Master struct {
	Tempo     int `xml:"Tempo"`
	DubMode   int `xml:"DubMode"`
	RecAction int `xml:"RecAction"`
	AutoRec   int `xml:"AutoRec"`
	FadeTime  int `xml:"FadeTime"`
	Level     int `xml:"Level"`
	LpMod     int `xml:"LpMod"`
	LpLen     int `xml:"LpLen"`
	TrkMod    int `xml:"TrkMod"`
	Sync      int `xml:"Sync"`
}

func LoadMemoryFile(filename string) (*Database, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	db := Database{}
	err = xml.Unmarshal(raw, &db)
	if err != nil {
		return nil, err
	}
	return &db, nil
}
