import React, {Component} from 'react';
import {withRouter} from 'react-router-dom';
import {connect} from 'react-redux';

class UserInfo extends Component {

  state = {
    isShowingUserMenu: false
  };

  render() {
    let {location} = this.props;
    let {user} = this.props;
    let {isShowingUserMenu} = this.state;
    return (
      <div className="user-wrapper">
        <div className="user">
          <img src={user.imageUrl} alt={user.name} className="photo"/>
          <div className="userinfo">
            <div className="username">
              {user.username}
            </div>
            <div className="title">{user.name}</div>
          </div>
          {/*<span*/}
          {/*onClick={() => this.setState({ isShowingUserMenu: !this.state.isShowingUserMenu })}*/}
          {/*className={cx("pe-7s-angle-down collapse-arrow", {*/}
          {/*active: isShowingUserMenu*/}
          {/*})}/>*/}
        </div>
        {/*<Collapse in={isShowingUserMenu}>*/}
        {/*<ul className="nav user-nav">*/}
        {/*<li className={this.isPathActive('/profile') ? 'active' : null}>*/}
        {/*<Link to="/profile">My Profile</Link>*/}
        {/*</li>*/}
        {/*<li className={this.isPathActive('/profile') ? 'active' : null}>*/}
        {/*<Link to="/profile">Edit Profile</Link>*/}
        {/*</li>*/}
        {/*<li className={this.isPathActive('/profile') ? 'active' : null}>*/}
        {/*<Link to="/profile">Settings</Link>*/}
        {/*</li>*/}
        {/*</ul>*/}
        {/*</Collapse>*/}
      </div>
    );
  }

  isPathActive(path) {
    return this.props.location.pathname.startsWith(path);
  }
}

const mapStateToProps = state => ({
  user: state.auth.user
});

export default withRouter(connect(mapStateToProps)(UserInfo));
