// Copyright (c) 2017 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package duration

import (
	"math"
	"time"

	"github.com/miniflux/miniflux/locale"
)

// Texts to be translated if necessary.
var (
	NotYet     = `not yet`
	JustNow    = `just now`
	LastMinute = `1 minute ago`
	Minutes    = `%d minutes ago`
	LastHour   = `1 hour ago`
	Hours      = `%d hours ago`
	Yesterday  = `yesterday`
	Days       = `%d days ago`
	Weeks      = `%d weeks ago`
	Months     = `%d months ago`
	Years      = `%d years ago`
)

// ElapsedTime returns in a human readable format the elapsed time
// since the given datetime.
func ElapsedTime(translator *locale.Language, timezone string, t time.Time) string {
	if t.IsZero() {
		return translator.Get(NotYet)
	}

	var now time.Time
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		now = time.Now()
	} else {
		now = time.Now().In(loc)

		// The provided date is already converted to the user timezone by Postgres,
		// but the timezone information is not set in the time struct.
		// We cannot use time.In() because the date will be converted a second time.
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
	}

	if now.Before(t) {
		return translator.Get(NotYet)
	}

	diff := now.Sub(t)
	// Duration in seconds
	s := diff.Seconds()
	// Duration in days
	d := int(s / 86400)
	switch {
	case s < 60:
		return translator.Get(JustNow)
	case s < 120:
		return translator.Get(LastMinute)
	case s < 3600:
		return translator.Get(Minutes, int(diff.Minutes()))
	case s < 7200:
		return translator.Get(LastHour)
	case s < 86400:
		return translator.Get(Hours, int(diff.Hours()))
	case d == 1:
		return translator.Get(Yesterday)
	case d < 7:
		return translator.Get(Days, d)
	case d < 31:
		return translator.Get(Weeks, int(math.Ceil(float64(d)/7)))
	case d < 365:
		return translator.Get(Months, int(math.Ceil(float64(d)/30)))
	default:
		return translator.Get(Years, int(math.Ceil(float64(d)/365)))
	}
}
