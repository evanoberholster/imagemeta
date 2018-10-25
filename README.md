# Exiftool wrapper for Go
Heavily based on: https://github.com/mostlygeek/go-exiftool/


Example usage
```golang
	f, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
	}
	// ExtractExif from Object
	resp, err := exiftool.ExtractWithReader(r)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp.ToExif())
	log.Println(resp.ToLocation())

```