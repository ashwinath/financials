import React from 'react';
import {
  EuiCallOut,
} from '@elastic/eui';

export function ErrorBar({ title, errorMessage }) {
  if (!errorMessage) {
    return null;
  }

  return (
    <EuiCallOut
      title={title}
      color="danger"
      iconType="alert"
    >
      <p>{errorMessage}</p>
    </EuiCallOut>
  );
}
