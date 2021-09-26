import React from 'react';

import {
  useHistory,
} from "react-router-dom";

import { 
  useDispatch,
  useSelector,
} from 'react-redux';
import { querySessionAsync } from "../redux/mainPageSlice";

export default function HomePage() {
  const history = useHistory();
  const dispatch = useDispatch();

  const { isLoggedIn, status, triedLoggingIn } = useSelector((state) => state.mainPage)
  if (!triedLoggingIn && status === "idle") {
    dispatch(querySessionAsync())
  }

  if (!isLoggedIn && triedLoggingIn) {
    history.push("/login");
  }

  return (
    <h1>Hello world</h1>
  );
}
