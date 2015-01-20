package cfg

import "os"

type media struct {
	Root string
}

var Media media

func init() {
	Media.Root = os.Getenv("MEDIA_ROOT")
}
