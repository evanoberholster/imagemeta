package exiftool
import (
	"time"
)

type File struct {
	FileSize		int
	MIMEType		string
	ImageWidth		int
	ImageHeight		int
} 

func (e ExifResponse)ToFile() (File) {
	filesize := e.GetExifItem(GroupFile, "FileSize").ToInt()
	mimeType := e.GetExifItem(GroupFile, "MIMEType").ToString()
	width := e.ExifPath(Path(GroupFile, "ImageWidth"), Path(GroupExif, "ImageWidth")).ToInt()
	height := e.ExifPath(Path(GroupFile, "ImageHeight"), Path(GroupExif, "ImageHeight")).ToInt()
	return File{
		FileSize: filesize,
		MIMEType: mimeType,
		ImageWidth: width,
		ImageHeight: height,
	}
}

type Exif struct {
	FileSize					int 		`json:"FileSize"`
	MIMEType					string 		`json:"MIMEType"`
	ImageWidth              	int     	`json:"ImageWidth"`
	ImageHeight             	int     	`json:"ImageHeight"`
	CameraMake                  string  	`json:"CameraMake"`
	CameraModel                 string  	`json:"CameraModel"`
	CameraSerial				string 		`json:"CameraSerial"`
	LensModel					string 		`json:"LensModel"`
	LensSerial					string 		`json:"LensSerial"`
	CreateDate 					time.Time  	`json:"CreateDate`
	ModifyDate                	time.Time  	`json:"ModifyDate"`
	Artist                    	string  	`json:"Artist"`
	Copyright                 	string  	`json:"Copyright"`
	Aperture 					string  	`json:"Aperture"`
	ShutterSpeed 				string 		`json:"ShutterSpeed"`
	ISO                       	int     	`json:"ISO"`
	ExposureProgram           	int 	  	`json:"ExposureProgram"`
	MeteringMode				int 		`json:"MeteringMode"`	
	Orientation             	int 	  	`json:"Orientation"`
	Flash						int 		`json:"Flash"`
	Software					string 		`json:"Software"`
	FocalLength 				float64		`json:"FocalLength"` // mm
	ScaleFactor35Efl       		float64 	`json:"ScaleFactor35efl"`
	CircleOfConfusion      		float64 	`json:"CircleOfConfusion"` // mm
	DOF                    		string  	`json:"DOF"`
	FOV                    		float64  	`json:"FOV"` // degrees
	HyperfocalDistance			float64		`json:"HyperfocalDistance"` // meters
}

func (e ExifResponse)ToExif() (Exif) {
	//dateTimeOriginal := e.ExifPath(Path(GroupComposite, "SubSecDateTimeOriginal"), Path(GroupExif, "DateTimeOriginal")).ToTime()
	return Exif{
		FileSize: e.GetExifItem(GroupFile, "FileSize").ToInt(),
		MIMEType: e.GetExifItem(GroupFile, "MIMEType").ToString(),
		ImageWidth: e.ExifPath(Path(GroupFile, "ImageWidth"), Path(GroupExif, "ImageWidth"), Path(GroupExif, "ExifImageWidth")).ToInt(),
		ImageHeight: e.ExifPath(Path(GroupFile, "ImageHeight"), Path(GroupExif, "ImageHeight"), Path(GroupExif, "ExifImageHeight")).ToInt(),
		CameraMake: e.ExifPath(Path(GroupExif, "Make")).ToString(),
		CameraModel: e.ExifPath(Path(GroupExif, "Model")).ToString(),
		LensModel: e.ExifPath(Path(GroupComposite, "LensID"), Path(GroupExif, "LensModel"), Path(GroupComposite, "Lens")).ToString(),
		CreateDate: e.ExifPath(Path(GroupComposite, "SubSecCreateDate"), Path(GroupExif, "CreateDate")).ToTime(),
		ModifyDate: e.ExifPath(Path(GroupComposite, "SubSecModifyDate"), Path(GroupExif, "ModifyDate")).ToTime(),
		Artist: e.GetExifItem(GroupExif, "Artist").ToString(),
		Copyright: e.GetExifItem(GroupExif, "Copyright").ToString(),
		Aperture: e.ExifPath(Path(GroupComposite, "Aperture"), Path(GroupExif, "ApertureValue"), Path(GroupExif, "FNumber")).ToString(),
		ShutterSpeed: e.ExifPath(Path(GroupComposite, "ShutterSpeed"), Path(GroupExif, "ShutterSpeedValue")).ToString(),
		ISO: e.ExifPath(Path(GroupComposite, "ISO"), Path(GroupExif, "ISO"), Path(GroupMakerNotes, "ISO")).ToInt(),
		ExposureProgram: e.GetExifItem(GroupExif, "ExposureProgram").ToInt(),
		Orientation: e.GetExifItem(GroupExif, "Orientation").ToInt(),
		Flash: e.GetExifItem(GroupExif, "Flash").ToInt(),
		MeteringMode: e.GetExifItem(GroupExif, "MeteringMode").ToInt(),
		Software: e.GetExifItem(GroupExif, "Software").ToString(),
		FocalLength: e.ExifPath(Path(GroupMakerNotes, "FocalLength"),Path(GroupExif, "FocalLength")).ToFloat(),
		ScaleFactor35Efl: e.GetExifItem(GroupComposite, "ScaleFactor35efl").ToFloat(),
		CircleOfConfusion: e.GetExifItem(GroupComposite, "CircleOfConfusion").ToFloat(),
		DOF: e.GetExifItem(GroupComposite, "DOF").ToString(),
		FOV: e.GetExifItem(GroupComposite, "FOV").ToFloat(),
		HyperfocalDistance: e.GetExifItem(GroupComposite, "HyperfocalDistance").ToFloat(),
		CameraSerial: e.ExifPath(Path(GroupExif, "SerialNumber"),Path(GroupExif, "CameraSerialNumber")).ToString(),
		LensSerial: e.ExifPath(Path(GroupMakerNotes, "LensSerialNumber")).ToString(),
	}
}

type MakerNotes struct {
	MacroMode                  	string  `json:"MacroMode"`
	SelfTimer                  	string  `json:"SelfTimer"`
	ContinuousDrive            	string  `json:"ContinuousDrive"`
	MeteringMode               	string  `json:"MeteringMode"`
	CanonExposureMode          	string  `json:"CanonExposureMode"`
	FocusMode                  	string  `json:"FocusMode"`
	WhiteBalance               	string  `json:"WhiteBalance"`
	ColorTemperature           	int     `json:"ColorTemperature"`
	ColorTempAsShot            	int     `json:"ColorTempAsShot"`
	CameraTemperature          	string  `json:"CameraTemperature"`
	OwnerName                  	string  `json:"OwnerName"`
	AEB 				     	string  `json:"AutoExposureBracketing"`
	AEBBracketValue            	int     `json:"AEBBracketValue"`
	AEBShotCount               	string  `json:"AEBShotCount"`
	BracketMode                	string  `json:"BracketMode"`
	BracketValue               	int     `json:"BracketValue"`
	BracketShotNumber          	int     `json:"BracketShotNumber"`
	ExposureCompensation       	string  `json:"ExposureCompensation"`
	LiveViewShooting           	string  `json:"LiveViewShooting"`
	ColorSpace                 	string  `json:"ColorSpace"`
	HighlightTonePriority      	string  `json:"HighlightTonePriority"`
	MaxFocalLength             	string  `json:"MaxFocalLength"`
	MinFocalLength             	string  `json:"MinFocalLength"`
	FocalLength                	string  `json:"FocalLength"`
	FocusDistanceUpper         	string  `json:"FocusDistanceUpper"`
	FocusDistanceLower         	string  `json:"FocusDistanceLower"`
	TimeZone                   	string  `json:"TimeZone"`
	TimeZoneCity               	string  `json:"TimeZoneCity"`
	DaylightSavings            	string  `json:"DaylightSavings"`
	AspectRatio                	string  `json:"AspectRatio"`
	LensModel                  	string  `json:"LensModel"`
}

func (e ExifResponse)ToMakerNotes() (MakerNotes) {
	return MakerNotes{
		MacroMode: e.ExifPath(Path(GroupMakerNotes, "MacroMode")).ToString(),
		SelfTimer: e.ExifPath(Path(GroupMakerNotes, "SelfTimer")).ToString(),
		ContinuousDrive: e.ExifPath(Path(GroupMakerNotes, "ContinuousDrive")).ToString(),
		MeteringMode: e.ExifPath(Path(GroupMakerNotes, "MeteringMode")).ToString(),
		CanonExposureMode: e.ExifPath(Path(GroupMakerNotes, "CanonExposureMode")).ToString(),
		FocusMode: e.ExifPath(Path(GroupMakerNotes, "FocusMode")).ToString(),
		WhiteBalance: e.ExifPath(Path(GroupMakerNotes, "WhiteBalance")).ToString(),
		ColorTemperature: e.ExifPath(Path(GroupMakerNotes, "ColorTemperature")).ToInt(),
		ColorTempAsShot: e.ExifPath(Path(GroupMakerNotes, "ColorTempAsShot")).ToInt(),
		CameraTemperature: e.ExifPath(Path(GroupMakerNotes, "CameraTemperature")).ToString(),
		OwnerName: e.ExifPath(Path(GroupMakerNotes, "OwnerName")).ToString(),
		AEB: e.ExifPath(Path(GroupMakerNotes, "AutoExposureBracketing")).ToString(),
		AEBBracketValue: e.ExifPath(Path(GroupMakerNotes, "AEBBracketValue")).ToInt(),
		AEBShotCount: e.ExifPath(Path(GroupMakerNotes, "AEBShotCount")).ToString(),
		BracketMode: e.ExifPath(Path(GroupMakerNotes, "BracketMode")).ToString(),
		BracketValue: e.ExifPath(Path(GroupMakerNotes, "BracketValue")).ToInt(),
		BracketShotNumber: e.ExifPath(Path(GroupMakerNotes, "BracketShotNumber")).ToInt(),
		ExposureCompensation: e.ExifPath(Path(GroupMakerNotes, "ExposureCompensation")).ToString(),
		LiveViewShooting: e.ExifPath(Path(GroupMakerNotes, "LiveViewShooting")).ToString(),
		ColorSpace: e.ExifPath(Path(GroupMakerNotes, "ColorSpace")).ToString(),
		HighlightTonePriority: e.ExifPath(Path(GroupMakerNotes, "HighlightTonePriority")).ToString(),
		MaxFocalLength: e.ExifPath(Path(GroupMakerNotes, "MaxFocalLength")).ToString(),
		MinFocalLength: e.ExifPath(Path(GroupMakerNotes, "MinFocalLength")).ToString(),
		FocalLength: e.ExifPath(Path(GroupMakerNotes, "FocalLength")).ToString(),
		FocusDistanceUpper: e.ExifPath(Path(GroupMakerNotes, "FocusDistanceUpper")).ToString(),
		FocusDistanceLower: e.ExifPath(Path(GroupMakerNotes, "FocusDistanceLower")).ToString(),
		TimeZone: e.ExifPath(Path(GroupMakerNotes, "TimeZone")).ToString(),
		TimeZoneCity: e.ExifPath(Path(GroupMakerNotes, "TimeZoneCity")).ToString(),
		DaylightSavings: e.ExifPath(Path(GroupMakerNotes, "DaylightSavings")).ToString(),
		AspectRatio: e.ExifPath(Path(GroupMakerNotes, "AspectRatio")).ToString(),
		LensModel: e.ExifPath(Path(GroupMakerNotes, "LensModel")).ToString(),
	}
}


type CanonAutoFocus struct {
	AFAreaMode                 string  `json:"AFAreaMode"`
	NumAFPoints                int     `json:"NumAFPoints"`
	ValidAFPoints              int     `json:"ValidAFPoints"`
	AFImageWidth               int     `json:"AFImageWidth"`
	AFImageHeight              int     `json:"AFImageHeight"`
	AFAreaWidths               string  `json:"AFAreaWidths"`
	AFAreaHeights              string  `json:"AFAreaHeights"`
	AFAreaXPositions           string  `json:"AFAreaXPositions"`
	AFAreaYPositions           string  `json:"AFAreaYPositions"`
	AFPointsInFocus            int     `json:"AFPointsInFocus"`
	AFPointsSelected           int     `json:"AFPointsSelected"`
	OrientationLinkedAFPoint   string  `json:"OrientationLinkedAFPoint"`
}

type Location struct {
	GPSAltitude 				float64
	GPSDateTime					time.Time
	GPSLatitude 				float64
	GPSLongitude 				float64	
	GPSPositionError 			float64	
}

func (e ExifResponse)ToLocation() (Location) {
	altitude := e.ExifPath(Path(GroupComposite, "GPSAltitude"),Path(GroupExif, "GPSAltitude")).ToFloat()
	latitude := e.ExifPath(Path(GroupComposite, "GPSLatitude"),Path(GroupExif, "GPSLatitude")).ToFloat()
	longitude := e.ExifPath(Path(GroupComposite, "GPSLongitude"),Path(GroupExif, "GPSLongitude")).ToFloat()
	gpsTime := e.ExifPath(Path(GroupComposite, "GPSDateTime"),Path(GroupExif, "GPSDateTime")).ToTime()
	gpsError := e.GetExifItem(GroupExif, "GPSHPositioningError").ToFloat()

	return Location{
		GPSAltitude: altitude,
		GPSDateTime: gpsTime,
		GPSLatitude: latitude,
		GPSLongitude: longitude,
		GPSPositionError: gpsError,
	}
}

type Composite struct {
	Aperture               	float64 	`json:"Aperture"`
	ShutterSpeed           	string  	`json:"ShutterSpeed"`
	ISO                    	int     	`json:"ISO"`
	Lens                   	string  	`json:"Lens"`
	LensID                 	string  	`json:"LensID"`
	Megapixels             	float64 	`json:"Megapixels"`
	CreateDate       		time.Time 	`json:"CreateDate"`
	DateTimeOriginal 		time.Time  	`json:"DateTimeOriginal"`
	ModifyDate       		time.Time  	`json:"ModifyDate"`
	CircleOfConfusion      	string  	`json:"CircleOfConfusion"`
	DOF                    	string  	`json:"DOF"`
	FOV                    	string  	`json:"FOV"`
	ScaleFactor35Efl       	float64 	`json:"ScaleFactor35efl"`
	FocalLength35Efl       	float64  	`json:"FocalLength35efl"` // mm
	HyperfocalDistance     	string  	`json:"HyperfocalDistance"`
	Lens35Efl              	string  	`json:"Lens35efl"`
	LightValue             	float64 	`json:"LightValue"`
	ShootingMode           	string  	`json:"ShootingMode"`
	DriveMode              	string  	`json:"DriveMode"`
}

func (e ExifResponse)ToComposite() (Composite) {	
	return Composite{
		Aperture: e.GetExifItem(GroupComposite, "Aperture").ToFloat(),
		ISO: e.ExifPath(Path(GroupComposite, "ISO"), Path(GroupExif, "ISO")).ToInt(),
		Lens: e.GetExifItem(GroupComposite, "Lens").ToString(),
		LensID: e.GetExifItem(GroupComposite, "LensID").ToString(),
		Megapixels: e.GetExifItem(GroupComposite, "Megapixels").ToFloat(),
		ScaleFactor35Efl: e.GetExifItem(GroupComposite, "ScaleFactor35efl").ToFloat(),
		ShutterSpeed: e.GetExifItem(GroupComposite, "ShutterSpeed").ToString(),
		CreateDate: e.ExifPath(Path(GroupComposite, "SubSecCreateDate"), Path(GroupExif, "CreateDate")).ToTime(),
		DateTimeOriginal: e.ExifPath(Path(GroupComposite, "SubSecDateTimeOriginal"), Path(GroupExif, "DateTimeOriginal")).ToTime(),
		ModifyDate: e.ExifPath(Path(GroupComposite, "SubSecModifyDate"), Path(GroupExif, "ModifyDate")).ToTime(),
		CircleOfConfusion: e.GetExifItem(GroupComposite, "CircleOfConfusion").ToString(),
		DOF: e.GetExifItem(GroupComposite, "DOF").ToString(),
		FOV: e.GetExifItem(GroupComposite, "FOV").ToString(),
		FocalLength35Efl: e.GetExifItem(GroupComposite, "FocalLength35Efl").ToFloat(),
		HyperfocalDistance: e.GetExifItem(GroupComposite, "HyperfocalDistance").ToString(),
		Lens35Efl: e.GetExifItem(GroupComposite, "Lens35Efl").ToString(),
		LightValue:  e.GetExifItem(GroupComposite, "LightValue").ToFloat(),	
		ShootingMode: e.ExifPath(Path(GroupComposite, "ShootingMode"), Path(GroupMakerNotes, "ShootingMode")).ToString(),
		DriveMode: e.ExifPath(Path(GroupComposite, "DriveMode"), Path(GroupMakerNotes, "DriveMode")).ToString(),
	}
}

