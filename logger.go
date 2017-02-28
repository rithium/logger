package logger

import (
	"sync"
	"os"
	"time"
	"path/filepath"
)

type RotatingFileWriter struct {
	lock		sync.Mutex
	filename	string
	fp		*os.File
	maxWrites	int
	writes		int
}

var Logger *RotatingFileWriter

func NewRotatingFileWriter(filename string, maxWrites int) (*RotatingFileWriter, error) {
	Logger = &RotatingFileWriter{filename: filename, maxWrites: maxWrites}

	if err := Logger.rotate(); err != nil {
		return nil, err
	}

	return Logger, nil
}

func (w *RotatingFileWriter) Write(output []byte) (int, error) {
	w.lock.Lock()

	total, err := w.fp.Write(output)

	w.lock.Unlock()

	if err != nil {
		return total, err
	}

	w.writes = w.writes + 1

	if w.writes >= w.maxWrites {
		err := w.rotate()

		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func (w *RotatingFileWriter) rotate() (error) {
	w.lock.Lock()

	defer w.lock.Unlock()

	if w.fp != nil {
		err := w.fp.Close()

		w.fp = nil

		if err != nil {
			return nil
		}
	}

	_, err := os.Stat(w.filename)

	if err == nil {
		err = os.Rename(w.filename, filepath.Dir(w.filename)+"/"+time.Now().Format(time.RFC3339)+"-"+filepath.Base(w.filename))

		if err != nil {
			return err
		}
	}

	w.fp, err = os.Create(w.filename)

	w.writes = 0

	return nil
}