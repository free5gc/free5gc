import React, {Component} from 'react';
import {Modal} from "react-bootstrap";
import Form from "react-jsonschema-form";
import PropTypes from 'prop-types';

class SubscriberModal extends Component {
  static propTypes = {
    open: PropTypes.bool.isRequired,
    setOpen: PropTypes.func.isRequired,
    onSubmit: PropTypes.func.isRequired,
  };

  state = {};

  schema = {
    // title: "A registration form",
    // "description": "A simple form example.",
    type: "object",
    required: [
      "plmnID",
      "ueId",
      "authenticationMethod",
      "K",
      "OPOPcSelect",
      "OPOPc",
    ],
    properties: {
      plmnID: {
        type: "string",
        title: "PLMN ID",
        default: "20893",
      },
      ueId: {
        type: "string",
        title: "UE ID (IMSI)",
        default: "imsi-208930000000003",
      },
      authenticationMethod: {
        type: "string",
        title: "Authentication Method",
        default: "5G_AKA",
        enum: ["5G_AKA", "EAP_AKA_PRIME"],
      },
      K: {
        type: "string",
        title: "K",
        default: "8baf473f2f8fd09487cccbd7097c6862",
      },
      OPOPcSelect: {
        type: "string",
        title: "Operator Code Type",
        enum: ["OP", "OPc"],
        default: "OPc",
      },
      OPOPc: {
        type: "string",
        title: "Operator Code Value",
        default: "8e27b6af0e692e750f32667a3b14605d",
      },
    }
  };

  uiSchema = {
    OPOPcSelect: {
      "ui:widget": "select",
    },
    authenticationMethod: {
      "ui:widget": "select",
    },
  };

  onSubmitClick(result) {
    console.log("onSubmitClick", result);

    const formData = result.formData;
    const OP = formData["OPOPcSelect"] === "OP" ? formData["OPOPc"] : "";
    const OPc = formData["OPOPcSelect"] === "OPc" ? formData["OPOPc"] : "";

    let subscriberData = {
      "plmnID": formData["plmnID"], // Change required
      "ueId": formData["ueId"], // Change required
      "AuthenticationSubscription": {
        "authenticationManagementField": "8000",
        "authenticationMethod": formData["authenticationMethod"], // "5G_AKA", "EAP_AKA_PRIME"
        "milenage": {
          "op": {
            "encryptionAlgorithm": 0,
            "encryptionKey": 0,
            "opValue": OP // Change required
          }
        },
        "opc": {
          "encryptionAlgorithm": 0,
          "encryptionKey": 0,
          "opcValue": OPc // Change required (one of OPc/OP should be filled)
        },
        "permanentKey": {
          "encryptionAlgorithm": 0,
          "encryptionKey": 0,
          "permanentKeyValue": formData["K"] // Change required
        },
        "sequenceNumber": "16f3b3f70fc2",
      },
      "AccessAndMobilitySubscriptionData": {
        "gpsis": [
          "msisdn-0900000000"
        ],
        "nssai": {
          "defaultSingleNssais": [
            {
              "sd": "010203",
              "sst": 1
            },
            {
              "sd": "112233",
              "sst": 1
            }
          ],
          "singleNssais": [
            {
              "sd": "010203",
              "sst": 1
            },
            {
              "sd": "112233",
              "sst": 1
            }
          ]
        },
        "subscribedUeAmbr": {
          "downlink": "1000 Kbps",
          "uplink": "1000 Kbps"
        },
      },
      "SmfSelectionSubscriptionData": {
        "subscribedSnssaiInfos": {
          "01010203": {
            "dnnInfos": [
              {
                "dnn": "internet"
              }
            ]
          },
          "01112233": {
            "dnnInfos": [
              {
                "dnn": "internet"
              }
            ]
          }
        },
      }
    };

    this.props.onSubmit(subscriberData);
  }

  render() {
    return (
      <Modal
        show={this.props.open}
        className={"fields__edit-modal theme-light"}
        backdrop={"static"}
        onHide={this.props.setOpen.bind(this, false)}>
        <Modal.Header closeButton>
          <Modal.Title id="example-modal-sizes-title-lg">
            New Subscriber
          </Modal.Title>
        </Modal.Header>

        <Modal.Body>
          <Form schema={this.schema}
                uiSchema={this.uiSchema}
                onSubmit={this.onSubmitClick.bind(this)}/>
        </Modal.Body>
      </Modal>
    );
  }
}

export default SubscriberModal;
