const month = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

const handleDateSelectChange = () => {
  const expensePeriodElement = window.document.getElementById("expense-period");
  expensePeriodElement.innerHTML = month[0];
};
