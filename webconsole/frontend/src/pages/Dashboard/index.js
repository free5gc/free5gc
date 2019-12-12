import React from 'react';
import {Button, Jumbotron} from "react-bootstrap";
import {Link} from "react-router-dom";

const Dashboard = () => (
  <div className="content">
    <div className="container-fluid">
      <div className="row">
        <div className="col-md-12">
          <div className="card">
            {/*<div className="header">*/}
            {/*  <h4>Dashboard</h4>*/}
            {/*</div>*/}
            <div className="content">
              <Jumbotron>
                <h2 style={{marginTop: 24}}>
                  free5GC Web Console
                </h2>
                <p>
                  Manage the 5G core network easily by this web console.<br/>
                  To add SIM card information, go to manage subscribers page.
                </p>
                <p style={{marginTop: 24}}>
                  <Button bsStyle="primary">
                    <Link to={"/subscriber"}>Manage Subscribers</Link>
                  </Button>
                </p>
              </Jumbotron>
            </div>

            <p>&nbsp;</p><p>&nbsp;</p>
            <p>&nbsp;</p><p>&nbsp;</p>
            <p>&nbsp;</p><p>&nbsp;</p>
            <p>&nbsp;</p><p>&nbsp;</p>
            <p>&nbsp;</p><p>&nbsp;</p>

          </div>
        </div>
      </div>
    </div>
  </div>
);

export default Dashboard;
