# xmlmeta

## Work In Progress

Example usage
```go 

package main

import (
    "bytes"
    "fmt"
    "os"
    "time"

    "github.com/evanoberholster/exiftool/meta/jpegmeta"
    "github.com/evanoberholster/exiftool/meta/xmlmeta"
)

const testFilename = "image.jpeg"

func main() {
    f, err := os.Open(testFilename)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    fmt.Println(testFilename)
    start := time.Now()
    m, err := jpegmeta.Scan(f)

    fmt.Println("Decode JPEG:", time.Since(start))
    rXML := bytes.NewReader([]byte(m.XML()))
    start = time.Now()
    pckt, err := xmlmeta.Walk(rXML)
    fmt.Println("Decode XML:", time.Since(start))
    fmt.Println(pckt.ModifyDate())
    fmt.Println("Create Date:", pckt.XMP.CreateDate)
    fmt.Println("MetadataDate:", pckt.XMP.MetadataDate)
    fmt.Println("ModifyDate:", pckt.XMP.ModifyDate)
    fmt.Println("Rights:", pckt.Rights())
    fmt.Println("Creator:", pckt.Creator())
    fmt.Println("Subject:", pckt.Subject())

    fmt.Println("Error:",err)
}

```