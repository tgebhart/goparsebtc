package blockchainbuilder

import (
   "errors"
    "fmt"
    "os"
    "github.com/tgebhart/goparsebtc/block"
    "github.com/tgebhart/goparsebtc/filefunctions"
    "github.com/tgebhart/goparsebtc/blockvalidation"
    "github.com/tgebhart/goparsebtc/btchashing"
    "encoding/hex"
    "encoding/csv"
    "strconv"
)

//Blockchain holds the BlockMap object
type Blockchain struct {
  BlockMap map[string]block.HashBlock
}

//NewBlockchain constructs a Blockchain instance
func NewBlockchain() *Blockchain {
  var b Blockchain
  b.BlockMap = make(map[string]block.HashBlock)
  return &b
}


func readExactMagicNumber(file *os.File) (uint32, error) {

  var magicNumber uint32
  b, err := filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, err
  }
  fmt.Println("Magic Number Bytes: ", b)
  err = filefunctions.ReadBinaryToUInt32(b, &magicNumber)
  if err != nil {
    fmt.Println("binary.Read failed:", err)
  }
  if blockvalidation.ValidateMagicNumber(magicNumber) {
    return magicNumber, nil
  }
  return magicNumber, ErrBadMagic
}

//ErrBadMagic is returned when magic number is unusual or can't be found. Can be used to trigger new file opening
var ErrBadMagic = errors.New("blockchainbuilder: unusual or invalid magic number")
//ErrBadOutputValue is thrown when output value isn't picked up correctly
var ErrBadOutputValue = errors.New("blockchainbuilder : unusual output value")
//ErrBadFormatVersion is thrown when block format version is odd
var ErrBadFormatVersion = errors.New("Unusual format version")
//ErrWriteToFile is thrown when an error occurs while writing main chain to file
var ErrWriteToFile = errors.New("WriteToFile: Could not locate previous block hash")
//ErrBadSequenceNumber is thrown when an errors occurs in reading sequence number
var ErrBadSequenceNumber = errors.New("blockchainbuilder: unusual sequence number")
//ErrBadTransactionVersion is thrown when an error occurs in reading transaction version
var ErrBadTransactionVersion = errors.New("blockchainbuilder: unusual transaction version")



func readMagicNumber(file *os.File) (uint32, error) {

  var magicNumber uint32
  b, err := filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, err
  }
  err = filefunctions.ReadBinaryToUInt32(b, &magicNumber)
  if err != nil {
    fmt.Println("binary.Read failed:", err)
  }
  if blockvalidation.ValidateMagicNumber(magicNumber) {
    return magicNumber, nil
  }
  fmt.Println("Looking for magic")
  magicNumber, err = filefunctions.DetailedLookForMagic(file)
  if err != nil {
    return 0, err
  }
  if blockvalidation.ValidateMagicNumber(magicNumber) {
    return magicNumber, nil
  }
  return magicNumber, ErrBadMagic
}

func readBlockLength(file *os.File) (uint32, error) {
  var blockLength uint32
  b, err := filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, err
  }
  err = filefunctions.ReadBinaryToUInt32(b, &blockLength)
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
  b, err := filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, nil, err
  }
  err = filefunctions.ReadBinaryToUInt32(b, &formatVersion)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateFormatVersion(formatVersion) {
      return formatVersion, b, nil
  }
  return 0, nil, ErrBadFormatVersion
}

func readPreviousBlockHash(file *os.File) (string, []byte, error) {
  var previousBlockHash string
  b, err := filefunctions.ReadNextBytes(file, 32)
  if err != nil {
    return "", nil, err
  }
  filefunctions.ReadUInt8ByteArrayLength32(b, &previousBlockHash)
  return previousBlockHash, b, nil
}

func readMerkleRoot(file *os.File) (string, []byte, error) {
  var merkleRoot string
  b, err := filefunctions.ReadNextBytes(file, 32)
  if err != nil {
    return "", nil, err
  }
  filefunctions.ReadUInt8ByteArrayLength32(b, &merkleRoot)
  return merkleRoot, b, nil
}

func readTimeStamp(file *os.File) (uint32, []byte, error) {
  var timeStamp uint32
  b, err := filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, nil, err
  }
  err = filefunctions.ReadBinaryToUInt32(b, &timeStamp)
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
  b, err := filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, nil, err
  }
  err = filefunctions.ReadBinaryToUInt32(b, &targetValue)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return targetValue, b, nil
}

func readNonce(file *os.File) (uint32, []byte, error) {
  var nonce uint32
  b, err := filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, nil, err
  }
  err = filefunctions.ReadBinaryToUInt32(b, &nonce)
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
  b, err := filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, nil, err
  }
  err = filefunctions.ReadBinaryToUInt32(b, &transactionVersion)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateTransactionVersion(transactionVersion) {
    return transactionVersion, b, nil
  }
  if transactionVersion == 16777216 {
    filefunctions.StepBack(5, file)
    b, err := filefunctions.ReadNextBytes(file, 4)
    if err != nil {
      return 0, nil, err
    }
    err = filefunctions.ReadBinaryToUInt32(b, &transactionVersion)
    if err != nil {
      fmt.Println("binary.Read failed: ", err)
    }
    if blockvalidation.ValidateTransactionVersion(transactionVersion) {
      return transactionVersion, b, nil
    }
  }
  return 0, nil, ErrBadTransactionVersion
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
  b, err := filefunctions.ReadNextBytes(file, 32)
  if err != nil {
    return "", nil, err
  }
  filefunctions.ReadUInt8ByteArrayLength32(b, &transactionHash)
  return transactionHash, b, nil
}

func readTransactionIndex(file *os.File) (uint32, []byte, error) {
  var transactionIndex uint32
  b, err := filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, nil, err
  }
  err = filefunctions.ReadBinaryToUInt32(b, &transactionIndex)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateTransactionIndex(transactionIndex) {
    return transactionIndex, b, nil
  }
  b, err = filefunctions.RewindAndRead32(b, file, &transactionIndex)
  if err != nil {
    fmt.Println("rewind read failed:  ", err)
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
  b, err := filefunctions.ReadNextBytes(file, inputScriptLength)
  if err != nil {
    return "", nil, err
  }
  err = filefunctions.ReadUInt8ByteArray(b, &inputScriptBytes)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  /*if inputScriptBytes[len(inputScriptBytes) - 1 ] == 0xFF {
    inputScriptBytes = inputScriptBytes[0:len(inputScriptBytes) - 1]
    filefunctions.StepBack(1, file)
    return hex.EncodeToString(inputScriptBytes), b, nil
  }*/
  return hex.EncodeToString(inputScriptBytes), b, nil
}

func readSequenceNumber(file *os.File) (uint32, []byte, error) {
  var sequenceNumber uint32
  b, err := filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, nil, err
  }
  err = filefunctions.ReadBinaryToUInt32(b, &sequenceNumber)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }

  if blockvalidation.ValidateSequenceNumber(b) {
    return sequenceNumber, b, nil
  }
  filefunctions.StepBack(5, file)
  b, err = filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, nil, err
  }
  err = filefunctions.ReadBinaryToUInt32(b, &sequenceNumber)
  if err != nil {
    fmt.Println("binary.Read failed on StepBack", err)
  }
  if blockvalidation.ValidateSequenceNumber(b) {
    return sequenceNumber, b, nil
  }
  return 0, nil, ErrBadSequenceNumber
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
  b, err := filefunctions.ReadNextBytes(file, 8)
  if err != nil {
    return 0, nil, err
  }
  err = filefunctions.ReadBinaryToUInt64(b, &outputValue)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateOutputValue(outputValue) {
    return outputValue, b, nil
  }
  b, err = filefunctions.RewindAndRead64(b, file, &outputValue)
  if err != nil {
    fmt.Println("rewind read failed: ", err)
  }
  if blockvalidation.ValidateOutputValue(outputValue) {
    return outputValue, b, nil
  }
  return 0, nil, ErrBadOutputValue
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
  b, err := filefunctions.ReadNextBytes(file, challengeScriptLength)
  if err != nil {
    return "", nil, err
  }
  err = filefunctions.ReadUInt8ByteArray(b, &challengeScriptBytes)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return hex.EncodeToString(challengeScriptBytes), b, nil
}

func readTransactionLockTime(file *os.File) (uint32, []byte, error) {
  var transactionLockTime uint32
  b, err := filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, nil, err
  }
  err = filefunctions.ReadBinaryToUInt32(b, &transactionLockTime)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if blockvalidation.ValidateTransactionLockTime(transactionLockTime) {
    return transactionLockTime, b, nil
  }
  filefunctions.StepBack(5, file)
  b, err = filefunctions.ReadNextBytes(file, 4)
  if err != nil {
    return 0, nil, err
  }
  err = filefunctions.ReadBinaryToUInt32(b, &transactionLockTime)
  if err != nil {
    fmt.Println("binary.Read failed on StepBack", err)
  }
  if blockvalidation.ValidateTransactionLockTime(transactionLockTime) {
    return transactionLockTime, b, nil
  }
  return 1, nil, errors.New("Invalid Lock Time on Transaction")
}



/***************************OUTER LOOPS****************************************/


//ParseIndividualBlock parses a block using the functions in blockchainbuilder
func (Blockchain) ParseIndividualBlock(Block *block.Block, file *os.File) (error) {

  bmagicNumber, err := readMagicNumber(file)
  if err != nil {
    fmt.Println("No magic number recovered", err)
    return err
  }
  Block.MagicNumber = bmagicNumber
  fmt.Println("Magic Number: ", Block.MagicNumber)

  offset, err := file.Seek(0, 1)
  if err != nil {
    return err
  }
  Block.HashBlock.ByteOffset = int(offset - 4)

  Block.BlockLength, err = readBlockLength(file)
  if err != nil {
    fmt.Println("No blocklength recovered", err)
    return err
  }
  fmt.Println("Block Length: ", Block.BlockLength)

  //Update ByteOffset and ParsedBlockLength fields to track where in the file the block ends
  Block.HashBlock.ParsedBlockLength = Block.BlockLength

  filefunctions.SetByteCount(0)

  Block.Header.FormatVersion, Block.Header.ByteFormatVersion, err = readFormatVersion(file)
  if err != nil {
    fmt.Println("Error reading format version", Block.Header.FormatVersion, err)
    return err
  }
  fmt.Println("Format Version: ", Block.Header.FormatVersion)

  Block.Header.PreviousBlockHash, Block.Header.BytePreviousBlockHash, err = readPreviousBlockHash(file)
  if err != nil {
    fmt.Println("Error reading previous block hash", err)
    return err
  }
  fmt.Println("Previous Block Hash: ", blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash))

  //Update HashBlock PreviousBlockHash field with parsed value. Will be used to build Main Chain
  Block.HashBlock.PreviousCompressedBlockHash = btchashing.ComputeCompressedBlockHash(blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash))
  Block.HashBlock.CompressedBlockHash = blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash)

  Block.Header.MerkleRoot, Block.Header.ByteMerkleRoot, err = readMerkleRoot(file)
  if err != nil {
    fmt.Println("Error reading merkle root", err)
    return err
  }
  fmt.Println("Merkle Root: ", blockvalidation.ReverseEndian(Block.Header.MerkleRoot))

  Block.Header.TimeStamp, Block.Header.ByteTimeStamp, err = readTimeStamp(file)
  if err != nil {
    fmt.Println("Error reading timestamp", err)
    return err
  }
  fmt.Println("Time Stamp: ", blockvalidation.ConvertUnixEpochToDate(Block.Header.TimeStamp))

  //Update TimeStamp of hashblock
  Block.HashBlock.TimeStamp = Block.Header.TimeStamp

  Block.Header.TargetValue, Block.Header.ByteTargetValue, err = readTargetValue(file)
  if err != nil {
    fmt.Println("Error reading target value", err)
    return err
  }
  fmt.Println("Target Value: ", Block.Header.TargetValue)

  Block.Header.Nonce, Block.Header.ByteNonce, err = readNonce(file)
  if err != nil {
    fmt.Println("Error reading nonce", err)
    return err
  }
  fmt.Println("Nonce: ", Block.Header.Nonce)

  Block.BlockHash, err = btchashing.ComputeBlockHash(Block)
  if err != nil {
    fmt.Println("Error computing block hash", err)
    return err
  }
  fmt.Println("Block Hash: ", blockvalidation.ReverseEndian(Block.BlockHash))

  //Add BlockHash field to HashBlock object in Block and compress hash to limit search space for BlockChain hashmap
  Block.HashBlock.CompressedBlockHash = btchashing.ComputeCompressedBlockHash(blockvalidation.ReverseEndian(Block.BlockHash))
  Block.HashBlock.BlockHash = blockvalidation.ReverseEndian(Block.BlockHash)

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

    Block.Transactions[transactionIndex].TransactionVersionNumber, Block.Transactions[transactionIndex].ByteTransactionVersionNumber, err = readTransactionVersion(file)
    if err != nil {
      fmt.Println("Error reading transaction version number", Block.Transactions[transactionIndex].TransactionVersionNumber, err)
      return err
    }
    fmt.Println("Transaction Version: ", Block.Transactions[transactionIndex].TransactionVersionNumber, Block.Transactions[transactionIndex].ByteTransactionVersionNumber)

    Block.Transactions[transactionIndex].InputCount, Block.Transactions[transactionIndex].ByteInputCount, err = readInputCount(file)
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

      Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionHash, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteTransactionHash, err = readTransactionHash(file)
      if err != nil {
        fmt.Println("Error reading transaction hash", err)
        return err
      }
      fmt.Println("Transaction Hash: ", blockvalidation.ReverseEndian(Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionHash))

      Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionIndex, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteTransactionIndex, err = readTransactionIndex(file)
      if err != nil {
        fmt.Println("Error reading transaction index", err)
        return err
      }
      fmt.Println("Transaction Index: ", Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionIndex)

      Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteInputScriptLength, err = readInputScriptLength(file)
      if err != nil {
        fmt.Println("Error reading script length", err)
        return err
      }
      fmt.Println("Script Length: ", Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteInputScriptLength)

      Block.Transactions[transactionIndex].Inputs[inputIndex].InputScript, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteInputScript, err = readInputScriptBytes(int(Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength), file)
      if err != nil {
        fmt.Println("Error reading script bytes", err)
        return err
      }
      fmt.Println("Input Script: ", Block.Transactions[transactionIndex].Inputs[inputIndex].InputScript)

      Block.Transactions[transactionIndex].Inputs[inputIndex].SequenceNumber, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteSequenceNumber, err = readSequenceNumber(file)
      if err != nil {
        fmt.Println("Error reading sequence number", err)
        return err
      }
      fmt.Println("Sequence Number: ", Block.Transactions[transactionIndex].Inputs[inputIndex].SequenceNumber, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteSequenceNumber)

    }

    Block.Transactions[transactionIndex].OutputCount, Block.Transactions[transactionIndex].ByteOutputCount, err = readOutputCount(file)
    if err != nil {
      fmt.Println("Error reading output count", err)
      return err
    }
    fmt.Println("Output Count: ", Block.Transactions[transactionIndex].OutputCount, Block.Transactions[transactionIndex].ByteOutputCount)

/**********************************Outputs*************************************
 ******************************************************************************/

    for outputIndex := 0; outputIndex < int(Block.Transactions[transactionIndex].OutputCount); outputIndex++ {

      fmt.Println("**** Output " , outputIndex + 1, " of ", int(Block.Transactions[transactionIndex].OutputCount), " ****")

      Block.Transactions[transactionIndex].Outputs = append(Block.Transactions[transactionIndex].Outputs, block.Output{})

      Block.Transactions[transactionIndex].Outputs[outputIndex].OutputValue, Block.Transactions[transactionIndex].Outputs[outputIndex].ByteOutputValue, err = readOutputValue(file)
      if err != nil {
        fmt.Println("Error reading output value", err)
        return err
      }
      fmt.Println("Output Value: ", Block.Transactions[transactionIndex].Outputs[outputIndex].OutputValue)

      Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength, Block.Transactions[transactionIndex].Outputs[outputIndex].ByteChallengeScriptLength, err = readChallengeScriptLength(file)
      if err != nil {
        fmt.Println("Error reading challenge script length", err)
        return err
      }
      fmt.Println("Challenge Script Length: ", Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength)

      Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScript, Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptBytes, err = readChallengeScriptBytes(int(Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength), file)
      if err != nil {
        fmt.Println("Error reading challenge script bytes", err)
        return err
      }
      fmt.Println("Challenge Script: ", Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScript)

      Block.Transactions[transactionIndex].Outputs[outputIndex].KeyType, err = blockvalidation.ParseOutputScript(&Block.Transactions[transactionIndex].Outputs[outputIndex])
      if err != nil {
        return err
      }

      fmt.Println("Hash160: ", Block.Transactions[transactionIndex].Outputs[outputIndex].Addresses[0].RipeMD160)
      fmt.Println("Address: ", Block.Transactions[transactionIndex].Outputs[outputIndex].Addresses[0].Address)
      fmt.Println("PublicKey: ", Block.Transactions[transactionIndex].Outputs[outputIndex].Addresses[0].PublicKey)

    }

    Block.Transactions[transactionIndex].TransactionLockTime, Block.Transactions[transactionIndex].ByteTransactionLockTime, err = readTransactionLockTime(file)
    if err != nil {
      fmt.Println("Error reading transaction lock time", err)
      return err
    }
    fmt.Println("Transaction Lock Time: ", Block.Transactions[transactionIndex].TransactionLockTime)

    Block.Transactions[transactionIndex].TransactionHash, err = btchashing.ComputeTransactionHash(&Block.Transactions[transactionIndex], Block.Transactions[transactionIndex].InputCount, Block.Transactions[transactionIndex].OutputCount)
    if err != nil {
      fmt.Println("Error in computing transaction hash", err)
      return err
    }
    fmt.Println("Transaction Hash: ", blockvalidation.ReverseEndian(Block.Transactions[transactionIndex].TransactionHash))
  }

  _, err = filefunctions.ResetBlockHeadPointer(Block.BlockLength, file)
  if err != nil {
    fmt.Println("Error in resetting block head pointer", err)
  }
  return nil

}

//WriteMainChainToFile writes the binary data of the compressed HashBlock to filename
func WriteMainChainToFile(chain *Blockchain, currentKey string, filename string) (error) {

  f, err := os.Create("" + filename + ".csv")
  if err != nil {
    return err
  }
  defer f.Close()

  writer := csv.NewWriter(f)

  fmt.Println("Writing to file...")
  //fmt.Println(chain.BlockMap)
  var thisHash string
  var nextKey string
  var nextHash string

  for chain.BlockMap[currentKey].PreviousCompressedBlockHash != "5fa1ced70304aad0f8c7de728c0b20cd" {

    fmt.Println(chain.BlockMap[currentKey])

    if chain.BlockMap[currentKey].PreviousCompressedBlockHash == "" {
      fmt.Println("search hash: ", thisHash)
      nextHash, thisHash, err = blockvalidation.GetReplacementKey(nextHash)
      if err != nil {
        return err
      }
      nextKey = btchashing.ComputeCompressedBlockHash(nextHash)
      currentKey = nextKey
    } else if chain.BlockMap[currentKey].BlockHash == "" {
      nextHash, thisHash, err = blockvalidation.GetReplacementKey(nextHash)
      if err != nil {
        return err
      }
      nextKey = btchashing.ComputeCompressedBlockHash(nextHash)
      currentKey = nextKey
      } else {
      nextKey = chain.BlockMap[currentKey].PreviousCompressedBlockHash
      thisHash = chain.BlockMap[currentKey].BlockHash
      nextHash = chain.BlockMap[currentKey].PreviousBlockHash
      currentKey = nextKey
    }
    err = writer.Write([]string{thisHash, chain.BlockMap[currentKey].FileEndpoint, strconv.Itoa(chain.BlockMap[currentKey].ByteOffset), strconv.Itoa(int(chain.BlockMap[currentKey].ParsedBlockLength)), strconv.Itoa(chain.BlockMap[currentKey].RawBlockNumber), fmt.Sprint(chain.BlockMap[currentKey].TimeStamp)})
    if err != nil {
      fmt.Println("Error writing file", err)
      return err
    }
    fmt.Println("Wrote", thisHash)
    fmt.Println("next hash", nextHash)
  }
  return nil
}































/***************************OUTER LOOPS****************************************/


//ParseIndividualBlockSuppressOutput parses a block using the functions in blockchainbuilder -- no output
func (Blockchain) ParseIndividualBlockSuppressOutput(Block *block.Block, file *os.File) (error) {

  filefunctions.SetByteCount(0)

  bmagicNumber, err := readMagicNumber(file)
  if err != nil {
    fmt.Println("No magic number recovered", err)
    return err
  }
  Block.MagicNumber = bmagicNumber

  offset, err := file.Seek(0, 1)
  if err != nil {
    return err
  }
  Block.HashBlock.ByteOffset = int(offset - 4)
  //fmt.Println("Magic Number: ", Block.MagicNumber)

  Block.BlockLength, err = readBlockLength(file)
  if err != nil {
    fmt.Println("No blocklength recovered", err)
    return err
  }
  //Update ByteOffset and ParsedBlockLength fields to track where in the file the block ends
  Block.HashBlock.ParsedBlockLength = Block.BlockLength

  filefunctions.SetByteCount(0)

  Block.Header.FormatVersion, Block.Header.ByteFormatVersion, err = readFormatVersion(file)
  if err != nil {
    fmt.Println("Error reading format version", Block.Header.FormatVersion, err)
    return err
  }
  //fmt.Println("Format Version: ", Block.Header.FormatVersion)

  Block.Header.PreviousBlockHash, Block.Header.BytePreviousBlockHash, err = readPreviousBlockHash(file)
  if err != nil {
    fmt.Println("Error reading previous block hash", err)
    return err
  }
  //fmt.Println("Previous Block Hash: ", blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash))

  //Update HashBlock PreviousBlockHash field with parsed value. Will be used to build Main Chain
  Block.HashBlock.PreviousCompressedBlockHash = btchashing.ComputeCompressedBlockHash(blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash))
  Block.HashBlock.PreviousBlockHash = blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash)

  Block.Header.MerkleRoot, Block.Header.ByteMerkleRoot, err = readMerkleRoot(file)
  if err != nil {
    fmt.Println("Error reading merkle root", err)
    return err
  }
  //fmt.Println("Merkle Root: ", blockvalidation.ReverseEndian(Block.Header.MerkleRoot))

  Block.Header.TimeStamp, Block.Header.ByteTimeStamp, err = readTimeStamp(file)
  if err != nil {
    fmt.Println("Error reading timestamp", err)
    return err
  }
  //fmt.Println("Time Stamp: ", blockvalidation.ConvertUnixEpochToDate(Block.Header.TimeStamp))
  Block.HashBlock.TimeStamp = Block.Header.TimeStamp

  Block.Header.TargetValue, Block.Header.ByteTargetValue, err = readTargetValue(file)
  if err != nil {
    fmt.Println("Error reading target value", err)
    return err
  }
  //fmt.Println("Target Value: ", Block.Header.TargetValue)

  Block.Header.Nonce, Block.Header.ByteNonce, err = readNonce(file)
  if err != nil {
    fmt.Println("Error reading nonce", err)
    return err
  }
  //fmt.Println("Nonce: ", Block.Header.Nonce)

  Block.BlockHash, err = btchashing.ComputeBlockHash(Block)
  if err != nil {
    fmt.Println("Error computing block hash", err)
    return err
  }
  //fmt.Println("Block Hash: ", blockvalidation.ReverseEndian(Block.BlockHash))

  //Add BlockHash field to HashBlock object in Block and compress hash to limit search space for BlockChain hashmap
  Block.HashBlock.CompressedBlockHash = btchashing.ComputeCompressedBlockHash(blockvalidation.ReverseEndian(Block.BlockHash))
  Block.HashBlock.BlockHash = blockvalidation.ReverseEndian(Block.BlockHash)

  Block.TransactionCount, err = readTransactionCount(file)
  if err != nil {
    fmt.Println("Error reading transaction length", err)
    return err
  }
  //fmt.Println("Transaction Length: ", Block.TransactionCount)

/*===============================Transactions=================================
 ============================================================================*/

  for transactionIndex := 0; transactionIndex < int(Block.TransactionCount); transactionIndex++ {

    //fmt.Println(" ========== Transaction ", transactionIndex + 1, " of ", int(Block.TransactionCount), " ============")

    Block.Transactions = append(Block.Transactions, block.Transaction{})

    Block.Transactions[transactionIndex].TransactionVersionNumber, Block.Transactions[transactionIndex].ByteTransactionVersionNumber, err = readTransactionVersion(file)
    if err != nil {
      fmt.Println("Error reading transaction version number", Block.Transactions[transactionIndex].TransactionVersionNumber, err)
      return err
    }
    //fmt.Println("Transaction Version: ", Block.Transactions[transactionIndex].TransactionVersionNumber, Block.Transactions[transactionIndex].ByteTransactionVersionNumber)

    Block.Transactions[transactionIndex].InputCount, Block.Transactions[transactionIndex].ByteInputCount, err = readInputCount(file)
    if err != nil {
      fmt.Println("Error reading input count", err)
      return err
    }
    //fmt.Println("Input Count: ", Block.Transactions[transactionIndex].InputCount)

/**********************************Inputs**************************************
 ******************************************************************************/

    for inputIndex := 0; inputIndex < int(Block.Transactions[transactionIndex].InputCount); inputIndex++ {

      //fmt.Println("**** Input ", inputIndex + 1, " of ", int(Block.Transactions[transactionIndex].InputCount), " ****")

      Block.Transactions[transactionIndex].Inputs = append(Block.Transactions[transactionIndex].Inputs, block.Input{})

      Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionHash, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteTransactionHash, err = readTransactionHash(file)
      if err != nil {
        fmt.Println("Error reading transaction hash", err)
        return err
      }
      //fmt.Println("Transaction Hash: ", blockvalidation.ReverseEndian(Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionHash))

      Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionIndex, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteTransactionIndex, err = readTransactionIndex(file)
      if err != nil {
        fmt.Println("Error reading transaction index", err)
        return err
      }
      //fmt.Println("Transaction Index: ", Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionIndex)

      Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteInputScriptLength, err = readInputScriptLength(file)
      if err != nil {
        fmt.Println("Error reading script length", err)
        return err
      }
      //fmt.Println("Script Length: ", Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteInputScriptLength)

      Block.Transactions[transactionIndex].Inputs[inputIndex].InputScript, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteInputScript, err = readInputScriptBytes(int(Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength), file)
      if err != nil {
        fmt.Println("Error reading script bytes", err)
        return err
      }
      //fmt.Println("Input Script: ", Block.Transactions[transactionIndex].Inputs[inputIndex].InputScript)

      Block.Transactions[transactionIndex].Inputs[inputIndex].SequenceNumber, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteSequenceNumber, err = readSequenceNumber(file)
      if err != nil {
        fmt.Println("Error reading sequence number", err)
        return err
      }
      //fmt.Println("Sequence Number: ", Block.Transactions[transactionIndex].Inputs[inputIndex].SequenceNumber, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteSequenceNumber)

    }

    Block.Transactions[transactionIndex].OutputCount, Block.Transactions[transactionIndex].ByteOutputCount, err = readOutputCount(file)
    if err != nil {
      fmt.Println("Error reading output count", err)
      return err
    }
    //fmt.Println("Output Count: ", Block.Transactions[transactionIndex].OutputCount, Block.Transactions[transactionIndex].ByteOutputCount)

/**********************************Outputs*************************************
 ******************************************************************************/

    for outputIndex := 0; outputIndex < int(Block.Transactions[transactionIndex].OutputCount); outputIndex++ {

      //fmt.Println("**** Output " , outputIndex + 1, " of ", int(Block.Transactions[transactionIndex].OutputCount), " ****")

      Block.Transactions[transactionIndex].Outputs = append(Block.Transactions[transactionIndex].Outputs, block.Output{})

      Block.Transactions[transactionIndex].Outputs[outputIndex].OutputValue, Block.Transactions[transactionIndex].Outputs[outputIndex].ByteOutputValue, err = readOutputValue(file)
      if err != nil {
        fmt.Println("Error reading output value", err)
        return err
      }
      //fmt.Println("Output Value: ", Block.Transactions[transactionIndex].Outputs[outputIndex].OutputValue)

      Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength, Block.Transactions[transactionIndex].Outputs[outputIndex].ByteChallengeScriptLength, err = readChallengeScriptLength(file)
      if err != nil {
        fmt.Println("Error reading challenge script length", err)
        return err
      }
      //fmt.Println("Challenge Script Length: ", Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength)

      Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScript, Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptBytes, err = readChallengeScriptBytes(int(Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength), file)
      if err != nil {
        fmt.Println("Error reading challenge script bytes", err)
        return err
      }
      //fmt.Println("Challenge Script: ", Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScript)

      Block.Transactions[transactionIndex].Outputs[outputIndex].KeyType, err = blockvalidation.ParseOutputScript(&Block.Transactions[transactionIndex].Outputs[outputIndex])
      if err != nil {
        return err
      }

      //fmt.Println("Hash160: ", Block.Transactions[transactionIndex].Outputs[outputIndex].Addresses[0].RipeMD160)
      //fmt.Println("Address: ", Block.Transactions[transactionIndex].Outputs[outputIndex].Addresses[0].Address)
      //fmt.Println("PublicKey: ", Block.Transactions[transactionIndex].Outputs[outputIndex].Addresses[0].PublicKey)

    }

    Block.Transactions[transactionIndex].TransactionLockTime, Block.Transactions[transactionIndex].ByteTransactionLockTime, err = readTransactionLockTime(file)
    if err != nil {
      fmt.Println("Error reading transaction lock time", err)
      return err
    }
    //fmt.Println("Transaction Lock Time: ", Block.Transactions[transactionIndex].TransactionLockTime)

    Block.Transactions[transactionIndex].TransactionHash, err = btchashing.ComputeTransactionHash(&Block.Transactions[transactionIndex], Block.Transactions[transactionIndex].InputCount, Block.Transactions[transactionIndex].OutputCount)
    if err != nil {
      fmt.Println("Error in computing transaction hash", err)
      return err
    }
    //fmt.Println("Transaction Hash: ", blockvalidation.ReverseEndian(Block.Transactions[transactionIndex].TransactionHash))
  }

  /*skipped, err := filefunctions.ResetBlockHeadPointer(Block.BlockLength, file)
  if err != nil {
    fmt.Println("Error in resetting block head pointer", err)
  }
  fmt.Println(Block.BlockLength)
  fmt.Println(skipped)*/
  return nil

}


//PrepareSkipBlock fills in block with as much information as possible then sets all other fields to null values
func (Blockchain) PrepareSkipBlock(Block *block.Block, fe string, rbn int, byteCount int, file *os.File) (error) {
  Block.HashBlock.FileEndpoint = fe
  Block.HashBlock.RawBlockNumber = rbn
  Block.HashBlock.ByteOffset = byteCount
  _, err := filefunctions.ResetBlockHeadPointer(Block.BlockLength, file)
  if err != nil {
    return err
  }
  return nil
}














/***************************OUTER LOOPS****************************************/


//ParseBlockOnly parses a single block from a given file location and does not include the hash block
func ParseBlockOnly(Block *block.Block, file *os.File) (error) {

  filefunctions.SetByteCount(0)

  bmagicNumber, err := readMagicNumber(file)
  if err != nil {
    fmt.Println("No magic number recovered", err)
    return err
  }
  Block.MagicNumber = bmagicNumber

  Block.BlockLength, err = readBlockLength(file)
  if err != nil {
    fmt.Println("No blocklength recovered", err)
    return err
  }

  filefunctions.SetByteCount(0)

  Block.Header.FormatVersion, Block.Header.ByteFormatVersion, err = readFormatVersion(file)
  if err != nil {
    fmt.Println("Error reading format version", Block.Header.FormatVersion, err)
    return err
  }

  Block.Header.PreviousBlockHash, Block.Header.BytePreviousBlockHash, err = readPreviousBlockHash(file)
  if err != nil {
    fmt.Println("Error reading previous block hash", err)
    return err
  }
  Block.Header.PreviousBlockHash = blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash)

  Block.Header.MerkleRoot, Block.Header.ByteMerkleRoot, err = readMerkleRoot(file)
  if err != nil {
    fmt.Println("Error reading merkle root", err)
    return err
  }

  Block.Header.TimeStamp, Block.Header.ByteTimeStamp, err = readTimeStamp(file)
  if err != nil {
    fmt.Println("Error reading timestamp", err)
    return err
  }

  Block.Header.TargetValue, Block.Header.ByteTargetValue, err = readTargetValue(file)
  if err != nil {
    fmt.Println("Error reading target value", err)
    return err
  }

  Block.Header.Nonce, Block.Header.ByteNonce, err = readNonce(file)
  if err != nil {
    fmt.Println("Error reading nonce", err)
    return err
  }

  Block.BlockHash, err = btchashing.ComputeBlockHash(Block)
  if err != nil {
    fmt.Println("Error computing block hash", err)
    return err
  }
  Block.BlockHash = blockvalidation.ReverseEndian(Block.BlockHash)

  Block.TransactionCount, err = readTransactionCount(file)
  if err != nil {
    fmt.Println("Error reading transaction length", err)
    return err
  }

/*===============================Transactions=================================
 ============================================================================*/

  for transactionIndex := 0; transactionIndex < int(Block.TransactionCount); transactionIndex++ {

    Block.Transactions = append(Block.Transactions, block.Transaction{})

    Block.Transactions[transactionIndex].TransactionVersionNumber, Block.Transactions[transactionIndex].ByteTransactionVersionNumber, err = readTransactionVersion(file)
    if err != nil {
      fmt.Println("Error reading transaction version number", Block.Transactions[transactionIndex].TransactionVersionNumber, err)
      return err
    }

    Block.Transactions[transactionIndex].InputCount, Block.Transactions[transactionIndex].ByteInputCount, err = readInputCount(file)
    if err != nil {
      fmt.Println("Error reading input count", err)
      return err
    }

/**********************************Inputs**************************************
 ******************************************************************************/

    for inputIndex := 0; inputIndex < int(Block.Transactions[transactionIndex].InputCount); inputIndex++ {

      Block.Transactions[transactionIndex].Inputs = append(Block.Transactions[transactionIndex].Inputs, block.Input{})

      Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionHash, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteTransactionHash, err = readTransactionHash(file)
      if err != nil {
        fmt.Println("Error reading transaction hash", err)
        return err
      }
      Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionHash = blockvalidation.ReverseEndian(Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionHash)


      Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionIndex, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteTransactionIndex, err = readTransactionIndex(file)
      if err != nil {
        fmt.Println("Error reading transaction index", err)
        return err
      }

      Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteInputScriptLength, err = readInputScriptLength(file)
      if err != nil {
        fmt.Println("Error reading script length", err)
        return err
      }

      Block.Transactions[transactionIndex].Inputs[inputIndex].InputScript, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteInputScript, err = readInputScriptBytes(int(Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength), file)
      if err != nil {
        fmt.Println("Error reading script bytes", err)
        return err
      }

      Block.Transactions[transactionIndex].Inputs[inputIndex].SequenceNumber, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteSequenceNumber, err = readSequenceNumber(file)
      if err != nil {
        fmt.Println("Error reading sequence number", err)
        return err
      }

    }

    Block.Transactions[transactionIndex].OutputCount, Block.Transactions[transactionIndex].ByteOutputCount, err = readOutputCount(file)
    if err != nil {
      fmt.Println("Error reading output count", err)
      return err
    }

/**********************************Outputs*************************************
 ******************************************************************************/

    for outputIndex := 0; outputIndex < int(Block.Transactions[transactionIndex].OutputCount); outputIndex++ {

      Block.Transactions[transactionIndex].Outputs = append(Block.Transactions[transactionIndex].Outputs, block.Output{})

      Block.Transactions[transactionIndex].Outputs[outputIndex].OutputValue, Block.Transactions[transactionIndex].Outputs[outputIndex].ByteOutputValue, err = readOutputValue(file)
      if err != nil {
        fmt.Println("Error reading output value", err)
        return err
      }

      Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength, Block.Transactions[transactionIndex].Outputs[outputIndex].ByteChallengeScriptLength, err = readChallengeScriptLength(file)
      if err != nil {
        fmt.Println("Error reading challenge script length", err)
        return err
      }

      Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScript, Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptBytes, err = readChallengeScriptBytes(int(Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength), file)
      if err != nil {
        fmt.Println("Error reading challenge script bytes", err)
        return err
      }

      Block.Transactions[transactionIndex].Outputs[outputIndex].KeyType, err = blockvalidation.ParseOutputScript(&Block.Transactions[transactionIndex].Outputs[outputIndex])
      if err != nil {
        return err
      }

    }

    Block.Transactions[transactionIndex].TransactionLockTime, Block.Transactions[transactionIndex].ByteTransactionLockTime, err = readTransactionLockTime(file)
    if err != nil {
      fmt.Println("Error reading transaction lock time", err)
      return err
    }

    Block.Transactions[transactionIndex].TransactionHash, err = btchashing.ComputeTransactionHash(&Block.Transactions[transactionIndex], Block.Transactions[transactionIndex].InputCount, Block.Transactions[transactionIndex].OutputCount)
    if err != nil {
      fmt.Println("Error in computing transaction hash", err)
      return err
    }
    Block.Transactions[transactionIndex].TransactionHash = blockvalidation.ReverseEndian(Block.Transactions[transactionIndex].TransactionHash)
  }

  return nil

}






//ParseBytesOnly takes a byte array and extracts block features
func ParseBytesOnly(b *block.Block, bytes []byte) (error) {

  var magicnumber uint32
  filefunctions.ReadBinaryToUInt32(bytes[0:4], &magicnumber)
  fmt.Println("Magic Number: ", magicnumber, bytes[0:4])
  b.MagicNumber = magicnumber

  var blocklength uint32
  filefunctions.ReadBinaryToUInt32(bytes[4:8], &blocklength)
  fmt.Println("Block Length: ", blocklength, bytes[4:8])
  b.BlockLength = blocklength

  var formatversion uint32
  filefunctions.ReadBinaryToUInt32(bytes[8:12], &formatversion)
  fmt.Println("Format Version: ", formatversion, bytes[8:12])
  b.Header.FormatVersion = formatversion

  var previousblockhash string
  filefunctions.ReadUInt8ByteArrayLength32(bytes[12:44], &previousblockhash)
  previousblockhash = blockvalidation.ReverseEndian(previousblockhash)
  fmt.Println("Previous Block Hash: ", previousblockhash, bytes[12:44])
  b.Header.PreviousBlockHash = previousblockhash

  var merkleroot string
  filefunctions.ReadUInt8ByteArrayLength32(bytes[44:76], &merkleroot)
  fmt.Println("Merkle Root: ", merkleroot, bytes[44:76])
  b.Header.MerkleRoot = merkleroot

  var timestamp uint32
  filefunctions.ReadBinaryToUInt32(bytes[76:80], &timestamp)
  fmt.Println("Time Stamp: ", timestamp, bytes[76:80])
  b.Header.TimeStamp = timestamp

  var targetvalue uint32
  filefunctions.ReadBinaryToUInt32(bytes[80:84], &targetvalue)
  fmt.Println("Target Value: ", targetvalue, bytes[80:84])
  b.Header.TargetValue = targetvalue

  var nonce uint32
  filefunctions.ReadBinaryToUInt32(bytes[84:88], &nonce)
  fmt.Println("Nonce: ", nonce, bytes[84:88])
  b.Header.Nonce = nonce

  transactioncount, index, err := filefunctions.ReadVarIntFromBytes(bytes, 88)
  if err != nil {
    return err
  }
  fmt.Println("Transaction Count: ", transactioncount, bytes[88:index])
  b.TransactionCount = transactioncount

  var txHolder []block.Transaction

  for t := 0; t < int(b.TransactionCount); t++ {

    var tx block.Transaction

    var txversionnumber uint32
    filefunctions.ReadBinaryToUInt32(bytes[index:index+4], &txversionnumber)
    fmt.Println("Transaction Version Number: ", txversionnumber, bytes[index:index+4])
    tx.TransactionVersionNumber = txversionnumber

    fmt.Println("preindex", index)
    var inputcount uint64
    inputcount, index, err = filefunctions.ReadVarIntFromBytes(bytes, index+4)
    if err != nil {
      return err
    }
    fmt.Println("Input Count: ", inputcount)
    tx.InputCount = inputcount

    /*var transactionhash string
    filefunctions.ReadUInt8ByteArrayLength32(bytes[index:index+32], &transactionhash)
    transactionhash = blockvalidation.ReverseEndian(transactionhash)
    fmt.Println("Transaction Hash: ", transactionhash, bytes[index:index+32])
    tx.TransactionHash = transactionhash*/


    var inHolder []block.Input

    for i := 0; i < int(tx.InputCount); i++ {

      var in block.Input

      fmt.Println("postindex", index)

      var intransactionhash string
      filefunctions.ReadUInt8ByteArrayLength32(bytes[index:index+32], &intransactionhash)
      intransactionhash = blockvalidation.ReverseEndian(intransactionhash)
      fmt.Println("Input Transaction Hash: ", intransactionhash, bytes[index:index+32])
      in.TransactionHash = intransactionhash

      var transactionindex uint32
      filefunctions.ReadBinaryToUInt32(bytes[index+32:index+36], &transactionindex)
      fmt.Println("Transaction Index: ", transactionindex, bytes[index+32:index+36])
      in.TransactionIndex = transactionindex

      var inscriptlength uint64
      inscriptlength, index, err = filefunctions.ReadVarIntFromBytes(bytes, index+36)
      if err != nil {
        return err
      }
      fmt.Println("Input Script Length: ", inscriptlength)
      in.InputScriptLength = inscriptlength

      var inputscript string
      filefunctions.ReadUInt8ByteArrayToString(bytes[index:index+int(in.InputScriptLength)], &inputscript)
      fmt.Println("Input Script: ", inputscript, bytes[index:index+int(in.InputScriptLength)])
      in.InputScript = inputscript

      var sequencenumber uint32
      filefunctions.ReadBinaryToUInt32(bytes[index+int(in.InputScriptLength):index+int(in.InputScriptLength)+4], &sequencenumber)
      fmt.Println("Sequence Number: ", sequencenumber, bytes[index+int(in.InputScriptLength):index+int(in.InputScriptLength)+4])
      in.SequenceNumber = sequencenumber

      index = index + int(in.InputScriptLength) + 4
      inHolder = append(inHolder, in)

    }

    var outputcount uint64
    outputcount, index, err = filefunctions.ReadVarIntFromBytes(bytes, index)
    if err != nil {
      return err
    }
    fmt.Println("Output Count: ", outputcount)
    tx.OutputCount = outputcount

    var outHolder []block.Output

    for o := 0; o < int(tx.OutputCount); o++ {

      var out block.Output

      var outvalue uint64
      filefunctions.ReadBinaryToUInt64(bytes[index:index+8], &outvalue)
      fmt.Println("Output Value: ", outvalue, bytes[index:index+8])
      out.OutputValue = outvalue

      var outlength uint64
      outlength, index, err = filefunctions.ReadVarIntFromBytes(bytes, index+8)
      if err != nil {
        return err
      }
      fmt.Println("Output Length: ", outlength)
      out.ChallengeScriptLength = outlength

      var outscript string
      filefunctions.ReadUInt8ByteArrayToString(bytes[index:index+int(out.ChallengeScriptLength)], &outscript)
      fmt.Println("Output Script: ", outscript, bytes[index:index+int(out.ChallengeScriptLength)])
      out.ChallengeScript = outscript
      out.ChallengeScriptBytes = bytes[index:index+int(out.ChallengeScriptLength)]

      blockvalidation.ParseOutputScript(&out)

      fmt.Println("Addresss: ", out.Addresses[0])

      index = index + int(out.ChallengeScriptLength)
      outHolder = append(outHolder, out)

    }

    var transactionlock uint32
    filefunctions.ReadBinaryToUInt32(bytes[index:index+4], &transactionlock)
    fmt.Println("Transaction Lock Time: ", transactionlock, bytes[index:index+4])
    tx.TransactionLockTime = transactionlock

    txHolder = append(txHolder, tx)

  }

  return nil
}







// ParseBlock parses a block using the functions in blockchainbuilder
func ParseBlock(Block *block.Block, file *os.File) (error) {

  bmagicNumber, err := readExactMagicNumber(file)
  if err != nil {
    fmt.Println("No magic number recovered", err, bmagicNumber)
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

  //Update ByteOffset and ParsedBlockLength fields to track where in the file the block ends
  Block.HashBlock.ParsedBlockLength = Block.BlockLength

  filefunctions.SetByteCount(0)

  Block.Header.FormatVersion, Block.Header.ByteFormatVersion, err = readFormatVersion(file)
  if err != nil {
    fmt.Println("Error reading format version", Block.Header.FormatVersion, err)
    return err
  }
  fmt.Println("Format Version: ", Block.Header.FormatVersion)

  Block.Header.PreviousBlockHash, Block.Header.BytePreviousBlockHash, err = readPreviousBlockHash(file)
  if err != nil {
    fmt.Println("Error reading previous block hash", err)
    return err
  }
  fmt.Println("Previous Block Hash: ", blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash))

  //Update HashBlock PreviousBlockHash field with parsed value. Will be used to build Main Chain
  Block.HashBlock.PreviousCompressedBlockHash = btchashing.ComputeCompressedBlockHash(blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash))
  Block.HashBlock.CompressedBlockHash = blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash)

  Block.Header.MerkleRoot, Block.Header.ByteMerkleRoot, err = readMerkleRoot(file)
  if err != nil {
    fmt.Println("Error reading merkle root", err)
    return err
  }
  fmt.Println("Merkle Root: ", blockvalidation.ReverseEndian(Block.Header.MerkleRoot))

  Block.Header.TimeStamp, Block.Header.ByteTimeStamp, err = readTimeStamp(file)
  if err != nil {
    fmt.Println("Error reading timestamp", err)
    return err
  }
  fmt.Println("Time Stamp: ", blockvalidation.ConvertUnixEpochToDate(Block.Header.TimeStamp))

  //Update TimeStamp of hashblock
  Block.HashBlock.TimeStamp = Block.Header.TimeStamp

  Block.Header.TargetValue, Block.Header.ByteTargetValue, err = readTargetValue(file)
  if err != nil {
    fmt.Println("Error reading target value", err)
    return err
  }
  fmt.Println("Target Value: ", Block.Header.TargetValue)

  Block.Header.Nonce, Block.Header.ByteNonce, err = readNonce(file)
  if err != nil {
    fmt.Println("Error reading nonce", err)
    return err
  }
  fmt.Println("Nonce: ", Block.Header.Nonce)

  Block.BlockHash, err = btchashing.ComputeBlockHash(Block)
  if err != nil {
    fmt.Println("Error computing block hash", err)
    return err
  }
  Block.BlockHash = blockvalidation.ReverseEndian(Block.BlockHash)
  fmt.Println("Block Hash: ", blockvalidation.ReverseEndian(Block.BlockHash))

  //Add BlockHash field to HashBlock object in Block and compress hash to limit search space for BlockChain hashmap
  Block.HashBlock.CompressedBlockHash = btchashing.ComputeCompressedBlockHash(blockvalidation.ReverseEndian(Block.BlockHash))
  Block.HashBlock.BlockHash = blockvalidation.ReverseEndian(Block.BlockHash)

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

    Block.Transactions[transactionIndex].TransactionVersionNumber, Block.Transactions[transactionIndex].ByteTransactionVersionNumber, err = readTransactionVersion(file)
    if err != nil {
      fmt.Println("Error reading transaction version number", Block.Transactions[transactionIndex].TransactionVersionNumber, err)
      return err
    }
    fmt.Println("Transaction Version: ", Block.Transactions[transactionIndex].TransactionVersionNumber, Block.Transactions[transactionIndex].ByteTransactionVersionNumber)

    Block.Transactions[transactionIndex].InputCount, Block.Transactions[transactionIndex].ByteInputCount, err = readInputCount(file)
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

      Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionHash, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteTransactionHash, err = readTransactionHash(file)
      if err != nil {
        fmt.Println("Error reading transaction hash", err)
        return err
      }
      fmt.Println("Transaction Hash: ", blockvalidation.ReverseEndian(Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionHash))

      Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionIndex, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteTransactionIndex, err = readTransactionIndex(file)
      if err != nil {
        fmt.Println("Error reading transaction index", err)
        return err
      }
      fmt.Println("Transaction Index: ", Block.Transactions[transactionIndex].Inputs[inputIndex].TransactionIndex)

      Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteInputScriptLength, err = readInputScriptLength(file)
      if err != nil {
        fmt.Println("Error reading script length", err)
        return err
      }
      fmt.Println("Script Length: ", Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteInputScriptLength)

      Block.Transactions[transactionIndex].Inputs[inputIndex].InputScript, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteInputScript, err = readInputScriptBytes(int(Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength), file)
      if err != nil {
        fmt.Println("Error reading script bytes", err)
        return err
      }
      fmt.Println("Input Script: ", Block.Transactions[transactionIndex].Inputs[inputIndex].InputScript)

      Block.Transactions[transactionIndex].Inputs[inputIndex].SequenceNumber, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteSequenceNumber, err = readSequenceNumber(file)
      if err != nil {
        fmt.Println("Error reading sequence number", err)
        return err
      }
      fmt.Println("Sequence Number: ", Block.Transactions[transactionIndex].Inputs[inputIndex].SequenceNumber, Block.Transactions[transactionIndex].Inputs[inputIndex].ByteSequenceNumber)

    }

    Block.Transactions[transactionIndex].OutputCount, Block.Transactions[transactionIndex].ByteOutputCount, err = readOutputCount(file)
    if err != nil {
      fmt.Println("Error reading output count", err)
      return err
    }
    fmt.Println("Output Count: ", Block.Transactions[transactionIndex].OutputCount, Block.Transactions[transactionIndex].ByteOutputCount)

/**********************************Outputs*************************************
 ******************************************************************************/

    for outputIndex := 0; outputIndex < int(Block.Transactions[transactionIndex].OutputCount); outputIndex++ {

      fmt.Println("**** Output " , outputIndex + 1, " of ", int(Block.Transactions[transactionIndex].OutputCount), " ****")

      Block.Transactions[transactionIndex].Outputs = append(Block.Transactions[transactionIndex].Outputs, block.Output{})

      Block.Transactions[transactionIndex].Outputs[outputIndex].OutputValue, Block.Transactions[transactionIndex].Outputs[outputIndex].ByteOutputValue, err = readOutputValue(file)
      if err != nil {
        fmt.Println("Error reading output value", err)
        return err
      }
      fmt.Println("Output Value: ", Block.Transactions[transactionIndex].Outputs[outputIndex].OutputValue)

      Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength, Block.Transactions[transactionIndex].Outputs[outputIndex].ByteChallengeScriptLength, err = readChallengeScriptLength(file)
      if err != nil {
        fmt.Println("Error reading challenge script length", err)
        return err
      }
      fmt.Println("Challenge Script Length: ", Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength)

      Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScript, Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptBytes, err = readChallengeScriptBytes(int(Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptLength), file)
      if err != nil {
        fmt.Println("Error reading challenge script bytes", err)
        return err
      }
      fmt.Println("Challenge Script: ", Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScript)

      Block.Transactions[transactionIndex].Outputs[outputIndex].KeyType, err = blockvalidation.ParseOutputScript(&Block.Transactions[transactionIndex].Outputs[outputIndex])
      if err != nil {
        return err
      }

      fmt.Println("Hash160: ", Block.Transactions[transactionIndex].Outputs[outputIndex].Addresses[0].RipeMD160)
      fmt.Println("Address: ", Block.Transactions[transactionIndex].Outputs[outputIndex].Addresses[0].Address)
      fmt.Println("PublicKey: ", Block.Transactions[transactionIndex].Outputs[outputIndex].Addresses[0].PublicKey)

    }

    Block.Transactions[transactionIndex].TransactionLockTime, Block.Transactions[transactionIndex].ByteTransactionLockTime, err = readTransactionLockTime(file)
    if err != nil {
      fmt.Println("Error reading transaction lock time", err)
      return err
    }
    fmt.Println("Transaction Lock Time: ", Block.Transactions[transactionIndex].TransactionLockTime)

    Block.Transactions[transactionIndex].TransactionHash, err = btchashing.ComputeTransactionHash(&Block.Transactions[transactionIndex], Block.Transactions[transactionIndex].InputCount, Block.Transactions[transactionIndex].OutputCount)
    if err != nil {
      fmt.Println("Error in computing transaction hash", err)
      return err
    }
    fmt.Println("Transaction Hash: ", blockvalidation.ReverseEndian(Block.Transactions[transactionIndex].TransactionHash))
  }

  return nil

}
