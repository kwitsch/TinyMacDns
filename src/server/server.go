package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/kwitsch/TinyMacDns/cache"
	"github.com/miekg/dns"
)

type Server struct {
	udp   *dns.Server
	cache *cache.Cache
	ttl   uint32
}

func New(cache *cache.Cache, ttl int) *Server {
	res := &Server{
		udp:   createUDPServer(),
		cache: cache,
		ttl:   uint32(ttl),
	}

	handler := res.udp.Handler.(*dns.ServeMux)
	handler.HandleFunc(".", res.OnRequest)

	return res
}

func (s *Server) Start() {
	go func() {
		err := s.udp.ListenAndServe()
		if err != nil {
			fmt.Println("Server.Start error", err)
		}
	}()
}

func (s *Server) Stop() error {
	return s.udp.Shutdown()
}

func createUDPServer() *dns.Server {
	return &dns.Server{
		Addr:    "53",
		Net:     "udp",
		Handler: dns.NewServeMux(),
		NotifyStartedFunc: func() {
			fmt.Println("UDP server is up and running")
		},
		UDPSize: 65535,
	}
}

const rdnsSuf string = ".in-addr.arpa"

func (s *Server) OnRequest(w dns.ResponseWriter, request *dns.Msg) {
	q := request.Question[0]

	m := new(dns.Msg)
	m.SetReply(request)

	if q.Qtype == dns.TypePTR || q.Qtype == dns.TypeA {
		cname := strings.TrimSuffix(strings.ToLower(q.Name), ".")
		exists := false
		val := ""

		if q.Qtype == dns.TypePTR {
			crname := strings.TrimSuffix(cname, rdnsSuf)

			val, exists = s.cache.GetHostname(crname)
			if exists {
				rr := new(dns.PTR)
				rr.Hdr = dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypePTR,
					Class:  dns.ClassINET,
					Ttl:    s.ttl,
				}

				rr.Ptr = val

				m.Answer = []dns.RR{rr}
			}
		} else if q.Qtype == dns.TypeA {
			val, exists = s.cache.GetIp(cname)
			if exists {
				rr := new(dns.A)
				rr.Hdr = dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    s.ttl,
				}

				rr.A = net.ParseIP(val)

				m.Answer = []dns.RR{rr}
			}
		}

		if !exists {
			m.SetRcode(request, dns.RcodeNameError)
		}
	} else {
		m.SetRcode(request, dns.RcodeRefused)
	}
	w.WriteMsg(m)
}
