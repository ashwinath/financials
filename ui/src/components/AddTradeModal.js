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
  EuiText,
  EuiTextArea,
  EuiSpacer,
  EuiCode,
} from '@elastic/eui';

import { useDispatch, useSelector } from 'react-redux';
import {
  toggleIsAddTradeModalOpen,
  toggleIsAddBulkTradeModalOpen,
  updateAddTradeForm,
  submitTrade,
  submitTrades,
  updateCSVTrades,
} from '../redux/investmentsSlice';

import {
  convertDateToString,
  convertStringToDate,
  formatMoney,
  marshalCSV,
  validateTradeCSV,
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

  const isReadyToSubmit = symbol !== "" && !!date_purchased && quantity > 0 && price_each > 0 && trade_type !== "";
  let saveButton = (
    <EuiButton
      type="submit"
      form="modalFormId"
      onClick={() => dispatch(submitTrade(addTradeForm))}
      isDisabled={!isReadyToSubmit}
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

export function AddBulkTradeModal() {
  const dispatch = useDispatch();
  const { isTradeFormSubmitting, tradeCSVRaw } = useSelector((state) => state.investments);

  const isValidCSV = validateTradeCSV(tradeCSVRaw);
  let saveButton = (
    <EuiButton
      type="submit"
      form="modalFormId"
      onClick={() => dispatch(submitTrades(marshalCSV(tradeCSVRaw)))}
      isDisabled={!isValidCSV}
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
      onClose={() => dispatch(toggleIsAddBulkTradeModalOpen())}
      initialFocus="[name=popswitch]"
    >

      <EuiModalHeader>
        <EuiModalHeaderTitle>
          <h1>Add bulk trade</h1>
        </EuiModalHeaderTitle>
      </EuiModalHeader>

      <EuiModalBody>
        <EuiText>
          Input your data in <EuiCode>csv</EuiCode> format, with the following headers:
        </EuiText>
        <EuiSpacer size="s" />
        <EuiText>
          <EuiCode>date_purchased</EuiCode>: Format in <EuiCode>yyyy-mm-dd</EuiCode>.
        </EuiText>
        <EuiText>
          <EuiCode>symbol</EuiCode>: Using Alpha Vantage symbols.
        </EuiText>
        <EuiText>
          <EuiCode>trade_type</EuiCode>: <EuiCode>buy</EuiCode> or <EuiCode>sell</EuiCode>.
        </EuiText>
        <EuiText>
          <EuiCode>price_each</EuiCode>: number
        </EuiText>
        <EuiText>
          <EuiCode>quantity</EuiCode>: number
        </EuiText>
        <EuiSpacer size="s" />
        <EuiTextArea
          placeholder=""
          fullWidth={true}
          value={tradeCSVRaw}
          onChange={(e) => dispatch(updateCSVTrades(e.target.value))}
        />
      </EuiModalBody>
      <EuiModalFooter>
        <EuiButtonEmpty
          onClick={() => dispatch(toggleIsAddBulkTradeModalOpen())}
        >
          Cancel
        </EuiButtonEmpty>
        {saveButton}
      </EuiModalFooter>
    </EuiModal>
  );
}
