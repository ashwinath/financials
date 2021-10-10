// We assume this is for Singapore only
import moment from "moment";

export function convertDateToString(date) {
  if (date === null || date === undefined) {
    return null;
  }

  const mDate = date.utcOffset('+0800')
    .set({hour:16,minute:0,second:0,millisecond:0})

  return mDate.format();
}

export function convertStringToDate(dateString) {
  if (dateString === null || dateString === undefined) {
    return null;
  }

  return moment(dateString).utcOffset('+0800').set({hour:0,minute:0,second:0,millisecond:0});
}

export function getDateFromPeriod(numberOfMonthsFromNow) {
  if (numberOfMonthsFromNow === null || numberOfMonthsFromNow === undefined) {
    return null;
  }

  return moment()
    .utcOffset('+0800')
    .set({hour:16,minute:0,second:0,millisecond:0})
    .subtract(numberOfMonthsFromNow, "months")
    .format();
}
