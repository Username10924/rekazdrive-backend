package storage

import (
	"bytes"
	"path"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTPBackend struct {
	host string
	user string
	pass string
	basePath string
}

func NewFTPBackend(host, user, pass, basePath string) *FTPBackend {
	if basePath == "" {
		basePath = "/"
	}
	return &FTPBackend{host: host, user: user, pass: pass, basePath: basePath}
}

func (f *FTPBackend) dial() (*ftp.ServerConn, error) {
	conn, err := ftp.Dial(f.host, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}
	
	if err := conn.Login(f.user, f.pass); err != nil {
		conn.Quit()
		return nil, err
	}

	return conn, nil // connection established
}

func sanitizeID(id string) string {
	id = strings.ReplaceAll(id, "/", "_") // prevent directory traversal

	if id == "" {
		id = "someblob"
	}

	return id
}

func (f *FTPBackend) Save(id string, data []byte) error {
	conn, err := f.dial()
	if err != nil {
		return err
	}
	defer conn.Quit() // closing connection after saving to avoid memory leaks

	safeID := sanitizeID(id)
	fullPath := path.Join(f.basePath, safeID)

	return conn.Stor(fullPath, bytes.NewReader(data))
}

func (f *FTPBackend) Load(id string) ([]byte, error) {
	conn, err := f.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Quit()

	safeID := sanitizeID(id)
	fullPath := path.Join(f.basePath, safeID)

	response, err := conn.Retr(fullPath)
	if err != nil {
		return nil, err
	}
	defer response.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(response); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (f *FTPBackend) Delete(id string) error {
	conn, err := f.dial()
	if err != nil {
		return err
	}
	defer conn.Quit()
	safeID := sanitizeID(id)
	fullPath := path.Join(f.basePath, safeID)

	return conn.Delete(fullPath)
}