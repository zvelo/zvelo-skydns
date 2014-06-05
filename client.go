// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
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
	machines, err := getMachines(s.client)
	if err != nil {
		s.config.log.Infof("could not get new etcd machines, keeping old: %s", err.Error())
		return
	}
	/*
		}
		var client *etcd.Client
		if strings.HasPrefix(machines[0], "https://") {
			// First one is https, assume they all have.
			if client, err = etcd.NewTLSClient(machines, tlspem, tlskey, ""); err != nil {
				s.config.log.Infof("could not connect to new etcd machines, keeping old: %s", err.Error())
				return
			}
		}
		client = etcd.NewClient(machines)
		client.SyncCluster()
		s.Lock()
		s.client = client
		s.Unlock()
	*/
	println(machines)
}

// getMachine get a list of the machines from Etcd.
func getMachines(c *etcd.Client) ([]string, error) {
	p := "/machines"
	px := []string{}
	// Can not access the default consitency used in the client
	//if c.config.Consistency == etcd.STRONG_CONSISTENCY {
	//	options["consistent"] = true
	//}
	//
	//str, err := options.toParameters(etcd.VALID_GET_OPTIONS)
	//if err != nil {
	//	return nil, err
	//}
	//p += str

	println("HIER")
	req := etcd.NewRawRequest("GET", p, nil, nil)
	raw, err := c.SendRequest(req)
	if err != nil {
		println("RETURN")
		return nil, err
	}
	fmt.Printf("%s\n", raw)
	resp, err := raw.Unmarshal()
	if err != nil {
		println("RETURN2")
		return nil, err
	}
	if err = json.Unmarshal([]byte(resp.Node.Value), &px); err != nil {
		return nil, err
	}
	return nil, nil
}
