import React from 'react';
import {connect} from 'react-redux';
import {withRouter} from 'react-router-dom';
/**
 * Pages
 */
import Main from '../Main';
import Auth from '../Auth';
import AuthHelper from "../../util/AuthHelper";

const App = (props) => {
  // Logout route
  if (props.location.pathname === '/logout') {
    if (props.user !== null) {
      AuthHelper.logout().then(success => {
        props.history.push('/');
      });
    } else {
      props.history.push('/');
    }
    return null;
  }

  // App
  return (
    <div>
      <Auth isLoggedIn={props.user != null}/>
      <Main isLoggedIn={props.user != null}/>
    </div>
  );
};

const mapStateToProp = state => ({
  user: state.auth.user
});

export default withRouter(connect(mapStateToProp)(App));
