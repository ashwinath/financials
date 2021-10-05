export function capitaliseFirstLetter(string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

export function capitaliseAll(string) {
  return string ? string.toUpperCase() : null;
}

export function formatMoney(number) {
  return `$${number.toFixed(2)}`;
}