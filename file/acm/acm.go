package acm

import (
	"fmt"
	"io"
	"bufio"
)

const (
	kACMFile = 0x01032897
)

var (
	ErrNotAnACM = fmt.Errorf("Not a sound file")
)

type header struct {
	Samples   int
	Channels  int
	Rate      int
	Levels    int
	SubBlocks int
}

type SoundStreamer struct {
	header   header
	r        *bitstream
	rawBlock []int16
	samples  []int32
	auxMem1  []int32
	auxMem2  []int32

	leftInBlock int
	currentPos  int
}

func (s *SoundStreamer) Start(r io.Reader) error {
	s.Reset()
	s.r = newBitstream(bufio.NewReader(r))
	return s.initHeader()
}

func (s *SoundStreamer) Reset() {
	s.header = header{}
	s.r = nil
	s.samples = nil
	s.samples = nil
	s.leftInBlock = 0
	s.currentPos = 0
}

func (s *SoundStreamer) SampleRate() int {
	return s.header.Rate
}

func (s *SoundStreamer) NextSample() (float32, float32) {
	left := s.nextSingleSample()
	right := s.nextSingleSample()
	return left, right
}

func (s *SoundStreamer) nextSingleSample() float32 {
	if s.leftInBlock == 0 {
		s.readNextBlock()
		s.leftInBlock = len(s.samples)
	}
	s.leftInBlock--
	s.currentPos++
	return float32(s.samples[len(s.samples)-s.leftInBlock-1] >> s.header.Levels) / 32768.0
}

func (s *SoundStreamer) initHeader() error {
	if s.r.bits(32) != kACMFile {
		return ErrNotAnACM
	}

	s.header.Samples = int(s.r.bits(32))
	s.header.Channels = int(s.r.bits(16))
	if s.header.Channels != 1 && s.header.Channels != 2 {
		return ErrNotAnACM
	}

	s.header.Rate = int(s.r.bits(16))
	s.header.Levels = int(s.r.bits(4))
	if s.header.Levels == 0 {
		return ErrNotAnACM
	}
	s.header.SubBlocks = int(s.r.bits(12))

	subBlockSize := 1 << s.header.Levels
	blockSize := subBlockSize * s.header.SubBlocks
	s.samples  = make([]int32, blockSize)
	s.auxMem2  = make([]int32, 3 * (subBlockSize >> 1) - 2)

	return nil
}

func (s *SoundStreamer) readNextBlock() {
	s.decompressBlock()

	memory16 := make([]int16, len(s.auxMem2) * 2)
	for i := 0; i < len(s.auxMem2); i++ {
		memory16[i*2+0] = int16(s.auxMem2[i])
		memory16[i*2+1] = int16(s.auxMem2[i]>>16)
	}

	sbSize := (1 << s.header.Levels) >> 1
	blocks := s.header.SubBlocks << 1
	decoder1(memory16, s.samples, sbSize, blocks)

	for i := 0; i < len(s.auxMem2); i++ {
		s.auxMem2[i] = int32(memory16[i*2]) | (int32(memory16[i*2 + 1]) << 16)
	}

	mem := s.auxMem2[sbSize:]
	for sbSize > 1 {
		sbSize >>= 1
		blocks <<= 1
		decoder2(mem, s.samples, sbSize, blocks)
		mem = mem[sbSize << 1:]
	}
}
