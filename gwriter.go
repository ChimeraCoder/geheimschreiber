package main

import (
	"errors"
	"log"
	"strconv"
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

//Brute force
//TODO do better
func invertAlphabet(i int) (string, error) {
	for key, val := range alphabet {
		if val == i {
			return key, nil
		}
	}
	return "", errors.New("matching key not found in alphabet")
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

func EncryptString(plaintext string) (string, error) {
	result := ""
	for _, character := range plaintext {
		char, _ := strconv.Unquote(strconv.QuoteRune(character))
		encrypted, err := encryptCharacter(char)
		if err != nil {
			log.Printf("error occured on character %s", char)
			return "", err
		}
		result += encrypted
	}
	return result, nil
}

//EncryptCharacter takes a single character and encrypts it with all ten wheels in Wheels
func encryptCharacter(char string) (string, error) {

	c, ok := alphabet[char]
	if !ok {
		return "", errors.New("error: character not in alphabet")
	}

	log.Printf("Character %s gives integer %d", char, c)
	var i uint8
	for i = 0; i < 5; i++ {
		current_bit := wheels[i].CurrentBit()
		log.Printf("Offset is %d", wheels[i].CurrentIndex)
		log.Printf("Current bit is %d", current_bit)
		c = (c ^ (current_bit << (4 - i))) //
	}
	log.Printf("After XORing, C is %d", c)

	if wheels[5].CurrentBit() == 1 {
		c = interchangeBits(c, 0, 4)
	}

	if wheels[6].CurrentBit() == 1 {
		c = interchangeBits(c, 0, 1)
	}

	if wheels[7].CurrentBit() == 1 {
		c = interchangeBits(c, 1, 2)
	}

	if wheels[8].CurrentBit() == 1 {
		c = interchangeBits(c, 2, 3)
	}

	if wheels[9].CurrentBit() == 1 {
		c = interchangeBits(c, 3, 4)
	}

	encrypted_character, err := invertAlphabet(c)
	return encrypted_character, err
}

func decryptCharacter(char string) (string, error) {
	c, ok := alphabet[char]
	if !ok {
		return "", errors.New("error: character not in alphabet")
	}

	if wheels[9].CurrentBit() == 1 {
		c = interchangeBits(c, 3, 4)
	}

	if wheels[8].CurrentBit() == 1 {
		c = interchangeBits(c, 2, 3)
	}

	if wheels[7].CurrentBit() == 1 {
		c = interchangeBits(c, 1, 2)
	}

	if wheels[6].CurrentBit() == 1 {
		c = interchangeBits(c, 0, 1)
	}
	if wheels[5].CurrentBit() == 1 {
		c = interchangeBits(c, 0, 4)
	}

	//Order of XOR doesn't matter
	var i uint8
	for i = 0; i < 5; i++ {
		current_bit := wheels[i].CurrentBit()
		c = (c ^ (current_bit << (4 - i))) //
	}

	decrypted_character, err := invertAlphabet(c)
	return decrypted_character, err

}

//interchangeBits takes a uint8 (c) with only FIVE significant bits
//and interchanges the ith and jth bit
//i must be less than j
func interchangeBits(c int, i uint8, j uint8) int {

	//Get the ith digit of c
	ci_tmp := c & (16 >> i)
	//Get the jth digit of c
	cj_tmp := c & (16 >> j)

	c = c &^ (16 >> i)
	c = c | (cj_tmp << (j - i))

	//Set cj to be the OLD value of ci
	c = c &^ (16 >> j)
	c = c | (ci_tmp >> (j - i))
	return c
}

func ResetWheels(order []int) {
	wheels = make([]*Wheel, 10)

	for i := 0; i < 10; i++ {
		wheels[i] = NewWheel(wheel_values[order[i]])
	}
}

func OffsetWheel(wheel_index int, offset int) {
	wheels[wheel_index].CurrentIndex = offset
}

func main() {

	//Initialize wheels
	ResetWheels([]int{8, 7, 2, 4, 3, 5, 6, 1, 0, 9})

	OffsetWheel(0, 44)
	OffsetWheel(1, 52)
	OffsetWheel(2, 35)
	OffsetWheel(3, 14)
	OffsetWheel(4, 19)
	OffsetWheel(5, 55)
	OffsetWheel(6, 6)
	OffsetWheel(7, 4)
	OffsetWheel(8, 3)
	OffsetWheel(9, 51)

	encrypted, err := EncryptString("UMUM4VEVE35KING4HENRY4IV35UMUM4VEVE")
	if err != nil {
		panic(err)
	}
	log.Print(encrypted)

}
