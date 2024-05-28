/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package gnb

import (
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	serviceNas "my5G-RANTester/internal/control_test_engine/gnb/nas/service"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/trigger"
	"my5G-RANTester/internal/monitoring"

	"os"
	"os/signal"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func InitGnb(conf config.Config, wg *sync.WaitGroup) *context.GNBContext {

	// instance new gnb.
	gnb := &context.GNBContext{}

	// new gnb context.
	gnb.NewRanGnbContext(
		conf.GNodeB.PlmnList.GnbId,
		conf.GNodeB.PlmnList.Mcc,
		conf.GNodeB.PlmnList.Mnc,
		conf.GNodeB.PlmnList.Tac,
		conf.GNodeB.SliceSupportList.Sst,
		conf.GNodeB.SliceSupportList.Sd,
		conf.GNodeB.ControlIF.Ip,
		conf.GNodeB.DataIF.Ip,
		conf.GNodeB.ControlIF.Port,
		conf.GNodeB.DataIF.Port,
	)

	amfPool := gnb.GetAmfPool()

	// start communication with AMF (server SCTP).
	for _, amfConfig := range conf.AMFs {

		amfExisted := false

		amfPool.Range(func(key, value any) bool {
			gnbAmf, ok := value.(*context.GNBAmf)
			if !ok {
				return true
			}
			if gnbAmf.GetAmfIp() == amfConfig.Ip && gnbAmf.GetAmfPort() == amfConfig.Port {
				amfExisted = true
				return false
			}
			return true
		})

		if amfExisted {
			log.Warn("TNLA existed")
			continue
		}

		// new AMF context.
		amf := gnb.NewGnBAmf(amfConfig.Ip, amfConfig.Port)

		// start communication with AMF(SCTP).
		if err := ngap.InitConn(amf, gnb); err != nil {
			log.Warn("Error in", err)
			continue

		} else {
			log.Info("[GNB] SCTP/NGAP service is running")
			// wg.Add(1)
		}

		trigger.SendNgSetupRequest(gnb, amf)
	}

	time.Sleep(time.Second)

	amfPool.Range(func(k, v any) bool {
		if amf, ok := v.(*context.GNBAmf); ok {
			if amf.GetState() == 0 {
				amfPool.Delete(k)
			}
		}

		return true
	})

	// start communication with UE (server UNIX sockets).
	serviceNas.InitServer(gnb)

	go func() {
		// control the signals
		sigGnb := make(chan os.Signal, 1)
		signal.Notify(sigGnb, os.Interrupt)

		// Block until a signal is received.
		<-sigGnb
		//gnb.Terminate()
		wg.Done()
	}()

	return gnb
}

func InitGnbForLoadSeconds(conf config.Config, wg *sync.WaitGroup,
	monitor *monitoring.Monitor) {

	// instance new gnb.
	gnb := &context.GNBContext{}

	// new gnb context.
	gnb.NewRanGnbContext(
		conf.GNodeB.PlmnList.GnbId,
		conf.GNodeB.PlmnList.Mcc,
		conf.GNodeB.PlmnList.Mnc,
		conf.GNodeB.PlmnList.Tac,
		conf.GNodeB.SliceSupportList.Sst,
		conf.GNodeB.SliceSupportList.Sd,
		conf.GNodeB.ControlIF.Ip,
		conf.GNodeB.DataIF.Ip,
		conf.GNodeB.ControlIF.Port,
		conf.GNodeB.DataIF.Port)

	// start communication with AMF (server SCTP).
	for _, amf := range conf.AMFs {
		// new AMF context.
		amf := gnb.NewGnBAmf(amf.Ip, amf.Port)

		// start communication with AMF(SCTP).
		ngap.InitConn(amf, gnb)

		trigger.SendNgSetupRequest(gnb, amf)

		// timeout is 1 second for receive NG Setup Response
		time.Sleep(1000 * time.Millisecond)

		// AMF responds message sends by Tester
		// means AMF is available
		if amf.GetState() == 0x01 {
			monitor.IncRqs()
		}
	}

	gnb.Terminate()
	wg.Done()
	// os.Exit(0)
}

func InitGnbForAvaibility(conf config.Config,
	monitor *monitoring.Monitor) {

	// instance new gnb.
	gnb := &context.GNBContext{}

	// new gnb context.
	gnb.NewRanGnbContext(
		conf.GNodeB.PlmnList.GnbId,
		conf.GNodeB.PlmnList.Mcc,
		conf.GNodeB.PlmnList.Mnc,
		conf.GNodeB.PlmnList.Tac,
		conf.GNodeB.SliceSupportList.Sst,
		conf.GNodeB.SliceSupportList.Sd,
		conf.GNodeB.ControlIF.Ip,
		conf.GNodeB.DataIF.Ip,
		conf.GNodeB.ControlIF.Port,
		conf.GNodeB.DataIF.Port)

	// start communication with AMF (server SCTP).
	for _, amf := range conf.AMFs {
		// new AMF context.
		amf := gnb.NewGnBAmf(amf.Ip, amf.Port)

		// start communication with AMF(SCTP).
		if err := ngap.InitConn(amf, gnb); err != nil {
			log.Info("Error in ", err)

			return

		} else {
			log.Info("[GNB] SCTP/NGAP service is running")

		}

		trigger.SendNgSetupRequest(gnb, amf)

		// timeout is 1 second for receive NG Setup Response
		time.Sleep(1000 * time.Millisecond)

		// AMF responds message sends by Tester
		// means AMF is available
		if amf.GetState() == 0x01 {
			monitor.IncAvaibility()

		}
	}

	gnb.Terminate()
	// os.Exit(0)
}
