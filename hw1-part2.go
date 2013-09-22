package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

var wheels []*Wheel

var learnedWheels = [][]*int{}

var WHEEL_SIZES = []int{47, 53, 59, 61, 64, 65, 67, 69, 71, 73}

var REMOVE_WHITESPACE_REGEX = regexp.MustCompile(`[\n\r]`)

var interestingCharacters = map[string]struct{}{"T": {},
	"3": {},
	"4": {},
	"5": {},
	"E": {},
	"K": {},
	"Q": {},
	"6": {},
	"X": {},
	"V": {}}

//FindUniqueBitIndex takes an integer which must be a permutation of "00001" or "11110"
//and returns the index (counting from the left) of the unique bit
func FindUniqueBitIndex(i int) (index int, err error) {
	switch i {
	case 1:
		index = 4
		return
	case 2:
		index = 3
		return
	case 4:
		index = 2
		return
	case 8:
		index = 1
		return
	case 16:
		index = 0
		return
	case 30:
		index = 4
		return
	case 29:
		index = 3
		return
	case 27:
		index = 2
		return
	case 23:
		index = 1
		return
	case 15:
		index = 0
		return
	}
	return -1, fmt.Errorf("error: FindUniqueBitIndex called with an invalid integer input")
}

//We need these only because we cannot take the address of an integer literal
var ZERO = 0
var ONE = 1

var TRANSPOSITION_PATTERN = [][][]*int{
	[][]*int{
		[]*int{&ZERO, &ZERO, nil, nil, nil},
		[]*int{&ZERO, &ONE, &ZERO, nil, nil},
		[]*int{&ZERO, &ONE, &ONE, &ZERO, nil},
		[]*int{nil, nil, nil, nil, nil},
		[]*int{nil, nil, nil, nil, nil},
	},

	[][]*int{
		[]*int{nil, &ONE, nil, nil, nil},
		[]*int{nil, &ZERO, &ZERO, nil, nil},
		[]*int{nil, &ZERO, &ONE, &ZERO, nil},
		[]*int{nil, &ZERO, &ONE, &ONE, &ZERO},
		[]*int{nil, &ZERO, &ONE, &ONE, &ONE},
	},
	[][]*int{
		[]*int{nil, nil, nil, nil, nil},
		[]*int{nil, nil, &ONE, nil, nil},
		[]*int{nil, nil, &ZERO, &ZERO, nil},
		[]*int{nil, nil, &ZERO, &ONE, &ZERO},
		[]*int{nil, nil, &ZERO, &ONE, &ONE},
	},
	[][]*int{
		[]*int{nil, nil, nil, nil, nil},
		[]*int{nil, nil, nil, nil, nil},
		[]*int{nil, nil, nil, &ONE, nil},
		[]*int{nil, nil, nil, &ZERO, &ZERO},
		[]*int{nil, nil, nil, &ZERO, &ONE},
	},
	[][]*int{
		[]*int{&ONE, &ZERO, nil, nil, nil},
		[]*int{&ONE, &ONE, &ZERO, nil, nil},
		[]*int{&ONE, &ONE, &ONE, &ZERO, nil},
		[]*int{nil, nil, nil, nil, nil},
		[]*int{nil, nil, nil, nil, nil},
	},
}

//inferTransposeBits takes a source integer and a destination integer corresponding to
//(plaintext XORed with wheels 0-4) and ciphertext, respectively
//It returns what we know about the transpose wheels (wheels 5-9)
func inferTransposeBits(source, dest int) []*int {
	return TRANSPOSITION_PATTERN[source][dest]
}

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

//CurrentBit gets the current bit on the given wheel AND ticks the wheel forward to the next value
func (w *Wheel) CurrentBit() (bit int) {
	bit = w.Items[w.CurrentIndex]
	w.CurrentIndex = (w.CurrentIndex + 1) % w.MaxSize
	return
}

//Tick increments the current wheel (but ignore the actual value read from the wheel)
func (w *Wheel) Tick() {
	w.CurrentBit()
}

//TickAll takes a slice of wheels and calls Tick() on all of them
func TickAll(ws []*Wheel) {
	for _, w := range ws {
		w.Tick()
	}
}

//xorCurrentCharacter takes an integer representation of a character
//and XORs it with the current bit on each of wheel b0 through b4
//This assumes that "wheels" is a valid array of wheels of length <= 5
func xorCurrentCharacter(input int) int {
	//Iterate over each of the wheels (from the left)
	for i := 0; i < 5; i++ {
		currentBit := wheels[i].CurrentBit()
		input = input ^ currentBit<<(4-uint(i))
	}
	return input
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

func ResetWheels() {
	for _, w := range wheels {
		w.CurrentIndex = 0
	}
}

//learnedWheelToWheel converts a learnedWheel array to a Wheel struct
//Assume that all *int values are non-nil; otherwise this will panic
func learnedWheelToWheel(learnedWheel []*int) *Wheel {
	items := make([]int, len(learnedWheel))
	for i, _ := range learnedWheel {
		items[i] = *learnedWheel[i]
	}
	w := NewWheel(items)
	return w
}

//TODO these don't really need to be separate functions, as long as the format is the same (which it currently is)

func parsePlaintext(filename string) string {
	bts, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	plaintext := string(bts)

	//Strip out newlines and carriage returns from both plaintext and ciphertext
	plaintext = REMOVE_WHITESPACE_REGEX.ReplaceAllString(plaintext, "")
	return plaintext
}

func parseCiphertext(filename string) string {
	bts, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	ciphertext := string(bts)

	ciphertext = REMOVE_WHITESPACE_REGEX.ReplaceAllString(ciphertext, "")
	return ciphertext

}

//func learnFirstFiveWheels learns all spoke values from the first five wheels
//This happens to work for the plaintext/ciphertext pair that we used for testing; it is not guaranteed to work for all texts, particularly shorter texts
func learnFirstFiveWheels(plaintext string, ciphertext string) {

	//TODO don't use a global variable (learnedWheels) to store the results

	// Iterate across plaintext. For each character:
	// For each bit c0-c4:
	// Examine all possible destination bits for ci
	// Add p to implied element of appropriate spoke's pair

	// Examine all observed spoke 0-1 pairs. For all that pass a threshold, declare it 0 or 1.
	// Abort if any spoke fails this threshold.

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

				if learnedWheels[i][index%WHEEL_SIZES[i]] != nil && *learnedWheels[i][index%WHEEL_SIZES[i]] != bi {
					panic(fmt.Errorf("error: inconsistent XOR bit saved at %d for wheel %d", index, i))
				}

				learnedWheels[i][index%WHEEL_SIZES[i]] = &bi
			}

		}
	}

}

//learnEasyTransposeBits learns all of the bits in wheels 5-8, and most (but not all) of the bits in wheel 9
func learnEasyTransposeBits(plaintext, ciphertext string) {

	//Iterate over the ciphertext. If the ciphercharacter is one of T,3,4,5,E,K,Q,6,X,V,
	//we XOR the plainInt with the current state of the XOR wheels (which is known)
	//This gives a permutation of "00001" or "11110"
	//The current cipherInt must also be a (potentially different) permutation of the same two bit sequences
	//Based on where the unique bit (the unique 0 or unique 1) started and ended, we can deduce at least 2 transposed bits
	for index, plainRune := range plaintext {
		plainChar := string(plainRune)

		plainInt := alphabet[plainChar]

		cipherRune := rune(ciphertext[index])
		cipherChar := string(cipherRune)

		cipherInt := alphabet[cipherChar]

		//Check if the cipherCharacter is one of the characters we care about
		if _, present := interestingCharacters[cipherChar]; present {
			//XOR the plainInt with the current state of the XOR wheels

			xoredValue := xorCurrentCharacter(plainInt)

			sourceIndex, err := FindUniqueBitIndex(xoredValue)
			if err != nil {
				panic(err)
			}

			destIndex, err := FindUniqueBitIndex(cipherInt)
			if err != nil {
				panic(err)
			}

			inferredBits := inferTransposeBits(sourceIndex, destIndex)
			for i, bitP := range inferredBits {
				if bitP != nil {
					bit := *bitP

					//Check that we are never overwriting an existing known value with a (conflicting) known value; this should never be possible
					if learnedWheels[5+i][index%WHEEL_SIZES[5+i]] != nil && *learnedWheels[5+i][index%WHEEL_SIZES[5+i]] != bit {
						panic(fmt.Errorf("error: inconsistent transposition bit saved at %d for wheel %d", index, 5+i))
					}

					//Store the bit in the collection of learned wheels
					learnedWheels[5+i][index%WHEEL_SIZES[5+i]] = &bit
				}
			}

		} else {
			TickAll(wheels)
		}
	}

}

//learnHardTransposeBits will learn the missing transpose bits in wheel 9, assuming all of wheels 5-8 are known
//It WILL call ResetWheels() as part of its execution, which will reset wheel state.
func learnHardTransposeBits(plaintext, ciphertext string) error {
	//Reset the wheels
	//This is VERY IMPORTANT, or the learning will fail.
	ResetWheels()

	//Assume that any remaining "nil" values appear on wheel 9
	//This lets us update the missing spoke values on wheel 9, given
	//the values on wheels 5-8
	unknownSpokeIndices := map[int]struct{}{}
	for _, wheel := range learnedWheels[5:10] {
		for i, w := range wheel {
			if w == nil {
				unknownSpokeIndices[i] = struct{}{}
			}
		}
	}

	for index, plainRune := range plaintext {
		plainChar := string(plainRune)

		plainInt := alphabet[plainChar]

		cipherRune := rune(ciphertext[index])
		cipherChar := string(cipherRune)

		cipherInt := alphabet[cipherChar]

		present := false

		if _, ok := interestingCharacters[cipherChar]; ok {
			if _, ok := unknownSpokeIndices[index%WHEEL_SIZES[9]]; ok {
				present = true
				xoredValue := xorCurrentCharacter(plainInt)
				sourceIndex, err := FindUniqueBitIndex(xoredValue)
				if err != nil {
					panic(err)
				}
				destIndex, err := FindUniqueBitIndex(cipherInt)
				if err != nil {
					panic(err)
				}

				if destIndex == 4 {
					if sourceIndex == 0 {
						//We have already assumed that we know every spoke for every wheel but wheel 9 at this point
						tmp := 1 - wheels[5].Items[index%WHEEL_SIZES[5]]
						learnedWheels[9][index%WHEEL_SIZES[9]] = &tmp

						//The value is no longer unknown, so remove it from the set of unknownSpokeIndices
						delete(unknownSpokeIndices, *learnedWheels[9][index%WHEEL_SIZES[9]])
					}
					if sourceIndex == 4 {
						tmp := wheels[5].Items[index%WHEEL_SIZES[5]]
						learnedWheels[9][index%WHEEL_SIZES[9]] = &tmp
						delete(unknownSpokeIndices, *learnedWheels[9][index%WHEEL_SIZES[9]])
					}
				} else if destIndex == 3 {
					if sourceIndex == 4 {
						tmp := 1 - wheels[5].Items[index%WHEEL_SIZES[5]]
						learnedWheels[9][index%WHEEL_SIZES[9]] = &tmp
						delete(unknownSpokeIndices, *learnedWheels[9][index%WHEEL_SIZES[9]])
					} else if sourceIndex == 0 {
						tmp := wheels[5].Items[index%WHEEL_SIZES[5]]
						learnedWheels[9][index%WHEEL_SIZES[9]] = &tmp
						delete(unknownSpokeIndices, *learnedWheels[9][index%WHEEL_SIZES[9]])
					}
				}
			}
			if !present {
				//The wheels have only been incremented if xoredValue was called
				TickAll(wheels)
			}
		} else {
			TickAll(wheels)
		}
	}

	//Check if all bits of all wheels have been learned
	//If all bits of all wheels have not been learned by this point, throw an error

	for wheelIndex, wheel := range learnedWheels {
		for i, w := range wheel {
			if w == nil {
				return fmt.Errorf("error: wheel %d spoke %d is unknown", wheelIndex, i)
			}
		}
	}

	wheels = append(wheels, learnedWheelToWheel(learnedWheels[9]))
	return nil

}

func main() {

	for i := 0; i < 10; i++ {

		tmp_wheel := make([]*int, WHEEL_SIZES[i])

		learnedWheels = append(learnedWheels, tmp_wheel)
	}

	ciphertext := parseCiphertext("gwriter/part_2/ciphertext.txt")
	plaintext := parsePlaintext("gwriter/part_2/plaintext.txt")

	//Learn all bits of the first five wheels
	//The results are stored in learnedWheels (global variable)
	learnFirstFiveWheels(plaintext, ciphertext)

	//TODO check that all bit values are now known

	//Convert learnedWheels to Wheel structs
	//Store the result in the global variable "wheels"
	for _, learnedWheel := range learnedWheels[:5] {
		w := learnedWheelToWheel(learnedWheel)
		wheels = append(wheels, w)
	}

	learnEasyTransposeBits(plaintext, ciphertext)

	//Convert the last five wheels of learnedWheels to Wheel structs
	//Append the result to the global variable "wheels"
	//We can ONLY do this for wheels 5-8 right now, because there
	//are still unknown values on wheel 9
	for _, learnedWheel := range learnedWheels[5:9] {
		w := learnedWheelToWheel(learnedWheel)
		wheels = append(wheels, w)
	}

	if err := learnHardTransposeBits(plaintext, ciphertext); err != nil {
		panic(err)
	}

	for wheelIndex, wheel := range wheels {
		fmt.Printf("Wheel %d: ", wheelIndex)
		for _, w := range wheel.Items {
			fmt.Printf(" %d ", w)
		}
		fmt.Printf("\n")
	}
}
