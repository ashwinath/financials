import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';

import { LoadingPage } from '../pages';
import { querySessionAsync } from "../redux/loginSlice";

export function PrivateRoute({component: Component, ...rest}) {
  const { isLoggedIn, status, triedLoggingIn } = useSelector((state) => state.login)
  const dispatch = useDispatch();

  if (!isLoggedIn && !triedLoggingIn && status === "idle") {
    dispatch(querySessionAsync())
  }

  if (status === "loading") {
    return <LoadingPage/>;
  }

  return (
    <Route
      {...rest}
      render={props => (
        isLoggedIn
        ? <Component {...props} />
        : <Redirect to="/login" />
      )}
    />
  );
};

export const PublicRoute = ({component: Component, restricted, ...rest}) => {
  const { isLoggedIn, status, triedLoggingIn } = useSelector((state) => state.login)
  const dispatch = useDispatch();

  if (!isLoggedIn && !triedLoggingIn && status === "idle") {
    dispatch(querySessionAsync())
  }

  if (status === "loading") {
    return <LoadingPage/>;
  }

  return (
    <Route
      {...rest}
      render={props => (
        isLoggedIn
        ? <Redirect to="/" />
        : <Component {...props} />
      )}
    />
  );
};
