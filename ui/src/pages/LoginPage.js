import React from 'react';
import {
  EuiForm,
  EuiFieldText,
  EuiFormRow,
  EuiSpacer,
  EuiPageTemplate,
  EuiButton,
  EuiFieldPassword,
  EuiEmptyPrompt,
  EuiCallOut,
  EuiText,
} from '@elastic/eui';

import { Link, useHistory } from "react-router-dom";
import { useDispatch, useSelector } from 'react-redux';
import { updateUsername, updatePassword, loginAsync } from '../redux/loginSlice';

import { LoadingPage } from ".";
import { useLoginHook } from "../hooks/login";

function LoginForm() {
  const {
    username,
    password,
    isLoggedIn,
  } = useSelector((state) => state.login)
  const user = {
    username,
    password,
  }
  const dispatch = useDispatch();
  const history = useHistory();

  if (isLoggedIn) {
    history.push("/");
  }

  return (
    <EuiForm component="form">
      <EuiFormRow label="Username">
        <EuiFieldText
          name="username"
          onChange={(e) => dispatch(updateUsername(e.target.value))}
        />
      </EuiFormRow>

      <EuiFormRow label="Password">
        <EuiFieldPassword
          name="password"
          onChange={(e) => dispatch(updatePassword(e.target.value))}
        />
      </EuiFormRow>

      <EuiSpacer />

      <Link to="/create">
        <EuiText>Don't have one? Create here.</EuiText>
      </Link>

      <EuiSpacer />

      <EuiButton
        type="submit"
        fill
        onClick={
          (e) => {
            e.preventDefault();
            dispatch(loginAsync(user));
          }
        }
      >
        Login
      </EuiButton>

    </EuiForm>
  );
}

export function LoginPage() {
  const {
    errorMessage,
    status,
  } = useSelector((state) => state.login)
  useLoginHook();

  if (status === "loading") {
    return <LoadingPage/>;
  }

  return (
    <>
      {
        !!errorMessage
        ? <EuiCallOut
            title="Sorry, there was an error logging you in"
            color="danger"
            iconType="alert"
          >
            <p>{errorMessage}</p>
          </EuiCallOut>
        : null
      }

      <EuiPageTemplate
        template="centeredBody"
        pageContentProps={{ paddingSize: 'l' }}
        minHeight="80vh"
      >
        <EuiEmptyPrompt
          title={<span>Login into Financials</span>}
          body={<LoginForm/>}
          titleSize="m"
        />
      </EuiPageTemplate>
    </>
  );
}
