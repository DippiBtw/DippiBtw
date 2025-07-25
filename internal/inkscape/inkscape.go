package inkscape

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	defaultBinaryPath = "/usr/bin/inkscape"
)

type Converter struct {
	bin string
}

// New creates a new Converter with the default Inkscape binary path.
func New() *Converter {
	return &Converter{bin: defaultBinaryPath}
}

// SetBinary sets a custom path to the Inkscape binary.
func (c *Converter) SetBinary(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("inkscape: binary path cannot be empty")
	}
	c.bin = path
	return nil
}

/*
Try to convert the input SVG to the PNG image
*/
func (c *Converter) Convert(in []byte) (out []byte, err error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(c.bin,
		"--pipe",
		"--export-type=png",
		"--export-filename=-",
		"--export-dpi=300",
		"--export-text-to-path",
	)

	cmd.Stdin = bytes.NewBuffer(in)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("inkscape conversion failed: %w\nstderr:\n%s", err, stderr.String())
	}

	if stdout.Len() == 0 {
		return nil, errors.New("inkscape: conversion returned no output")
	}

	return stdout.Bytes(), nil
}
