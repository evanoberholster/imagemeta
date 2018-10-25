# Exiftool wrapper for Go
Heavily based on: https://github.com/mostlygeek/go-exiftool/

See example.go for more information

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
log.Println(exiftool.ToExif(resp))
log.Println(exiftool.ToLocation(resp))
log.Println(exiftool.ToComposite(resp))
log.Println(exiftool.ToMakerNotes(resp))

```