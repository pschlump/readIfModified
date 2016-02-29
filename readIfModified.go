package ReadIfModified

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var fileTimeStamp map[string]time.Time
var mutexFileTimeStamp = &sync.Mutex{}

var ErrNoSuchFile = errors.New("Error - no such file")

// This is like ioutil.ReadFile - but it will only re-read the content if the file has been modified.
// Also an un-changed flag is returned
func ReadFile(fileName string) (rv []byte, err error, exists, hasChanged bool) {
	var newFileInfo os.FileInfo
	hasChanged = true
	mutexFileTimeStamp.Lock()
	oldTs, ok := fileTimeStamp[fileName]
	mutexFileTimeStamp.Unlock()
	if !ok {
		exists, newFileInfo = ExistsGetFileInfo(fileName)
		if exists {
			fileTimeStamp[fileName] = newFileInfo.ModTime()
			rv, err = ioutil.ReadFile(fileName)
			return
		}
		// no stuch file
		err = ErrNoSuchFile
		return
	}
	exists, newFileInfo = ExistsGetFileInfo(fileName)
	if !exists {
		err = ErrNoSuchFile
		return
	}
	newModTime := newFileInfo.ModTime()
	if newModTime.After(oldTs) {
		fileTimeStamp[fileName] = newModTime
		rv, err = ioutil.ReadFile(fileName)
		return
	}
	hasChanged = false
	return
}

func ExistsGetFileInfo(name string) (bool, os.FileInfo) {
	fi, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
	}
	return true, fi
}
