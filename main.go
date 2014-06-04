// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"
)

var (
	machines    = strings.Split(os.Getenv("ETCD_MACHINES"), ",")      // list of URLs to etcd
	nameservers = strings.Split(os.Getenv("SKYDNS_NAMESERVERS"), ",") // list of nameservers
	tlskey      = os.Getenv("ETCD_TLSKEY")                            // TLS private key path
	tlspem      = os.Getenv("ETCD_TLSPEM")                            // X509 certificate
	config      = &Config{ReadTimeout: 0, Domain: "", DnsAddr: "", DNSSEC: ""}
	nameserver  = ""
	machine     = ""
	discover    = false
)

func init() {
	flag.StringVar(&config.Domain, "domain",
		func() string {
			if x := os.Getenv("SKYDNS_DOMAIN"); x != "" {
				return x
			}
			return "skydns.local."
		}(), "domain to anchor requests to (SKYDNS_DOMAIN)")
	flag.StringVar(&config.DnsAddr, "addr",
		func() string {
			if x := os.Getenv("SKYDNS_ADDR"); x != "" {
				return x
			}
			return "127.0.0.1:53"
		}(), "ip:port to bind to (SKYDNS_ADDR)")

	flag.StringVar(&nameserver, "nameserver", "", "nameserver address(es) to forward (non-local) queries to e.g. 8.8.8.8:53,8.8.4.4:53")
	flag.StringVar(&machine, "machines", "", "machine address(es) running etcd")
	flag.StringVar(&config.DNSSEC, "dnssec", "", "basename of DNSSEC key file e.q. Kskydns.local.+005+38250")
	flag.StringVar(&tlskey, "tls-key", "", "TLS Private Key path")
	flag.StringVar(&tlspem, "tls-pem", "", "X509 Certificate")
	flag.DurationVar(&config.ReadTimeout, "rtimeout", 2*time.Second, "read timeout")
	flag.BoolVar(&config.RoundRobin, "round-robin", true, "round robin A/AAAA replies")
	flag.BoolVar(&discover, "discover", false, "discover new machines running etcd by querying /v2/machines on startup")
	// TTl
	// Minttl
	flag.StringVar(&config.Hostmaster, "hostmaster", "hostmaster@skydns.local.", "hostmaster email address to use")
}

func main() {
	flag.Parse()
	client := newClient()
	if nameserver != "" {
		config.Nameservers = strings.Split(nameserver, ",")
	}
	config, err := LoadConfig(client, config)
	if err != nil {
		log.Fatal(err)
	}
	s := New(config)
	s.client = client
	if discover {
		s.updateClient()
	}
	statsCollect()
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
