package blockchainbuilder

import (
   "errors"
    "fmt"
    "os"
    "github.com/tgebhart/goparsebtc/block"
    "github.com/tgebhart/goparsebtc/filefunctions"
    "github.com/tgebhart/goparsebtc/blockvalidation"
    "crypto/sha256"
    "encoding/hex"
)

//Blockchain interface holds the top-level blockchain functions
type Blockchain interface {

  readMagicNumber(file *os.File) (uint32, error)
  validateMagicNumber(pmagicNumber uint32) (bool)

  readBlockLength(file *os.File) (uint32, error)
  validateBlockLength(blockLength uint32) (bool)

  readFormatVersion(file *os.File) (uint32, error)
  validateFormatVersion(formatVersion uint32) (bool)

  readPreviousBlockHash(file *os.File) ([]uint8, error)


}


func readMagicNumber(file *os.File) (uint32, error) {

  var magicNumber uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &magicNumber)
  if err != nil {
    fmt.Println("binary.Read failed:", err)
  }
  if blockvalidation.ValidateMagicNumber(magicNumber) {
    return magicNumber, nil
  }
  return 0, errors.New("unusual or invalid magic number value")
}

func readBlockLength(file *os.File) (uint32, error) {
  var blockLength uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &blockLength)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateBlockLength(blockLength) {
      return blockLength, nil
  }
  return 0, errors.New("Very large (or no) block length")
}

func readFormatVersion(file *os.File) (uint32, []byte, error) {
  var formatVersion uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &formatVersion)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateFormatVersion(formatVersion) {
      return formatVersion, b, nil
  }
  return 0, nil, errors.New("Unusual format version")
}

func readPreviousBlockHash(file *os.File) (string, []byte, error) {
  var previousBlockHash string
  b := filefunctions.ReadNextBytes(file, 32)
  filefunctions.ReadUInt8ByteArrayLength32(b, &previousBlockHash)
  return previousBlockHash, b, nil
}

func readMerkleRoot(file *os.File) (string, []byte, error) {
  var merkleRoot string
  b := filefunctions.ReadNextBytes(file, 32)
  filefunctions.ReadUInt8ByteArrayLength32(b, &merkleRoot)
  return merkleRoot, b, nil
}

func readTimeStamp(file *os.File) (uint32, []byte, error) {
  var timeStamp uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &timeStamp)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateTimeStamp(timeStamp) {
    return timeStamp, b, nil
  }
  return 0, nil, errors.New("Unexpected timestamp value")
}

func readTargetValue(file *os.File) (uint32, []byte, error) {
  var targetValue uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &targetValue)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return targetValue, b, nil
}

func readNonce(file *os.File) (uint32, []byte, error) {
  var nonce uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &nonce)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return nonce, b, nil
}

func readTransactionCount(file *os.File) (uint64, error) {
  var transactionLength uint64
  transactionLength, err := filefunctions.ReadVariableLengthInteger(file)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return transactionLength, nil
}

func readTransactionVersion(file *os.File) (uint32, error) {
  var transactionVersion uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &transactionVersion)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateTransactionVersion(transactionVersion) {
    return transactionVersion, nil
  }
  return 0, errors.New("Unexpected transaction version number")
}

func readInputCount(file *os.File) (uint64, error) {
  var inputCount uint64
  inputCount, err := filefunctions.ReadVariableLengthInteger(file)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return inputCount, nil
}

func readTransactionHash(file *os.File) (string, error) {
  var transactionHash string
  b := filefunctions.ReadNextBytes(file, 32)
  filefunctions.ReadUInt8ByteArrayLength32(b, &transactionHash)

  return transactionHash, nil
}

func readTransactionIndex(file *os.File) (uint32, error) {
  var transactionIndex uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &transactionIndex)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return transactionIndex, nil
}

func readInputScriptLength(file *os.File) (uint64, error) {
  var inputScriptLength uint64
  inputScriptLength, err := filefunctions.ReadVariableLengthInteger(file)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return inputScriptLength, nil
}

func readInputScriptBytes(inputScriptLength int, file *os.File) ([]uint8, error) {
  var inputScriptBytes = make([]uint8, inputScriptLength)
  b := filefunctions.ReadNextBytes(file, inputScriptLength)
  err := filefunctions.ReadUInt8ByteArray(b, &inputScriptBytes)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return inputScriptBytes, nil
}

func readSequenceNumber(file *os.File) (uint32, error) {
  var sequenceNumber uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &sequenceNumber)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateSequenceNumber(sequenceNumber) {
    return sequenceNumber, nil
  }
  return 0, errors.New("Invalid sequence number")
}

func readOutputCount(file *os.File) (uint64, error) {
  var outputCount uint64
  outputCount, err := filefunctions.ReadVariableLengthInteger(file)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return outputCount, nil
}

func readOutputValue(file *os.File) (uint64, error) {
  var outputValue uint64
  b := filefunctions.ReadNextBytes(file, 8)
  err := filefunctions.ReadBinaryToUInt64(b, &outputValue)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return outputValue, nil
}

func readChallengeScriptLength(file *os.File) (uint64, error) {
  var challengeScriptLength uint64
  challengeScriptLength, err := filefunctions.ReadVariableLengthInteger(file)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return challengeScriptLength, nil
}

func readChallengeScriptBytes(challengeScriptLength int, file *os.File) ([]uint8, error) {
  var challengeScriptBytes = make([]uint8, challengeScriptLength)
  b := filefunctions.ReadNextBytes(file, challengeScriptLength)
  err := filefunctions.ReadUInt8ByteArray(b, &challengeScriptBytes)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return challengeScriptBytes, nil
}

func readTransactionLockTime(file *os.File) (uint32, error) {
  var transactionLockTime uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &transactionLockTime)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateTransactionLockTime(transactionLockTime) {
    return transactionLockTime, nil
  }
  return 1, errors.New("Invalid Lock Time on Transaction")
}

//ComputeBlockHash computes the SHA256 double-hash of the block header
func ComputeBlockHash(Block *block.Block) (string, error) {
  hasher := sha256.New()
  slicetwo := append(Block.Header.HPreviousBlockHash[:], Block.Header.HMerkleRoot[:] ...)
  slicethree := append(Block.Header.HTimeStamp[:], Block.Header.HTargetValue[:] ...)
  slicefour := append(slicethree[:], Block.Header.HNonce[:] ...)
  slicetwofour := append(slicetwo[:], slicefour[:] ...)
  kimbo := append(Block.Header.HFormatVersion, slicetwofour ...)
  hasher.Write(kimbo)
  slasher := hasher.Sum(nil)
  hasherTwo := sha256.New()
  hasherTwo.Write(slasher)
  return hex.EncodeToString(hasherTwo.Sum(nil)), nil
}

//ComputeTransactionHash computes the dual-SHA256 hash of a given transaction
func ComputeTransactionHash(Block *block.Block) (string, error) {

}


/***************************OUTER LOOPS****************************************/


//ParseIndividualBlock parses a block using the functions in blockchainbuilder
func ParseIndividualBlock(Block *block.Block, file *os.File) (error) {

  bmagicNumber, err := readMagicNumber(file)
  if err != nil {
    fmt.Println("No magic number recovered", err)
    return err
  }
  Block.BMagicNumber = bmagicNumber
  fmt.Println("Magic Number: ", Block.BMagicNumber)

  Block.BBlockLength, err = readBlockLength(file)
  if err != nil {
    fmt.Println("No blocklength recovered", err)
    return err
  }
  fmt.Println("Block Length: ", Block.BBlockLength)

  filefunctions.SetByteCount(0)

  Block.BFormatVersion, Block.Header.HFormatVersion, err = readFormatVersion(file)
  if err != nil {
    fmt.Println("Error reading format version", err)
    return err
  }
  fmt.Println("Format Version: ", Block.BFormatVersion)

  Block.BPreviousBlockHash, Block.Header.HPreviousBlockHash, err = readPreviousBlockHash(file)
  if err != nil {
    fmt.Println("Error reading previous block hash", err)
    return err
  }
  fmt.Println("Previous Block Hash: ", blockvalidation.ReverseEndian(Block.BPreviousBlockHash))

  Block.BMerkleRoot, Block.Header.HMerkleRoot, err = readMerkleRoot(file)
  if err != nil {
    fmt.Println("Error reading merkle root", err)
    return err
  }
  fmt.Println("Merkle Root: ", blockvalidation.ReverseEndian(Block.BMerkleRoot))

  Block.BTimeStamp, Block.Header.HTimeStamp, err = readTimeStamp(file)
  if err != nil {
    fmt.Println("Error reading timestamp", err)
    return err
  }
  fmt.Println("Time Stamp: ", Block.BTimeStamp)

  Block.BTargetValue, Block.Header.HTargetValue, err = readTargetValue(file)
  if err != nil {
    fmt.Println("Error reading target value", err)
    return err
  }
  fmt.Println("Target Value: ", Block.BTargetValue)

  Block.BNonce, Block.Header.HNonce, err = readNonce(file)
  if err != nil {
    fmt.Println("Error reading nonce", err)
    return err
  }
  fmt.Println("Nonce: ", Block.BNonce)

  Block.BBlockHash, err = ComputeBlockHash(Block)
  if err != nil {
    fmt.Println("Error computing block hash", err)
    return err
  }
  fmt.Println("Block Hash: ", blockvalidation.ReverseEndian(Block.BBlockHash))

  Block.BTransactionCount, err = readTransactionCount(file)
  if err != nil {
    fmt.Println("Error reading transaction length", err)
    return err
  }
  fmt.Println("Transaction Length: ", Block.BTransactionCount)

/*===============================Transactions=================================
 ============================================================================*/

  for transactionIndex := 1; transactionIndex <= int(Block.BTransactionCount); transactionIndex++ {

    fmt.Println(" ========== Transaction ", transactionIndex, " of ", int(Block.BTransactionCount), " ============")

    Block.BTransactionVersionNumber, err = readTransactionVersion(file)
    if err != nil {
      fmt.Println("Error reading transaction version number", err)
      return err
    }
    fmt.Println("Transaction Version: ", Block.BTransactionVersionNumber)

    Block.BInputCount, err = readInputCount(file)
    if err != nil {
      fmt.Println("Error reading input count", err)
      return err
    }
    fmt.Println("Input Count: ", Block.BInputCount)

/**********************************Inputs**************************************
 ******************************************************************************/

    for inputIndex := 1; inputIndex <= int(Block.BInputCount); inputIndex++ {

      fmt.Println("**** Input ", inputIndex, " of ", int(Block.BInputCount), " ****")

      Block.BTransactionHash, err = readTransactionHash(file)
      if err != nil {
        fmt.Println("Error reading transaction hash", err)
        return err
      }
      fmt.Println("Transaction Hash: ", blockvalidation.ReverseEndian(Block.BTransactionHash))

      Block.BTransactionIndex, err = readTransactionIndex(file)
      if err != nil {
        fmt.Println("Error reading transaction index", err)
        return err
      }
      fmt.Println("Transaction Index: ", Block.BTransactionIndex)

      Block.BInputScriptLength, err = readInputScriptLength(file)
      if err != nil {
        fmt.Println("Error reading script length", err)
        return err
      }
      fmt.Println("Script Length: ", Block.BInputScriptLength)

      Block.BInputScriptBytes, err = readInputScriptBytes(int(Block.BInputScriptLength), file)
      if err != nil {
        fmt.Println("Error reading script bytes", err)
        return err
      }
      fmt.Println("Script Bytes: ", Block.BInputScriptBytes)

      Block.BSequenceNumber, err = readSequenceNumber(file)
      if err != nil {
        fmt.Println("Error reading sequence number", err)
        return err
      }
      fmt.Println("Sequence Number: ", Block.BSequenceNumber)

    }

    Block.BOutputCount, err = readOutputCount(file)
    if err != nil {
      fmt.Println("Error reading output count", err)
      return err
    }
    fmt.Println("Output Count: ", Block.BOutputCount)

/**********************************Outputs*************************************
 ******************************************************************************/

    for outputIndex := 1; outputIndex <= int(Block.BOutputCount); outputIndex++ {

      fmt.Println("**** Output " , outputIndex, " of ", int(Block.BOutputCount), " ****")

      Block.BOutputValue, err = readOutputValue(file)
      if err != nil {
        fmt.Println("Error reading output value", err)
        return err
      }
      fmt.Println("Output Value: ", Block.BOutputValue)

      Block.BChallengeScriptLength, err = readChallengeScriptLength(file)
      if err != nil {
        fmt.Println("Error reading challenge script length", err)
        return err
      }
      fmt.Println("Challenge Script Length: ", Block.BChallengeScriptLength)

      Block.BChallengeScriptBytes, err = readChallengeScriptBytes(int(Block.BChallengeScriptLength), file)
      if err != nil {
        fmt.Println("Error reading challenge script bytes", err)
        return err
      }
      fmt.Println("Challenge Script Bytes: ", Block.BChallengeScriptBytes)

    }
    Block.BTransactionLockTime, err = readTransactionLockTime(file)
    if err != nil {
      fmt.Println("Error reading transaction lock time", err)
      return err
    }
    fmt.Println("Transaction Lock Time: ", Block.BTransactionLockTime)
  }

  err = filefunctions.ResetBlockHeadPointer(Block.BBlockLength, file)
  if err != nil {
    fmt.Println("Error in resetting block head pointer", err)
  }
  return nil


}
