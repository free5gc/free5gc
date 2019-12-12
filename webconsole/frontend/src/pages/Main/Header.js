import React, { Component } from 'react';
import { connect } from 'react-redux';
import { toggleMobileNavVisibility } from '../../redux/reducers/layout';
import { Navbar, Nav, NavItem, NavDropdown, MenuItem, FormGroup, FormControl } from 'react-bootstrap';
import { withRouter } from 'react-router-dom';
import AuthHelper from "../../util/AuthHelper";

class Header extends Component {
  static async logout() {
    await AuthHelper.logout();
  }

  render() {
    return (
      <Navbar fluid={true}>
        <Navbar.Header>
          <button type="button" className="navbar-toggle" data-toggle="collapse" onClick={this.props.toggleMobileNavVisibility}>
            <span className="sr-only">Toggle navigation</span>
            <span className="icon-bar"/>
            <span className="icon-bar"/>
            <span className="icon-bar"/>
          </button>
        </Navbar.Header>

        <Navbar.Collapse>
          {/*<Nav>*/}
          {/*  <NavDropdown title={<i className="fa fa-globe" />} id="basic-nav-dropdown">*/}
          {/*    <MenuItem>Action</MenuItem>*/}
          {/*    <MenuItem>Another action</MenuItem>*/}
          {/*    <MenuItem>Something else here</MenuItem>*/}
          {/*    <MenuItem divider />*/}
          {/*    <MenuItem>Separated link</MenuItem>*/}
          {/*  </NavDropdown>*/}
          {/*</Nav>*/}

          {/*<div className="separator"/>*/}

          {/*<Navbar.Form pullLeft>*/}
            {/*<FormGroup>*/}
              {/*<span className="input-group-addon"><i className="fa fa-search"/></span>*/}
              {/*<FormControl type="text" placeholder="Type to search" />*/}
            {/*</FormGroup>*/}
          {/*</Navbar.Form>*/}

          <Nav pullRight>
            {/*<NavItem>Account</NavItem>*/}
            {/*<NavDropdown title="Dropdown" id="right-nav-bar">*/}
            {/*  <MenuItem>Action</MenuItem>*/}
            {/*  <MenuItem>Another action</MenuItem>*/}
            {/*  <MenuItem>Something else here</MenuItem>*/}
            {/*  <MenuItem divider />*/}
            {/*  <MenuItem>Separated link</MenuItem>*/}
            {/*</NavDropdown>*/}
            <NavItem onClick={Header.logout}>Log out</NavItem>
          </Nav>
        </Navbar.Collapse>
      </Navbar>
    );
  }
}

const mapDispatchToProp = dispatch => ({
  toggleMobileNavVisibility: () => dispatch(toggleMobileNavVisibility()),
});

export default withRouter(connect(null, mapDispatchToProp)(Header));
