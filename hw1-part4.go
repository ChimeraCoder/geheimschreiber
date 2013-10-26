package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
)

var wheels []*Wheel

var learnedWheels = [][]*int{}

var WHEEL_SIZES = []int{47, 53, 59, 61, 64, 65, 67, 69, 71, 73}
var LARGE_WHEEL_SIZES = []int{26996, 26996, 26996, 26996, 26996, 26996, 26996, 26996, 26996, 26996}

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

func (w Wheel) Equals(other Wheel) bool {
    if w.MaxSize != other.MaxSize {
        return false
    }
    for i, item := range w.Items{
        if item != other.Items[i]{
            return false
        }
    }
    return true
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

//Via http://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func parseCiphertext(filename string) (string, string) {

	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		panic(err)
	}

	ciphertext := ""
	plaintext := ""
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		currentLine := scanner.Text()
		//Assume every line is at least 13 characters

		plaintext += "UMUM4VEVE35"
		ciphertext += currentLine[:11]

		for _, _ = range currentLine[12 : len(currentLine)-1] {
			ciphertext += "-"
			plaintext += "-"
		}
		ciphertext += currentLine[len(currentLine)-2:]
		plaintext += "35"

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	ciphertext = REMOVE_WHITESPACE_REGEX.ReplaceAllString(ciphertext, "")
	return ciphertext, plaintext

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
		if plainChar == "-" {
			continue
		}

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

				if learnedWheels[i][index%LARGE_WHEEL_SIZES[i]] != nil && *learnedWheels[i][index%LARGE_WHEEL_SIZES[i]] != bi {
					panic(fmt.Errorf("error: inconsistent XOR bit saved at %d for wheel %d", index, i))
				}

				learnedWheels[i][index%LARGE_WHEEL_SIZES[i]] = &bi
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
		if plainChar == "-" {
			TickAll(wheels)
			continue
		}

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
					//Unlike Part 2, we don't need to take the modulus because the wheel size has been set to the length of the text
					if learnedWheels[5+i][index] != nil && *learnedWheels[5+i][index] != bit {
						panic(fmt.Errorf("error: inconsistent transposition bit saved at %d for wheel %d", index, 5+i))
					}

					//Store the bit in the collection of learned wheels
					learnedWheels[5+i][index] = &bit
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
			if _, ok := unknownSpokeIndices[index%LARGE_WHEEL_SIZES[9]]; ok {
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
						tmp := 1 - wheels[5].Items[index%LARGE_WHEEL_SIZES[5]]
						learnedWheels[9][index%LARGE_WHEEL_SIZES[9]] = &tmp

						//The value is no longer unknown, so remove it from the set of unknownSpokeIndices
						delete(unknownSpokeIndices, *learnedWheels[9][index%LARGE_WHEEL_SIZES[9]])
					}
					if sourceIndex == 4 {
						tmp := wheels[5].Items[index%LARGE_WHEEL_SIZES[5]]
						learnedWheels[9][index%LARGE_WHEEL_SIZES[9]] = &tmp
						delete(unknownSpokeIndices, *learnedWheels[9][index%LARGE_WHEEL_SIZES[9]])
					}
				} else if destIndex == 3 {
					if sourceIndex == 4 {
						tmp := 1 - wheels[5].Items[index%LARGE_WHEEL_SIZES[5]]
						learnedWheels[9][index%LARGE_WHEEL_SIZES[9]] = &tmp
						delete(unknownSpokeIndices, *learnedWheels[9][index%LARGE_WHEEL_SIZES[9]])
					} else if sourceIndex == 0 {
						tmp := wheels[5].Items[index%LARGE_WHEEL_SIZES[5]]
						learnedWheels[9][index%LARGE_WHEEL_SIZES[9]] = &tmp
						delete(unknownSpokeIndices, *learnedWheels[9][index%LARGE_WHEEL_SIZES[9]])
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

//removePossibleWheelState will remove the impossibleSize from the list of possible sizes for wheel with index wheelIndex in the set of possible sizes
func removePossibleWheelState(possibleSizes []map[int]struct{}, wheelIndex, impossibleSize int) []map[int]struct{} {

	//Remove it from the list of possible sizes for the specified wheel iff it is present
	if _, ok := possibleSizes[wheelIndex][impossibleSize]; ok {
		delete(possibleSizes[wheelIndex], impossibleSize)

		//If there is only one possibility left, we know this wheel with certainty
		//Delete this size from all other wheels
		if len(possibleSizes[wheelIndex]) == 1 {

			//TODO figure out better hack
			var actualSize int
			for k, _ := range possibleSizes[wheelIndex] {
				actualSize = k
				break
			}

			for i, _ := range possibleSizes {
				if i != wheelIndex {
					possibleSizes = removePossibleWheelState(possibleSizes, i, actualSize)
				}
			}
		}
	}
	return possibleSizes
}

//TODO rename this
func crackMessage(filename string) []*Wheel{

	//Use the LARGE wheel sizes instead of the regular wheel sizes this time
	for i := 0; i < 10; i++ {
		tmp_wheel := make([]*int, LARGE_WHEEL_SIZES[i])

		learnedWheels = append(learnedWheels, tmp_wheel)
	}

	ciphertext, plaintext := parseCiphertext(filename)

	//Learn all bits of the first five wheels
	//The results are stored in learnedWheels (global variable)
	learnFirstFiveWheels(plaintext, ciphertext)

	POSSIBLE_SIZES := make([]map[int]struct{}, 10)
	for i, _ := range POSSIBLE_SIZES {
		POSSIBLE_SIZES[i] = map[int]struct{}{
			47: {},
			53: {},
			59: {},
			61: {},
			64: {},
			65: {},
			67: {},
			69: {},
			71: {},
			73: {},
		}
	}

	for i := 0; i < 5; i++ {
		for size := range POSSIBLE_SIZES[i] {
			for spokeIndex := 0; spokeIndex < size; spokeIndex++ {
				impliedBit := -1
				for index := spokeIndex; index < len(plaintext); index += size {
					currentBit := learnedWheels[i][index]
					if currentBit != nil {
						if impliedBit == -1 {
							impliedBit = *currentBit
						} else if impliedBit != *currentBit {
							//This means we have found a conflict
							//Remove this wheel size from the pool of possible wheel sizes for this wheel
							POSSIBLE_SIZES = removePossibleWheelState(POSSIBLE_SIZES, i, size)
						}
					}
				}
			}
		}
	}

	//At this point, all of the first five wheels should be known
	//This is not always the case, but it will be the case for the current input

	//Now, we know the bits of wheels 0-4, but they are in the wrong locations
	//Set those bits to the correct locations
	for i, _ := range POSSIBLE_SIZES[:5] {
		//TODO figure out better hack
		var wheelSize int
		for k, _ := range POSSIBLE_SIZES[i] {
			wheelSize = k
			break
		}

		for j, _ := range learnedWheels[i] {
			if learnedWheels[i][j] != nil {
				learnedWheels[i][j%wheelSize] = learnedWheels[i][j]
			}
		}
	}

	//Truncate the first five wheels to their correct sizes
	for i, _ := range learnedWheels[:5] {
		//TODO figure out better hack
		var wheelSize int
		for k, _ := range POSSIBLE_SIZES[i] {
			wheelSize = k
			break
		}
		learnedWheels[i] = learnedWheels[i][:wheelSize]
	}

	//At this point, we know all the bits on the first five wheels
	//AND they are in the correct locations

	//We are ready to create Wheel structs for the first five wheels
	wheels = make([]*Wheel, 5)
	//Create actual Wheel structs for these fully-learned wheels and append them to "wheels" (global variable)
	for wheelIndex, lw := range learnedWheels[:5] {
		items := make([]int, len(lw))
		for i, item := range lw {
			items[i] = *item
		}
		wheel := NewWheel(items)
		wheels[wheelIndex] = wheel
	}

	for _, w := range wheels[:5] {
		log.Print(*w)
	}

	//Now, we need to do the same thing, but for wheels 5-9

	//Learn the transpose bits, though they will not be in the correct locations
	learnEasyTransposeBits(plaintext, ciphertext)

	//Figure out the actual sizes for wheels 5-9, so we can set the learned bits to be in the correct locations
	for i := 5; i < 10; i++ {
		for size := range POSSIBLE_SIZES[i] {
			for spokeIndex := 0; spokeIndex < size; spokeIndex++ {
				impliedBit := -1
				for index := spokeIndex; index < len(plaintext); index += size {
					currentBit := learnedWheels[i][index]
					if currentBit != nil {
						if impliedBit == -1 {
							impliedBit = *currentBit
						} else if impliedBit != *currentBit {
							//This means we have found a conflict
							//Remove this wheel size from the pool of possible wheel sizes for this wheel
							POSSIBLE_SIZES = removePossibleWheelState(POSSIBLE_SIZES, i, size)
						}
					}
				}
			}
		}
	}

	//At this point, all sizes for all wheels are known

	//Now, we know the bits of wheels 5-9, but they are in the wrong locations
	//Set those bits to the correct locations
	for i, _ := range POSSIBLE_SIZES[5:] {
		//TODO figure out better hack
		var wheelSize int
		for k, _ := range POSSIBLE_SIZES[5+i] {
			wheelSize = k
			break
		}

		for j, _ := range learnedWheels[5+i] {
			if learnedWheels[5+i][j] != nil {
				learnedWheels[5+i][j%wheelSize] = learnedWheels[5+i][j]
			}
		}
	}

	//Truncate the first five wheels to their correct sizes
	for i, _ := range learnedWheels[5:] {
		//TODO figure out better hack
		var wheelSize int
		for k, _ := range POSSIBLE_SIZES[5+i] {
			wheelSize = k
			break
		}
		learnedWheels[5+i] = learnedWheels[5+i][:wheelSize]
	}

	//We are ready to create Wheel structs for wheels 5-9
	//Create actual Wheel structs for these fully-learned wheels and append them to "wheels" (global variable)
	for _, lw := range learnedWheels[5:] {
		items := make([]int, len(lw))
		for i, item := range lw {
			items[i] = *item
		}
		wheel := NewWheel(items)
		wheels = append(wheels, wheel)
	}

	for _, w := range wheels[5:] {
		log.Print(w.Items)
	}

	fmt.Printf("package main\nvar PART_4_SOLVED_WHEELS = []*Wheel{")
	for _, wheel := range wheels {
		fmt.Printf("NewWheel([]int{")
		for _, w := range wheel.Items {
			fmt.Printf("%d, ", w)
		}
		fmt.Printf("}),\n")
	}
	fmt.Printf("}\n")
    return wheels
}

func main(){
    crackMessage("test_ciphertext.txt")
}
