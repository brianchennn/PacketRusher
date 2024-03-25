/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package context

import (
	"fmt"

	"github.com/free5gc/aper"
	"github.com/free5gc/openapi/models"
	"github.com/ishidawataru/sctp"
	log "github.com/sirupsen/logrus"
)

// AMF main states in the GNB Context.
const Inactive = 0x00
const Active = 0x01
const Overload = 0x02

type GNBAmf struct {
	amfIp               string         // AMF ip
	amfPort             int            // AMF port
	amfId               int64          // AMF id
	tnlaList                []*TNLAssociation // AMF sctp associations
	relativeAmfCapacity int64          // AMF capacity
	state               int
	name                string // amf name.
	regionId            byte
	setId               byte
	pointer             byte
	plmns               *PlmnSupported
	slices              *SliceSupported
	lenSlice            int
	lenPlmn             int
	// TODO implement the other fields of the AMF Context
}

type TNLAssociation struct {
	sctpConn         *sctp.SCTPConn
	amfName			 string
	backupAmfName	 string
	Guami			 *models.Guami
	tnlaWeightFactor int64
	usage            aper.Enumerated
	streams          uint16
}

type SliceSupported struct {
	sst    string
	sd     string
	status string
	next   *SliceSupported
}

type PlmnSupported struct {
	mcc  string
	mnc  string
	next *PlmnSupported
}

func (amf *GNBAmf) GetSliceSupport(index int) (string, string) {

	mov := amf.slices
	for i := 0; i < index; i++ {
		mov = mov.next
	}

	return mov.sst, mov.sd
}

func (amf *GNBAmf) GetPlmnSupport(index int) (string, string) {

	mov := amf.plmns
	for i := 0; i < index; i++ {
		mov = mov.next
	}

	return mov.mcc, mov.mnc
}

func convertMccMnc(plmn string) (mcc string, mnc string) {
	if plmn[2] == 'f' {
		mcc = fmt.Sprintf("%c%c%c", plmn[3], plmn[0], plmn[1])
		mnc = fmt.Sprintf("%c%c", plmn[5], plmn[4])
	} else {
		mcc = fmt.Sprintf("%c%c%c", plmn[3], plmn[0], plmn[1])
		mnc = fmt.Sprintf("%c%c%c", plmn[2], plmn[5], plmn[4])
	}

	return mcc, mnc
}

func (amf *GNBAmf) AddedPlmn(plmn string) {

	if amf.lenPlmn == 0 {
		newElem := &PlmnSupported{}

		// newElem.info = plmn
		newElem.next = nil
		newElem.mcc, newElem.mnc = convertMccMnc(plmn)
		// update list
		amf.plmns = newElem
		amf.lenPlmn++
		return
	}

	mov := amf.plmns
	for i := 0; i < amf.lenPlmn; i++ {

		// end of the list
		if mov.next == nil {

			newElem := &PlmnSupported{}
			newElem.mcc, newElem.mnc = convertMccMnc(plmn)
			newElem.next = nil

			mov.next = newElem

		} else {
			mov = mov.next
		}
	}

	amf.lenPlmn++
}

func (amf *GNBAmf) AddedSlice(sst string, sd string) {

	if amf.lenSlice == 0 {
		newElem := &SliceSupported{}
		newElem.sst = sst
		newElem.sd = sd
		newElem.next = nil

		// update list
		amf.slices = newElem
		amf.lenSlice++
		return
	}

	mov := amf.slices
	for i := 0; i < amf.lenSlice; i++ {

		// end of the list
		if mov.next == nil {

			newElem := &SliceSupported{}
			newElem.sst = sst
			newElem.sd = sd
			newElem.next = nil

			mov.next = newElem

		} else {
			mov = mov.next
		}
	}
	amf.lenSlice++
}

func (amf *GNBAmf) GetTNLAs() []*TNLAssociation {
	return amf.tnlaList
}

func (amf *GNBAmf) SetStateInactive() {
	amf.state = Inactive
}

func (amf *GNBAmf) SetStateActive() {
	amf.state = Active
}

func (amf *GNBAmf) SetStateOverload() {
	amf.state = Overload
}

func (amf *GNBAmf) GetState() int {
	return amf.state
}

func (amf *GNBAmf) GetSCTPConn(amfName string) *sctp.SCTPConn {
	if len(amf.tnlaList) == 0 {
		log.Error("[AMF][TNLA] TNLA list is empty")
		return nil
	}
	for _, tnla := range amf.tnlaList {
		if tnla.amfName == amfName {
			return tnla.sctpConn
		}
	}

	return nil
}

func (amf *GNBAmf) SetSCTPConn(amfName string, conn *sctp.SCTPConn) {
	if len(amf.tnlaList) == 0 {
		log.Error("[AMF][TNLA] TNLA list is empty")
		return
	}

	for _, tnla := range amf.tnlaList {
		if tnla.amfName == amfName {
			tnla.sctpConn = conn
			break
		}
	}
}

func (amf *GNBAmf) AddTNLA() {
	log.Info("[GNB] Add first TNLA")
	amf.tnlaList = append(amf.tnlaList, new(TNLAssociation))
}

func (tnla *TNLAssociation) SetAmfName(AmfName string) error {
	tnla.amfName = AmfName
	return nil
}

func (tnla *TNLAssociation) GetAmfName() string {
	return tnla.amfName
}

func (amf *GNBAmf) SetTNLAWeight(amfName string, weight int64) {
	if len(amf.tnlaList) == 0 {
		log.Error("[AMF][TNLA] TNLA list is empty")
		return
	}

	for _, tnla := range amf.tnlaList {
		if tnla.amfName == amfName {
			tnla.tnlaWeightFactor = weight
			break
		}
	}
}

func (amf *GNBAmf) SetTNLAUsage(amfName string, usage aper.Enumerated) {
	if len(amf.tnlaList) == 0 {
		log.Error("[AMF][TNLA] TNLA list is empty")
		return
	}

	for _, tnla := range amf.tnlaList {
		if tnla.amfName == amfName {
			tnla.usage = usage
			break
		}
	}
}

func (amf *GNBAmf) SetTNLAStreams(amfName string, streams uint16) {
	if len(amf.tnlaList) == 0 {
		log.Error("[AMF][TNLA] TNLA list is empty")
		return
	}

	for _, tnla := range amf.tnlaList {
		if tnla.amfName == amfName {
			tnla.streams = streams 
			break
		}
	}
}

func (amf *GNBAmf) GetTNLAStreams(amfName string) uint16 {
	if len(amf.tnlaList) == 0 {
		log.Error("[AMF][TNLA] TNLA list is empty")
		return 0
	}

	for _, tnla := range amf.tnlaList {
		if tnla.amfName == amfName {
			return tnla.streams
		}
	}

	return 0
}

func (amf *GNBAmf) GetAmfIp() string {
	return amf.amfIp
}

func (amf *GNBAmf) SetAmfIp(ip string) {
	amf.amfIp = ip
}

func (amf *GNBAmf) GetAmfPort() int {
	return amf.amfPort
}

func (amf *GNBAmf) setAmfPort(port int) {
	amf.amfPort = port
}

func (amf *GNBAmf) GetAmfId() int64 {
	return amf.amfId
}

func (amf *GNBAmf) setAmfId(id int64) {
	amf.amfId = id
}

func (amf *GNBAmf) GetAmfName() string {
	return amf.name
}

func (amf *GNBAmf) SetAmfName(name string) {
	amf.name = name
}

func (amf *GNBAmf) GetAmfCapacity() int64 {
	return amf.relativeAmfCapacity
}

func (amf *GNBAmf) SetAmfCapacity(capacity int64) {
	amf.relativeAmfCapacity = capacity
}

func (amf *GNBAmf) GetLenPlmns() int {
	return amf.lenPlmn
}

func (amf *GNBAmf) GetLenSlice() int {
	return amf.lenSlice
}

func (amf *GNBAmf) SetLenPlmns(value int) {
	amf.lenPlmn = value
}

func (amf *GNBAmf) SetLenSlice(value int) {
	amf.lenSlice = value
}
