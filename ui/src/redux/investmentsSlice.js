import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';

const DEFAULT_PAGE_SIZE = 10;

export const queryTrades = createAsyncThunk(
  'investments/queryTrades',
  async (payload) => {
    try {
      const {page, pageSize, orderBy, order} = payload;
      const url = `/api/v1/trades?page=${page}&page_size=${pageSize}&order_by=${orderBy}&order=${order}`;
      const response = await axios.get(url);
      return response;
    } catch (error) {
      return error.response;
    }
  }
);

export const investmentsSlice = createSlice({
  name: "investments",
  initialState: {
    pageSize: DEFAULT_PAGE_SIZE,
    page: 1,
    orderBy: "date_purchased",
    order: "desc",
    status: "idle",
    payload: null,
    errorMessage: "",
    shouldReload: false,
    init: false,
  },
  reducers: {
    updateTableInfo: (state, action) => {
      console.log(action.payload.sort)
      state.page = action.payload.page.index + 1;
      state.pageSize = action.payload.page.size;
      state.orderBy = action.payload.sort.field;
      state.order = state.order === "desc" ? "asc" : "desc";
      state.shouldReload = true;
    },
    resetShouldReload: (state) => {
      state.shouldReload = false;
    },
    setInitialState: (state, action) => {
      state.page = action.payload.get("page")
      state.pageSize = action.payload.get("page_size")
      state.orderBy = action.payload.get("order_by")
      state.order = action.payload.get("order")
      state.init = true;
      state.shouldReload = true;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(queryTrades.pending, (state) => {
        state.status = "loading";
        state.errorMessage = "";
      })
      .addCase(queryTrades.fulfilled, (state, action) => {
        state.status = "idle";
        if (action.payload.status === 200) {
          state.payload = action.payload;
          state.errorMessage = "";
        } else {
          state.payload = {};
          state.errorMessage = action.payload.data.message;
        }
      })
      .addCase(queryTrades.rejected, (state) => {
        state.status = "idle";
        state.errorMessage = "Had some trouble fetching transactions.";
      });
  },
});

export const { updateTableInfo, resetShouldReload, setInitialState } = investmentsSlice.actions;

export default investmentsSlice.reducer;
