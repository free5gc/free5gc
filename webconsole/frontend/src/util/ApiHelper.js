import Http from './Http';
import {store} from '../index';
import subscriberActions from "../redux/actions/subscriberActions";
import Subscriber from "../models/Subscriber";

class ApiHelper {

  static async fetchSubscribers() {
    try {
      let response = await Http.get('subscriber');
      if (response.status === 200 && response.data) {
        const subscribers = response.data.map(val => new Subscriber(val['ueId'], val['plmnID']));
        store.dispatch(subscriberActions.setSubscribers(subscribers));
        return true;
      }
    } catch (error) {
    }

    return false;
  }

  static async createSubscriber(subscriberData) {
    try {
      let response = await Http.post(
        `subscriber/${subscriberData["ueId"]}/${subscriberData["plmnID"]}`, subscriberData);
      if (response.status === 201)
        return true;
    } catch (error) {
      console.error(error);
    }

    return false;
  }

  static async updateSubscriber(subscriberData) {
    try {
      let response = await Http.patch(
        `subscriber/${subscriberData["ueId"]}/${subscriberData["plmnID"]}`, subscriberData);
      if (response.status === 201)
        return true;
    } catch (error) {
      console.error(error);
    }

    return false;
  }

  static async deleteSubscriber(id, plmn) {
    try {
      let response = await Http.delete(`subscriber/${id}/${plmn}`);
      if (response.status === 204)
        return true;
    } catch (error) {
      console.error(error);
    }

    return false;
  }
}

export default ApiHelper;
