import { convertDateToString, convertStringToDate } from ".";

export function capitaliseFirstLetter(string) {
  return string ? string.charAt(0).toUpperCase() + string.slice(1) : null;
}

export function capitaliseAll(string) {
  return string ? string.toUpperCase() : null;
}

export function formatMoney(number) {
	const numberString = number.toFixed(2).toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  return `$${numberString}`;
}

export function formatMoneyGraph(num, digits) {
  const lookup = [
    { value: 1, symbol: "" },
    { value: 1e3, symbol: "k" },
    { value: 1e6, symbol: "M" },
    { value: 1e9, symbol: "G" },
    { value: 1e12, symbol: "T" },
    { value: 1e15, symbol: "P" },
    { value: 1e18, symbol: "E" }
  ];
  const rx = /\.0+$|(\.[0-9]*[1-9])0+$/;
  var item = lookup.slice().reverse().find(function(item) {
    return num >= item.value;
  });
  return item ? (num / item.value).toFixed(digits).replace(rx, "$1") + item.symbol : "0";
}

export function validateTradeCSV(rawCSV) {
  if (rawCSV === "") {
    return false;
  }

  const arrays = CSVToArray(rawCSV)
  if (arrays.length < 2) {
    return false;
  }

  const datePurchasedIndex = arrays[0].indexOf("date_purchased")
  if (datePurchasedIndex === -1) {
    return false;
  }

  const symbolIndex = arrays[0].indexOf("symbol")
  if (symbolIndex === -1) {
    return false;
  }

  const tradeTypeIndex = arrays[0].indexOf("trade_type")
  if (tradeTypeIndex === -1) {
    return false;
  }

  const priceEachIndex = arrays[0].indexOf("price_each")
  if (priceEachIndex === -1) {
    return false;
  }

  const quantityIndex = arrays[0].indexOf("quantity")
  if (quantityIndex === -1) {
    return false;
  }

  for (let i = 1; i < arrays.length; ++i) {
    if (arrays[i].length !== 5) {
      return false;
    }
  }

  return true;
}

export function marshalCSV(rawCSV) {
  const arrays = CSVToArray(rawCSV)
  const datePurchasedIndex = arrays[0].indexOf("date_purchased")
  const symbolIndex = arrays[0].indexOf("symbol")
  const tradeTypeIndex = arrays[0].indexOf("trade_type")
  const priceEachIndex = arrays[0].indexOf("price_each")
  const quantityIndex = arrays[0].indexOf("quantity")

  const trades = []
  for (let i = 1; i < arrays.length; ++i) {
    trades.push({
      "date_purchased": convertDateToString(convertStringToDate(arrays[i][datePurchasedIndex])),
      "symbol": arrays[i][symbolIndex].trim(),
      "trade_type": arrays[i][tradeTypeIndex].trim(),
      "price_each": parseFloat(arrays[i][priceEachIndex]),
      "quantity": parseFloat(arrays[i][quantityIndex]),
    });
  }

  return trades;
}

function CSVToArray( strData, strDelimiter ){
  // Credits to: https://www.bennadel.com/blog/1504-ask-ben-parsing-csv-strings-with-javascript-exec-regular-expression-command.htm
  // Check to see if the delimiter is defined. If not,
  // then default to comma.
  strDelimiter = (strDelimiter || ",");
  // Create a regular expression to parse the CSV values.
  var objPattern = new RegExp(
    (
      // Delimiters.
      "(\\" + strDelimiter + "|\\r?\\n|\\r|^)" +
      // Quoted fields.
      "(?:\"([^\"]*(?:\"\"[^\"]*)*)\"|" +
      // Standard fields.
      "([^\"\\" + strDelimiter + "\\r\\n]*))"
    ),
    "gi"
    );
  // Create an array to hold our data. Give the array
  // a default empty first row.
  var arrData = [[]];
  // Create an array to hold our individual pattern
  // matching groups.
  var arrMatches = null;
  // Keep looping over the regular expression matches
  // until we can no longer find a match.
  // eslint-disable-next-line
  while (arrMatches = objPattern.exec( strData )){
    // Get the delimiter that was found.
    var strMatchedDelimiter = arrMatches[ 1 ];
    // Check to see if the given delimiter has a length
    // (is not the start of string) and if it matches
    // field delimiter. If id does not, then we know
    // that this delimiter is a row delimiter.
    if (
      strMatchedDelimiter.length &&
      (strMatchedDelimiter !== strDelimiter)
      ){
      // Since we have reached a new row of data,
      // add an empty row to our data array.
      arrData.push( [] );
    }
    // Now that we have our delimiter out of the way,
    // let's check to see which kind of value we
    // captured (quoted or unquoted).
    let strMatchedValue;
    if (arrMatches[ 2 ]){
      // We found a quoted value. When we capture
      // this value, unescape any double quotes.
      strMatchedValue = arrMatches[ 2 ].replace(
        new RegExp( "\"\"", "g" ),
        "\""
        );
    } else {
      // We found a non-quoted value.
      strMatchedValue = arrMatches[ 3 ];
    }
    // Now that we have our value string, let's add
    // it to the data array.
    arrData[ arrData.length - 1 ].push( strMatchedValue );
  }
  // Return the parsed data.
  return( arrData );
}

export function formatPercent(d) {
  return `${Number(d * 100).toFixed(2)}%`;

}
