package cmts

import (
	"fmt"
	"math"

	metro "github.com/dgryski/go-metro"
	bitset "github.com/willf/bitset"
)

// Calculate the position of the leftmost 1-bit.
func lsb(val uint64) (r uint64) {
	for r = 0; val&0x8000000000000000 == 0 && r < 64; r++ {
		val <<= 1
	}
	return 64 - r
}

// Sketch is Count-Min-Tree Sketch
type Sketch struct {
	barrier  []*bitset.BitSet
	counting []*bitset.BitSet
	nblayer  uint64
	baseSize uint64
	spire    uint64
}

// New returns a Count-Min-Tree Sketch for a given baseSize in bits
// the baseSize results in a sketch of the size 2*((2*baseSize)-1) due to the barrier bits
func New(baseSize uint) *Sketch {
	nblayer := uint64(math.Log2(float64(baseSize))) + 1
	barrier := make([]*bitset.BitSet, nblayer, nblayer)
	counting := make([]*bitset.BitSet, nblayer, nblayer)
	origBaseSize := uint64(baseSize)

	for i := uint64(0); i < nblayer; i++ {
		barrier[i] = bitset.New(baseSize)
		counting[i] = bitset.New(baseSize)
		baseSize /= 2
	}

	return &Sketch{
		barrier:  barrier,
		counting: counting,
		nblayer:  nblayer,
		baseSize: origBaseSize,
	}
}

func (sketch *Sketch) getPos(val []byte) uint {
	x := metro.Hash64(val, 0)
	pos := uint(x % sketch.baseSize)
	return pos
}

// Increment the counter for val in the sketch
func (sketch *Sketch) Increment(val []byte) {
	pos := sketch.getPos(val)
	sketch.encode(pos)
}

func (sketch *Sketch) encode(pos uint) {
	nblayer := sketch.nblayer

	nv := sketch.decode(pos) + 1 // Increment nv = value + 1
	x := (nv + 2) / 4            // x = (nv+2)/4
	lsb := lsb(x)                // lsb = lsb(x)
	var nb uint64                // nb = min(nblayer, lsb)
	if nblayer < lsb {
		nb = nblayer
	} else {
		nb = lsb
	}
	nc := nv - 2*uint64((math.Pow(2, float64(nb))-1))

	tpos := pos
	for i := uint64(0); i < nb; i++ {
		sketch.barrier[i].Set(tpos)
		tpos /= 2
	}

	nb++
	tpos = pos

	for i := uint64(0); i < nb && i < nblayer; i++ {
		sketch.counting[i].SetTo(tpos, nc%2 != 0)
		nc >>= 1
		tpos /= 2
	}

	nb = nb >> sketch.nblayer << sketch.nblayer
	if sketch.spire < nc {
		sketch.spire = nc
	}

	//fmt.Println("encode")
	//fmt.Printf("\tx: %d\tnblayer: %d\n", x, nblayer)
	//fmt.Printf("\tnv: %d\tnb: %d\tnc: %d\n", nv, nb, nc)
	//fmt.Println("spire:", sketch.spire)
}

// Get returns the frequency of val in the sketch
func (sketch *Sketch) Get(val []byte) uint64 {
	pos := sketch.getPos(val)
	return sketch.decode(pos)
}

func (sketch *Sketch) decode(pos uint) uint64 {
	var (
		b uint64
		c uint64
	)

	// The binary values of the barrier are gathered from the bottom cell of the counter, up to the first zero barrier
	// If there are 2 barrier bits contiguously set (like counter 0 in Figure 2), then b = 2.
	tpos := pos
	for _, layer := range sketch.barrier {
		if !layer.Test(tpos) {
			break
		}
		b++
		tpos /= 2
	}

	// (b + 1) value bits are gathered in counter c. If the bottom value bit is 0 and the two others are 1,
	// then the counter is c = 110b = 6
	tpos = pos
	for i := uint64(0); i < b+1 && i != sketch.nblayer; i++ {
		if sketch.counting[i].Test(tpos) {
			c |= (1 << uint64(i))
		}
		tpos /= 2
	}

	if b == sketch.nblayer {
		c = (sketch.spire << b) + c
	}

	// The real value can finally be computed: v= c + 2*(2^b − 1) = 12
	v := uint64(float64(c) + 2*(math.Pow(2, float64(b))-1) + 0.5)

	//fmt.Printf("Decode\t\tb: %d\tc: %d\n", uint64(b), c)
	return v
}

func bitreverse(x uint32, num uint) uint32 {
	number := x
	rNumber := number - number // reserve type
	for i := uint(0); i < num; i++ {
		rNumber <<= 1
		rNumber |= number & 1
		number >>= 1
	}
	return rNumber
}

func (sketch *Sketch) printSketch() {
	formatBase := "%%s %%0%db\n"
	format := fmt.Sprintf(formatBase, 64)
	fmt.Printf(format, "s", sketch.spire)
	for i := len(sketch.counting) - 1; i >= 0; i-- {
		l := sketch.barrier[i].Len()
		format := fmt.Sprintf(formatBase, l)
		for _, b := range sketch.barrier[i].Bytes() {
			fmt.Printf(format, "b", bitreverse(uint32(b), l))
		}
		for _, c := range sketch.counting[i].Bytes() {
			fmt.Printf(format, "c", bitreverse(uint32(c), l))
		}
	}
}
