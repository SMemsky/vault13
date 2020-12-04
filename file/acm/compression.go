package acm

var (
	Table1 = [...]int{
		0, 1, 2, 4, 5, 6, 8, 9, 10, 16, 17, 18, 20, 21, 22, 24, 25, 26, 32, 33, 34, 36, 37, 38, 40, 41, 42,
	}

	Table2 = [...]int{
		0, 1, 2, 3, 4, 8, 9, 10, 11, 12, 16, 17, 18, 19, 20, 24, 25, 26, 27, 28,
		32, 33, 34, 35, 36, 64, 65, 66, 67, 68, 72, 73, 74, 75, 76, 80, 81, 82,
		83, 84, 88, 89, 90, 91, 92, 96, 97, 98, 99, 100, 128, 129, 130, 131, 132,
		136, 137, 138, 139, 140, 144, 145, 146, 147, 148, 152, 153, 154, 155, 156,
		160, 161, 162, 163, 164, 192, 193, 194, 195, 196, 200, 201, 202, 203, 204,
		208, 209, 210, 211, 212, 216, 217, 218, 219, 220, 224, 225, 226, 227, 228,
		256, 257, 258, 259, 260, 264, 265, 266, 267, 268, 272, 273, 274, 275, 276,
		280, 281, 282, 283, 284, 288, 289, 290, 291, 292,
	}

	Table3 = [...]int{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x20, 0x21,
		0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x30, 0x31, 0x32,
		0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x40, 0x41, 0x42, 0x43,
		0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x50, 0x51, 0x52, 0x53, 0x54,
		0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x60, 0x61, 0x62, 0x63, 0x64, 0x65,
		0x66, 0x67, 0x68, 0x69, 0x6A, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76,
		0x77, 0x78, 0x79, 0x7A, 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87,
		0x88, 0x89, 0x8A, 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98,
		0x99, 0x9A, 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9,
		0xAA,
	}
)

func (s *SoundStreamer) prepareAux1() int {
	count := 1 << int(s.r.bits(4))
	s.auxMem1 = make([]int32, count * 2)

	{
		value := int16(s.r.bits(16))
		t1, t2 := int16(0), -value
		for i := 0; i < count; i++ {
			s.auxMem1[count+i  ] = int32(t1)
			s.auxMem1[count-i-1] = int32(t2)
			t1 += value
			t2 -= value
		}
	}

	return count
}

func (s *SoundStreamer) decompressBlock() {
	count := s.prepareAux1()

	subBlockSize := 1 << s.header.Levels
	for i := 0; i < subBlockSize; i++ {
		op := int(s.r.bits(5))
		switch op {
		case 0: // zeros
			for j := 0; j < s.header.SubBlocks; j++ {
				s.samples[j * subBlockSize + i] = 0
			}
		case 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16: // copy from aux buffer
			for j := 0; j < s.header.SubBlocks; j++ {
				s.samples[j * subBlockSize + i] = s.auxMem1[count + int(s.r.bits(op)) - (1 << (op - 1))]
			}
		case 17:
			for j := 0; j < s.header.SubBlocks; j++ {
				if s.r.bit() == 1 {
					if s.r.bit() == 1 {
						if s.r.bit() == 1 {
							s.samples[j * subBlockSize + i] = s.auxMem1[count + 1]
						} else {
							s.samples[j * subBlockSize + i] = s.auxMem1[count - 1]
						}
					} else {
						s.samples[j * subBlockSize + i] = 0
					}
				} else {
					s.samples[j * subBlockSize + i] = 0
					j++; if j >= s.header.SubBlocks { break }
					s.samples[j * subBlockSize + i] = 0
				}
			}
		case 18:
			for j := 0; j < s.header.SubBlocks; j++ {
				if s.r.bit() == 1 {
					if s.r.bit() == 1 {
						s.samples[j * subBlockSize + i] = s.auxMem1[count + 1]
					} else {
						s.samples[j * subBlockSize + i] = s.auxMem1[count - 1]
					}
				} else {
					s.samples[j * subBlockSize + i] = 0
				}
			}
		case 19:
			for j := 0; j < s.header.SubBlocks; j++ {
				x := Table1[s.r.bits(5)]
				s.samples[j * subBlockSize + i] = s.auxMem1[count + (x & 3) - 1]
				j++; if j >= s.header.SubBlocks { break }
				x >>= 2
				s.samples[j * subBlockSize + i] = s.auxMem1[count + (x & 3) - 1]
				j++; if j >= s.header.SubBlocks { break }
				x >>= 2
				s.samples[j * subBlockSize + i] = s.auxMem1[count + x - 1]
			}
		case 20:
			for j := 0; j < s.header.SubBlocks; j++ {
				if s.r.bit() == 1 {
					if s.r.bit() == 1 {
						if s.r.bit() == 1 {
							if s.r.bit() == 1 {
								s.samples[j * subBlockSize + i] = s.auxMem1[count + 2]
							} else {
								s.samples[j * subBlockSize + i] = s.auxMem1[count - 1]
							}
						} else {
							if s.r.bit() == 1 {
								s.samples[j * subBlockSize + i] = s.auxMem1[count + 1]
							} else {
								s.samples[j * subBlockSize + i] = s.auxMem1[count - 2]
							}
						}
					} else {
						s.samples[j * subBlockSize + i] = 0
					}
				} else {
					s.samples[j * subBlockSize + i] = 0
					j++; if j >= s.header.SubBlocks { break }
					s.samples[j * subBlockSize + i] = 0
				}
			}
		case 21:
			for j := 0; j < s.header.SubBlocks; j++ {
				if s.r.bit() == 1 {
					if s.r.bit() == 1 {
						if s.r.bit() == 1 {
							s.samples[j * subBlockSize + i] = s.auxMem1[count + 2]
						} else {
							s.samples[j * subBlockSize + i] = s.auxMem1[count - 1]
						}
					} else {
						if s.r.bit() == 1 {
							s.samples[j * subBlockSize + i] = s.auxMem1[count + 1]
						} else {
							s.samples[j * subBlockSize + i] = s.auxMem1[count - 2]
						}
					}
				} else {
					s.samples[j * subBlockSize + i] = 0
				}
			}
		case 22:
			for j := 0; j < s.header.SubBlocks; j++ {
				x := Table2[s.r.bits(7)]
				s.samples[j * subBlockSize + i] = s.auxMem1[count + (x & 7) - 2]
				j++; if j >= s.header.SubBlocks { break }
				x >>= 3
				s.samples[j * subBlockSize + i] = s.auxMem1[count + (x & 7) - 2]
				j++; if j >= s.header.SubBlocks { break }
				x >>= 3
				s.samples[j * subBlockSize + i] = s.auxMem1[count + x - 2]
			}
		case 23:
			for j := 0; j < s.header.SubBlocks; j++ {
				if s.r.bit() == 1 {
					if s.r.bit() == 1 {
						if s.r.bit() == 1 {
							b := int(s.r.bits(2))
							if b >= 2 {
								b += 3
							}
							s.samples[j * subBlockSize + i] = s.auxMem1[count - 3 + b]
						} else {
							if s.r.bit() == 1 {
								s.samples[j * subBlockSize + i] = s.auxMem1[count + 1]
							} else {
								s.samples[j * subBlockSize + i] = s.auxMem1[count - 1]
							}
						}
					} else {
						s.samples[j * subBlockSize + i] = 0
					}
				} else {
					s.samples[j * subBlockSize + i] = 0
					j++; if j >= s.header.SubBlocks { break }
					s.samples[j * subBlockSize + i] = 0
				}
			}
		case 24:
			for j := 0; j < s.header.SubBlocks; j++ {
				if s.r.bit() == 1 {
					if s.r.bit() == 1 {
						b := int(s.r.bits(2))
						if b >= 2 {
							b += 3
						}
						s.samples[j * subBlockSize + i] = s.auxMem1[count - 3 + b]
					} else {
						if s.r.bit() == 1 {
							s.samples[j * subBlockSize + i] = s.auxMem1[count + 1]
						} else {
							s.samples[j * subBlockSize + i] = s.auxMem1[count - 1]
						}
					}
				} else {
					s.samples[j * subBlockSize + i] = 0
				}
			}
		case 26:
			for j := 0; j < s.header.SubBlocks; j++ {
				if s.r.bit() == 1 {
					if s.r.bit() == 1 {
						b := int(s.r.bits(3))
						if b >= 4 {
							b++
						}
						s.samples[j * subBlockSize + i] = s.auxMem1[count - 4 + b]
					} else {
						s.samples[j * subBlockSize + i] = 0
					}
				} else {
					s.samples[j * subBlockSize + i] = 0
					j++; if j >= s.header.SubBlocks { break }
					s.samples[j * subBlockSize + i] = 0
				}
			}
		case 27:
			for j := 0; j < s.header.SubBlocks; j++ {
				if s.r.bit() == 1 {
					b := int(s.r.bits(3))
					if b >= 4 {
						b++
					}
					s.samples[j * subBlockSize + i] = s.auxMem1[count - 4 + b]
				} else {
					s.samples[j * subBlockSize + i] = 0
				}
			}
		case 29:
			for j := 0; j < s.header.SubBlocks; j++ {
				x := Table3[s.r.bits(7)]
				s.samples[j * subBlockSize + i] = s.auxMem1[count + (x & 0xf) - 5]
				j++; if j >= s.header.SubBlocks { break }
				x >>= 4
				s.samples[j * subBlockSize + i] = s.auxMem1[count + x - 5]
			}
		default:
			// todo: what shall we do?
		}
	}
}
