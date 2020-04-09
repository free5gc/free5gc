import Http from './Http';
import {store} from '../index';
import ueinfoActions from "../redux/actions/ueinfoActions";
import UEInfo from "../models/UEInfo";
import UEInfoDetail from "../models/UEInfoDetail";

class UeInfoApiHelper {

  static async fetchRegisteredUE() {
    try {
      let url =  "http://localhost:5000/api/registered-ue-context"
      console.log("Making request to ", url, " ....")
      let response = await Http.get(url);
      if (response.status === 200 && response.data) {

        const registered_users = response.data.map(
            function(ue_context) {

              return new UEInfo(
                ue_context.Supi, 
                ue_context.CmState
              );
            });
            
        store.dispatch(ueinfoActions.setRegisteredUE(registered_users));
        return true;
      } else {

        console.log("Request fail")
        console.log("Response Status: ", response.status)
        console.log("Response Status: ", response.data)

        let err_msg;
        if (response.data !== undefined){
          err_msg = response.data
        } else {
          err_msg = "Fetch Registered UE Request fail"
        }
        store.dispatch(ueinfoActions.setRegisteredUEError(err_msg));
      }
    } catch (error) {
        console.log(error)
        let err_msg;
        if (error.response !== undefined){
          err_msg = error.response.data.cause
        } else {
          err_msg = "Fetch Registered UE Request fail"
        }
        store.dispatch(ueinfoActions.setRegisteredUEError(err_msg));
    }

    return false;
  }

  static async fetchUEInfoDetail(supi) {
    try {
      let url = `http://localhost:5000/api/registered-ue-context/${supi}`
      console.log("Making request to ", url, " ....")

      let response = await Http.get(url);
      if (response.status === 200 && response.data) {
        //To do: implement set rgistered ue action
        
        console.log(response.data)

        let ue_context = response.data[0]
        store.dispatch(ueinfoActions.setUEInfoDetailAMF(ue_context));

        let smContextRef = ue_context.PduSessions[0].SmContextRef
    
        return [true, smContextRef];
      } else {

        console.log("Request fail")
        console.log("Response Status: ", response.status)
      }
    } catch (error) {
        console.log(error)
    }

    return [false, ""];
  }

  static async fetchUEInfoDetailSMF(smContextRef) {
    try {
      let  url = `http://localhost:5000/api/ue-pdu-session-info/${smContextRef}`
      console.log("Making request to ", url, " ....")

      let response = await Http.get(url);
      if (response.status === 200 && response.data) {
        //To do: implement set rgistered ue action

        let smContext = response.data
        store.dispatch(ueinfoActions.setUEInfoDetailSMF(smContext));

       
        return true;
      } else {

        console.log("Request fail")
        console.log("Response Status: ", response.status)
      }
    } catch (error) {
        console.log(error)
    }

    return false;
  }
}

export default UeInfoApiHelper;
