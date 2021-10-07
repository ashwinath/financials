import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios from 'axios';
import moment from "moment";
import {
  convertDateToString,
} from "../utils";

const DEFAULT_PAGE_SIZE = 20;

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

export const submitTrade = createAsyncThunk(
  'investments/submitTrade',
  async (singleTrade) => {
    try {
      const url = "/api/v1/trades";
      const response = await axios.post(url, {
        transactions: [singleTrade],
      });
      return response;
    } catch (error) {
      return error.response;
    }
  }
);

export const submitTrades = createAsyncThunk(
  'investments/submitTrades',
  async (trades) => {
    try {
      const url = "/api/v1/trades";
      const response = await axios.post(url, {
        transactions: trades,
      });
      return response;
    } catch (error) {
      return error.response;
    }
  }
);

const addTradeFormInitState = {
  date_purchased: convertDateToString(moment()),
  symbol: "",
  price_each: 0,
  quantity: 0,
  trade_type: "buy",
}

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
    isAddTradeModalOpen: false,
    isAddBulkTradeModalOpen: false,
    addTradeForm: addTradeFormInitState,
    isTradeFormSubmitting: false,
    submitSuccess: "none", // can be none/success/failure
    tradeCSVRaw: "",
  },
  reducers: {
    updateTableInfo: (state, action) => {
      state.page = action.payload.page.index + 1;
      state.pageSize = action.payload.page.size;
      state.orderBy = action.payload.sort.field;
      state.order = action.payload.sort.direction;
      state.shouldReload = true;
    },
    resetShouldReload: (state) => {
      state.shouldReload = false;
    },
    setInitialState: (state, action) => {
      // This function gets the query params and sets the intial state so it renders right
      const page = action.payload.get("page");
      state.page = page ? page : state.page;

      const pageSize = action.payload.get("page_size");
      state.pageSize = pageSize ? pageSize : state.pageSize;

      const orderBy = action.payload.get("order_by");
      state.orderBy = orderBy ? orderBy : state.orderBy;

      const order = action.payload.get("order");
      state.order = order ? order : state.order;

      state.init = true;
      state.shouldReload = true;
    },
    toggleIsAddTradeModalOpen: (state) => {
      state.submitSuccess = "none";
      state.isAddTradeModalOpen = !state.isAddTradeModalOpen;
    },
    toggleIsAddBulkTradeModalOpen: (state) => {
      state.submitSuccess = "none";
      state.isAddBulkTradeModalOpen = !state.isAddBulkTradeModalOpen;
    },
    updateAddTradeForm: (state, action) => {
      const payload = action.payload;
      state.addTradeForm[payload.name] = payload.value;
    },
    updateCSVTrades: (state, action) => {
      state.tradeCSVRaw = action.payload
    },
  },
  extraReducers: (builder) => {
    // Get trades
    builder
      .addCase(queryTrades.pending, (state) => {
        state.status = "loading";
        state.errorMessage = "";
      })
      .addCase(queryTrades.fulfilled, (state, action) => {
        state.status = "idle";
        if (action.payload.status === 200) {
          state.payload = action.payload.data;
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

    // Submit single trade
    builder
      .addCase(submitTrade.pending, (state) => {
        state.isTradeFormSubmitting = true;
        state.errorMessage = "";
      })
      .addCase(submitTrade.fulfilled, (state, action) => {
        state.isTradeFormSubmitting = false;
        if (action.payload.status === 201) {
          state.submitSuccess = "success";
          state.isAddTradeModalOpen = false;
          state.errorMessage = "";
          state.shouldReload = true;
          state.addTradeForm = addTradeFormInitState
          state.tradeCSVRaw = ""
        } else {
          state.submitSuccess = "failure";
          state.isAddTradeModalOpen = false;
          state.errorMessage = action.payload.data.message;
        }
      })
      .addCase(submitTrade.rejected, (state) => {
        state.isTradeFormSubmitting = false;
        state.submitSuccess = "failure";
        state.isAddTradeModalOpen = false;
        state.errorMessage = "Had some trouble submitting your trades.";
      });

    // Submit trades
    builder
      .addCase(submitTrades.pending, (state) => {
        state.isTradeFormSubmitting = true;
        state.errorMessage = "";
      })
      .addCase(submitTrades.fulfilled, (state, action) => {
        state.isTradeFormSubmitting = false;
        state.isAddBulkTradeModalOpen = false;
        if (action.payload.status === 201) {
          state.submitSuccess = "success";
          state.errorMessage = "";
          state.shouldReload = true;
          state.addTradeForm = addTradeFormInitState
          state.tradeCSVRaw = ""
        } else {
          state.submitSuccess = "failure";
          state.errorMessage = action.payload.data.message;
        }
      })
      .addCase(submitTrades.rejected, (state) => {
        state.isTradeFormSubmitting = false;
        state.submitSuccess = "failure";
        state.isAddTradeModalOpen = false;
        state.errorMessage = "Had some trouble submitting your trades.";
      });
  },
});

export const { 
  updateTableInfo,
  resetShouldReload,
  setInitialState,
  toggleIsAddTradeModalOpen,
  updateAddTradeForm,
  toggleIsAddBulkTradeModalOpen,
  updateCSVTrades,
} = investmentsSlice.actions;

export default investmentsSlice.reducer;
