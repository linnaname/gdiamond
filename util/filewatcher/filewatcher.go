/**
 */

package filewatcher

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

/**
code from https://github.com/farmergreg/rfsnotify,  I need to change something about it so I copy it and make some difference
*/

//RWatcher Recursive file or directory watcher
type RWatcher struct {
	//you can for select all events and errors
	Events   chan fsnotify.Event
	Errors   chan error
	done     chan struct{}
	fsnotify *fsnotify.Watcher
	isClosed bool
}

//NewWatcher create RWatcher
func NewWatcher() (*RWatcher, error) {
	fsWatch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &RWatcher{}
	w.fsnotify = fsWatch
	w.Events = make(chan fsnotify.Event)
	w.Errors = make(chan error)
	w.done = make(chan struct{})
	go w.start()
	return w, nil
}

// AddRecursive  watch the  directory and all sub-directories
func (rw *RWatcher) AddRecursive(name string) error {
	if rw.isClosed {
		return errors.New("rfsnotify instance already closed")
	}
	if err := rw.watchRecursive(name, false); err != nil {
		return err
	}
	return nil
}

// RemoveRecursive stop watch the  directory and all sub-directories
func (rw *RWatcher) RemoveRecursive(name string) error {
	if err := rw.watchRecursive(name, true); err != nil {
		return err
	}
	return nil
}

// Close remove all watche and close the events channel
func (rw *RWatcher) Close() error {
	if rw.isClosed {
		return nil
	}
	close(rw.done)
	rw.isClosed = true
	return nil
}

func (rw *RWatcher) start() {
	for {
		select {

		case e := <-rw.fsnotify.Events:
			s, err := os.Stat(e.Name)
			if err == nil && s != nil && s.IsDir() {
				//adds all directory under the new file or directory
				if e.Op&fsnotify.Create != 0 {
					rw.watchRecursive(e.Name, false)
				}
			}
			//Can't stat a deleted directory and sub-directories
			if e.Op&fsnotify.Remove != 0 {
				rw.fsnotify.Remove(e.Name)
			}
			rw.Events <- e

		case e := <-rw.fsnotify.Errors:
			rw.Errors <- e

		case <-rw.done:
			rw.fsnotify.Close()
			close(rw.Events)
			close(rw.Errors)
			return
		}
	}
}

// watchRecursive adds all directories under the given one to the watch list.
func (rw *RWatcher) watchRecursive(path string, unWatch bool) error {
	err := filepath.Walk(path, func(walkPath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			if unWatch {
				if err = rw.fsnotify.Remove(walkPath); err != nil {
					return err
				}
			} else {
				if err = rw.fsnotify.Add(walkPath); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}
