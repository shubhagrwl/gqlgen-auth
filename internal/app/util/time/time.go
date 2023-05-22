package time

import (
	"time"
)

//ParseTime : func to parse various times in string format
func ParseTime(dateString string) (time.Time, error) {
	var parsedDate time.Time
	var err error
	parsedDate, err = time.Parse("2006-01-02T15:04:05.000Z", dateString)
	if err != nil {
		parsedDate, err = time.Parse("2006-01-02T15:04:05.000-07:00", dateString)
		if err != nil {
			parsedDate, err = time.Parse("2006-01-02T15:04:05Z", dateString)
			if err != nil {
				parsedDate, err = time.Parse("2006-01-02T15:04:05.9+02:00", dateString)
				if err != nil {
					parsedDate, err = time.Parse("2006-01-02T15:04:05.99+02:00", dateString)
					if err != nil {
						parsedDate, err = time.Parse("2006-01-02T15:04:05.999-07:00", dateString)
						if err != nil {
							parsedDate, err = time.Parse("2006-01-02T15:04:05.999999", dateString)
							if err != nil {
								parsedDate, err = time.Parse("2006-01-02T15:04:05-07:00", dateString)
								if err != nil {
									parsedDate, err = time.Parse("2006-01-02T15:04:05+07:00", dateString)
									if err != nil {
										parsedDate, err = time.Parse("2006-01-02 15:04:05.000+07", dateString)
										if err != nil {
											parsedDate, err = time.Parse("2006-01-02", dateString)
											if err != nil {
												parsedDate, err = time.Parse("2006-01-02 15:04:05.999999Z", dateString)
												if err != nil {
													parsedDate, err = time.Parse("2006-01-02 15:04:05.999999", dateString)
													if err != nil {
														parsedDate, err = time.Parse("2006/01/02", dateString)
														if err != nil {
															parsedDate, err = time.Parse("02/01/2006", dateString)
															if err != nil {
																return time.Time{}, err
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return parsedDate, nil
}

//ParseTime : func to parse various times in string format
func ParseTimeToISO(dateString string, locationString string) (string, error) {
	var parsedDate time.Time
	var err error
	loc, _ := time.LoadLocation(locationString)
	parsedDate, err = time.ParseInLocation("2006-01-02T03:04:05Z", dateString, loc)
	if err != nil {
		parsedDate, err = time.ParseInLocation("2006-01-02T03:04Z", dateString, loc)
		if err != nil {
			parsedDate, err = time.ParseInLocation("2006-01-02T3:04Z", dateString, loc)
			if err != nil {
				parsedDate, err = time.ParseInLocation("2006-01-02T15:04Z", dateString, loc)
				if err != nil {
					parsedDate, err = time.ParseInLocation("2006-01-02T15:04:05Z", dateString, loc)
					if err != nil {
						parsedDate, err = time.ParseInLocation("2006-01-02", dateString, loc)
						if err != nil {
							parsedDate, err = time.ParseInLocation("2006-01-02T15:04:05.000Z", dateString, loc)
							if err != nil {
								return "", err
							}
						}
					}
				}
			}
		}
	}
	return parsedDate.Format("2006-01-02T15:04:05-0700"), nil
}

func DateRange(week, month, year int) (startDate, endDate time.Time) {

	timeBenchmark := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	weekStartBenchmark := timeBenchmark.AddDate(0, 0, -(int(timeBenchmark.Weekday())+6)%7)

	startDate = weekStartBenchmark.AddDate(0, 0, (week-1)*7)
	endDate = startDate.AddDate(0, 0, 6)

	return startDate, endDate
}
