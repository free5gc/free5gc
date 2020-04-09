import React from 'react';
import {Route} from 'react-router-dom';
import SubscriberOverview from "./SubscriberOverview";

const Subscribers = ({match}) => (
  //console.log("In Subscribers")
  <div className="content">
    <Route exact path={`${match.url}/`} component={SubscriberOverview} />
    {/*<Route path={`${match.url}/:uuid`} component={} />*/}
  </div>
);

export default Subscribers;
