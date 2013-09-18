package main

import (
	"errors"
	"log"
)

var wheels []*Wheel

var wheel_values = [][]int{
	[]int{1, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 1, 1, 1, 0, 1, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0},
	[]int{1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	[]int{0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 1, 1},
	[]int{0, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
	[]int{1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 0, 0, 1, 0},
	[]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 1, 1, 0, 1, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 0, 1, 0},
	[]int{1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 0},
	[]int{1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0},
	[]int{1, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1},
	[]int{0, 1, 0, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1},
}

var alphabet = map[string]int{
	"2": 0,
	"T": 1,
	"3": 2,
	"O": 3,
	"4": 4,
	"H": 5,
	"N": 6,
	"M": 7,
	"5": 8,
	"L": 9,
	"R": 10,
	"G": 11,
	"I": 12,
	"P": 13,
	"C": 14,
	"V": 15,
	"E": 16,
	"Z": 17,
	"D": 18,
	"B": 19,
	"S": 20,
	"Y": 21,
	"F": 22,
	"X": 23,
	"A": 24,
	"W": 25,
	"J": 26,
	"6": 27,
	"U": 28,
	"Q": 29,
	"K": 30,
	"7": 31,
}

type Wheel struct {
	Items        []int
	CurrentIndex int //Starts at zero
	MaxSize      int
}

func NewWheel(items []int) (w *Wheel) {
	w = new(Wheel)
	w.MaxSize = len(items)

	w.Items = items

	w.CurrentIndex = 0
	return w

}

func (w *Wheel) CurrentBit() (bit int) {
	bit = w.Items[w.CurrentIndex]
	w.CurrentIndex = (w.CurrentIndex + 1) % w.MaxSize
	return
}

//EncryptCharacter takes a single character and encrypts it with all ten wheels in Wheels
func EncryptCharacter(char string) error {

	c, ok := alphabet[char]
	if !ok {
		return errors.New("error: character not in alphabet")
	}
	return nil

	var i uint8
	for i = 0; i < 5; i++ {
		c = (c ^ wheels[i].CurrentBit()) << (4 - i) //
	}

    
    if wheels[5].CurrentBit() == 1 {
        //Interchange c0 and c4

        //Bitwise AND with "10000" tells us if c0 is 1

        //c0_tmp is 16 if c0 is 1 (if 16ths place is 1)
        
        c0_tmp := c & (1 << 4)
        c4_tmp := c & (1 << 0)



        /**
        10101 original
     &  01111 mask       (c4_tmp << 4)
        00101

        00000 new value
     &  10000 mask
        00000
        **/

        //Set c0 to be the OLD value of c4
        c = c &^ (1 << 4)
        c = c | (c4_tmp << 4)

        //Set c4 to be the OLD value of c0
        c = c &^ (1 << 0)
        c = c | (c0_tmp >> 4)
    }

    if wheels[5].CurrentBit() == 1 {

    }

    if wheels[5].CurrentBit() == 1 {
        
    }

    if wheels[5].CurrentBit() == 1 {
        
    }

    if wheels[5].CurrentBit() == 1 {
        
    }



	return nil
}

func main() {
	//TODO initialize Wheels

	wheels = make([]*Wheel, 10)

    for i := 0; i < 10; i++{
        wheels[i] = NewWheel(wheel_values[i])
    }


	for j := 0; j < 10; j++ {
		log.Print(wheels[0].CurrentBit())
	}

    for _, wheel := range wheels {
        log.Print(wheel.Items)
        
    }

	log.Print(1 ^ 1)
}
