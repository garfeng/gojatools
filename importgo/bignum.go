package importgo

import "math"

func BigNumber(a any, b any) any {
	const div = 0x100000000
	switch a.(type) {
	case int:
		a1 := a.(int)
		b1 := b.(int)
		a2 := float64(a1)
		b2 := float64(b1)

		v := a2*div + b2

		if v > math.MaxInt {
			a3 := uint64(a1)
			b3 := uint64(b1)
			return a3*div + b3
		}
		return a1*div + b1

		/*
			if a1 > (div >> 1) {
				a2 := uint64(a1)
				b2 := uint64(b1)

				return b2 + (a2 * div)
			}
			return b1 + (a1 * div)
		*/

	case float32:
		a1 := a.(float32)
		b1 := b.(float32)
		return b1 + (a1 * div)
	case float64:
		a1 := a.(float64)
		b1 := b.(float64)
		return b1 + (a1 * div)
	default:
		panic("invalid number type")
	}
}
