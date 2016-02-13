package btchashing


import (
  "errors"
   "fmt"
   "github.com/tgebhart/goparsebtc/block"
   "github.com/tgebhart/goparsebtc/base58"
   "crypto/sha256"
   "golang.org/x/crypto/ripemd160"
   "encoding/hex"
   //"bytes"
   //"encoding/binary"
   //"strconv"
   //"github.com/tv42/base58"
)

//ComputeBlockHash computes the SHA256 double-hash of the block header
func ComputeBlockHash(Block *block.Block) (string, error) {
  hasher := sha256.New()
  slicetwo := append(Block.Header.BytePreviousBlockHash[:], Block.Header.ByteMerkleRoot[:] ...)
  slicethree := append(Block.Header.ByteTimeStamp[:], Block.Header.ByteTargetValue[:] ...)
  slicefour := append(slicethree[:], Block.Header.ByteNonce[:] ...)
  slicetwofour := append(slicetwo[:], slicefour[:] ...)
  kimbo := append(Block.Header.ByteFormatVersion, slicetwofour ...)
  hasher.Write(kimbo)
  slasher := hasher.Sum(nil)
  hasherTwo := sha256.New()
  hasherTwo.Write(slasher)
  return hex.EncodeToString(hasherTwo.Sum(nil)), nil
}

//ComputeTransactionHash computes the dual-SHA256 hash of a given transaction
func ComputeTransactionHash(Transaction *block.Transaction, inputCount uint64, outputCount uint64) (string, error) {
  hasher := sha256.New()
  var inputBytes []byte
  var outputBytes []byte
  for i := 0; i < int(inputCount); i++ {
    inputBytes = append(inputBytes[:], combineInputBytes(&Transaction.Inputs[i])[:] ...)
  }
  for o := 0; o < int(outputCount); o++ {
    outputBytes = append(outputBytes[:], combineOutputBytes(&Transaction.Outputs[o])[:] ...)
  }
  sliceone := append(Transaction.ByteTransactionVersionNumber[:], Transaction.ByteInputCount[:] ...)
  slicetwo := append(inputBytes[:], Transaction.ByteOutputCount[:] ...)
  slicethree := append(slicetwo[:], outputBytes[:] ...)
  sliceonethree := append(sliceone[:], slicethree[:] ...)
  kimbo := append(sliceonethree[:], Transaction.ByteTransactionLockTime[:] ...)
  hasher.Write(kimbo)
  slasher := hasher.Sum(nil)
  hasherTwo := sha256.New()
  hasherTwo.Write(slasher)
  return hex.EncodeToString(hasherTwo.Sum(nil)), nil
}

//ComputeCompressedBlockHash truncates a block hash to just the last half of its SHA256 hash
func ComputeCompressedBlockHash(hash string) (string) {
  ret := hash[int(len(hash)/2):]
  return ret
}

//BitcoinPublicKeyToAddress takes a 65 byte public key found in parsing addresses
//and converts it to the 20 byte form
func BitcoinPublicKeyToAddress(pubKey []byte, address *block.Address) ([]byte, []byte, error) {
  if pubKey[0] != 0x04 {
    return nil, nil, errors.New("Beginning of 65 byte public key does not match expected format")
  }
  sha1 := sha256.New()
  sha1.Write(pubKey)
  hash1 := sha1.Sum(nil)
  ripemd := ripemd160.New()
  ripemd.Write(hash1)
  hash160 := ripemd.Sum(nil)
  ret := BitcoinRipeMD160ToAddress(hash160, address)
  address.RipeMD160 = hex.EncodeToString(hash160)
  address.PublicKey = hex.EncodeToString(pubKey)
  return ret, hash160, nil
}

//BitcoinRipeMD160ToAddress takes 20 byte RipeMD160 hash and returns the 25-byte address as well as updates the address representation of the output
func BitcoinRipeMD160ToAddress(hash160 []byte, address *block.Address) ([]byte) {
  ret := append([]byte{0}, hash160[:] ...) //append network byte of 0 (main network) as checksum
  sha2 := sha256.New()
  sha2.Write(ret) //sha256 on ripemd hash
  hash3 := sha2.Sum(nil)
  sha3 := sha256.New()
  sha3.Write(hash3) // compute second hash to get checksum to store at end of output
  hash4 := sha3.Sum(nil)
  ret = append(ret[:], hash4[0:4] ...)
  address.Address = BitcoinToASCII(ret)
  address.RipeMD160 = hex.EncodeToString(hash160)
  address.PublicKey = hex.EncodeToString(hash160)
  return ret

}

//BitcoinCompressedPublicKeyToAddress takes a compressed ECDSA key and converts it to 25-byte address
func BitcoinCompressedPublicKeyToAddress(key []byte, address *block.Address) ([]byte) {
  if key[0] == 0x02 || key[0] == 0x03 {
    address.PublicKey = hex.EncodeToString(key)
    sha1 := sha256.New()
    sha1.Write(key)
    hash1 := sha1.Sum(nil)
    return BitcoinRipeMD160ToAddress(hash1, address)
  }
  fmt.Println("Invalid Compressed Public Key")
  return nil
}

//FormatAddress is the top-level method for formatting the various address types encountered while parsing the blockchain
/*func FormatAddress(address []byte) (string) {
  if len(address) == 65 {
    twentyfive, twenty, err := BitcoinPublicKeyToAddress(address)
    if err != nil {
      fmt.Println("Public key to address failed in FormatAddress", err)
    }
    fmt.Println("RipeMD160: ", hex.EncodeToString(twenty))
    return BitcoinRipeMD160ToASCII(twentyfive)
  }
  fmt.Println("RipeMD160: ", hex.EncodeToString(address))
  return BitcoinRipeMD160ToASCII(address)
}
*/

//BitcoinToASCII returns the ASCII representation of the 25 byte public key
func BitcoinToASCII(address []byte) (string) {
  return base58.HexToBase58(address)
}

func combineInputBytes(Input *block.Input) ([]byte) {
  var inputBytes []byte
  sliceone := append(Input.ByteTransactionHash[:], Input.ByteTransactionIndex[:] ...)
  slicetwo := append(Input.ByteInputScriptLength[:], Input.ByteInputScript[:] ...)
  slicethree := append(sliceone[:], slicetwo[:] ...)
  inputBytes = append(slicethree[:], Input.ByteSequenceNumber[:] ...)
  return inputBytes
}

func combineOutputBytes(Output *block.Output) ([]byte) {
  var outputBytes []byte
  sliceone := append(Output.ByteOutputValue[:], Output.ByteChallengeScriptLength[:] ...)
  outputBytes = append(sliceone[:], Output.ChallengeScriptBytes[:] ...)
  return outputBytes
}

//ParseAddressFromOutputScript parses an output's script to determine address hash
/*func ParseAddressFromOutputScript(ByteOutput block.ByteOutput, Output block.Output) (string, error) {
  var ret []byte
  var err error
  if Output.ChallengeScriptLength == 67 {
    ret, err = sixSevenByteScript(ByteOutput.ChallengeScript)
    if err != nil {
      ret, err = twentyFivePlusScript(ByteOutput.ChallengeScript)
      if err != nil {
        ret, err = searchForAddress(ByteOutput.ChallengeScript)
        if err != nil {
          ret, err = fiveErrorScript(ByteOutput.ChallengeScript)
          if err != nil {
            return "", err
          }
        }
      }
    }
    return FormatAddress(ret), nil
  }
  if Output.ChallengeScriptLength == 66 {
    ret, err = sixSixByteScript(ByteOutput.ChallengeScript)
    if err != nil {
      ret, err = twentyFivePlusScript(ByteOutput.ChallengeScript)
      if err != nil {
        ret, err = searchForAddress(ByteOutput.ChallengeScript)
        if err != nil {
          ret, err = fiveErrorScript(ByteOutput.ChallengeScript)
          if err != nil {
            return "", err
          }
        }
      }
    }
    return FormatAddress(ret), nil
  }

  if Output.ChallengeScriptLength == 20 {
    return FormatAddress(twentyScript(ByteOutput.ChallengeScript)), nil
  }

  if Output.ChallengeScriptLength == 0 {
    return block.NullHash, nil
  }

  ret, err = twentyFivePlusScript(ByteOutput.ChallengeScript)
  if err != nil {
    ret, err = searchForAddress(ByteOutput.ChallengeScript)
    if err != nil {
      return "", err
    }
  }
  return FormatAddress(ret), nil
}
*/

func sixSevenByteScript(script []byte) ([]byte, error) {
  ret := script[1:66]
  if script[0] != 65 { //hex value of coming address length does not match the raw 65 byte address hash
    return nil, errors.New("Incorrect address length at beginning of 67 byte block")
  }
  fmt.Println("----SixSevenByte----")
  return ret, nil
}

func sixSixByteScript(script []byte) ([]byte, error) {
  ret := script[0:64]
  if script[65] != 0xac { //hex value of OP_CHECKSIG
    return nil, errors.New("Incorrect OP_CHECKSIG at end of 66 byte block")
  }
  fmt.Println("----SixSixByte----")
  return ret, nil
}

func twentyFivePlusScript(script []byte) ([]byte, error) {
  ret := script[3:23]
  if script[0] != 0x76 { //hex value of OP_DUP
    return nil, errors.New("Incorrect OP_DUP at beginning of 25+ byte block")
  }
  if script[1] != 0xa9 { //hex value of OP_HASH160
    return nil, errors.New("Incorrect OP_HASH160 at beginning of 25+ byte block")
  }
  if script[2] != 0x14 { //hex value of 20
    return nil, errors.New("Possibly incorrect script length")
  }
  fmt.Println("----TwentyFivePlusByte----")
  return ret, nil
}

func fiveErrorScript(script []byte) ([]byte, error) {
  var ret []byte
  if script[2] == 0 {
    return ret, nil
  }
  fmt.Println("----FiveError----")
  return ret, errors.New("Not a five script error")
}

func twentyScript(script []byte) ([]byte) {
  return script
}

func searchForAddress(script []byte) ([]byte, error) {
  for i := 0; i < len(script) - 22; i++ {
    if script[i] == 0xA9 && script[i+1] == 0x14 {
      return script[i+2:i+22], nil
    }
  }
  fmt.Println("----SearchedAddress----")
  return nil, errors.New("Could not find address in output script")
}
