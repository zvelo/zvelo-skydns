// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package main

import (
	"log"
	"strings"
	"net/url"

	"github.com/coreos/go-etcd/etcd"
)

func newClient() (client *etcd.Client) {
	// set default if not specified in env
	if len(machines) == 1 && machines[0] == "" {
		machines[0] = "http://127.0.0.1:4001"

	}
	// override if we have a commandline flag as well
	if machine != "" {
		machines = strings.Split(machine, ",")
	}
	var err error
	if strings.HasPrefix(machines[0], "https://") {
		// TODO(miek): this probably does not work.
		if client, err = etcd.NewTLSClient(machines, tlspem, tlskey, ""); err != nil {
			log.Fatal(err)
		}
	}
	client = etcd.NewClient(machines)
	client.SyncCluster()
	return client
}

// updateClient updates the client with the machines found in v2/_etcd/machines.
func (s *server) updateClient() {
	resp, err := s.client.Get("/_etcd/machines/", false, true)
	if err != nil {
		s.config.log.Infof("could not get new etcd machines, keeping old: %s", err.Error())
		return
	}
	machine := make([]string, 0)
	for _, m := range resp.Node.Nodes {
		u, e := url.Parse(m.Value)
		if e != nil {
			continue
		}
		// etcd=bla&raft=bliep
		// TODO(miek): surely there is a better way to do this
		ms := strings.Split(u.String(), "&")
		if len(ms) == 0 {
			continue
		}
		if len(ms[0]) < 5 {
			continue
		}
		machine = append(machine, ms[0][5:])
	}
	s.config.log.Infof("setting new etcd cluster to %v", machines)
	s.Lock()
	s.client.SetCluster(machines) // TODO(miek): return value
	s.client.SyncCluster()
	s.Unlock()
}
