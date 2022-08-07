package network

import "encoding/xml"

type Conntrack struct {
	XMLName xml.Name `xml:"conntrack"`
	Text    string   `xml:",chardata"`
	Flow    Flow     `xml:"flow"`
}
type Flow struct {
	Text string `xml:",chardata"`
	Meta []Meta `xml:"meta"`
}
type Meta struct {
	Text      string `xml:",chardata"`
	Direction string `xml:"direction,attr"`
	Layer3    Layer3 `xml:"layer3"`
	Layer4    Layer4 `xml:"layer4"`
	State     string `xml:"state"`
	Timeout   string `xml:"timeout"`
	Mark      string `xml:"mark"`
	Use       string `xml:"use"`
	ID        string `xml:"id"`
	Assured   string `xml:"assured"`
}
type Layer3 struct {
	Text      string `xml:",chardata"`
	Protonum  string `xml:"protonum,attr"`
	Protoname string `xml:"protoname,attr"`
	Src       string `xml:"src"`
	Dst       string `xml:"dst"`
}

type Layer4 struct {
	Text      string `xml:",chardata"`
	Protonum  string `xml:"protonum,attr"`
	Protoname string `xml:"protoname,attr"`
	Sport     string `xml:"sport"`
	Dport     string `xml:"dport"`
}
