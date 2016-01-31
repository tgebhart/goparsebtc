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
  transactionLength, _, err := filefunctions.ReadVariableLengthInteger(file)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return transactionLength, nil
}

func readTransactionVersion(file *os.File) (uint32, []byte, error) {
  var transactionVersion uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &transactionVersion)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateTransactionVersion(transactionVersion) {
    return transactionVersion, b, nil
  }
  return 0, nil, errors.New("Unexpected transaction version number")
}

func readInputCount(file *os.File) (uint64, []byte, error) {
  var inputCount uint64
  inputCount, b, err := filefunctions.ReadVariableLengthInteger(file)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return inputCount, b, nil
}

func readTransactionHash(file *os.File) (string, []byte, error) {
  var transactionHash string
  b := filefunctions.ReadNextBytes(file, 32)
  filefunctions.ReadUInt8ByteArrayLength32(b, &transactionHash)

  return transactionHash, b, nil
}

func readTransactionIndex(file *os.File) (uint32, []byte, error) {
  var transactionIndex uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &transactionIndex)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return transactionIndex, b, nil
}

func readInputScriptLength(file *os.File) (uint64, []byte, error) {
  var inputScriptLength uint64
  inputScriptLength, b, err := filefunctions.ReadVariableLengthInteger(file)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return inputScriptLength, b, nil
}

func readInputScriptBytes(inputScriptLength int, file *os.File) (string, []byte, error) {
  var inputScriptBytes = make([]uint8, inputScriptLength)
  b := filefunctions.ReadNextBytes(file, inputScriptLength)
  err := filefunctions.ReadUInt8ByteArray(b, &inputScriptBytes)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return hex.EncodeToString(inputScriptBytes), b, nil
}

func readSequenceNumber(file *os.File) (uint32, []byte, error) {
  var sequenceNumber uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &sequenceNumber)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateSequenceNumber(sequenceNumber) {
    return sequenceNumber, b, nil
  }
  return 0, nil, errors.New("Invalid sequence number")
}

func readOutputCount(file *os.File) (uint64, []byte, error) {
  var outputCount uint64
  outputCount, b, err := filefunctions.ReadVariableLengthInteger(file)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return outputCount, b, nil
}

func readOutputValue(file *os.File) (uint64, []byte, error) {
  var outputValue uint64
  b := filefunctions.ReadNextBytes(file, 8)
  err := filefunctions.ReadBinaryToUInt64(b, &outputValue)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return outputValue, b, nil
}

func readChallengeScriptLength(file *os.File) (uint64, []byte, error) {
  var challengeScriptLength uint64
  challengeScriptLength, b, err := filefunctions.ReadVariableLengthInteger(file)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return challengeScriptLength, b, nil
}

func readChallengeScriptBytes(challengeScriptLength int, file *os.File) (string, []byte, error) {
  var challengeScriptBytes = make([]uint8, challengeScriptLength)
  b := filefunctions.ReadNextBytes(file, challengeScriptLength)
  err := filefunctions.ReadUInt8ByteArray(b, &challengeScriptBytes)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return hex.EncodeToString(challengeScriptBytes), b, nil
}

func readTransactionLockTime(file *os.File) (uint32, []byte, error) {
  var transactionLockTime uint32
  b := filefunctions.ReadNextBytes(file, 4)
  err := filefunctions.ReadBinaryToUInt32(b, &transactionLockTime)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateTransactionLockTime(transactionLockTime) {
    return transactionLockTime, b, nil
  }
  return 1, nil, errors.New("Invalid Lock Time on Transaction")
}

//ComputeBlockHash computes the SHA256 double-hash of the block header
func ComputeBlockHash(Block *block.Block) (string, error) {
  hasher := sha256.New()
  slicetwo := append(Block.ByteHeader.PreviousBlockHash[:], Block.ByteHeader.MerkleRoot[:] ...)
  slicethree := append(Block.ByteHeader.TimeStamp[:], Block.ByteHeader.TargetValue[:] ...)
  slicefour := append(slicethree[:], Block.ByteHeader.Nonce[:] ...)
  slicetwofour := append(slicetwo[:], slicefour[:] ...)
  kimbo := append(Block.ByteHeader.FormatVersion, slicetwofour ...)
  hasher.Write(kimbo)
  slasher := hasher.Sum(nil)
  hasherTwo := sha256.New()
  hasherTwo.Write(slasher)
  return hex.EncodeToString(hasherTwo.Sum(nil)), nil
}

//ComputeTransactionHash computes the dual-SHA256 hash of a given transaction
func ComputeTransactionHash(Transaction *block.ByteTransaction, inputCount uint64, outputCount uint64) (string, error) {
  hasher := sha256.New()
  var inputBytes []byte
  var outputBytes []byte
  for i := 0; i < int(inputCount); i++ {
    inputBytes = append(inputBytes[:], combineInputBytes(&Transaction.Inputs[i])[:] ...)
  }
  for o := 0; o < int(outputCount); o++ {
    fmt.Println("Outputcount: ", int(outputCount))
    outputBytes = append(outputBytes[:], combineOutputBytes(&Transaction.Outputs[o])[:] ...)
  }
  sliceone := append(Transaction.TransactionVersionNumber[:], Transaction.InputCount[:] ...)
  slicetwo := append(inputBytes[:], Transaction.OutputCount[:] ...)
  slicethree := append(slicetwo[:], outputBytes[:] ...)
  sliceonethree := append(sliceone[:], slicethree[:] ...)
  kimbo := append(sliceonethree[:], Transaction.TransactionLockTime[:] ...)
  hasher.Write(kimbo)
  slasher := hasher.Sum(nil)
  hasherTwo := sha256.New()
  hasherTwo.Write(slasher)
  return hex.EncodeToString(hasherTwo.Sum(nil)), nil
}

func combineInputBytes(Input *block.ByteInput) ([]byte) {
  var inputBytes []byte
  sliceone := append(Input.TransactionHash[:], Input.TransactionIndex[:] ...)
  slicetwo := append(Input.InputScriptLength[:], Input.InputScriptBytes[:] ...)
  slicethree := append(sliceone[:], slicetwo[:] ...)
  inputBytes = append(slicethree[:], Input.SequenceNumber[:] ...)
  return inputBytes
}

func combineOutputBytes(Output *block.ByteOutput) ([]byte) {
  var outputBytes []byte
  sliceone := append(Output.OutputValue[:], Output.ChallengeScriptLength[:] ...)
  outputBytes = append(sliceone[:], Output.ChallengeScriptBytes[:] ...)
  return outputBytes
}

//ParseAddressFromOutputScript parses an output's script to determine address hash
func ParseAddressFromOutputScript(Output *block.ByteOutput) (string, error) {

}


/***************************OUTER LOOPS****************************************/


//ParseIndividualBlock parses a block using the functions in blockchainbuilder
func ParseIndividualBlock(Block *block.Block, file *os.File) (error) {

  bmagicNumber, err := readMagicNumber(file)
  if err != nil {
    fmt.Println("No magic number recovered", err)
    return err
  }
  Block.MagicNumber = bmagicNumber
  fmt.Println("Magic Number: ", Block.MagicNumber)

  Block.BlockLength, err = readBlockLength(file)
  if err != nil {
    fmt.Println("No blocklength recovered", err)
    return err
  }
  fmt.Println("Block Length: ", Block.BlockLength)

  filefunctions.SetByteCount(0)

  Block.Header.FormatVersion, Block.ByteHeader.FormatVersion, err = readFormatVersion(file)
  if err != nil {
    fmt.Println("Error reading format version", err)
    return err
  }
  fmt.Println("Format Version: ", Block.Header.FormatVersion)

  Block.Header.PreviousBlockHash, Block.ByteHeader.PreviousBlockHash, err = readPreviousBlockHash(file)
  if err != nil {
    fmt.Println("Error reading previous block hash", err)
    return err
  }
  fmt.Println("Previous Block Hash: ", blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash))

  Block.Header.MerkleRoot, Block.ByteHeader.MerkleRoot, err = readMerkleRoot(file)
  if err != nil {
    fmt.Println("Error reading merkle root", err)
    return err
  }
  fmt.Println("Merkle Root: ", blockvalidation.ReverseEndian(Block.Header.MerkleRoot))

  Block.Header.TimeStamp, Block.ByteHeader.TimeStamp, err = readTimeStamp(file)
  if err != nil {
    fmt.Println("Error reading timestamp", err)
    return err
  }
  fmt.Println("Time Stamp: ", Block.Header.TimeStamp)

  Block.Header.TargetValue, Block.ByteHeader.TargetValue, err = readTargetValue(file)
  if err != nil {
    fmt.Println("Error reading target value", err)
    return err
  }
  fmt.Println("Target Value: ", Block.Header.TargetValue)

  Block.Header.Nonce, Block.ByteHeader.Nonce, err = readNonce(file)
  if err != nil {
    fmt.Println("Error reading nonce", err)
    return err
  }
  fmt.Println("Nonce: ", Block.Header.Nonce)

  Block.BlockHash, err = ComputeBlockHash(Block)
  if err != nil {
    fmt.Println("Error computing block hash", err)
    return err
  }
  fmt.Println("Block Hash: ", blockvalidation.ReverseEndian(Block.BlockHash))

  Block.TransactionCount, err = readTransactionCount(file)
  if err != nil {
    fmt.Println("Error reading transaction length", err)
    return err
  }
  fmt.Println("Transaction Length: ", Block.TransactionCount)

/*===============================Transactions=================================
 ============================================================================*/

  for transactionIndex := 0; transactionIndex < int(Block.TransactionCount); transactionIndex++ {

    fmt.Println(" ========== Transaction ", transactionIndex + 1, " of ", int(Block.TransactionCount), " ============")

    Block.Transactions = append(Block.Transactions, block.Transaction{})
    Block.ByteTransactions = append(Block.ByteTransactions, block.ByteTransaction{})

    Block.Transactions[transactionIndex].TransactionVersionNumber, Block.ByteTransactions[transactionIndex].TransactionVersionNumber, err = readTransactionVersion(file)
    if err != nil {
      fmt.Println("Error reading transaction version number", err)
      return err
    }
    fmt.Println("Transaction Version: ", Block.Transactions[transactionIndex].TransactionVersionNumber)

    Block.Transactions[transactionIndex].InputCount, Block.ByteTransactions[transactionIndex].InputCount, err = readInputCount(file)
    if err != nil {
      fmt.Println("Error reading input count", err)
      return err
    }
    fmt.Println("Input Count: ", Block.Transactions[transactionIndex].InputCount)

/**********************************Inputs**************************************
 ******************************************************************************/

    for inputIndex := 0; inputIndex < int(Block.Transactions[transactionIndex].InputCount); inputIndex++ {

      fmt.Println("**** Input ", inputIndex + 1, " of ", int(Block.Transactions[transactionIndex].InputCount), " ****")

      Block.Transactions[transactionIndex].Inputs = append(Block.Transactions[transactionIndex].Inputs, block.Input{})
      Block.ByteTransactions[transactionIndex].Inputs = append(Block.ByteTransactions[transactionIndex].Inputs, block.ByteInput{})

      Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionHash, Block.ByteTransactions[transactionIndex].Inputs[inputIndex].TransactionHash, err = readTransactionHash(file)
      if err != nil {
        fmt.Println("Error reading transaction hash", err)
        return err
      }
      fmt.Println("Transaction Hash: ", blockvalidation.ReverseEndian(Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionHash))

      Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionIndex, Block.ByteTransactions[transactionIndex].Inputs[inputIndex].TransactionIndex, err = readTransactionIndex(file)
      if err != nil {
        fmt.Println("Error reading transaction index", err)
        return err
      }
      fmt.Println("Transaction Index: ", Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionIndex)

      Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength, Block.ByteTransactions[transactionIndex].Inputs[inputIndex].InputScriptLength, err = readInputScriptLength(file)
      if err != nil {
        fmt.Println("Error reading script length", err)
        return err
      }
      fmt.Println("Script Length: ", Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength)

      Block.Transactions[transactionIndex].Inputs[inputIndex].InputScript, Block.ByteTransactions[transactionIndex].Inputs[inputIndex].InputScriptBytes, err = readInputScriptBytes(int(Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength), file)
      if err != nil {
        fmt.Println("Error reading script bytes", err)
        return err
      }
      fmt.Println("Input Script: ", Block.Transactions[transactionIndex].Inputs[inputIndex].InputScript)

      Block.Transactions[transactionIndex].Inputs[inputIndex].SequenceNumber, Block.ByteTransactions[transactionIndex].Inputs[inputIndex].SequenceNumber, err = readSequenceNumber(file)
      if err != nil {
        fmt.Println("Error reading sequence number", err)
        return err
      }
      fmt.Println("Sequence Number: ", Block.Transactions[transactionIndex].Inputs[inputIndex].SequenceNumber)

    }

    Block.Transactions[transactionIndex].OutputCount, Block.ByteTransactions[transactionIndex].OutputCount, err = readOutputCount(file)
    if err != nil {
      fmt.Println("Error reading output count", err)
      return err
    }
    fmt.Println("Output Count: ", Block.Transactions[transactionIndex].OutputCount)

/**********************************Outputs*************************************
 ******************************************************************************/

    for outputIndex := 0; outputIndex < int(Block.Transactions[transactionIndex].OutputCount); outputIndex++ {

      fmt.Println("**** Output " , outputIndex + 1, " of ", int(Block.Transactions[transactionIndex].OutputCount), " ****")

      Block.Transactions[transactionIndex].Outputs = append(Block.Transactions[transactionIndex].Outputs, block.Output{})
      Block.ByteTransactions[transactionIndex].Outputs = append(Block.ByteTransactions[transactionIndex].Outputs, block.ByteOutput{})

      Block.Transactions[transactionIndex].Outputs[outputIndex].OutputValue, Block.ByteTransactions[transactionIndex].Outputs[outputIndex].OutputValue, err = readOutputValue(file)
      if err != nil {
        fmt.Println("Error reading output value", err)
        return err
      }
      fmt.Println("Output Value: ", Block.Transactions[transactionIndex].Outputs[outputIndex].OutputValue)

      Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength, Block.ByteTransactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength, err = readChallengeScriptLength(file)
      if err != nil {
        fmt.Println("Error reading challenge script length", err)
        return err
      }
      fmt.Println("Challenge Script Length: ", Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength)

      Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScript, Block.ByteTransactions[transactionIndex].Outputs[outputIndex].ChallengeScriptBytes, err = readChallengeScriptBytes(int(Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength), file)
      if err != nil {
        fmt.Println("Error reading challenge script bytes", err)
        return err
      }
      fmt.Println("Challenge Script: ", Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScript)

    }
    Block.Transactions[transactionIndex].TransactionLockTime, Block.ByteTransactions[transactionIndex].TransactionLockTime, err = readTransactionLockTime(file)
    if err != nil {
      fmt.Println("Error reading transaction lock time", err)
      return err
    }
    fmt.Println("Transaction Lock Time: ", Block.Transactions[transactionIndex].TransactionLockTime)

    Block.Transactions[transactionIndex].TransactionHash, err = ComputeTransactionHash(&Block.ByteTransactions[transactionIndex], Block.Transactions[transactionIndex].InputCount, Block.Transactions[transactionIndex].OutputCount)
    if err != nil {
      fmt.Println("Error in computing transaction hash", err)
      return err
    }
    fmt.Println("Transaction Hash: ", blockvalidation.ReverseEndian(Block.Transactions[transactionIndex].TransactionHash))
  }

  err = filefunctions.ResetBlockHeadPointer(Block.BlockLength, file)
  if err != nil {
    fmt.Println("Error in resetting block head pointer", err)
  }
  return nil


}
