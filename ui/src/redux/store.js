import { configureStore } from '@reduxjs/toolkit';
import loginReducer from './loginSlice';
import createAccountReducer from './createAccountSlice';

export const store = configureStore({
  reducer: {
    login: loginReducer,
    createAccount: createAccountReducer,
  },
});
