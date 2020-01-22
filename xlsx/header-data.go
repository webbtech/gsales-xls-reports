package xlsx

// BankCardJSON function
func BankCardJSON() string {
	return `[
		{
			"label": "Station",
			"width": 22.00
		},
		{
			"label": "Record Number",
			"width": 15.00
		},
		{
			"label": "AMEX",
			"width": 15.00
		},
		{
			"label": "Discover",
			"width": 15.00
		},
		{
			"label": "Gales",
			"width": 15.00
		},
		{
			"label": "MC",
			"width": 15.00
		},
		{
			"label": "Visa",
			"width": 15.00
		},
		{
			"label": "Debit",
			"width": 15.00
		},
		{
			"label": "Other",
			"width": 15.00
		}
	]`
}

// EmployeeOSJSON function
func EmployeeOSJSON() string {
	return `[
		{
			"label": "Employee",
			"width": 20.00
		},
		{
			"label": "Record Number",
			"width": 15.00
		},
		{
			"label": "Station",
			"width": 22.00
		},
		{
			"label": "Shift Overshort",
			"width": 14.00
		},
		{
			"label": "Attendant Overshort",
			"width": 18.00
		},
		{
			"label": "Overshort Diff",
			"width": 13.00
		},
		{
			"label": "Discrepancy Description",
			"width": 35.00
		}
	]`
}

// MonthlySalesJSON function
func MonthlySalesJSON() string {
	return `[
		{
			"label": "Station",
			"width": 22.00
		},
		{
			"label": "Record Number",
			"width": 15.00
		},
		{
			"label": "Employee",
			"width": 20.00
		},
		{
			"label": "Shift Overshort",
			"width": 14.5
		},
		{
			"label": "Fuel Sales",
			"width": 10.5
		},
		{
			"label": "Fuel Sales HST",
			"width": 14.5
		},
		{
			"label": "Fuel Sales Total",
			"width": 14.75
		},
		{
			"label": "Other Fuel Sales",
			"width": 15.5
		},
		{
			"label": "Gift Certificates",
			"width": 15.0
		},
		{
			"label": "Bobs Sales",
			"width": 11.5
		},
		{
			"label": "Bobs Gift Certs.",
			"width": 15.0
		},
		{
			"label": "Bobs Fuel Adjust.",
			"width": 16.5
		},
		{
			"label": "Total Non-Fuel",
			"width": 14.0
		},
		{
			"label": "Cash Bills",
			"width": 11.0
		},
		{
			"label": "Cash Debit",
			"width": 11.0
		},
		{
			"label": "Cash Diesel Discount",
			"width": 20.0
		},
		{
			"label": "Cash Drive Off NSF",
			"width": 18.0
		},
		{
			"label": "Gales Loyalty Redeem",
			"width": 21.0
		},
		{
			"label": "Cash Gift Cert Redeem",
			"width": 21.0
		},
		{
			"label": "Cash Lottery Payout",
			"width": 19.0
		},
		{
			"label": "Cash Other",
			"width": 11.5
		},
		{
			"label": "Cash OS Adjusted",
			"width": 17.0
		},
		{
			"label": "Cash Payout",
			"width": 12.0
		},
		{
			"label": "Cash Write Off",
			"width": 14.0
		},
		{
			"label": "Bank Amex",
			"width": 11.0
		},
		{
			"label": "Bank Discover",
			"width": 14.0
		},
		{
			"label": "Bank Gales",
			"width": 11.5
		},
		{
			"label": "Bank MC",
			"width": 10.0
		},
		{
			"label": "Bank VISA",
			"width": 11.5
		},
		{
			"label": "Cigarette Sales",
			"width": 14.5
		},
		{
			"label": "Cigarette Qty",
			"width": 12.5
		},
		{
			"label": "Oil Sales",
			"width": 9.0
		},
		{
			"label": "Oil Qty",
			"width": 7.5
		},
		{
			"label": "Carwash Qty",
			"width": 12.5
		},
		{
			"label": "Loyalty Funded Qty",
			"width": 16.5
		}
	]`
}

// PayPeriodJSON function
func PayPeriodJSON() string {
	return `[
		{
			"label": "Employee",
			"width": 20.00
		},
		{
			"label": "Record Number",
			"width": 15.00
		},
		{
			"label": "Station",
			"width": 22.00
		},
		{
			"label": "Shift Overshort",
			"width": 15.00
		},
		{
			"label": "Total NonFuel Sales",
			"width": 19.00
		},
		{
			"label": "Product Sales",
			"width": 14.00
		},
		{
			"label": "Commission Eligible Qty",
			"width": 21.00
		},
		{
			"label": "Commission Eligible Sales",
			"width": 21.00
		},
		{
			"label": "Commission Amount",
			"width": 18.00
		},
		{
			"label": "Number Car Washes",
			"width": 19.00
		},
		{
			"label": "Attendant Adjustment",
			"width": 21.00
		}
	]`
}

// ProductNumbersJSON function
func ProductNumbersJSON() string {
	return `[
		{
			"label": "Product",
			"width": 20.00
		},
		{
			"label": "Quantity",
			"width": 15.00
		}
	]`
}
