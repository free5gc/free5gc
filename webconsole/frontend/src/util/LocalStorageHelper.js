import User from '../models/User'

export default class LocalStorageHelper {

  /**
   * @param user {models/User}
   */
  static setUserInfo(user) {
    localStorage.setItem('user_info', user ? user.serialize() : null);
  }

  /**
   * @return {(models/User|null)}
   */
  static getUserInfo() {
    let json = localStorage.getItem('user_info');
    return json === null ? null : ApiTokens.deserialize(json);
  }
}
