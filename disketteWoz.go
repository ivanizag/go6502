package apple2

type disketteWoz struct {
	data    *fileWoz
	cycleOn uint64 // Cycle when the disk was last turned on
	turning bool

	latch     uint8
	position  uint32
	cycle     uint64
	trackSize uint32

	visibleLatch          uint8
	visibleLatchCountDown int8 // The visible latch stores a valid latch reading for 2 bit timings
}

func (d *disketteWoz) powerOn(cycle uint64) {
	d.turning = true
	d.cycleOn = cycle
}

func (d *disketteWoz) powerOff(_ uint64) {
	d.turning = false
}

func (d *disketteWoz) read(quarterTrack int, cycle uint64) uint8 {
	// Count cycles to know how many bits have been read
	cycles := cycle - d.cycle
	deltaBits := cycles / cyclesPerBit // TODO: Use Woz optimal bit timing

	// Process bits from woz
	// TODO: avoid processing too many bits if delta is big
	for i := uint64(0); i < deltaBits; i++ {
		d.position++
		bit := d.data.getBit(d.position, quarterTrack)
		d.latch = (d.latch << 1) + bit
		if d.latch >= 0x80 {
			// Valid byte, store value a bit longer and clear the internal latch
			//fmt.Printf("Valid 0x%.2x\n", d.latch)
			d.visibleLatch = d.latch
			d.visibleLatchCountDown = 1
			d.latch = 0
		} else if d.visibleLatchCountDown > 0 {
			// Continue showing the valid byte
			d.visibleLatchCountDown--
		} else {
			// The valid byte is lost, show the internal latch
			d.visibleLatch = d.latch
		}
	}

	//fmt.Printf("Visible: 0x%.2x, latch: 0x%.2x, bits: %v, cycles: %v\n", d.visibleLatch, d.latch, deltaBits, cycle-d.cycle)

	// Update the internal last cycle without losing the remainder not processed
	d.cycle += deltaBits * cyclesPerBit

	return d.visibleLatch
}

func (d *disketteWoz) write(quarterTrack int, value uint8, _ uint64) {
	panic("Write not implemented on woz disk implementation")
}

func loadDisquetteWoz(filename string) (*disketteWoz, error) {
	var d disketteWoz

	f, err := loadFileWoz(filename)
	if err != nil {
		return nil, err
	}
	d.data = f

	return &d, nil
}
