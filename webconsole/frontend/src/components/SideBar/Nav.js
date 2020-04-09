import React, {Component} from 'react';
import {Link, withRouter} from 'react-router-dom';

class Nav extends Component {
  state = {};

  render() {
    let {location} = this.props;
    /* Icons:
     *  - https://fontawesome.com/icons
     *  - http://themes-pixeden.com/font-demos/7-stroke/
     */
    return (
      <ul className="nav">
        <li className={location.pathname === '/' ? 'active' : null}>
          <Link to="/ueinfo">
            <i className="pe-7s-network"/>
            <p>Dashboard</p>
          </Link>
        </li>

        <li className={this.isPathActive('/subscriber') ? 'active' : null}>
          <Link to="/subscriber">
            <i className="fa fa-mobile-alt"/>
            <p>Subscribers</p>
          </Link>
        </li>

        <li className={this.isPathActive('/tasks') ? 'active' : null}>
          <Link to="/tasks">
            <i className="pe-7s-graph1"/>
            <p>Analytics</p>
          </Link>
        </li>

      </ul>
    );
  }

  isPathActive(path) {
    return this.props.location.pathname.startsWith(path);
  }
}

export default withRouter(Nav);
