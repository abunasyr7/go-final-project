package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"


func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date %q: %w", dstart, err)
	}

	if repeat == "" {
		return "", fmt.Errorf("repeat rule is empty")
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		return nextDateDay(date, now, parts)
	case "y":
		return nextDateYear(date, now), nil
	case "w":
		return nextDateWeek(date, now, parts)
	case "m":
		return nextDateMonth(date, now, parts)
	default:
		return "", fmt.Errorf("unsupported repeat rule %q", repeat)
	}
}


func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	d := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	n := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
	return d.After(n)
}

func nextDateDay(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("day interval is not specified")
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil || days <= 0 || days > 400 {
		return "", fmt.Errorf("invalid day interval %q: must be from 1 to 400", parts[1])
	}

	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func nextDateYear(date, now time.Time) string {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat)
}

func nextDateWeek(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("weekdays are not specified")
	}

	var weekdays [8]bool 
	for _, s := range strings.Split(parts[1], ",") {
		d, err := strconv.Atoi(s)
		if err != nil || d < 1 || d > 7 {
			return "", fmt.Errorf("invalid weekday %q: must be from 1 to 7", s)
		}
		weekdays[d] = true
	}

	for {
		date = date.AddDate(0, 0, 1)
		wd := int(date.Weekday())
		if wd == 0 {
			wd = 7
		}
		if weekdays[wd] && afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func nextDateMonth(date, now time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("invalid month rule format")
	}

	days := map[int]bool{}
	for _, s := range strings.Split(parts[1], ",") {
		d, err := strconv.Atoi(s)
		if err != nil || d == 0 || d < -2 || d > 31 {
			return "", fmt.Errorf("invalid day of month %q: must be from 1 to 31, -1 or -2", s)
		}
		days[d] = true
	}

	var months map[int]bool
	if len(parts) == 3 {
		months = map[int]bool{}
		for _, s := range strings.Split(parts[2], ",") {
			m, err := strconv.Atoi(s)
			if err != nil || m < 1 || m > 12 {
				return "", fmt.Errorf("invalid month %q: must be from 1 to 12", s)
			}
			months[m] = true
		}
	}

	matches := func(d time.Time) bool {
		if months != nil && !months[int(d.Month())] {
			return false
		}
		if days[d.Day()] {
			return true
		}
		lastDay := lastDayOfMonth(d)
		if days[-1] && d.Day() == lastDay {
			return true
		}
		if days[-2] && d.Day() == lastDay-1 {
			return true
		}
		return false
	}

	limit := date
	if afterNow(limit, now) {
		limit = now
	}
	limit = limit.AddDate(2, 0, 0)

	for {
		date = date.AddDate(0, 0, 1)
		if date.After(limit) {
			return "", fmt.Errorf("no matching date found for rule %q", strings.Join(parts, " "))
		}
		if matches(date) && afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func lastDayOfMonth(d time.Time) int {
	firstOfNext := time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location()).AddDate(0, 1, 0)
	return firstOfNext.AddDate(0, 0, -1).Day()
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	now := time.Now()
	if nowParam != "" {
		parsed, err := time.Parse(dateFormat, nowParam)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid now %q: %v", nowParam, err), http.StatusBadRequest)
			return
		}
		now = parsed
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
