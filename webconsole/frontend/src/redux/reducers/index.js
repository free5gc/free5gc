import {reducer as formReducer} from 'redux-form'
import auth from './auth';
import layout from './layout';
import subscriber from "./subscriber";

export default {
  auth,
  layout,
  subscriber,
  form: formReducer,
};
