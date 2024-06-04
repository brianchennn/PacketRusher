/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package interface_management

import (
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"sync"

	"github.com/free5gc/ngap"

	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
)

func AmfConfigurationUpdateAcknowledge(amfPool *sync.Map) ([]byte, error) {
	message := BuildAmfConfigurationUpdateAcknowledge(amfPool)

	return ngap.Encoder(message)
}

func BuildAmfConfigurationUpdateAcknowledge(amfPool *sync.Map) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeAMFConfigurationUpdate
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentAMFConfigurationUpdateAcknowledge
	successfulOutcome.Value.AMFConfigurationUpdateAcknowledge = new(ngapType.AMFConfigurationUpdateAcknowledge)

	amfTNLAssociationSetupList := ngapType.AMFTNLAssociationSetupList{}

	amfPool.Range(
		func(k, v any) bool {
			gnbAmf, ok := v.(*context.GNBAmf)
			if !ok || gnbAmf.GetState() == 0 {
				return true
			}

			endpointIpAddress_Ngap := ngapConvert.IPAddressToNgap(gnbAmf.GetAmfIp(), "")
			amfTNLAssociationSetupItem := ngapType.AMFTNLAssociationSetupItem{
				AMFTNLAssociationAddress: ngapType.CPTransportLayerInformation{
					Present:           ngapType.CPTransportLayerInformationPresentEndpointIPAddress,
					EndpointIPAddress: &endpointIpAddress_Ngap,
				},
			}
			amfTNLAssociationSetupList.List = append(amfTNLAssociationSetupList.List, amfTNLAssociationSetupItem)
			return true
		},
	)

	ie1 := ngapType.AMFConfigurationUpdateAcknowledgeIEs{
		Id: ngapType.ProtocolIEID{
			Value: ngapType.ProtocolIEIDAMFTNLAssociationSetupList,
		},
		Criticality: ngapType.Criticality{
			Value: ngapType.CriticalityPresentReject,
		},
		Value: ngapType.AMFConfigurationUpdateAcknowledgeIEsValue{
			Present:                    ngapType.AMFConfigurationUpdateAcknowledgeIEsPresentAMFTNLAssociationSetupList,
			AMFTNLAssociationSetupList: &amfTNLAssociationSetupList,
		},
	}

	/*criticalityDiagnostics := ngapType.CriticalityDiagnostics{
		ProcedureCode: &ngapType.ProcedureCode{
			Value: ngapType.ProcedureCodeAMFConfigurationUpdate,
		},
		TriggeringMessage: &ngapType.TriggeringMessage{
			Value: ngapType.TriggeringMessagePresentSuccessfulOutcome,
		},
		ProcedureCriticality: &ngapType.Criticality{
			Value: ngapType.CriticalityPresentReject,
		},
		IEExtensions: &ngapType.ProtocolExtensionContainerCriticalityDiagnosticsExtIEs{},
	}

	/*ie2 := ngapType.AMFConfigurationUpdateAcknowledgeIEs{
		Id: ngapType.ProtocolIEID{
			Value: ngapType.ProtocolIEIDCriticalityDiagnostics,
		},
		Criticality: ngapType.Criticality{
			Value: ngapType.CriticalityPresentIgnore,
		},
		Value: ngapType.AMFConfigurationUpdateAcknowledgeIEsValue{
			Present:                ngapType.AMFConfigurationUpdateAcknowledgeIEsPresentCriticalityDiagnostics,
			CriticalityDiagnostics: &criticalityDiagnostics,
		},
	}*/

	successfulOutcome.Value.AMFConfigurationUpdateAcknowledge.ProtocolIEs.List = append(successfulOutcome.Value.AMFConfigurationUpdateAcknowledge.ProtocolIEs.List, ie1)
	//successfulOutcome.Value.AMFConfigurationUpdateAcknowledge.ProtocolIEs.List = append(successfulOutcome.Value.AMFConfigurationUpdateAcknowledge.ProtocolIEs.List, ie2)

	return
}
