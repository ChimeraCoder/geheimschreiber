package geheimschreiber

import (
	"io/ioutil"
	"testing"
)

const TEST_CIPHERTEXT_FILE = "test_ciphertext.txt"
const TEST_PLAINTEXT_FILE = "test_plaintext.txt"

var TEST_CIPHERTEXT_SOLVED_WHEELS = []*Wheel{NewWheel([]int{0, 0, 1, 1, 1, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 1}),
	NewWheel([]int{0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0}),
	NewWheel([]int{1, 1, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 1, 1, 0, 0}),
	NewWheel([]int{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1, 1}),
	NewWheel([]int{1, 1, 1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0, 1, 0, 1, 1, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 0}),
	NewWheel([]int{0, 1, 1, 0, 0, 1, 1, 0, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 1}),
	NewWheel([]int{0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 1, 1, 1, 0, 0, 0}),
	NewWheel([]int{1, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1}),
	NewWheel([]int{1, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 0}),
	NewWheel([]int{1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1}),
}

func Test_LoadingCiphertext(t *testing.T) {

	ciphertext, plaintext := parseCiphertext(TEST_CIPHERTEXT_FILE)

	//Plaintext is generated from reading ciphertext
	//We know it starts off with UMUM4VEVE35
	if plaintext[:11] != "UMUM4VEVE35" {
		t.Error("Error reading test ciphertext")
	}
	if ciphertext[:11] != "BTEVUIO7WGR" {
		t.Error("Error reading test ciphertext")
	}
}

func Test_WheelsEqual(t *testing.T) {
	wheel := NewWheel([]int{0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0})
	wheel2 := NewWheel([]int{0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0})
	wheel3 := NewWheel([]int{1, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0})
	if !wheel.Equals(*wheel) {
		t.Error("Wheel equality failed - wheel is not equal to itself")
	}

	if !wheel.Equals(*wheel2) {
		t.Error("Wheel equality failed - wheel is not equal to identical other wheel")
	}

	if wheel.Equals(*wheel3) || wheel2.Equals(*wheel3) {
		t.Error("Wheel equality failed - wheel is equal to a non-identical other wheel")
	}

}

func Test_Decryption(t *testing.T) {

	wheels := crackMessage(TEST_CIPHERTEXT_FILE)

	//Plaintext is generated from reading ciphertext
	//We know it starts off with UMUM4VEVE35
	for i, wheel := range wheels {
		if !wheel.Equals(*TEST_CIPHERTEXT_SOLVED_WHEELS[i]) {
			t.Error("Error decoding ciphertext: wheel %d does not match expected result")
		}
	}
}

func Test_Encryption(t *testing.T) {
	bts, err := ioutil.ReadFile(TEST_PLAINTEXT_FILE)
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	plaintext := string(bts)

	bts, err = ioutil.ReadFile(TEST_CIPHERTEXT_FILE)
	if err != nil {
		t.Errorf("Error reading file: %s", err.Error())
	}
	ciphertext := string(bts)

	result, err := EncryptString(TEST_CIPHERTEXT_SOLVED_WHEELS, plaintext)
	if err != nil {
		t.Errorf("Error encrypting string: %s", err.Error())
	}

	if result != ciphertext {
		diff := -1
		for i := 0; i < len(result); i++ {
			if result[i] != ciphertext[i] {
				diff = i
				break
			}
		}
		t.Errorf("Encrypted plaintext (length %d) does not match target ciphertext (length %d) - first error at char %d ", len(result), len(ciphertext), diff)
	}
}
