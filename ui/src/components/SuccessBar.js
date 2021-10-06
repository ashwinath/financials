import React from 'react';
import {
  EuiCallOut,
} from '@elastic/eui';

export function SuccessBar({ title, message }) {
  if (!message) {
    return null;
  }

  return (
    <EuiCallOut
      title={title}
      color="success"
      iconType="user"
    >
      <p>{message}</p>
    </EuiCallOut>
  );
}

