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
} from '@elastic/eui';

import { useHistory } from "react-router-dom";

import { 
  useDispatch,
  useSelector,
} from 'react-redux';

import {
  updateUsername,
  updatePassword,
  createAsync,
} from '../redux/createAccountSlice';

function LoginForm() {
  const {
    username,
    password,
    isLoggedIn
  } = useSelector((state) => state.createAccount)
  const user = {
    username,
    password,
  }
  const dispatch = useDispatch()
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

      <EuiButton
        type="submit"
        fill
        onClick={
          (e) => {
            e.preventDefault();
            dispatch(createAsync(user));
          }
        }
      >
        Create
      </EuiButton>

    </EuiForm>
  );
}

function LoginPage() {
  const {
    errorMessage,
  } = useSelector((state) => state.login)
  return (
    <>
      {
        !!errorMessage
        ? <EuiCallOut
            title="Sorry, there was an error creating an account for you."
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

export default LoginPage;
