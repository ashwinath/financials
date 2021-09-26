import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

export const querySessionAsync = createAsyncThunk(
  'mainPage/querySession',
  async (user) => {
    try {
      const response = await axios.get(
        '/api/v1/session',
      );
      return response;
    } catch (error) {
      return error.response;
    }
  }
);

export const mainPageSlice = createSlice({
  name: "mainPage",
  initialState: {
    triedLoggingIn: false,
    status: "idle",
    isLoggedIn: true,
  },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(querySessionAsync.pending, (state) => {
        state.status = "loading";
      })
      .addCase(querySessionAsync.fulfilled, (state, action) => {
        state.status = "idle";
        state.triedLoggingIn = true;
        if (action.payload.status === 200) {
          state.isLoggedIn = true;
        } else {
          state.isLoggedIn = false;
        }
      })
      .addCase(querySessionAsync.rejected, (state) => {
        state.status = "idle";
        state.isLoggedIn = false;
      })
  },
});

export default mainPageSlice.reducer;
