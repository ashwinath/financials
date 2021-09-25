import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

export const loginAsync = createAsyncThunk(
  'login/loginAsync',
  async (user) => {
    const response = await axios.post(
      '/api/v1/login',
      user,
    );
    return response;
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
      state.errorMessage = ""
    },
    updatePassword: (state, action) => {
      state.password = action.payload;
      state.errorMessage = ""
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(loginAsync.pending, (state) => {
        state.status = "loading";
        state.errorMessage = ""
      })
      .addCase(loginAsync.fulfilled, (state, action) => {
        state.status = "idle";
          state.isLoggedIn = true;
          state.errorMessage = ""
      })
      .addCase(loginAsync.rejected, (state) => {
        state.errorMessage = "Wrong credentials provided."
      });
  },
});

export const { updateUsername, updatePassword } = loginSlice.actions

export default loginSlice.reducer;
