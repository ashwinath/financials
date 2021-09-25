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

import { useDispatch } from 'react-redux';
import { updateUsername, updatePassword } from '../redux/loginSlice';

function LoginForm() {
  const dispatch = useDispatch()

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

      <EuiButton type="submit" fill>
        Login
      </EuiButton>

    </EuiForm>
  );
}

function LoginPage() {
  return (
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
  );
}

export default LoginPage;
