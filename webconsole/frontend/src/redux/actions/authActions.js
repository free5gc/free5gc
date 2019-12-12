
export default class authActions {
  static SET_USER = 'AUTH/SET_USER';
  static LOGOUT = 'AUTH/LOGOUT';

  /**
   * @param user  {User}
   */
  static setUser(user) {
    return {
      type: this.SET_USER,
      user: user,
    };
  }

  static logout() {
    return {
      type: this.LOGOUT,
    };
  }
}
