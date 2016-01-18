package main

import (
   "errors"
    "bytes"
    "encoding/binary"
    "fmt"
    "log"
    "os"
    "flag"
    "strconv"
)

type Block struct {
  BMagicNumber uint32
  BBlockLength uint32
  BFormatVersion uint32
  BPreviousBlockHash [32]uint8
  BMerkleRoot [32]uint8
  BTimeStamp uint32
  BTargetValue uint32
  BNonce uint32
  BTransactionLength uint64
  BTransactionVersionNumber uint32
  BInputCount uint64
  BTransactionHash [32]uint8
  BTransactionIndex uint32
  BInputScriptLength uint64
  BInputScriptBytes []uint8
  BSequenceNumber uint32
  BOutputCount uint64
  BOutputValue uint64
  BChallengeScriptLength uint64
  BChallengeScriptBytes []uint8
  BTransactionLockTime uint32
}

type Header struct {
  HMagicNumber uint32
  HBlockLength uint32
  HFormatVersion uint32
  HPreviousBlockHash [32]uint8
  HMerkleRoot [32]uint8
  HTimeStamp uint32
  HTargetValue uint32
  HNonce uint32
}

type Transaction struct {
  TTransactionLength uint64
  TTransactionVersionNumber uint32
  TInputCount uint64
  TOutputCount uint64
  TTransactionLockTime uint32
}

type TransactionInterface interface {
}

type Input struct {
  ITransactionHash uint8
  ITransactionIndex uint32
  IInputScriptLength uint64
  IInputScriptBytes []uint8
  ISequenceNumber uint32
}

type Output struct {
  OOutputValue uint64
  OChallengeScriptLength uint64
  OChallengeScriptBytes []uint8
}

type InputInterface interface {
}

type Blockchain interface {

  readMagicNumber(file *os.File) (uint32, error)
  validateMagicNumber(pmagicNumber uint32) (bool)

  readBlockLength(file *os.File) (uint32, error)
  validateBlockLength(blockLength uint32) (bool)

  readFormatVersion(file *os.File) (uint32, error)
  validateFormatVersion(formatVersion uint32) (bool)

  readPreviousBlockHash(file *os.File) ([]uint8, error)


}

var byteCount int

func readMagicNumber(file *os.File) (uint32, error) {

  var magicNumber uint32
  b := readNextBytes(file, 4)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &magicNumber)
  if err != nil {
    fmt.Println("binary.Read failed:", err)
  }
  if validateMagicNumber(magicNumber) {
    return magicNumber, nil
  }
  return 0, errors.New("unusual or invalid magic number value")
}

func validateMagicNumber(magicNumber uint32) (bool) {
  if magicNumber == 3652501241 || magicNumber == 4190024921 {
    return true
  }
  return false
}

func readBlockLength(file *os.File) (uint32, error) {
  var blockLength uint32
  b := readNextBytes(file, 4)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &blockLength)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if validateBlockLength(blockLength) {
      return blockLength, nil
  }
  log.Print(b)
  return 0, errors.New("Very large (or no) block length")
}

func validateBlockLength(blockLength uint32) (bool) {
  if blockLength <= 4294967295 && blockLength > 0 { //2^32 -1 ~ 4GB or maximum possible block length
    return true
  }
  return false
}

func readFormatVersion(file *os.File) (uint32, error) {
  var formatVersion uint32
  b := readNextBytes(file, 4)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &formatVersion)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if validateFormatVersion(formatVersion) {
      return formatVersion, nil
  }
  return 0, errors.New("Unusual format version")
}

func validateFormatVersion(formatVersion uint32) (bool) {
  if formatVersion == 1 {  //format version should still be 1 for BTC
    return true
  }
  return false
}


func readPreviousBlockHash(file *os.File) ([32]uint8, error) {
  var previousBlockHash [32]uint8
  b := readNextBytes(file, 32)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &previousBlockHash)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return previousBlockHash, nil
}

func readMerkleRoot(file *os.File) ([32]uint8, error) {
  var merkleRoot [32]uint8
  b := readNextBytes(file, 32)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &merkleRoot)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return merkleRoot, nil
}

func readTimeStamp(file *os.File) (uint32, error) {
  var timeStamp uint32
  b := readNextBytes(file, 4)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &timeStamp)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if validateTimeStamp(timeStamp) {
    return timeStamp, nil
  }
  return 0, errors.New("Unexpected timestamp value")
}

func validateTimeStamp(timeStamp uint32) (bool) {
  if timeStamp >= 1231006505 && timeStamp <= 4294967295 {  //genesis block UNIX epoch time && maximum value for unsigned integer
    return true
  }
  return false
}

func readTargetValue(file *os.File) (uint32, error) {
  var targetValue uint32
  b := readNextBytes(file, 4)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &targetValue)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return targetValue, nil
}

func readNonce(file *os.File) (uint32, error) {
  var nonce uint32
  b := readNextBytes(file, 4)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &nonce)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return nonce, nil
}

func readTransactionLength(file *os.File) (uint64, error) {
  var transactionLength uint64
  b := readNextBytes(file, 1)
  buf := bytes.NewReader(b)
  transactionLength, err := binary.ReadUvarint(buf)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return transactionLength, nil
}

func readTransactionVersion(file *os.File) (uint32, error) {
  var transactionVersion uint32
  b := readNextBytes(file, 4)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &transactionVersion)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if validateTransactionVersion(transactionVersion) {
    return transactionVersion, nil
  }
  return 0, errors.New("Unexpected transaction version number")
}

func validateTransactionVersion(transactionVersion uint32) (bool) {
  if transactionVersion == 1 {  //current transaction version
    return true
  }
  return false
}

func readInputCount(file *os.File) (uint64, error) {
  var inputCount uint64
  b := readNextBytes(file, 1)
  buf := bytes.NewReader(b)
  inputCount, err := binary.ReadUvarint(buf)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return inputCount, nil
}

func readTransactionHash(file *os.File) ([32]uint8, error) {
  var transactionHash [32]uint8
  b := readNextBytes(file, 32)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &transactionHash)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return transactionHash, nil
}

func readTransactionIndex(file *os.File) (uint32, error) {
  var transactionIndex uint32
  b := readNextBytes(file, 4)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &transactionIndex)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return transactionIndex, nil
}

func readInputScriptLength(file *os.File) (uint64, error) {
  var inputScriptLength uint64
  b := readNextBytes(file, 1)
  buf := bytes.NewReader(b)
  inputScriptLength, err := binary.ReadUvarint(buf)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return inputScriptLength, nil
}

func readInputScriptBytes(inputScriptLength int, file *os.File) ([]uint8, error) {
  var inputScriptBytes = make([]uint8, inputScriptLength)
  b := readNextBytes(file, inputScriptLength)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &inputScriptBytes)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return inputScriptBytes, nil
}

func readSequenceNumber(file *os.File) (uint32, error) {
  var sequenceNumber uint32
  b := readNextBytes(file, 4)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &sequenceNumber)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if validateSequenceNumber(sequenceNumber) {
    return sequenceNumber, nil
  }
  fmt.Println("sequence b: ", b)
  fmt.Println("sequence: ", sequenceNumber)
  return 0, errors.New("Invalid sequence number")
}

func validateSequenceNumber(sequenceNumber uint32) (bool) {
  if sequenceNumber <= 4294967295 {  //current largest sequence number
    return true
  }
  return false
}

func readOutputCount(file *os.File) (uint64, error) {
  var outputCount uint64
  b := readNextBytes(file, 1)
  buf := bytes.NewReader(b)
  outputCount, err := binary.ReadUvarint(buf)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return outputCount, nil
}

func readOutputValue(file *os.File) (uint64, error) {
  var outputValue uint64
  b := readNextBytes(file, 8)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &outputValue)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return outputValue, nil
}

func readChallengeScriptLength(file *os.File) (uint64, error) {
  var challengeScriptLength uint64
  b := readNextBytes(file, 1)
  buf := bytes.NewReader(b)
  challengeScriptLength, err := binary.ReadUvarint(buf)
  if err != nil {
    fmt.Println("binary.ReadUvarint failed: ", err)
  }
  return challengeScriptLength, nil
}

func readChallengeScriptBytes(challengeScriptLength int, file *os.File) ([]uint8, error) {
  var challengeScriptBytes = make([]uint8, challengeScriptLength)
  b := readNextBytes(file, challengeScriptLength)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &challengeScriptBytes)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  return challengeScriptBytes, nil
}

func readTransactionLockTime(file *os.File) (uint32, error) {
  var transactionLockTime uint32
  b := readNextBytes(file, 4)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &transactionLockTime)
  if err != nil {
    fmt.Println("binary.Read failed: ", err)
  }
  if validateTransactionLockTime(transactionLockTime) {
    return transactionLockTime, nil
  }
  return 1, errors.New("Invalid Lock Time on Transaction")
}

func validateTransactionLockTime(transactionLockTime uint32) (bool) {
  if transactionLockTime == 0 {
    return true
  }
  return false
}


func getTransactions(transactionLength uint64, file *os.File) (error) {
  var i uint64
  for ; i < transactionLength; {

  }
  return nil
}

func resetBlockHeadPointer(blockLength uint32, bytesUsed int, file *os.File) (error) {
  fmt.Print(bytesUsed, " ", blockLength)
  if bytesUsed <= int(blockLength) {
    readNextBytes(file, int(blockLength) - bytesUsed)
    return nil
  }
  return errors.New("used more bytes than listed in blocklength")
}




/***************************OUTER LOOPS****************************************/

func parseIndividualBlock(Block Block, file *os.File) error {

  bmagicNumber, err := readMagicNumber(file)
  if err != nil {
    fmt.Println("No magic number recovered", err)
    return err
  }
  fmt.Println("Magic Number: ", Block.BMagicNumber)
  Block.BMagicNumber = bmagicNumber

  Block.BBlockLength, err = readBlockLength(file)
  if err != nil {
    fmt.Println("No blocklength recovered", err)
    return err
  }
  fmt.Println("Block Length: ", Block.BBlockLength)

  byteCount = 0

  Block.BFormatVersion, err = readFormatVersion(file)
  if err != nil {
    fmt.Println("Error reading format version", err)
    return err
  }
  fmt.Println("Format Version: ", Block.BFormatVersion)

  Block.BPreviousBlockHash, err = readPreviousBlockHash(file)
  if err != nil {
    fmt.Println("Error reading previous block hash", err)
    return err
  }
  fmt.Println("Previous Block Hash: ", Block.BPreviousBlockHash)

  Block.BMerkleRoot, err = readMerkleRoot(file)
  if err != nil {
    fmt.Println("Error reading merkle root", err)
    return err
  }
  fmt.Println("Merkle Root: ", Block.BMerkleRoot)

  Block.BTimeStamp, err = readTimeStamp(file)
  if err != nil {
    fmt.Println("Error reading timestamp", err)
    return err
  }
  fmt.Println("Time Stamp: ", Block.BTimeStamp)

  Block.BTargetValue, err = readTargetValue(file)
  if err != nil {
    fmt.Println("Error reading target value", err)
    return err
  }
  fmt.Println("Target Value: ", Block.BTargetValue)

  Block.BNonce, err = readNonce(file)
  if err != nil {
    fmt.Println("Error reading nonce", err)
    return err
  }
  fmt.Println("Nonce: ", Block.BNonce)

  Block.BTransactionLength, err = readTransactionLength(file)
  if err != nil {
    fmt.Println("Error reading transaction length", err)
    return err
  }
  fmt.Println("Transaction Length: ", Block.BTransactionLength)

/*===============================Transactions=================================
 ============================================================================*/

  for transactionIndex := 1; transactionIndex <= int(Block.BTransactionLength); transactionIndex++ {

    fmt.Println(" ========== Transaction ", transactionIndex, " of ", int(Block.BTransactionLength), " ============")

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
      fmt.Println("Transaction Hash: ", Block.BTransactionHash)

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

  err = resetBlockHeadPointer(Block.BBlockLength, byteCount, file)
  if err != nil {
    fmt.Println("Error in resetting block head pointer", err)
  }
  return nil


}




/******************************MAIN********************************************/

func main() {
  path := "/Users/tgebhart/Library/Application Support/Bitcoin/blocks/"
  flag.Parse()
    s := flag.Arg(0)
    numberOfFiles, err := strconv.Atoi(s)
    if err != nil {
        fmt.Println(err)
        os.Exit(2)
    }
  for j := 0; j < numberOfFiles; j++ {
    path = "/Users/tgebhart/Library/Application Support/Bitcoin/blocks/"
    e := strconv.Itoa(j)
    tempString := e
    for k := len(e); k < 5; k++ {
      tempString = "0" + tempString
    }
    path = path + "blk" + tempString + ".dat"

    file, err := os.Open(path)
    if err != nil {
        log.Fatal("Error while opening file", err)
    }
    fmt.Printf("%s opened\n", path)
    Block := Block{}
    err = nil
    for err == nil {
      err = parseIndividualBlock(Block, file)
      if err != nil {
        log.Println("error in parseIndividualBlock ", err)
      }
    }
    defer file.Close()
  }
}


/**************************FILE FUNCTIONS**************************************/


func readNextBytes(file *os.File, number int) []byte {
  bytes := make([]byte, number)
  byteCount = byteCount + number

  _, err := file.Read(bytes)
  if err != nil {
      log.Fatal(err)
  }
  return bytes
}

func readBinaryToUInt8(b []byte, passedVariable uint32) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &passedVariable)
  return err
}

func readBinaryToUInt32(b []byte, passedVariable uint32) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &passedVariable)
  return err
}

func readBinaryToUInt64(b []byte, passedVariable uint32) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &passedVariable)
  return err
}



func readVariableLengthInteger(file *os.File) (uint64, error) {

  var ret uint64
  var v uint32

  bytes := make([]byte, 1)
  _, err := file.Read(bytes)
  if err != nil {
    return ret, err
  }
  if bytes < 0xFD {       // If it's less than 0xFD use this value as the unsigned integer
    ret = uint64(bytes)
  } else {
      bytes = make([]byte, 2)
      _, err = file.Read(bytes)
      if err != nil {
        return ret, err
      }
      if bytes < 0xFFFF {
        ret = uint64(bytes)
      } else {
          bytes = make([]byte, 4)
          v = file.Read(bytes)
          if v < 0xFFFFFFFF {
            ret = uint64(v)
          } else {      // never expect to actually encounter a 64bit integer in the block-chain stream; it's outside of any reasonable expected value
              bytes = make([]byte, 8)
              v = file.Read(bytes)
              ret = uint64(v)
            }
          }
      }
  return ret
}
