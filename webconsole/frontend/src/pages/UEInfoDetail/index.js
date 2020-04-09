import React, {Component} from 'react';
import {withRouter} from "react-router-dom";
import {connect} from "react-redux";
import {Button, Jumbotron, Table} from "react-bootstrap";
import {Link} from "react-router-dom";
import ApiHelper from "../../util/ApiHelper";
import UEInfoApiHelper from "../../util/UEInfoApiHelper"

class UEInfoDetail extends Component {

    constructor(props) {

        super(props);

        this.getAMFUEContexts = this.getAMFUEContexts.bind(this)  
        this.getSMFInfo = this.getSMFInfo.bind(this)
        this.getPCFInfo = this.getPCFInfo.bind(this)

    }

    async deleteSubscriber(subscriber) {
        if (!window.confirm(`Delete subscriber ${subscriber.id}?`))
          return;
    
        const result = await ApiHelper.deleteSubscriber(subscriber.id, subscriber.plmn);
        ApiHelper.fetchSubscribers().then();
        if (!result) {
          alert("Error deleting subscriber: " + subscriber.id);
        }
      }

      componentDidMount() {

        // console.log("In UEInfoDetail")
        // console.log("In componentDidMount")

      }

    getAMFUEContexts() {

        var UEContexts = this.props.amfInfo
        var PduSessions = this.props.amfInfo.PduSessions

        if (PduSessions === undefined) {
            return Arr
        }

        var Arr = []
        Object.getOwnPropertyNames(UEContexts).forEach(
            function (key, idx, array) { 

                if (key !== "PduSessions") {
                    Arr.push(
                        <tr key={key}>
                            <td>{key}</td>
                            <td>{UEContexts[key]}</td>
                        </tr>
                    )
                }
                
            });
        
        PduSessions.map( (obj) => { 
            
            for (var key in obj) {

                Arr.push(
                    <tr key={key}>
                        <td>{key}</td>
                        <td>{obj[key]}</td>
                    </tr>
                )
            }  
        });
          return Arr
    }

    getSMFInfo() {
        var smfInfo = this.props.smfInfo
        var Arr = []

        if (smfInfo === undefined) {
            return Arr
        }

        let smContext = {
            AnType: smfInfo.AnType,
            Dnn: smfInfo.Dnn,
            LocalSEID: smfInfo.LocalSEID,
            PDUAddress: smfInfo.PDUAddress,
            PDUSessionID: smfInfo.PDUSessionID,
            RemoteSEID: smfInfo.RemoteSEID,
            Sd: smfInfo.Sd,
            Sst: smfInfo.Sst
        };

        Object.getOwnPropertyNames(smContext).forEach(
            function (key) { 


                //if key 
                Arr.push(
                    <tr key={key}>
                        <td>{key}</td>
                        <td>{smContext[key]}</td>
                    </tr>
                )   
        });

        return Arr;
    }

    getPCFInfo() {
        var AmPolicyData = this.props.pcfInfo.AmPolicyData
        var Arr = []

        Object.getOwnPropertyNames(AmPolicyData).forEach(
            function (obj) { 

                switch (obj) {
                    case "Triggers": 
                    
                    AmPolicyData[obj].forEach(
                            function(value, index, array) {
                                var key = "Trigger " + (index+1).toString()
                                Arr.push(
                                    <tr key={key}>
                                        <td>{key}</td>
                                        <td>{value}</td>
                                    </tr>
                                )
                            }
                        )
                    break;
                    case "Areas":
                        AmPolicyData[obj].forEach(
                            function(value, index, array) {
                                var key = "Area " + (index+1).toString()
                                Arr.push(
                                    <tr key={key}>
                                        <td>{key}</td>
                                        <td>{value}</td>
                                    </tr>
                                )
                            }
                        )
                            
                    break;

                    default:
                        Arr.push(
                            <tr key={obj}>
                                <td>{obj}</td>
                                <td>{AmPolicyData[obj]}</td>
                            </tr>
                        )   
                }
        });
        return Arr;

    }

    render() {
        return (
                <div className="container-fluid">
                    <div className="row">
                        <div className="col-md-12">
                            <div className="card">
                                <div className="header subscribers__header">
                                    <h4>{`AMF Information [SUPI:${this.props.amfInfo.Supi}]`}</h4><br></br>
                                </div>
                                <div className="content subscribers__content">
                                    <Table className="subscribers__table" striped bordered condensed hover>
                                    <thead>
                                    <tr>
                                        <th colSpan={1}>Information Entity</th>
                                        <th colSpan={2}>Value</th>
                                    </tr>
                                    </thead>
                                    <tbody>
                                        {this.getAMFUEContexts()}

                                    </tbody>
                                    </Table>
                                </div>
                                <div className="pdu__Sessions">
                            
                                </div>
                            </div>
                            <div className="card">
                                <div className="header subscribers__header">
                                    <h4>{`SMF Information [SUPI:${this.props.amfInfo.Supi}]`}</h4>
                                </div>
                                <div className="content subscribers__content">
                                    <Table className="subscribers__table" striped bordered condensed hover>
                                    <thead>
                                    <tr>
                                        <th colSpan={1}>Information Entity</th>
                                        <th colSpan={2}>Value</th>
                                    </tr>
                                    </thead>
                                    <tbody>
                                        {this.getSMFInfo()}
                                    </tbody>
                                    </Table>
                                </div>
                                    <p>&nbsp;</p><p>&nbsp;</p>
                                    <p>&nbsp;</p><p>&nbsp;</p>
                                    <p>&nbsp;</p><p>&nbsp;</p>
                            </div>
                        </div>
                    </div>
                </div>
        );
    }

}

const mapStateToProps = state => ({
    amfInfo: state.ueinfo.amfInfo,
    smfInfo: state.ueinfo.smfInfo,
    pcfInfo: state.ueinfo.ueInfoDetail.pcfInfo

  });

export default withRouter(connect(mapStateToProps)(UEInfoDetail));

