package models

import (
	"encoding/json"
	"time"
)

type MonthYear time.Time

const layoutMonthYear = "01-2006"

func (m *MonthYear) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse(layoutMonthYear, s)
	if err != nil {
		return err
	}
	*m = MonthYear(t)
	return nil
}

func (m MonthYear) MarshalJSON() ([]byte, error) {
	t := time.Time(m)
	return json.Marshal(t.Format(layoutMonthYear))
}

func (m MonthYear) ToTime() time.Time {
	return time.Time(m)
}

func FromTime(t time.Time) MonthYear {
	return MonthYear(t)
}
