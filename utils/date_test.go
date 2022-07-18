package utils

import (
	"testing"
	"time"
)

func TestFromMonth(t *testing.T) {

	month := "2020-02"
	s, e, err := DatesFromMonth(month)
	if err != nil {
		t.Fatalf("Expected null error, got: %s", err)
	}

	strStr := month + "-01"
	expectedStartDate, _ := time.Parse(timeDayFormat, strStr)
	if s != expectedStartDate {
		t.Fatalf("expected start date: %s, got: %s", expectedStartDate, s)
	}

	endStr := month + "-29"
	expectedEndDate, _ := time.Parse(timeDayFormat, endStr)
	if e != expectedEndDate {
		t.Fatalf("expected end date: %s, got: %s", expectedEndDate, e)
	}
}

func TestFromDays(t *testing.T) {

	startStr := "2022-04-15"
	endStr := "2022-05-09"

	s, e, err := DatesFromDays(startStr, endStr)
	if err != nil {
		t.Fatalf("Expected null error, got: %s", err)
	}

	expectedStartDate, _ := time.Parse(timeDayFormat, startStr)
	if s != expectedStartDate {
		t.Fatalf("expected start date: %s, got: %s", expectedStartDate, s)
	}

	expectedEndDate, _ := time.Parse(timeDayFormat, endStr)
	if e != expectedEndDate {
		t.Fatalf("expected end date: %s, got: %s", expectedEndDate, e)
	}
}
