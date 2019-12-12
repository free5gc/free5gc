import actions from '../actions/authActions';

const initialState = {
  user: null, // also used to identify if user is logged in
};

export default function reducer(state = initialState, action) {
  let nextState = {...state};

  switch (action.type) {
    case actions.SET_USER:
      nextState.user = action.user;
      return nextState;

    case actions.LOGOUT:
      nextState.user = null;
      return nextState;

    default:
      return state;
  }
}
