# meta/exif/tags — ExifTool-style tag data for imagemeta

This directory holds **build-time data, not runtime code**. Each YAML file is a
simplified projection of one ExifTool tag table (`%Main`), intended as the
single source of truth for a future `taggen` code generator that will emit
idiomatic Go `switch` dispatch (no runtime tables, no func pointers, zero
allocs preserved).

## Sources

Extracted from ExifTool **13.55** (production release; latest known development
release at time of writing is **13.59**, May 27 2026 — same `%Main` shape, only
new lenses/tags added):

| YAML | ExifTool module | Table | Docs |
|---|---|---|---|
| `exif.yaml` | `Image::ExifTool::Exif` | `%Main` | https://exiftool.org/TagNames/EXIF.html |
| `gps.yaml` | `Image::ExifTool::GPS` | `%Main` | https://exiftool.org/TagNames/GPS.html |
| `canon.yaml` | `Image::ExifTool::Canon` | `%Main` | https://exiftool.org/TagNames/Canon.html |
| `nikon.yaml` | `Image::ExifTool::Nikon` | `%Main` | https://exiftool.org/TagNames/Nikon.html |
| `sony.yaml` | `Image::ExifTool::Sony` | `%Main` | https://exiftool.org/TagNames/Sony.html |
| `panasonic.yaml` | `Image::ExifTool::Panasonic` | `%Main` | https://exiftool.org/TagNames/Panasonic.html |
| `apple.yaml` | `Image::ExifTool::Apple` | `%Main` | https://exiftool.org/TagNames/Apple.html |

`exiftool_version`, `source_tag_count` (numeric ExifTool entries) and `tag_count`
(named entries exported) are recorded in each file header.

## Schema

```yaml
exiftool_version: "13.55"
source: {module, table, file, docs}
groups: {0: EXIF, 1: GPS, ...}   # ExifTool GROUPS, informational
source_tag_count: 590            # numeric keys in %Main
tag_count: 414                   # named tags exported (binary/data aliases skipped)
tags:
  - id: 0x010f
    name: "Make"
    writable: "string"           # ExifTool Writable (format hint)
    group: "IFD0"                # WriteGroup
    count: ""                    # ExifTool Count (-1 = variable, "" = default 1)
    description: "..."           # Description/Notes, single-line, truncated
    print_conv: "hash|code|scalar:...|none"  # summary only, never CODE refs
    value_conv: "none"
    subdirectory: false
    subdirectory_table: "Nikon::ShotInfo"    # only when present
```

`print_conv`/`value_conv` are intentionally lossy summaries (`hash`, `code`,
`scalar:…`, `none`). Full Perl converters are not portable to Go; the Go
`taggen` step maps them to the existing `parse*/Print` helpers in
`meta/exif` (e.g. `apexAperture`, `GPSCoord`).

## Regenerating

Requires the ExifTool Perl modules (Homebrew `exiftool` ships them):

```sh
PERL5LIB=/opt/homebrew/Cellar/exiftool/13.55/libexec/lib/perl5 \
TAGS_OUT=meta/exif/tags \
perl meta/exif/tags/dump_exiftool_tables.pl
```

The committed generator emits byte-identical YAML when run against ExifTool
13.55 (verified in CI by hand; see PR). Re-running against 13.59+ will show up
as a clean data-only diff.

## Codegen plan (not implemented in this PR)

`cmd/taggen` (future) reads these YAML files and emits per-IFD
`zz_generated_*.go` containing plain Go `switch t.ID` bodies plus cold-only
`Name()` maps — the same shape as today's hand-written
`parseIFD0Tag`/`parseExifTag`/`parseGPSTag`, but generated. Runtime stays
`func`-table free: `BulkTrusted` header scan and `state` queue are untouched,
`b.ReportAllocs()` budgets (`JPG 64B/2allocs`) hold.

## Non-goals of this PR

No parser, model, or public API changes. No new Go dependencies. These YAML
files are not loaded at runtime.
