package bgp

import (
	"context"
	"errors"
	"log"
	"net"

	api "github.com/osrg/gobgp/v3/api"
	"github.com/osrg/gobgp/v3/pkg/server"
	apb "google.golang.org/protobuf/types/known/anypb"
)

var (
	ErrorInvalidASN = errors.New("invalid ASN")
)

// CreateBgpServer creates a BGP server
// It needs the ASN and the IP address of the BGP server
// TODO: check that the ASN and the IP address are valid
func CreateBgpServer(asn int, ipAddress net.IP) (*server.BgpServer, error) {
	// Check for the correct ASN
	if asn < 1 || asn > 2147483647 {
		return nil, ErrorInvalidASN
	}
	s := server.NewBgpServer()
	go s.Serve()

	if err := s.StartBgp(context.Background(), &api.StartBgpRequest{
		Global: &api.Global{
			Asn:      uint32(asn),
			RouterId: ipAddress.String(),
		},
	}); err != nil {
		return nil, err
	}

	if err := s.WatchEvent(context.Background(), &api.WatchEventRequest{Peer: &api.WatchEventRequest_Peer{}}, func(r *api.WatchEventResponse) {
		if p := r.GetPeer(); p != nil && p.Type == api.WatchEventResponse_PeerEvent_STATE {
			log.Println(p)
		}
	}); err != nil {
		return nil, err
	}

	return s, nil
}

func CreateBgpPeer(s *server.BgpServer) error {
	n := &api.Peer{
		Conf: &api.PeerConf{
			NeighborAddress: "20.127.114.9",
			PeerAsn:         65002,
		},
		EbgpMultihop: &api.EbgpMultihop{Enabled: true, MultihopTtl: 250},
	}

	if err := s.AddPeer(context.Background(), &api.AddPeerRequest{
		Peer: n,
	}); err != nil {
		return err
	}

	nlri, _ := apb.New(&api.IPAddressPrefix{
		Prefix:    "10.0.0.0",
		PrefixLen: 24,
	})

	a1, _ := apb.New(&api.OriginAttribute{
		Origin: 0,
	})
	a2, _ := apb.New(&api.NextHopAttribute{
		NextHop: "10.0.0.1",
	})
	a3, _ := apb.New(&api.AsPathAttribute{
		Segments: []*api.AsSegment{
			{
				Type:    2,
				Numbers: []uint32{6762, 39919, 65000, 35753, 65000},
			},
		},
	})
	attrs := []*apb.Any{a1, a2, a3}

	_, err := s.AddPath(context.Background(), &api.AddPathRequest{
		Path: &api.Path{
			Family: &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST},
			Nlri:   nlri,
			Pattrs: attrs,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// GetPrimaryIP returns the primary IP address of the machine
// It leaves the decision to the OS to choose the primary IP address.
func GetPrimaryIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, err
}
