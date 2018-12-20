package raspi

import "errors"

//ErrNorAValidPin When a given pin is invalid,this error will be generated.
var ErrNorAValidPin = errors.New("Given pin is invalid")

var pins = map[uint8]map[string]uint8{
	3: {
		"1": 0,
		"2": 2,
		"3": 2,
	},
	5: {
		"1": 1,
		"2": 3,
		"3": 3,
	},
	7: {
		"*": 4,
	},
	8: {
		"*": 14,
	},
	10: {
		"*": 15,
	},
	11: {
		"*": 17,
	},
	12: {
		"*": 18,
	},
	13: {
		"1": 21,
		"2": 27,
		"3": 27,
	},
	15: {
		"*": 22,
	},
	16: {
		"*": 23,
	},
	18: {
		"*": 24,
	},
	19: {
		"*": 10,
	},
	21: {
		"*": 9,
	},
	22: {
		"*": 25,
	},
	23: {
		"*": 11,
	},
	24: {
		"*": 8,
	},
	26: {
		"*": 7,
	},
	29: {
		"3": 5,
	},
	31: {
		"3": 6,
	},
	32: {
		"3": 12,
	},
	33: {
		"3": 13,
	},
	35: {
		"3": 19,
	},
	36: {
		"3": 16,
	},
	37: {
		"3": 26,
	},
	38: {
		"3": 20,
	},
	40: {
		"3": 21,
	},
}

func translatePin(pin uint8) (res uint8, err error) {
	if val, ok := pins[pin][GetBoardVersion()]; ok {
		res = val
		return
	} else if val, ok := pins[pin]["*"]; ok {
		res = val
		return
	} else {
		return emptyResult, ErrNorAValidPin
	}
}
