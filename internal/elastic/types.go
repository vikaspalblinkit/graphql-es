package elastic

// Compensation represents an employee's compensation details.
type Compensation struct {
	ID            string  `json:"id,omitempty"`
	Timestamp     string  `json:"timestamp"`
	AgeRange      string  `json:"age_range"`
	Industry      string  `json:"industry"`
	JobTitle      string  `json:"job_title"`
	AnnualSalary  float64 `json:"annual_salary"`
	Currency      string  `json:"currency"`
	Location      string  `json:"location"`
	Experience    string  `json:"experience"`
	JobContext    string  `json:"job_context"`
	OtherCurrency string  `json:"other_currency"`
}
