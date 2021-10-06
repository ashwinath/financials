import React from 'react';
import {
  EuiModal,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiModalBody,
  EuiModalFooter,
  EuiButtonEmpty,
  EuiButton,
  EuiForm,
  EuiFormRow,
  EuiFieldText,
  EuiDatePicker,
  EuiRadioGroup,
} from '@elastic/eui';

import { useDispatch, useSelector } from 'react-redux';
import {
  toggleIsAddTradeModalOpen,
  updateAddTradeForm,
  submitTrade,
} from '../redux/investmentsSlice';

import {
  convertDateToString,
  convertStringToDate,
  formatMoney,
} from "../utils";

export function AddTradeModal() {
  const dispatch = useDispatch();
  const { addTradeForm, isTradeFormSubmitting } = useSelector((state) => state.investments);
  const {
    symbol,
    date_purchased,
    quantity,
    price_each,
    trade_type,
  } = addTradeForm;

  let saveButton = (
    <EuiButton
      type="submit"
      form="modalFormId"
      onClick={() => dispatch(submitTrade(addTradeForm))}
      fill
    >
      Save
    </EuiButton>
  );

  if (isTradeFormSubmitting) {
    saveButton = (
      <EuiButton isLoading={true}>Loading</EuiButton>
    );
  }

  return (
    <EuiModal
      onClose={() => dispatch(toggleIsAddTradeModalOpen())}
      initialFocus="[name=popswitch]"
    >

      <EuiModalHeader>
        <EuiModalHeaderTitle>
          <h1>Add trade</h1>
        </EuiModalHeaderTitle>
      </EuiModalHeader>

      <EuiModalBody>
        <EuiForm component="form">
          <EuiFormRow label="Date">
            <EuiDatePicker
              name="date_purchased"
              selected={convertStringToDate(date_purchased)}
              onChange={(date) => dispatch(updateAddTradeForm({
                name: "date_purchased",
                value: convertDateToString(date),
              }))}
            />
          </EuiFormRow>

          <EuiFormRow label="Symbol">
            <EuiFieldText
              name="symbol"
              value={symbol}
              onChange={(e) => dispatch(updateAddTradeForm({
                name: "symbol",
                value: e.target.value,
              }))}
            />
          </EuiFormRow>

          <EuiFormRow label="Trade Type">
            <EuiRadioGroup
              options={[
                {
                  id: "buy",
                  label: "Buy",
                },
                {
                  id: "sell",
                  label: "Sell",
                },
              ]}
              idSelected={trade_type}
              onChange={(id) => dispatch(updateAddTradeForm({
                name: "trade_type",
                value: id,
              })) }
              name="Trade Type"
            />
          </EuiFormRow>

          <EuiFormRow label="Price">
            <EuiFieldText
              name="price_each"
              value={price_each}
              onChange={(e) => dispatch(updateAddTradeForm({
                name: "price_each",
                value: parseFloat(e.target.value),
              }))}
            />
          </EuiFormRow>

          <EuiFormRow label="Quantity">
            <EuiFieldText
              name="quantity"
              value={quantity}
              onChange={(e) => dispatch(updateAddTradeForm({
                name: "quantity",
                value: parseFloat(e.target.value),
              }))}
            />
          </EuiFormRow>

          <EuiFormRow label="Total">
            <EuiFieldText
              name="quantity"
              value={formatMoney(quantity * price_each)}
              disabled
            />
          </EuiFormRow>
        </EuiForm>
      </EuiModalBody>

      <EuiModalFooter>
        <EuiButtonEmpty
          onClick={() => dispatch(toggleIsAddTradeModalOpen())}
        >
          Cancel
        </EuiButtonEmpty>
        {saveButton}
      </EuiModalFooter>
    </EuiModal>
  );
}
