package dat

import (
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

var (
	ErrNotADat  = fmt.Errorf("Unable to parse DAT file")
	ErrNotFound = fmt.Errorf("Requested file is missing")
)

type file struct {
	name       string // always lower case
	compressed bool
	unpacked   uint32
	packed     uint32
	offset     uint32
}

type Archive struct {
	path  string
	files []*file // sorted alphabetically
}

func Load(path string) (*Archive, error) {
	f, err := os.Open(os.Args[1])
	if err != nil {
		return nil, err
	}
	defer f.Close()

	offset, err := f.Seek(-4, io.SeekEnd)
	if err != nil {
		return nil, ErrNotADat
	}

	var size uint32
	if binary.Read(f, binary.LittleEndian, &size) != nil {
		return nil, ErrNotADat
	}
	if size != uint32(offset + 4) || size < 8 {
		return nil, ErrNotADat
	}

	if _, err := f.Seek(-8, io.SeekCurrent); err != nil {
		return nil, ErrNotADat
	}
	if binary.Read(f, binary.LittleEndian, &size) != nil {
		return nil, ErrNotADat
	}

	if _, err := f.Seek(-int64(uint64(size))-4, io.SeekCurrent); err != nil {
		return nil, ErrNotADat
	}
	if binary.Read(f, binary.LittleEndian, &size) != nil {
		return nil, ErrNotADat
	}

	itemCount := int(size)
	a := &Archive{
		path:  path,
		files: make([]*file, itemCount),
	}
	for i := 0; i < itemCount; i++ {
		if binary.Read(f, binary.LittleEndian, &size) != nil {
			return nil, ErrNotADat
		}
		nameData := make([]byte, size)
		if binary.Read(f, binary.LittleEndian, nameData) != nil {
			return nil, ErrNotADat
		}

		var compressed uint8
		var unpackedSize, packedSize, offset uint32

		if binary.Read(f, binary.LittleEndian, &compressed) != nil {
			return nil, ErrNotADat
		}
		if binary.Read(f, binary.LittleEndian, &unpackedSize) != nil {
			return nil, ErrNotADat
		}
		if binary.Read(f, binary.LittleEndian, &packedSize) != nil {
			return nil, ErrNotADat
		}
		if binary.Read(f, binary.LittleEndian, &offset) != nil {
			return nil, ErrNotADat
		}

		a.files[i] = new(file)
		a.files[i].name = strings.ToLower(string(nameData))
		a.files[i].compressed = compressed != 0
		a.files[i].unpacked = unpackedSize
		a.files[i].packed = packedSize
		a.files[i].offset = offset
	}

	sort.Slice(a.files, func(i, j int) bool {
		return a.files[i].name < a.files[j].name
	})

	return a, nil
}

func (a *Archive) HasFile(path string) bool {
	path = strings.ToLower(path)
	i := sort.Search(len(a.files), func(i int) bool {
		return a.files[i].name >= path
	})
	return i < len(a.files) && a.files[i].name == path
}

type compressibleWrapper struct {
	limit int64
	r io.ReadCloser
	z io.ReadCloser	
}

func (c *compressibleWrapper) Read(p []byte) (int, error) {
	if c.limit == 0 {
		return 0, io.EOF
	}
	limit := int64(len(p))
	if limit > c.limit {
		limit = c.limit
	}
	c.limit -= limit
	return c.z.Read(p[:limit])
}

func (c *compressibleWrapper) Close() error {
	if c.r != nil {
		c.r.Close()
	}
	return c.z.Close()
}

func (a *Archive) LoadFile(path string) (io.ReadCloser, error) {
	path = strings.ToLower(path)
	i := sort.Search(len(a.files), func(i int) bool {
		return a.files[i].name >= path
	})
	if i == len(a.files) || a.files[i].name != path {
		return nil, ErrNotFound
	}
	file := a.files[i]

	f, err := os.Open(a.path)
	if err != nil {
		return nil, err
	}
	f.Seek(int64(uint64(file.offset)), io.SeekStart)
	wr := new(compressibleWrapper)
	wr.limit = int64(uint64(file.unpacked))
	wr.z = f
	if file.compressed {
		wr.r = f
		wr.z, err = zlib.NewReader(wr.r)
		if err != nil {
			f.Close()
			return nil, err
		}
	}
	return wr, nil
}
