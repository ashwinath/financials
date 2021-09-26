import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

export const createAsync = createAsyncThunk(
  'createAccount/createAsync',
  async (user) => {
    try {
      const response = await axios.post(
        '/api/v1/users',
        user,
      );
      return response;
    } catch (error) {
      return error.response;
    }
  }
);

export const createAccountSlice = createSlice({
  name: "createAccount",
  initialState: {
    username: "",
    password: "",
    status: "idle",
    isLoggedIn: false,
    errorMessage: "",
  },
  reducers: {
    updateUsername: (state, action) => {
      state.username = action.payload;
      state.errorMessage = "";
    },
    updatePassword: (state, action) => {
      state.password = action.payload;
      state.errorMessage = "";
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(createAsync.pending, (state) => {
        state.status = "loading";
        state.errorMessage = "";
      })
      .addCase(createAsync.fulfilled, (state, action) => {
        state.status = "idle";
        if (action.payload.status === 201) {
          state.isLoggedIn = true;
          state.errorMessage = "";
        } else {
          state.errorMessage = action.payload.data.message;
        }
      })
      .addCase(createAsync.rejected, (state) => {
        state.status = "idle";
        state.errorMessage = "Had some trouble creating an account";
      });
  },
});

export const { updateUsername, updatePassword } = createAccountSlice.actions;

export default createAccountSlice.reducer;
