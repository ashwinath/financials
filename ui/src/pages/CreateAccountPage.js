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
} from '@elastic/eui';

import { useHistory } from "react-router-dom";
import { useDispatch, useSelector } from 'react-redux';
import { updateUsername, updatePassword, createAsync } from '../redux/createAccountSlice';
import { LoadingPage } from ".";
import { ErrorBar } from "../components"
import { useLoginHook } from "../hooks/login";

function CreateAccountForm() {
  const {
    username,
    password,
    isLoggedIn,
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

export function CreateAccountPage() {
  const {
    errorMessage,
    status,
  } = useSelector((state) => state.createAccount);
  useLoginHook();

  if (status === "loading") {
    return <LoadingPage/>;
  }

  return (
    <>
      <ErrorBar 
        title="Sorry, there was an error creating an account for you."
        errorMessage={errorMessage}
      />

      <EuiPageTemplate
        template="centeredBody"
        pageContentProps={{ paddingSize: 'l' }}
        minHeight="80vh"
      >
        <EuiEmptyPrompt
          title={<span>Create an account</span>}
          body={<CreateAccountForm/>}
          titleSize="m"
        />
      </EuiPageTemplate>
    </>
  );
}
