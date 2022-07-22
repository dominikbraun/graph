package draw

import "io"

const (
	defaultDirectory = "."
	defaultFilename  = "graph.gv"
)

type config struct {
	directory string
	filename  string
	writer    io.Writer
}

func defaultConfig() config {
	return config{
		directory: defaultDirectory,
		filename:  defaultFilename,
	}
}

func Directory(directory string) func(*config) {
	return func(c *config) {
		c.directory = directory
	}
}

func Filename(filename string) func(*config) {
	return func(c *config) {
		c.filename = filename
	}
}

func Writer(writer io.Writer) func(*config) {
	return func(c *config) {
		c.writer = writer
	}
}
