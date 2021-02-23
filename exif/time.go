package exif

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta/exif/ifds"
	"github.com/evanoberholster/imagemeta/exif/ifds/exififd"
	"github.com/evanoberholster/imagemeta/exif/tag"
)

// Errors for Parsing of Time
var (
	ErrParseYear  = fmt.Errorf("error parsing Year")
	ErrParseMonth = fmt.Errorf("error parsing Month")
	ErrParseDay   = fmt.Errorf("error parsing Day")
	ErrParseHour  = fmt.Errorf("error parsing Hour")
	ErrParseMin   = fmt.Errorf("error parsing Min")
)

// ModifyDate - the date and time at which the Exif file was modified
func (e *Data) ModifyDate() (time.Time, error) {
	// "IFD" DateTime
	// "IFD/Exif" SubSecTime
	return e.getDateTags(ifds.RootIFD, ifds.DateTime, ifds.ExifIFD, exififd.SubSecTime)
}

// DateTime - the date and time at which the EXIF file was created
// with sub-second precision
func (e *Data) DateTime() (time.Time, error) {
	// "IFD/Exif" DateTimeOriginal
	// "IFD/Exif" SubSecTimeOriginal
	// TODO: "IFD/Exif" OffsetTimeOriginal
	if t, err := e.getDateTags(ifds.ExifIFD, exififd.DateTimeOriginal, ifds.ExifIFD, exififd.SubSecTimeOriginal); err == nil {
		return t, err
	}

	// "IFD/Exif" DateTimeDigitized
	// "IFD/Exif" SubSecTimeDigitized
	// TODO: "IFD/Exif" OffsetTimeDigitized
	if t, err := e.getDateTags(ifds.ExifIFD, exififd.DateTimeDigitized, ifds.ExifIFD, exififd.SubSecTimeDigitized); err == nil {
		return t, err
	}
	return time.Time{}, ErrEmptyTag
}

func (e *Data) getDateTags(dateIFD ifds.IFD, dateTagID tag.ID, subSecIFD ifds.IFD, subSecTagID tag.ID) (time.Time, error) {
	// "IFD" DateTime
	t, err := e.GetTag(dateIFD, 0, dateTagID)
	if err != nil {
		return time.Time{}, ErrEmptyTag
	}
	if dateRaw, err := t.ASCIIValue(e.exifReader); err == nil && dateRaw != "" {
		var subSecRaw string
		// "IFD/Exif" SubSecTime
		if t, err := e.GetTag(subSecIFD, 0, subSecTagID); err != nil {
			subSecRaw, _ = t.ASCIIValue(e.exifReader)
		}
		if dateTime, err := parseExifFullTimestamp(dateRaw, subSecRaw); err == nil && !dateTime.IsZero() {
			return dateTime, nil
		}
	}
	return time.Time{}, ErrEmptyTag
}

// parseExifFullTimestamp parses dates like "2018:11:30 13:01:49" into a UTC
// `time.Time` struct.
func parseExifFullTimestamp(fullTimestampPhrase string, subSecString string) (timestamp time.Time, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	parts := strings.Split(fullTimestampPhrase, " ")
	datestampValue, timestampValue := parts[0], parts[1]

	dateParts := strings.Split(datestampValue, ":")

	year, err := strconv.ParseUint(dateParts[0], 10, 16)
	if err != nil {
		err = ErrParseYear
		return
	}

	month, err := strconv.ParseUint(dateParts[1], 10, 8)
	if err != nil {
		err = ErrParseMonth
		return
	}

	day, err := strconv.ParseUint(dateParts[2], 10, 8)
	if err != nil {
		err = ErrParseDay
		return
	}

	timeParts := strings.Split(timestampValue, ":")

	hour, err := strconv.ParseUint(timeParts[0], 10, 8)
	if err != nil {
		err = ErrParseHour
		return
	}

	minute, err := strconv.ParseUint(timeParts[1], 10, 8)
	if err != nil {
		err = ErrParseMin
		return
	}

	second, err := strconv.ParseUint(timeParts[2], 10, 8)
	if err != nil {
		err = ErrParseMin
		return
	}

	subSec, err := strconv.ParseUint(subSecString, 10, 16)
	if err != nil {
		subSec = 0
	}

	timestamp = time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), int(subSec*1000000), time.UTC)
	return timestamp, nil
}

// ParseTimestamp parses dates like "2018:11:30" into a UTC `time.Time` struct.
func parseTimestamp(dateStamp string, hour, min, sec int) (timestamp time.Time, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	dateParts := strings.Split(dateStamp, ":")

	year, err := strconv.ParseUint(dateParts[0], 10, 16)
	if err != nil {
		err = ErrParseYear
		return
	}

	month, err := strconv.ParseUint(dateParts[1], 10, 8)
	if err != nil {
		err = ErrParseMonth
		return
	}

	day, err := strconv.ParseUint(dateParts[2], 10, 8)
	if err != nil {
		err = ErrParseDay
		return
	}

	timestamp = time.Date(int(year), time.Month(month), int(day), hour, min, sec, 0, time.UTC)
	return timestamp, nil
}
