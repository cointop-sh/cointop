package filecache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DefaultCacheDir ...
var DefaultCacheDir = ":PREFERRED_CACHE_HOME:/cointop"

// FileCache ...
type FileCache struct {
	mapLock  sync.Mutex
	muts     map[string]*sync.Mutex
	prefix   string
	cacheDir string
}

// Config ...
type Config struct {
	CacheDir string
	Prefix   string
}

// NewFileCache ...
func NewFileCache(config *Config) (*FileCache, error) {
	if config == nil {
		config = &Config{}
	}

	cacheDir := DefaultCacheDir
	if config.CacheDir != "" {
		cacheDir = config.CacheDir
	}
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0700); err != nil {
			return nil, err
		}
	}

	return &FileCache{
		muts:     make(map[string]*sync.Mutex),
		cacheDir: cacheDir,
		prefix:   config.Prefix,
	}, nil
}

// Set writes item to cache
func (f *FileCache) Set(key string, data interface{}, expire time.Duration) error {
	var mu *sync.Mutex
	var ok bool
	f.mapLock.Lock()
	if mu, ok = f.muts[key]; !ok {
		f.muts[key] = new(sync.Mutex)
		mu = f.muts[key]
	}
	f.mapLock.Unlock()

	mu.Lock()
	defer mu.Unlock()

	key = regexp.MustCompile("[^a-zA-Z0-9_-]").ReplaceAllLiteralString(key, "")
	if f.prefix != "" {
		key = fmt.Sprintf("%s.%s", f.prefix, key)
	}
	ts := strconv.FormatInt(time.Now().Add(expire).Unix(), 10)
	file := fmt.Sprintf("fcache.%s.%v", key, ts)
	fpath := filepath.Join(f.cacheDir, file)

	f.clean(key)

	serialized, err := serialize(data)
	if err != nil {
		return err
	}

	fp, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer fp.Close()
	if _, err = fp.Write(serialized); err != nil {
		return err
	}

	return nil
}

// Get reads item from cache
func (f *FileCache) Get(key string, dst interface{}) error {
	key = regexp.MustCompile("[^a-zA-Z0-9_-]").ReplaceAllLiteralString(key, "")
	if f.prefix != "" {
		key = fmt.Sprintf("%s.%s", f.prefix, key)
	}
	pattern := filepath.Join(f.cacheDir, fmt.Sprintf("fcache.%s.*", key))
	files, err := filepath.Glob(pattern)
	if len(files) < 1 || err != nil {
		return errors.New("fcache: no cache file found")
	}

	if _, err = os.Stat(files[0]); err != nil {
		return err
	}

	fp, err := os.OpenFile(files[0], os.O_RDONLY, 0400)
	if err != nil {
		return err
	}
	defer fp.Close()

	var serialized []byte
	buf := make([]byte, 1024)
	for {
		var n int
		n, err = fp.Read(buf)
		serialized = append(serialized, buf[0:n]...)
		if err != nil || err == io.EOF {
			break
		}
	}

	if err = deserialize(serialized, dst); err != nil {
		return err
	}

	for _, file := range files {
		parts := strings.Split(file, ".")
		ts := parts[len(parts)-1]
		exptime, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return err
		}

		if exptime < time.Now().Unix() {
			if _, err = os.Stat(file); err == nil {
				os.Remove(file)
			}
		}
	}

	return nil
}

// clean removes item from cache
func (f *FileCache) clean(key string) error {
	pattern := filepath.Join(f.cacheDir, fmt.Sprintf("fcache.%s.*", key))
	files, _ := filepath.Glob(pattern)
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			os.Remove(file)
		}
	}

	return nil
}

// serialize encodes a value using binary
func serialize(src interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// deserialize decodes a value using binary
func deserialize(src []byte, dst interface{}) error {
	buf := bytes.NewReader(src)
	if err := gob.NewDecoder(buf).Decode(dst); err != nil {
		return err
	}

	return nil
}
