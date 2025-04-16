package database

import (
	"encoding/xml"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// struct representation of the XML Data.

// The toplevel database. One per file.
type Database struct {
	Name     string       `xml:"name,attr"`
	Revision string       `xml:"revision,attr"`
	Mem      []MemorySlot `xml:"mem"`
}

// Each memory slot has track1,master,rhythm that are just k/v pairs.

// Memory slot, one per <mem> tag.
type MemorySlot struct {
	XmlId  string   `xml:"id,attr"`
	Name   SlotName `xml:"NAME"`
	Master Master   `xml:"MASTER"`
	Track1 Track1   `xml:"TRACK1"`
	Rhythm Rhythm   `xml:"RHYTHM"`
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

// Contents of the <MASTER> tag.
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

type Track1 struct {
	Rev      int `xml:"Rev"`
	PlyLvl   int `xml:"PlyLvl"`
	Pan      int `xml:"Pan"`
	One      int `xml:"One"`
	StrtMod  int `xml:"StrtMod"`
	StpMod   int `xml:"StpMod"`
	Measure  int `xml:"Measure"`
	MeasMod  int `xml:"MeasMod"`
	MeasLen  int `xml:"MeasLen"`
	MeasBtLp int `xml:"MeasBtLp"`
	RecTmp   int `xml:"RecTmp"`
	WavStat  int `xml:"WavStat"`
	WavLen   int `xml:"WavLen"`
}

type Rhythm struct {
	Level           int `xml:"Level"`
	Reverb          int `xml:"Reverb"`
	Pattern         int `xml:"Pattern"`
	Variation       int `xml:"Variation"`
	VariationChange int `xml:"VariationChange"`
	Kit             int `xml:"Kit"`
	Beat            int `xml:"Beat"`
	Fill            int `xml:"Fill"`
	Part1           int `xml:"Part1"`
	Part2           int `xml:"Part2"`
	Part3           int `xml:"Part3"`
	Part4           int `xml:"Part4"`
	RecCount        int `xml:"RecCount"`
	PlayCount       int `xml:"PlayCount"`
	Start           int `xml:"Start"`
	Stop            int `xml:"Stop"`
	ToneLow         int `xml:"ToneLow"`
	ToneHigh        int `xml:"ToneHigh"`
	State           int `xml:"State"`
}

// Utility Functions for the above structs.

// Get a memory slot by its number (zero-padded two-char string)
func (d *Database) GetMemorySlotByNumber(num string) *MemorySlot {
	for i := range d.Mem {
		if d.Mem[i].Number() == num {
			return &d.Mem[i]
		}
	}
	return nil
}

func (d Database) WriteXML(filename string, xml string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(xml)
	if err != nil {
		return err
	}
	return nil
}

func (d Database) ToXML() (string, error) {
	ret := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"
	ret += "<database name=\"RC-5\" revision=\"0\">\n"

	for _, mem := range d.Mem {
		ret += fmt.Sprintf("<mem id=\"%v\">\n", mem.XmlId)
		ret += fmt.Sprintf("%v\n", mem.Name.XmlString())
		ret += fmt.Sprintf("%v\n", mem.Track1.XmlString())
		ret += fmt.Sprintf("%v\n", mem.Master.XmlString())
		ret += fmt.Sprintf("%v\n", mem.Rhythm.XmlString())
		ret += "</mem>\n"
	}
	ret += "</database>"
	return ret, nil
}

// The XML string representation of the master struct.
func (m Master) XmlString() string {
	ret := "<MASTER>\n"
	value := reflect.ValueOf(m)
	for i := 0; i < value.NumField(); i++ {
		ret += fmt.Sprintf("\t<%v>%v</%v>\n", value.Type().Field(i).Name, value.Field(i).Int(), value.Type().Field(i).Name)
	}
	ret += "</MASTER>"
	return ret
}

func (t Track1) XmlString() string {
	ret := "<TRACK1>\n"
	value := reflect.ValueOf(t)
	for i := 0; i < value.NumField(); i++ {
		ret += fmt.Sprintf("\t<%v>%v</%v>\n", value.Type().Field(i).Name, value.Field(i).Int(), value.Type().Field(i).Name)
	}
	ret += "</TRACK1>"
	return ret
}

func (r Rhythm) XmlString() string {
	ret := "<RHYTHM>\n"
	value := reflect.ValueOf(r)
	for i := 0; i < value.NumField(); i++ {
		ret += fmt.Sprintf("\t<%v>%v</%v>\n", value.Type().Field(i).Name, value.Field(i).Int(), value.Type().Field(i).Name)
	}
	ret += "</RHYTHM>"
	return ret
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
	d, err := DescribeMemoryData(m.Master)
	if err != nil {
		return fmt.Sprintf("Error describing memory data: %v", err)
	}
	ret += d
	d, err = DescribeMemoryData(m.Track1)
	if err != nil {
		return fmt.Sprintf("Error describing memory data: %v", err)
	}
	ret += d
	d, err = DescribeMemoryData(m.Rhythm)
	if err != nil {
		return fmt.Sprintf("Error describing memory data: %v", err)
	}
	ret += d

	return ret
}

func (m MemorySlot) GetAttributeByName(name string) (int, error) {
	// The name is the struct name, then the name of the field in it (since some fields appear in multiple structs).
	// e.g. "Master.Tempo" or "Rhythm.Level"
	fragments := strings.Split(name, ".")
	if len(fragments) != 2 {
		return 0, fmt.Errorf("invalid attribute name: %v", name)
	}
	attr_name := fragments[1]
	var value reflect.Value
	switch fragments[0] {
	case "Master":
		value = reflect.ValueOf(m.Master)
	case "Track1":
		value = reflect.ValueOf(m.Track1)
	case "Rhythm":
		value = reflect.ValueOf(m.Rhythm)
	default:
		return 0, fmt.Errorf("invalid top-level sttribute (must be [Master|Track1|Rhythm]): %v", fragments[0])
	}

	for i := 0; i < value.NumField(); i++ {
		if value.Type().Field(i).Name == attr_name {
			return int(value.Field(i).Int()), nil
		}
	}

	return 0, fmt.Errorf("attribute not found: %v", name)
}

func (m *MemorySlot) GetReflectedElementByName(name string) (string, reflect.Value, error) {
	// Returns the reflected value of the attribute, and the string name of the attribute.
	// i.e. ("Track1", <reflect.Value>, nil)

	// The name is the struct name, then the name of the field in it (since some fields appear in multiple structs).
	// e.g. "Master.Tempo" or "Rhythm.Level"
	fragments := strings.Split(name, ".")
	if len(fragments) != 2 {
		return "", reflect.Value{}, fmt.Errorf("invalid attribute name: %v (must be [Master|Track1|Rhythm].Attribute (or Name))", name)
	}
	switch fragments[0] {
	case "Master":
		return fragments[0], reflect.ValueOf(&m.Master).Elem(), nil
	case "Track1":
		return fragments[0], reflect.ValueOf(&m.Track1).Elem(), nil
	case "Rhythm":
		return fragments[0], reflect.ValueOf(&m.Rhythm).Elem(), nil
	default:
		return "", reflect.Value{}, fmt.Errorf("invalid top-level attribute (must be [Master|Track1|Rhythm]): %v", fragments[0])
	}
}

func (m *MemorySlot) SetAttributeByName(name string, new string) error {
	valueInt, err := strconv.ParseInt(new, 10, 64)
	if err != nil {
		return fmt.Errorf("value must be an integer: %v", new)
	}
	_, value, err := m.GetReflectedElementByName(name)
	if err != nil {
		return err
	}

	fragments := strings.Split(name, ".")

	for i := 0; i < value.NumField(); i++ {
		if value.Type().Field(i).Name == fragments[1] {
			value.Field(i).SetInt(valueInt)
			return nil
		}
	}
	return fmt.Errorf("attribute not found: %v", name)
}

func (m *MemorySlot) SetNameFromString(str string) error {
	new, err := SlotNameFromString(str)
	if err != nil {
		return err
	}
	m.Name = *new
	return nil
}

func DescribeMemoryData(d interface{}) (string, error) {
	var ret string
	value := reflect.ValueOf(d)
	if value.Type().Kind() != reflect.Struct {
		return "", fmt.Errorf("expected struct, got %v", value.Type().Kind())
	}
	for i := 0; i < value.NumField(); i++ {
		ret += fmt.Sprintf("%v: %v\n", value.Type().Field(i).Name, value.Field(i).Interface())
	}
	return ret, nil
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

func (s SlotName) XmlString() string {
	// The looper uses \v as an indenter, because of course.
	var ret string
	ret += "<NAME>\n"
	ret += fmt.Sprintf("\t<C01>%v</C01>\n", s.C01)
	ret += fmt.Sprintf("\t<C02>%v</C02>\n", s.C02)
	ret += fmt.Sprintf("\t<C03>%v</C03>\n", s.C03)
	ret += fmt.Sprintf("\t<C04>%v</C04>\n", s.C04)
	ret += fmt.Sprintf("\t<C05>%v</C05>\n", s.C05)
	ret += fmt.Sprintf("\t<C06>%v</C06>\n", s.C06)
	ret += fmt.Sprintf("\t<C07>%v</C07>\n", s.C07)
	ret += fmt.Sprintf("\t<C08>%v</C08>\n", s.C08)
	ret += fmt.Sprintf("\t<C09>%v</C09>\n", s.C09)
	ret += fmt.Sprintf("\t<C10>%v</C10>\n", s.C10)
	ret += fmt.Sprintf("\t<C11>%v</C11>\n", s.C11)
	ret += fmt.Sprintf("\t<C12>%v</C12>\n", s.C12)
	ret += "</NAME>"
	return ret
}

func SlotNameFromString(str string) (*SlotName, error) {
	ret := SlotName{}
	if len(str) > 12 {
		return nil, fmt.Errorf("memory slot name must be <12 characters long")
	}
	// Pad shorter stirngs to 12.
	if len(str) < 12 {
		str = fmt.Sprintf("%-12s", str)
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
