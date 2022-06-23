package main

import "io/fs"

const (
	ModDir   fs.FileMode = fs.ModeDir + 0775
	ModeCode fs.FileMode = fs.ModeDevice + 0664
)
