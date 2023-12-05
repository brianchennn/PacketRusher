/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package context

import (
	"errors"
	"net"

	"github.com/free5gc/openapi/models"
)

type SessionContext struct {
	dataNetworks    []DataNetwork
	sessionsRules   []*models.SessionRule
	lastAllocatedIP net.IP
	n3              net.IP
}

type DataNetwork struct {
	Dnn string
	Dns DNS
}

type DNS struct {
	IPv4Addr net.IP
	IPv6Addr net.IP
}

func (s *SessionContext) NewSessionContext() {
	//TODO parametrize session data
	s.dataNetworks = []DataNetwork{
		{
			Dnn: "internet",
			Dns: DNS{
				IPv4Addr: net.ParseIP("8.8.8.8"),
				IPv6Addr: net.ParseIP("2001:4860:4860::8888"),
			},
		},
	}

	s.sessionsRules = []*models.SessionRule{{
		AuthSessAmbr: &models.Ambr{
			Uplink:   "1 Gbps",
			Downlink: "1 Gbps",
		},
		AuthDefQos: &models.AuthorizedDefaultQos{
			Var5qi: 6,
			Arp: &models.Arp{
				PriorityLevel: 8,
			},
			PriorityLevel: 8,
		},
		SessRuleId: "SessRuleId-1",
	}}

	s.lastAllocatedIP = net.ParseIP("10.0.0.1")
	s.n3 = net.ParseIP("127.0.0.1")
}

func (s *SessionContext) GetN3() net.IP {
	return s.n3
}

func (s *SessionContext) GetDnnList() []string {
	dnn := []string{}
	for _, dn := range s.dataNetworks {
		dnn = append(dnn, dn.Dnn)
	}
	return dnn
}

func (s *SessionContext) GetSessionRules() []*models.SessionRule {
	return s.sessionsRules
}

func (s *SessionContext) GetDataNetwork(dnn string) (DataNetwork, error) {
	for _, dn := range s.dataNetworks {
		if dn.Dnn == dnn {
			return dn, nil
		}
	}
	return DataNetwork{}, errors.New("[5GC] Could not find requested datanetwork")
}
