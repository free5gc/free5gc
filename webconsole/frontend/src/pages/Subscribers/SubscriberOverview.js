import React, {Component} from 'react';
import {withRouter} from "react-router-dom";
import {connect} from "react-redux";
import {Button, Table} from "react-bootstrap";
import SubscriberModal from "./components/SubscriberModal";
import ApiHelper from "../../util/ApiHelper";

class SubscriberOverview extends Component {
  state = {
    subscriberModalOpen: false,
  };

  componentDidMount() {
    ApiHelper.fetchSubscribers().then();
  }

  openAddSubscriber() {
    this.setState({subscriberModalOpen: true});
  }

  async addSubscriber(subscriberData) {
    this.setState({subscriberModalOpen: false});

    if (!await ApiHelper.createSubscriber(subscriberData)) {
      alert("Error creating new subscriber");
    }
    ApiHelper.fetchSubscribers().then();
  }

  /**
   * @param subscriber  {Subscriber}
   */
  async deleteSubscriber(subscriber) {
    if (!window.confirm(`Delete subscriber ${subscriber.id}?`))
      return;

    const result = await ApiHelper.deleteSubscriber(subscriber.id, subscriber.plmn);
    ApiHelper.fetchSubscribers().then();
    if (!result) {
      alert("Error deleting subscriber: " + subscriber.id);
    }
  }

  render() {
    return (
      <div className="container-fluid">
        <div className="row">
          <div className="col-md-12">
            <div className="card">
              <div className="header subscribers__header">
                <h4>Subscribers</h4>
                <Button bsStyle={"primary"} className="subscribers__button"
                        onClick={this.openAddSubscriber.bind(this)}>
                  New Subscriber
                </Button>
              </div>
              <div className="content subscribers__content">
                <Table className="subscribers__table" striped bordered condensed hover>
                  <thead>
                  <tr>
                    <th style={{width: 80}}>PLMN</th>
                    <th colSpan={2}>UE ID</th>
                  </tr>
                  </thead>
                  <tbody>
                  {this.props.subscribers.map(subscriber => (
                    <tr key={subscriber.id}>
                      <td>{subscriber.plmn}</td>
                      <td>{subscriber.id}</td>
                      <td style={{textAlign: 'center'}}>
                        <i className="fa fa-trash-alt" onClick={this.deleteSubscriber.bind(this, subscriber)}/>
                      </td>
                    </tr>
                  ))}
                  </tbody>
                </Table>

                <p>&nbsp;</p><p>&nbsp;</p>
                <p>&nbsp;</p><p>&nbsp;</p>
                <p>&nbsp;</p><p>&nbsp;</p>
              </div>
            </div>
          </div>
        </div>

        <SubscriberModal open={this.state.subscriberModalOpen}
                         setOpen={val => this.setState({subscriberModalOpen: val})}
                         onSubmit={this.addSubscriber.bind(this)}/>
      </div>
    );
  }
}

const mapStateToProps = state => ({
  subscribers: state.subscriber.subscribers,
});

export default withRouter(connect(mapStateToProps)(SubscriberOverview));
