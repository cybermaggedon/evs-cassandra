package main

import (
	"github.com/cybermaggedon/evs-golang-api"
	pb "github.com/cybermaggedon/evs-golang-api/protos"
	"github.com/golang/protobuf/ptypes"
	"net"
	"time"
)

type DnsQuery struct {
	Name  string `cql:"name"`
	Type  string `cql:"type"`
	Class string `cql:"cls"`
}

type DnsAnswer struct {
	Name    string `cql:"name"`
	Type    string `cql:"type"`
	Class   string `cql:"cls"`
	Address net.IP `cql:"address"`
}

type HttpRequest struct {
	Method string `cql:"method"`
}

type HttpResponse struct {
	Status string `cql:"status"`
	Code   int32  `cql:"code"`
}

type Indicator struct {
	Id          string  `cql:"id"`
	Type        string  `cql:"type"`
	Value       string  `cql:"value"`
	Category    string  `cql:"category"`
	Source      string  `cql:"source"`
	Author      string  `cql:"author"`
	Description string  `cql:"description"`
	Probability float32 `cql:"probability"`
}

type Event struct {
	Id           string
	Time         time.Time
	Action       string
	Device       string
	Network      string
	SrcIp        net.IP
	DestIp       net.IP
	SrcPort      int
	DestPort     int
	Protocol     string
	Url          string
	DnsType      string
	DnsQuery     []DnsQuery
	DnsAnswer    []DnsAnswer
	HttpRequest  HttpRequest
	HttpResponse HttpResponse
	Header       map[string]string
	Indicator    []Indicator
}

func Convert(ev *pb.Event) *Event {

	ob := &Event{}

	ob.Header = map[string]string{}

	ob.Id = ev.Id
	ob.Action = ev.Action.String()
	ob.Device = ev.Device
	ob.Network = ev.Network
	tm, _ := ptypes.Timestamp(ev.Time)
	ob.Time = tm
	ob.Url = ev.Url

	switch d := ev.Detail.(type) {
	case *pb.Event_DnsMessage:

		msg := d.DnsMessage

		if len(msg.Query) > 0 {

			qs := []DnsQuery{}

			for _, val := range msg.Query {
				q := DnsQuery{
					Name:  val.Name,
					Type:  val.Type,
					Class: val.Class,
				}
				qs = append(qs, q)
			}
			ob.DnsQuery = qs
		}

		if len(msg.Answer) > 0 {
			as := []DnsAnswer{}
			for _, val := range msg.Answer {
				a := DnsAnswer{
					Name:  val.Name,
					Type:  val.Type,
					Class: val.Class,
				}
				if val.Address != nil {
					a.Address = evs.AddressToIp(
						val.Address,
					)
				}
				as = append(as, a)
			}
			ob.DnsAnswer = as
		}

	case *pb.Event_HttpRequest:
		ob.HttpRequest = HttpRequest{
			Method: d.HttpRequest.Method,
		}
		ob.Header = d.HttpRequest.Header
		break

	case *pb.Event_HttpResponse:
		ob.HttpResponse = HttpResponse{
			Status: d.HttpResponse.Status,
			Code:   d.HttpResponse.Code,
		}
		ob.Header = d.HttpResponse.Header
		break

	}

	// Indicators.
	if len(ev.Indicators) > 0 {
		is := []Indicator{}
		for _, val := range ev.Indicators {
			i := Indicator{
				Id:          val.Id,
				Type:        val.Type,
				Value:       val.Value,
				Category:    val.Category,
				Source:      val.Source,
				Author:      val.Author,
				Description: val.Description,
				Probability: val.Probability,
			}
			is = append(is, i)
		}
		ob.Indicator = is
	}

	for _, addr := range ev.Src {
		if addr.Protocol == pb.Protocol_ipv4 {
			ob.SrcIp = evs.Int32ToIp(addr.Address.GetIpv4())
		}
		if addr.Protocol == pb.Protocol_ipv6 {
			ob.SrcIp = evs.BytesToIp(addr.Address.GetIpv6())
		}
		if addr.Protocol == pb.Protocol_tcp {
			ob.SrcPort = int(addr.Address.GetPort())
			ob.Protocol = "tcp"
		}
		if addr.Protocol == pb.Protocol_udp {
			ob.SrcPort = int(addr.Address.GetPort())
			ob.Protocol = "udp"
		}
	}

	for _, addr := range ev.Dest {
		if addr.Protocol == pb.Protocol_ipv4 {
			ob.DestIp = evs.Int32ToIp(addr.Address.GetIpv4())
		}
		if addr.Protocol == pb.Protocol_ipv6 {
			ob.DestIp = evs.BytesToIp(addr.Address.GetIpv6())
		}
		if addr.Protocol == pb.Protocol_tcp {
			ob.DestPort = int(addr.Address.GetPort())
		}
		if addr.Protocol == pb.Protocol_udp {
			ob.DestPort = int(addr.Address.GetPort())
		}
	}

	return ob

}
