# Apple Metadata Parser for Go

A Go package for parsing and representing Apple device EXIF metadata and MakerNotes, particularly from iPhone/iPad images.

## Overview

This package provides Go types and functions for handling Apple-specific metadata found in images taken with iOS devices. The implementation is based on Phil Harvey's ExifTool documentation.

## Features

- Parse Apple MakerNote limited metadata
- Support for common Apple metadata tags:
  - Auto Focus state
  - Camera type detection (Wide/Normal/Front)
  - Image capture mode (ProRAW/Portrait/Manual/Scene)
  - Acceleration vector data
  - Focus distance ranges

## Todo

- Implement PList parsing