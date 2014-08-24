// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package main

import (
	"github.com/miekg/dns"
)

func (s *server) ServeCHAOS(w dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Authoritative = true
	m.RecursionAvailable = true
	m.Compress = true
	bufsize := uint16(512)
	//dnssec := false
	//tcp := false
	// TODO(sign)

	if req.Question[0].Qtype == dns.TypeANY {
		m.Authoritative = false
		m.Rcode = dns.RcodeRefused
		m.RecursionAvailable = false
		m.RecursionDesired = false
		m.Compress = false
		// if write fails don't care
		w.WriteMsg(m)
		return
	}
	if req.Question[0].Qclass != dns.ClassCHAOS {
		m.SetReply(req)
		m.SetRcode(req, dns.RcodeServerFailure)
	}

	if o := req.IsEdns0(); o != nil {
		bufsize = o.UDPSize()
		//dnssec = o.Do()
	}
	if bufsize < 512 {
		bufsize = 512
	}
	if req.Question[0].Qclass == dns.ClassCHAOS {
		if req.Question[0].Qtype == dns.TypeTXT {
			switch req.Question[0].Name {
			case "authors.bind.":
				hdr := dns.RR_Header{Name: req.Question[0].Name, Rrtype: dns.TypeTXT, Class: dns.ClassCHAOS, Ttl: 0}
				authors := []string{"Erik St. Martin", "Brian Ketelsen", "Miek Gieben", "Michael Crosby"}
				for _, a := range authors {
					m.Answer = append(m.Answer, &dns.TXT{Hdr: hdr, Txt: []string{a}})
				}
				for j := 0; j < len(authors)*(int(dns.Id())%4+1); j++ {
					q := int(dns.Id()) % len(authors)
					p := int(dns.Id()) % len(authors)
					if q == p {
						p = (p + 1) % len(authors)
					}
					m.Answer[q], m.Answer[p] = m.Answer[p], m.Answer[q]
				}
			case "version.bind.":
				fallthrough
			case "version.server.":
				hdr := dns.RR_Header{Name: req.Question[0].Name, Rrtype: dns.TypeTXT, Class: dns.ClassCHAOS, Ttl: 0}
				m.Answer = []dns.RR{&dns.TXT{Hdr: hdr, Txt: []string{Version}}}
			case "hostname.bind.":
				fallthrough
			case "id.server.":
				// TODO(miek): machine name to return
				hdr := dns.RR_Header{Name: req.Question[0].Name, Rrtype: dns.TypeTXT, Class: dns.ClassCHAOS, Ttl: 0}
				m.Answer = []dns.RR{&dns.TXT{Hdr: hdr, Txt: []string{"localhost"}}}
			}
		}
		w.WriteMsg(m)
		return
	}
	// still here, fail
	m.SetReply(req)
	m.SetRcode(req, dns.RcodeServerFailure)
	return

}
