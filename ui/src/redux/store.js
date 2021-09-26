import { configureStore } from '@reduxjs/toolkit';
import loginReducer from './loginSlice';
import createAccountReducer from './createAccountSlice';
import mainPageReducer from './mainPageSlice';

export const store = configureStore({
  reducer: {
    login: loginReducer,
    createAccount: createAccountReducer,
    mainPage: mainPageReducer,
  },
});
