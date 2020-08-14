package apple2

type card interface {
	loadRom(data []uint8)
	assign(a *Apple2, slot int)
}

type cardBase struct {
	a        *Apple2
	romCsxx  *memoryRangeROM
	romC8xx  *memoryRangeROM
	romCxxx  *memoryRangeROM
	slot     int
	_ssr     [16]softSwitchR
	_ssw     [16]softSwitchW
	_ssrName [16]string
	_sswName [16]string
}

func (c *cardBase) loadRom(data []uint8) {
	if c.a != nil {
		panic("Assert failed. Rom must be loaded before inserting the card in the slot")
	}
	if len(data) == 0x100 {
		// Just 256 bytes in Cs00
		c.romCsxx = newMemoryRangeROM(0, data)
	} else if len(data) == 0x800 {
		// The file has C800 to C8FF
		// The 256 bytes in Cx00 are copied from the first page in C800
		c.romCsxx = newMemoryRangeROM(0, data)
		c.romC8xx = newMemoryRangeROM(0xc800, data)
	} else if len(data) == 0x1000 {
		// The file covers the full Cxxx range. Only showing the page
		// corresponding to the slot used.
		c.romCxxx = newMemoryRangeROM(0xc000, data)
	} else {
		panic("Invalid ROM size")
	}
}

func (c *cardBase) assign(a *Apple2, slot int) {
	c.a = a
	c.slot = slot
	if slot != 0 {
		if c.romCsxx != nil {
			// Relocate to the assigned slot
			//c.romCsxx.base = uint16(0xc000 + slot*0x100)
			c.romCsxx.setBase(uint16(0xc000 + slot*0x100))
			a.mmu.setCardROM(slot, c.romCsxx)
		}
		if c.romC8xx != nil {
			a.mmu.setCardROMExtra(slot, c.romC8xx)
		}
		if c.romCxxx != nil {
			a.mmu.setCardROM(slot, c.romCxxx)
			a.mmu.setCardROMExtra(slot, c.romCxxx)
		}
	}

	for i := 0; i < 0x10; i++ {
		a.io.addSoftSwitchR(uint8(0xC80+slot*0x10+i), c._ssr[i], c._ssrName[i])
		a.io.addSoftSwitchW(uint8(0xC80+slot*0x10+i), c._ssw[i], c._sswName[i])
	}
}

func (c *cardBase) addCardSoftSwitchR(address uint8, ss softSwitchR, name string) {
	c._ssr[address] = ss
	c._ssrName[address] = name
}

func (c *cardBase) addCardSoftSwitchW(address uint8, ss softSwitchW, name string) {
	c._ssw[address] = ss
	c._sswName[address] = name
}
