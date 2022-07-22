package draw

import (
	"bytes"
	"io"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	tests := map[string]struct {
		expected *config
	}{
		"default config": {
			expected: &config{
				directory: defaultDirectory,
				filename:  defaultFilename,
			},
		},
	}

	for name, test := range tests {
		c := defaultConfig()

		if !configsAreEqual(test.expected, &c) {
			t.Errorf("%s: config expectation doesn't match: expected %v, got %v", name, test.expected, c)
		}
	}
}

func TestDirectory(t *testing.T) {
	tests := map[string]struct {
		directory string
		expected  *config
	}{
		"default directory": {
			directory: defaultDirectory,
			expected: &config{
				directory: defaultDirectory,
			},
		},
		"custom directory": {
			directory: "./my/directory",
			expected: &config{
				directory: "./my/directory",
			},
		},
	}

	for name, test := range tests {
		c := &config{}

		Directory(test.directory)(c)

		if !configsAreEqual(test.expected, c) {
			t.Errorf("%s: config expectation doesn't match: expected %v, got %v", name, test.expected, c)
		}
	}
}

func TestFilename(t *testing.T) {
	tests := map[string]struct {
		filename string
		expected *config
	}{
		"default filename": {
			filename: defaultFilename,
			expected: &config{
				filename: defaultFilename,
			},
		},
		"custom filename": {
			filename: "myfile.dot",
			expected: &config{
				filename: "myfile.dot",
			},
		},
	}

	for name, test := range tests {
		c := &config{}

		Filename(test.filename)(c)

		if !configsAreEqual(test.expected, c) {
			t.Errorf("%s: config expectation doesn't match: expected %v, got %v", name, test.expected, c)
		}
	}
}

func TestWriter(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})

	tests := map[string]struct {
		writer   io.Writer
		expected *config
	}{
		"default writer": {
			writer: nil,
			expected: &config{
				writer: nil,
			},
		},
		"custom filename": {
			writer: buffer,
			expected: &config{
				writer: buffer,
			},
		},
	}

	for name, test := range tests {
		c := &config{}

		Writer(test.writer)(c)

		if !configsAreEqual(test.expected, c) {
			t.Errorf("%s: config expectation doesn't match: expected %v, got %v", name, test.expected, c)
		}
	}
}

func configsAreEqual(a, b *config) bool {
	return a.directory == b.directory &&
		a.filename == b.filename
}
