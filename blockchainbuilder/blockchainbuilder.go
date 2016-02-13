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
  if magicNumber == 0 {
    b = filefunctions.LookForMagic(file)
    err := filefunctions.ReadBinaryToUInt32(b, &magicNumber)
    if err != nil {
      fmt.Println("failed to find MagicNumber from zeros: ", err)
    }
    if blockvalidation.ValidateMagicNumber(magicNumber) {
      return magicNumber, nil
    }
  }
  fmt.Println("bad magic number", magicNumber)
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
  if blockvalidation.ValidateTransactionIndex(transactionIndex) {
    return transactionIndex, b, nil
  }
  b, err = filefunctions.RewindAndRead32(b, file, &transactionIndex)
  if err != nil {
    fmt.Println("rewind read failed: ", err)
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
  if blockvalidation.ValidateOutputValue(outputValue) {
    return outputValue, b, nil
  }
  b, err = filefunctions.RewindAndRead64(b, file, &outputValue)
  if err != nil {
    fmt.Println("rewind read failed: ", err)
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



/***************************OUTER LOOPS****************************************/


//ParseIndividualBlock parses a block using the functions in blockchainbuilder
func (Blockchain) ParseIndividualBlock(Block *block.Block, file *os.File) (error) {

  Block.HashBlock.FilePointer = file

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

  Block.Header.FormatVersion, Block.Header.ByteFormatVersion, err = readFormatVersion(file)
  if err != nil {
    fmt.Println("Error reading format version", err)
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
  Block.HashBlock.PreviousBlockHash = blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash)
  Block.HashBlock.PreviousCompressedBlockHash = btchashing.ComputeCompressedBlockHash(blockvalidation.ReverseEndian(Block.Header.PreviousBlockHash))

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
      fmt.Println("Error reading transaction version number", err)
      return err
    }
    fmt.Println("Transaction Version: ", Block.Transactions[transactionIndex].TransactionVersionNumber)

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
      fmt.Println("Script Length: ", Block.Transactions[transactionIndex].Inputs[inputIndex].InputScriptLength)

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
      fmt.Println("Sequence Number: ", Block.Transactions[transactionIndex].Inputs[inputIndex].SequenceNumber)

    }

    Block.Transactions[transactionIndex].OutputCount, Block.Transactions[transactionIndex].ByteOutputCount, err = readOutputCount(file)
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
      fmt.Println("Bytes: ", Block.Transactions[transactionIndex].Outputs[outputIndex].ChallengeScriptBytes)

      _, err = blockvalidation.ParseOutputScript(&Block.Transactions[transactionIndex].Outputs[outputIndex])

    /*  Block.Transactions[transactionIndex].Outputs[outputIndex].Address, err = btchashing.ParseAddressFromOutputScript(Block.ByteTransactions[transactionIndex].Outputs[outputIndex], Block.Transactions[transactionIndex].Outputs[outputIndex])
      if err != nil {
        fmt.Println("Error reading address from output script", err)
        return err
      }
      fmt.Println("Address: ", Block.Transactions[transactionIndex].Outputs[outputIndex].Address)
*/
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

  //Update ByteOffset and ParsedBlockLength fields to track where in the file the block ends
  Block.HashBlock.ByteOffset = filefunctions.ByteCount
  Block.HashBlock.ParsedBlockLength = Block.BlockLength

  _, err = filefunctions.ResetBlockHeadPointer(Block.BlockLength, file)
  if err != nil {
    fmt.Println("Error in resetting block head pointer", err)
  }
  return nil

}


//WriteMainChainToFile writes the binary data of the compressed HashBlock to filename
func (Blockchain) WriteMainChainToFile(lastKey string, chain *Blockchain, filename string) (error) {

  f, err := os.Create("" + filename + ".dat")
  if err != nil {
    return err
  }

  fmt.Println("Writing to file...")

  for chain.BlockMap[lastKey].CompressedBlockHash != "" {
    lastKey = chain.BlockMap[lastKey].PreviousCompressedBlockHash
    _, err = f.Write([]byte(chain.BlockMap[lastKey].BlockHash))
    if err != nil {
      return err
    }
    _, err = f.WriteString("\n")
    if err != nil {
      return err
    }
    fmt.Println("Wrote", chain.BlockMap[lastKey].BlockHash)
  }
  defer f.Close()
  return nil
}
