package acm

func decoder1(memory []int16, buffer []int32, subBandSize int, blockCount int) {
	for i := 0; i < subBandSize; i++ {
		db0 := int32(memory[0])
		db1 := int32(memory[1])

		workArea := buffer[:]
		for j := 0; j < (blockCount >> 2); j++ {
			row0 := int32(workArea[0])
			workArea[0] = db0 + 2 * db1 + row0
			workArea = workArea[subBandSize:]

			row1 := int32(workArea[0])
			workArea[0] = -db1 + 2 * row0 - row1
			workArea = workArea[subBandSize:]

			row2 := int32(workArea[0])
			workArea[0] = row0 + 2 * row1 + row2
			workArea = workArea[subBandSize:]

			row3 := int32(workArea[0])
			workArea[0] = -row1 + 2 * row2 - row3
			if j+1 != (blockCount >> 2) {
				workArea = workArea[subBandSize:]
			}

			db0 = row2
			db1 = row3
		}

		memory[0] = int16(db0)
		memory[1] = int16(db1)
		memory = memory[2:]
		buffer = buffer[1:]
	}
}

func decoder2(memory []int32, buffer []int32, subBandSize int, blockCount int) {
	for i := 0; i < subBandSize; i++ {
		db0 := memory[0]
		db1 := memory[1]

		workArea := buffer[:]
		for j := 0; j < (blockCount >> 2); j++ {
			row0 := workArea[0]
			workArea[0] = db0 + 2 * db1 + row0
			workArea = workArea[subBandSize:]

			row1 := workArea[0]
			workArea[0] = -db1 + 2 * row0 - row1
			workArea = workArea[subBandSize:]

			row2 := workArea[0]
			workArea[0] = row0 + 2 * row1 + row2
			workArea = workArea[subBandSize:]

			row3 := workArea[0]
			workArea[0] = -row1 + 2 * row2 - row3
			if j+1 != (blockCount >> 2) {
				workArea = workArea[subBandSize:]
			}

			db0 = row2
			db1 = row3
		}

		memory[0] = db0
		memory[1] = db1
		memory = memory[2:]
		buffer = buffer[1:]
	}
}
