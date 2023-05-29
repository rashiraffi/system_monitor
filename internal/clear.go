package internal

import "os"

func Clear() {
	os.Remove(dbPath)
}
