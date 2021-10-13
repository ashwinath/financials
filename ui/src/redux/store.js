import { configureStore, combineReducers } from '@reduxjs/toolkit';
import loginReducer from './loginSlice';
import createAccountReducer from './createAccountSlice';
import investmentsReducer from './investmentsSlice';

const appReducer = combineReducers({
    login: loginReducer,
    createAccount: createAccountReducer,
    investments: investmentsReducer,
});


const rootReducer = (state, action) => {
  let newState = state;
  if (action.type === "login/logout") {
    newState = undefined;
  }
  return appReducer(newState, action);
}

export const store = configureStore({
  reducer: rootReducer,
});
