import React from 'react';
import {Route, Switch} from 'react-router-dom';
import Login from './Login';

const Auth = ({isLoggedIn}) => {
  if (isLoggedIn) {
    return null;
  }

  return (
    <div className="wrapper">
      <Switch>
        <Route component={Login}/>
      </Switch>
    </div>
  )
};

export default Auth;
