/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package ngap

import (
	contexxt "context"
	"fmt"
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"time"

	"github.com/ishidawataru/sctp"
	log "github.com/sirupsen/logrus"
)

var ConnCount int

func DialSCTPExtWithTimeout(network, local, remote string, options sctp.InitMsg, timeout time.Duration) (*sctp.SCTPConn, error) {
	ctx, cancel := contexxt.WithTimeout(contexxt.Background(), timeout)
	defer cancel()

	resultChan := make(chan struct {
		conn *sctp.SCTPConn
		err  error
	}, 1)

	go func() {
		laddr, err := sctp.ResolveSCTPAddr(network, local)
		if err != nil {
			resultChan <- struct {
				conn *sctp.SCTPConn
				err  error
			}{nil, err}
			return
		}
		raddr, err := sctp.ResolveSCTPAddr(network, remote)
		if err != nil {
			resultChan <- struct {
				conn *sctp.SCTPConn
				err  error
			}{nil, err}
			return
		}

		conn, err := sctp.DialSCTPExt(network, laddr, raddr, sctp.InitMsg{NumOstreams: 2, MaxInstreams: 2})
		resultChan <- struct {
			conn *sctp.SCTPConn
			err  error
		}{conn, err}
	}()

	select {
	case result := <-resultChan:
		return result.conn, result.err
	case <-ctx.Done():
		return nil, fmt.Errorf("dial SCTP timed out")
	}
}

func InitConn(amf *context.GNBAmf, gnb *context.GNBContext) error {

	// check AMF IP and AMF port.
	remote := fmt.Sprintf("%s:%d", amf.GetAmfIp(), amf.GetAmfPort())
	local := fmt.Sprintf("%s:%d", gnb.GetGnbIp(), gnb.GetGnbPort()+ConnCount)
	ConnCount++

	// streams := amf.GetTNLAStreams()
	timeout := 200 * time.Millisecond

	conn, err := DialSCTPExtWithTimeout(
		"sctp",
		local,
		remote,
		sctp.InitMsg{NumOstreams: 2, MaxInstreams: 2},
		timeout,
	)

	if err != nil {
		amf.SetSCTPConn(nil)
		return err
	}

	// set streams and other information about TNLA

	// successful established SCTP (TNLA - N2)
	amf.SetSCTPConn(conn)
	gnb.SetN2(conn)

	conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)

	go GnbListen(amf, gnb)

	return nil
}

func GnbListen(amf *context.GNBAmf, gnb *context.GNBContext) {

	buf := make([]byte, 65535)
	conn := amf.GetSCTPConn()

	/*
		defer func() {
			err := conn.Close()
			if err != nil {
				log.Info("[GNB][SCTP] Error in closing SCTP association for %d AMF\n", amf.GetAmfId())
			}
		}()
	*/

	for {

		n, info, err := conn.SCTPRead(buf[:])
		if err != nil {
			break
		}

		log.Info("[GNB][SCTP] Receive message in ", info.Stream, " stream\n")

		forwardData := make([]byte, n)
		copy(forwardData, buf[:n])

		// handling NGAP message.
		go Dispatch(amf, gnb, forwardData)

	}

}
