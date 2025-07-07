package graph

import (
	"compensation-api/graph/model"
	"compensation-api/internal/elastic"
)

func elasticToModelCompensation(c *elastic.Compensation) *model.Compensation {
	if c == nil {
		return nil
	}
	return &model.Compensation{
		ID:            c.ID,
		Timestamp:     &c.Timestamp,
		AgeRange:      &c.AgeRange,
		Industry:      &c.Industry,
		JobTitle:      &c.JobTitle,
		AnnualSalary:  &c.AnnualSalary,
		Currency:      &c.Currency,
		Location:      &c.Location,
		Experience:    &c.Experience,
		JobContext:    &c.JobContext,
		OtherCurrency: &c.OtherCurrency,
	}
}
