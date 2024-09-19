package eventcallback

import (
	"bufio"
	"errors"
	"os"
	"seatalk-bot/internal/constants"
	"sort"
	"strings"
	"time"
)

// Define constants

type Schedule struct {
	PIC   string    // Name of the Person In Charge
	Date  time.Time // Date of the schedule
	Email string    // Email of the Person In Charge
}

func ParseDate(dateStr string) (time.Time, error) {
	date, err := time.Parse(constants.DateFormat, dateStr)
	if err != nil {
		return time.Time{}, errors.New(constants.ErrorDateParse + ": " + err.Error())
	}
	return date, nil
}

// Read schedules from file
func ReadSchedules(filename string) ([]Schedule, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New(constants.ErrorFileOpen + ": " + err.Error())
	}
	defer file.Close()

	var schedules []Schedule
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 { // Ensure there are exactly 3 parts
			continue
		}

		pic := strings.TrimSpace(parts[0])
		dateStr := strings.TrimSpace(parts[1])
		email := strings.TrimSpace(parts[2])
		date, err := ParseDate(dateStr)
		if err != nil {
			return nil, errors.New(constants.ErrorInvalidDateFormat + ": " + pic)
		}

		schedules = append(schedules, Schedule{
			PIC:   pic,
			Date:  date,
			Email: email,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.New(constants.ErrorFileRead + ": " + err.Error())
	}

	return schedules, nil
}

// UpdatePreviousPICDate updates the date of the previous PIC
func UpdatePreviousPICDate(filename, currentPIC string) error {
	schedules, err := ReadSchedules(filename)
	if err != nil {
		return err
	}

	var currentPICDate, lastDate time.Time
	var previousPIC *Schedule

	for _, schedule := range schedules {
		if schedule.PIC == currentPIC {
			currentPICDate = schedule.Date
		}
		if schedule.Date.After(lastDate) {
			lastDate = schedule.Date
		}
	}

	if currentPICDate.IsZero() {
		return errors.New(constants.ErrorPICNotFound)
	}

	for i := len(schedules) - 1; i >= 0; i-- {
		if schedules[i].Date.Before(currentPICDate) {
			previousPIC = &schedules[i]
			break
		}
	}

	if previousPIC == nil {
		return errors.New(constants.ErrorPreviousPICNotFound)
	}

	// Update the date of the previous PIC
	previousPIC.Date = lastDate.Add(7 * 24 * time.Hour)

	// Sort schedules by date
	sort.Slice(schedules, func(i, j int) bool {
		return schedules[i].Date.Before(schedules[j].Date)
	})

	// Write the updated schedules back to the file
	file, err := os.Create(filename)
	if err != nil {
		return errors.New(constants.ErrorFileCreate + ": " + err.Error())
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, schedule := range schedules {
		_, err := writer.WriteString(strings.Join([]string{
			schedule.PIC,
			schedule.Date.Format(constants.DateFormat),
			schedule.Email,
		}, ",") + "\n")
		if err != nil {
			return errors.New(constants.ErrorWriteSchedule + ": " + schedule.PIC + ": " + err.Error())
		}
	}

	return writer.Flush()
}

// Get the start and end date of the current week
func GetCurrentWeekRange() (time.Time, time.Time) {
	loc := time.FixedZone("Asia/Jakarta", 7*60*60)
	now := time.Now().In(loc)

	// Calculate the start of the current week (Tuesday)
	offset := int(time.Tuesday) - int(now.Weekday())
	if offset > 0 {
		offset -= 7
	}
	startOfWeek := now.AddDate(0, 0, offset).Truncate(24 * time.Hour).In(loc)

	// Calculate the end of the current week (Monday of the next week)
	endOfWeek := startOfWeek.AddDate(0, 0, 6).Truncate(24 * time.Hour).In(loc)
	endOfWeek = endOfWeek.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return startOfWeek, endOfWeek
}

// Display PICs within a date range
func DisplayPICsWithinRange(schedules []Schedule, startDate, endDate time.Time) string {
	var result strings.Builder

	for _, schedule := range schedules {
		if (schedule.Date.After(startDate) || schedule.Date.Equal(startDate)) &&
			(schedule.Date.Before(endDate) || schedule.Date.Equal(endDate)) {
			result.WriteString(strings.Join([]string{
				"Date: " + schedule.Date.Format("2006-01-02"),
				"- PIC: " + schedule.PIC,
				" <mention-tag target=\"seatalk://user?email=" + schedule.Email + "\"/>",
			}, "") + "\n")
			if err := UpdatePreviousPICDate(constants.StockInventoryScheduleFile, schedule.PIC); err != nil {
				// Optionally log the error or handle it as needed
			}
		}
	}
	return result.String()
}

// Display the full schedule
func DisplayFullSchedule(schedules []Schedule) string {
	var result strings.Builder

	// Sort schedules by date
	sort.Slice(schedules, func(i, j int) bool {
		return schedules[i].Date.Before(schedules[j].Date)
	})

	result.WriteString("\nStock Inventory Schedule for Following Weeks:\n")
	for _, schedule := range schedules {
		result.WriteString("Date: " + schedule.Date.Format("2006-01-02") + " - PIC: " + schedule.PIC + "\n")
	}

	return result.String()
}
