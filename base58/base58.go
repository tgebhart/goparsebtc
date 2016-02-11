package base58


import "math/big"

//Useful materials:
//https://en.bitcoin.it/wiki/Base_58_Encoding

//alphabet used by Bitcoins
var alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

//Base58 holds the Base58 string
type Base58 string

//reverse alphabet used for quckly converting base58 strings into numbers
var revalp = map[string]int{
	"1": 0, "2": 1, "3": 2, "4": 3, "5": 4, "6": 5, "7": 6, "8": 7, "9": 8, "A": 9,
	"B": 10, "C": 11, "D": 12, "E": 13, "F": 14, "G": 15, "H": 16, "J": 17, "K": 18, "L": 19,
	"M": 20, "N": 21, "P": 22, "Q": 23, "R": 24, "S": 25, "T": 26, "U": 27, "V": 28, "W": 29,
	"X": 30, "Y": 31, "Z": 32, "a": 33, "b": 34, "c": 35, "d": 36, "e": 37, "f": 38, "g": 39,
	"h": 40, "i": 41, "j": 42, "k": 43, "m": 44, "n": 45, "o": 46, "p": 47, "q": 48, "r": 49,
	"s": 50, "t": 51, "u": 52, "v": 53, "w": 54, "x": 55, "y": 56, "z": 57,
}

//ToBig converts base58 to big.Int
func (b Base58) ToBig() *big.Int {
	answer := new(big.Int)
	for i := 0; i < len(b); i++ {
		answer.Mul(answer, big.NewInt(58))                              //multiply current value by 58
		answer.Add(answer, big.NewInt(int64(revalp[string(b[i:i+1])]))) //add value of the current letter
	}
	return answer
}

//ToHex converts base58 to hex bytes
func (b Base58) ToHex() []byte {
	value := b.ToBig() //convert to big.Int
	oneCount := 0
	for string(b)[oneCount] == '1' {
		oneCount++
	}
	return append(make([]byte, oneCount), value.Bytes()...) //convert big.Int to bytes
}

//ToHex convert base58 to hex bytes
func ToHex(b string) []byte {
	return Base58(b).ToHex()
}

//BitHex converts base58 to hexes used by Bitcoins (keeping the zeroes on the front, 25 bytes long)
func (b Base58) BitHex() []byte {
	value := b.ToBig() //convert to big.Int

	tmp := value.Bytes() //convert to hex bytes
	if len(tmp) == 25 {  //if it is exactly 25 bytes, return
		return tmp
	} else if len(tmp) > 25 { //if it is longer than 25, return nothing
		return nil
	}
	answer := make([]byte, 25)      //make 25 byte container
	for i := 0; i < len(tmp); i++ { //copy converted bytes
		answer[24-i] = tmp[len(tmp)-1-i]
	}
	return answer
}

//BigToBase58 encodes big.Int to base58 string
func BigToBase58(val *big.Int) Base58 {
	answer := ""
	valCopy := new(big.Int).Abs(val) //copies big.Int

	if val.Cmp(big.NewInt(0)) <= 0 { //if it is less than 0, returns empty string
		return Base58("")
	}

	tmpStr := ""
	tmp := new(big.Int)
	for valCopy.Cmp(big.NewInt(0)) > 0 { //converts the number into base58
		tmp.Mod(valCopy, big.NewInt(58))                //takes modulo 58 value
		valCopy.Div(valCopy, big.NewInt(58))            //divides the rest by 58
		tmpStr += alphabet[tmp.Int64() : tmp.Int64()+1] //encodes
	}
	for i := (len(tmpStr) - 1); i > -1; i-- {
		answer += tmpStr[i : i+1] //reverses the order
	}
	return Base58(answer) //returns
}

//HexToBase58 encodes hex bytes into base58
func HexToBase58(val []byte) (string) {
	tmp := BigToBase58(HexToBig(val)) //encoding of the number without zeroes in front

	//looking for zeros at the beginning
	i := 0
	for i = 0; val[i] == 0 && i < len(val); i++ {
	}
	answer := ""
	for j := 0; j < i; j++ { //adds zeroes from the front
		answer += alphabet[0:1]
	}
	answer += string(tmp) //concatenates

	return string(answer) //returns
}

//HexToBig encodes hex representation to big.int
func HexToBig(b []byte) *big.Int{
	answer:=big.NewInt(0)

	for i:=0; i<len(b); i++{
		answer.Lsh(answer, 8)
		answer.Add(answer, big.NewInt(int64(b[i])))
	}

	return answer
}
/*

//Alphabet holds the reference alphabet for the bas58 packages
const Alphabet string = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
//EncodedZero holds the assumed zero element for string representation in base58
const EncodedZero string = "1"
//IndexLength holds the length of base58's indexing array
const IndexLength int = 128
//holds index array for alphabet conversion
var indexes [IndexLength]int

func init() {
  instantiateIndexes()
}

func instantiateIndexes() {
  var indexes [IndexLength]int
  for i := 0; i < len(indexes); i++ {
    indexes[i] = -1
  }
  for i := 0; i < len(Alphabet); i++ {
    indexes[Alphabet[i]] = i
  }
}





//EncodeBase58 formats a byte array into its ASCII base58 equivalent. Can be used to decode addresses
func EncodeBase58(input []byte) (string) {
  if len(input) == 0 {
    return ""
  }
  //Count leading zeros
  var zeros int
  for zeros < len(input) && input[zeros] == 0 {
    zeros++
  }
  //Convert base256 digits to base58 digits and convert to ASCII
  var inputcopy []byte
  copy(inputcopy, input)
  var encoded []string
  ostart := len(input) * 2 //upper bound
  for istart := zeros; istart < len(input); istart++ {
    ostart = ostart - 1
    encoded[ostart] = Alphabet[divmod(input, istart, 256, 58)]
    if inputcopy[istart] == 0 {
      istart++ //skip leading zeros
    }
  }
  //keep number of leading encoded zeros in output as there were in input
  for ostart < len(encoded) && encoded[ostart] == EncodedZero {
    ostart++
  }
  zeros = zeros - 1
  for zeros >= 0 {
    ostart = ostart - 1
    encoded[ostart] = EncodedZero
    zeros = zeros - 1
  }
  return string(encoded[ostart : len(encoded) - ostart])
}

//Divides a number, represented as an array of bytes each containing a
//single digit in the specifed base, by the given divisor. The given number is
//modified in place to contain the quotient, and the return value is the remainder.
func divmod(number []byte, firstDigit int, base int, divisor int) (byte) {
  //implementation of long division
  remainder := 0
  for i := firstDigit; i < len(number); i++ {
    digit := int(number[i] & 0xFF)
    temp := remainder * base + digit
    number[i] = byte(temp / divisor)
    remainder = temp % divisor
  }
  return byte(remainder)
}
*/
