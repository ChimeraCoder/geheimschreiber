package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

var wheels []*Wheel

var spokeWeights = [][]*int{}

var WHEEL_SIZES = []int{47, 53, 59, 61, 64, 65, 67, 69, 71, 73}

// TRANSPOSE_PROBS[i][j] is prob ci ends up in j'th bit
var TRANSPOSE_PROBS = [][]float64{
	[]float64{0.25, 0.125, 0.0625, 0.25 + 0.03125, 0.25 + 0.03125},
	[]float64{0.5, 0.25, 0.125, 0.0625, 0.0625},
	[]float64{0, 0.5, 0.25, 0.125, 0.125},
	[]float64{0, 0, 0.5, 0.25, 0.25},
	[]float64{0.25, 0.125, 0.0625, 0.03125 + 0.25, 0.03125 + 0.025}}

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

		char := string(character)
		if char == "\n" || char == "\r" {

			result += char
			continue
		}
		encrypted, err := encryptCharacter(char)
		if err != nil {
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
		log.Printf("Cannot find character %s adf", char)
		return "", errors.New("error: character not in alphabet")
	}

	var i uint8
	for i = 0; i < 5; i++ {
		current_bit := wheels[i].CurrentBit()
		c = (c ^ (current_bit << (4 - i))) //
	}

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

func DecryptString(ciphertext string) (string, error) {
	result := ""
	for _, character := range ciphertext {

		char := string(character)
		if char == "\n" || char == "\r" {
			result += char
			continue
		}
		decrypted, err := decryptCharacter(char)
		if err != nil {
			return "", err
		}
		result += decrypted
	}
	return result, nil
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

//getNthBit returns the nth bit from the right (ie, place value 2^n)
//For example, getNthBit(2, 0) should return 0
func getNthBit(i int, n int) int {
	return (i & (1 << uint(n))) >> uint(n)
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

func OffsetWheel(wheel_index int, offset int) {
	wheels[wheel_index].CurrentIndex = offset
}

func main() {

	// Initialize wheel slice to empty pairs

	for i := 0; i < 10; i++ {

		tmp_wheel := make([]*int, WHEEL_SIZES[i])

		spokeWeights = append(spokeWeights, tmp_wheel)
	}

	// Iterate across plaintext. For each character:
	// For each bit c0-c4:
	// Examine all possible destination bits for ci
	// Add p to implied element of appropriate spoke's pair

	// Examine all observed spoke 0-1 pairs. For all that pass a threshold, declare it 0 or 1.
	// Abort if any spoke fails this threshold.

	//Initialize wheels

	bts, err := ioutil.ReadFile("gwriter/part_2/plaintext.txt")
	if err != nil {
		panic(err)
	}

	plaintext := string(bts)

	bts, err = ioutil.ReadFile("gwriter/part_2/ciphertext.txt")
	if err != nil {
		print(err)
	}

	ciphertext := string(bts)

	//Strip out newlines and carriage returns from both plaintext and ciphertext
	r := regexp.MustCompile(`[\n\r]`)
	ciphertext = r.ReplaceAllString(ciphertext, "")
	plaintext = r.ReplaceAllString(plaintext, "")

	//Iterate over plaintext and ciphertext in lockstep
	// For each cipherchar 2 or 7 encountered:
	// Learn b0-b4 and save to appropriate slot on each wheel
	// If all b0-b4 learned:
	// Else:
	for index, plainRune := range plaintext {
		plainChar := string(plainRune)

		plainInt := alphabet[plainChar]

		cipherRune := rune(ciphertext[index])
		cipherChar := string(cipherRune)

		cipherInt := alphabet[cipherChar]

		if cipherInt == 0 || cipherInt == 31 {
			//All output bits were 0, so we know that EVERY plaintext bit XORed with b_{i} to 0, for all i
			//To learn the bits b_{i}, we XOR again and store into the appropriate b_{i} slot

			mask := cipherInt ^ plainInt

			//Store each bit of mask in the appropriate b_{i} slot

			for i := 0; i < 5; i++ {
				bi := getNthBit(mask, 4-i)

				if spokeWeights[i][index%WHEEL_SIZES[i]] != nil && *spokeWeights[i][index%WHEEL_SIZES[i]] != bi {
					panic(fmt.Errorf("error: inconsistent XOR bit saved at %d for wheel %d", index, i))
				}

				spokeWeights[i][index%WHEEL_SIZES[i]] = &bi
			}

		}
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < WHEEL_SIZES[i]; j++ {
			if spokeWeights[i][j] != nil {
				fmt.Printf("%d ", *spokeWeights[i][j])
			} else {
				fmt.Printf("X")
			}
		}
		fmt.Print("\n")
	}

}
