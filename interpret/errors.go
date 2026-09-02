package interpret

import "errors"

var (
	// ErrEmptyInput indicates that an interpreter received no usable input.
	ErrEmptyInput = errors.New("empty input")
	// ErrUnrecognizedInput indicates that no supported interpretation was found.
	ErrUnrecognizedInput = errors.New("unrecognized input")
	// ErrAmbiguousInput indicates that multiple interpretations are equally plausible.
	ErrAmbiguousInput = errors.New("ambiguous input")
)
