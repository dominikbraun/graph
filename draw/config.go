package draw

const (
	defaultDirectory = "."
	defaultFilename  = "graph.dot"
)

type config struct {
	directory string
	filename  string
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
