/* eslint-disable no-useless-constructor */
import React, {Component} from "react";
import {connect} from 'react-redux';
import {withRouter} from 'react-router-dom';
import {Button, FormControl, FormGroup} from "react-bootstrap";
import AuthHelper from "../../util/AuthHelper"
import Free5gcLogo from "../../assets/images/free5gc_logo.png";

class Login extends Component {
  state = {
    submitDisabled: false,
    errorMsg: "",

    // Form
    username: "",
    password: "",
  };

  conponentWillMount() {
    this.setState({
      submitDisabled: false,
      errorMsg: "",
    });
  }

  validateForm() {
    return this.state.username.length > 0 && this.state.password.length > 0;
  }

  async handleSubmit(event) {
    event.preventDefault();

    if (!this.validateForm()) {
      this.setState({
        errorMsg: "Invalid inputs",
      });
      return;
    }

    this.setState({
      submitDisabled: true,
      errorMsg: "",
    });

    let result = await AuthHelper.login(this.state.username, this.state.password);

    if (result === true) {
      console.log('login successful');
    } else {
      this.setState({
        submitDisabled: false,
        errorMsg: "Wrong credentials",
      });
    }
  };

  render() {
    return (
      <div className="Login">
        <div className="LoginForm">
          <img src={Free5gcLogo} alt="free5GC"/>

          <form onSubmit={this.handleSubmit.bind(this)}>
            <span className="error-msg"><p>{this.state.errorMsg}&nbsp;</p></span>

            <FormGroup controlId="username" bsSize="large">
              <FormControl
                autoFocus
                type="text"
                placeholder="Username"
                value={this.state.username}
                onChange={e => this.setState({username: e.target.value})}
              />
            </FormGroup>

            <FormGroup controlId="password" bsSize="large">
              <FormControl
                type="password"
                placeholder="Password"
                value={this.state.password}
                onChange={e => this.setState({password: e.target.value})}
              />
            </FormGroup>

            <Button block type="submit" className="btn-login" disabled={this.state.submitDisabled}>
              Login
            </Button>
          </form>
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => ({});

export default withRouter(connect(mapStateToProps)(Login));
