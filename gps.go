package exiftool

import (
	"errors"
	"fmt"
	"time"

	"github.com/evanoberholster/exiftool/ifds"
	"github.com/evanoberholster/exiftool/ifds/gpsifd"
	"github.com/evanoberholster/exiftool/tag"
	"github.com/golang/geo/s2"
)

var (
	// ErrGpsCoordsNotValid means that some part of the geographic data were unparseable.
	ErrGpsCoordsNotValid = errors.New("error GPS coordinates not valid")
	// ErrGPSRationalNotValid means that the rawCoordinates were not long enough.
	ErrGPSRationalNotValid = errors.New("error GPS Coords requires a raw-coordinate with exactly three rationals")
)

// gpsCoordsFromRationals returns a decimal given the EXIF-encoded information.
// The refValue is the N/E/S/W direction that this position is relative to.
func gpsCoordsFromRationals(refValue string, rawCoordinate []tag.Rational) (decimal float64, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	if len(rawCoordinate) != 3 {
		err = ErrGPSRationalNotValid
		return
	}

	decimal = (float64(rawCoordinate[0].Numerator) / float64(rawCoordinate[0].Denominator))
	decimal += (float64(rawCoordinate[1].Numerator) / float64(rawCoordinate[1].Denominator) / 60.0)
	decimal += (float64(rawCoordinate[2].Numerator) / float64(rawCoordinate[2].Denominator) / 3600.0)

	// Decimal is a negative value for a South or West Orientation
	if refValue[0] == 'S' || refValue[0] == 'W' {
		decimal = -decimal
	}

	return
}

// GpsInfo encapsulates all of the geographic information in one place.
type GpsInfo struct {
	Latitude, Longitude float64
	Altitude            int
	Timestamp           time.Time
}

// String returns a descriptive string.
func (gi *GpsInfo) String() string {
	return fmt.Sprintf("GpsInfo | LAT=(%.05f) LON=(%.05f) ALT=(%d) TIME=[%s] |",
		gi.Latitude, gi.Longitude, gi.Altitude, gi.Timestamp)
}

// S2CellID returns the cell-ID of the geographic location on the earth.
func (gi *GpsInfo) S2CellID() (cellID s2.CellID, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()

	latLng := s2.LatLngFromDegrees(gi.Latitude, gi.Longitude)
	cellID = s2.CellIDFromLatLng(latLng)

	if !cellID.IsValid() {
		err = ErrGpsCoordsNotValid
		return
	}

	return cellID, nil
}

// GPSAltitude convenience func. "IFD/GPS" GPSAltitude
// WIP
func (e *ExifData) GPSAltitude() (alt float32, err error) {
	// Altitude
	t, err := e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSAltitude)
	if err == nil {
		rats, err := t.RationalValues(e.exifReader)
		if err == nil && len(rats) == 1 {
			alt = float32(rats[0].Numerator) / float32(rats[0].Denominator)
			return alt, nil
		}
	}

	// Altitude Ref - GPSAltitudeRef uint16
	return 0.0, err
}

// GPSCellID convenience func. retrieves "IFD/GPS" GPSLatitude and GPSLongitude
// converts them into an S2 CellID and returns the CellID.
//
// If the CellID is not valid it returns ErrGpsCoordsNotValid.
func (e *ExifData) GPSCellID() (cellID s2.CellID, err error) {
	lat, lng, err := e.GPSInfo()
	if err != nil {
		return
	}

	latLng := s2.LatLngFromDegrees(lat, lng)
	cellID = s2.CellIDFromLatLng(latLng)

	if !cellID.IsValid() {
		err = ErrGpsCoordsNotValid
		return
	}

	return cellID, nil
}

// GPSInfo convenience func. "IFD/GPS" GPSLatitude and GPSLongitude
func (e *ExifData) GPSInfo() (lat, lng float64, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()
	var ref string
	var raw []tag.Rational

	// Latitude
	t, err := e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSLatitudeRef)
	if err != nil {
		return
	}
	ref, err = t.ASCIIValue(e.exifReader)
	if err != nil {
		return
	}

	t, err = e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSLatitude)
	if err != nil {
		return
	}
	raw, err = t.RationalValues(e.exifReader)
	if err != nil {
		return
	}

	lat, err = gpsCoordsFromRationals(ref, raw)
	if err != nil {
		return
	}

	// Longitude
	t, err = e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSLongitudeRef)
	if err != nil {
		return
	}
	ref, err = t.ASCIIValue(e.exifReader)
	if err != nil {
		return
	}

	t, err = e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSLongitude)
	if err != nil {
		return
	}
	raw, err = t.RationalValues(e.exifReader)
	if err != nil {
		return
	}

	lng, err = gpsCoordsFromRationals(ref, raw)
	if err != nil {
		return
	}

	return
}

// GPSTime convenience func. "IFD/GPS" GPSDateStamp and GPSTimeStamp
func (e *ExifData) GPSTime() (timestamp time.Time, err error) {
	t, err := e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSDateStamp)
	if err != nil {
		return
	}
	dateRaw, err := t.ASCIIValue(e.exifReader)
	if err != nil {
		return
	}
	t, err = e.GetTag(ifds.GPSIFD, 0, gpsifd.GPSTimeStamp)
	if err != nil {
		return
	}
	timeRaw, err := t.RationalValues(e.exifReader)
	if err != nil {
		return
	}
	hour := int(timeRaw[0].Numerator / timeRaw[0].Denominator)
	min := int(timeRaw[1].Numerator / timeRaw[1].Denominator)
	sec := int(timeRaw[2].Numerator / timeRaw[2].Denominator)

	timestamp, err = parseTimestamp(dateRaw, hour, min, sec)
	return
}
