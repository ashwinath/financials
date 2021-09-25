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

function LoginForm() {
  return (
    <EuiForm component="form">
      <EuiFormRow label="Username">
        <EuiFieldText name="username" />
      </EuiFormRow>

      <EuiFormRow label="Password">
        <EuiFieldPassword name="password" />
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
