// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"
	"strings"

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

// updateClient updates the client with the machines found
// in v2/machines.
func (s *server) updateClient() {
	machines := make([]string, 0)
	n, err := s.client.Get("v2/machines", false, false)
	if err != nil {
		s.config.log.Info("could not read /machines from etcd, keeping old: ", err)
		return
	}
	if err := json.Unmarshal([]byte(n.Node.Value), &machines); err != nil {
		s.config.log.Infof("failed to parse json: %s", err.Error())
		return
	}
	var client *etcd.Client
	if strings.HasPrefix(machines[0], "https://") {
		// First one is https, assume they all have.
		if client, err = etcd.NewTLSClient(machines, tlspem, tlskey, ""); err != nil {
			s.config.log.Info("could not connect to new etcd machines, keeping old: ", err)
			return
		}
	}
	client = etcd.NewClient(machines)
	client.SyncCluster()
	s.client = client
}
