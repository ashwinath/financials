import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

export const loginAsync = createAsyncThunk(
  'login/loginAsync',
  async (user) => {
    try {
      const response = await axios.post(
        '/api/v1/login',
        user,
      );
      return response;
    } catch (error) {
      return error.response;
    }
  }
);

export const loginSlice = createSlice({
  name: "login",
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
      .addCase(loginAsync.pending, (state) => {
        state.status = "loading";
        state.errorMessage = "";
      })
      .addCase(loginAsync.fulfilled, (state, action) => {
        state.status = "idle";
        if (action.payload.status === 200) {
          state.isLoggedIn = true;
          state.errorMessage = ""
        } else {
          state.errorMessage = action.payload.data.message;
        }
      })
      .addCase(loginAsync.rejected, (state) => {
        state.status = "idle";
        state.errorMessage = "Something went wrong logging you in.";
      });
  },
});

export const { updateUsername, updatePassword } = loginSlice.actions

export default loginSlice.reducer;
