import { useDispatch, useSelector } from 'react-redux';

import { useHistory } from "react-router-dom";

import { querySessionAsync } from "../redux/loginSlice";

export function useLoginHook() {
  const history = useHistory();
  const dispatch = useDispatch();
  const { isLoggedIn, status, triedLoggingIn } = useSelector((state) => state.login)

  if (!isLoggedIn && !triedLoggingIn && status === "idle") {
    dispatch(querySessionAsync())
  }

  if (!isLoggedIn && triedLoggingIn) {
    history.push("/login");
  }

  return status;
}
