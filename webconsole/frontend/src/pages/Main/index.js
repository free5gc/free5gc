import React from 'react';
import {Route, withRouter} from 'react-router-dom';
import {connect} from 'react-redux';
import cx from 'classnames';
import {setMobileNavVisibility} from '../../redux/reducers/layout';

import Header from './Header';
import Footer from './Footer';
import SideBar from '../../components/SideBar';
/**
 * Pages
 */
import Dashboard from '../Dashboard';
import Subscribers from '../Subscribers';
import Tasks from '../Tasks';

const Main = ({
                mobileNavVisibility,
                hideMobileMenu,
                isLoggedIn,
                history
              }) => {
  if (!isLoggedIn) {
    return null;
  }

  history.listen(() => {
    if (mobileNavVisibility === true) {
      hideMobileMenu();
    }
  });

  return (
    <div className={cx({
      'nav-open': mobileNavVisibility === true
    })}>
      <div className="wrapper">
        <div className="close-layer" onClick={hideMobileMenu}/>
        <SideBar/>

        <div className="main-panel">
          <Header/>

          <Route exact path="/" component={Dashboard}/>
          <Route exact path="/subscriber" component={Subscribers}/>
          <Route path="/tasks" component={Tasks}/>

          <Footer/>
        </div>
      </div>
    </div>
  )
};

const mapStateToProp = state => ({
  mobileNavVisibility: state.layout.mobileNavVisibility
});

const mapDispatchToProps = (dispatch, ownProps) => ({
  hideMobileMenu: () => dispatch(setMobileNavVisibility(false))
});

export default withRouter(connect(mapStateToProp, mapDispatchToProps)(Main));
