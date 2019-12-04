
export default class subscriberActions {
  static SET_SUBSCRIBERS = 'SUBSCRIBER/SET_SUBSCRIBERS';

  /**
   * @param subscribers  {Subscriber}
   */
  static setSubscribers(subscribers) {
    return {
      type: this.SET_SUBSCRIBERS,
      subscribers: subscribers,
    };
  }
}
