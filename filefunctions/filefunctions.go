package filefunctions

import (
    "bytes"
    "encoding/binary"
    "encoding/hex"
    "log"
    "os"
    "errors"
)

var byteCount int

//Possible64ByteErrorFlag tracks whether we've hit a missed byte in parsing in the main method output count
//var Possible64ByteErrorFlag bool

//ReadNextBytes reads number of bytes from binary file
func ReadNextBytes(file *os.File, number int) []byte {
  bytes := make([]byte, number)
  byteCount = byteCount + number

  _, err := file.Read(bytes)
  if err != nil {
      log.Fatal(err)
  }
  return bytes
}

//RewindAndRead64 is called when Possible64ByteErrorFlag is raised. The function moves the file pointer back and re-reads with fewer bytes included in the read
func RewindAndRead64(b []byte, file *os.File, outputValue *uint64) ([]byte, error) {
  var secondTryLen int64 = 7
    if (b[0] != byte(0) || b[1] != byte(0) || b[2] != byte(0)) && (b[7] != byte(0) && b[6] == byte(0) && b[5] == byte(0)) {
      bytesTwo := make([]byte, secondTryLen)
      byteCount = byteCount - (len(b) - int(secondTryLen))

      a, b := file.Seek(-(secondTryLen + 1), 1)
      log.Println("seeking...", a, b)
      _, err := file.Read(bytesTwo)
      if err != nil {
        log.Fatal(err)
      }
      ReadBinaryToUInt64(bytesTwo, outputValue)
      return bytesTwo, nil
    }
  return nil, errors.New("Could not rewind and read smaller uint64 integer")
}

//RewindAndRead32 is called when we fail validation of unsigned 32 bit integer and want to skip back a bit and restart parsing
func RewindAndRead32(b []byte, file *os.File, transactionIndex *uint32) ([]byte, error) {
  var secondTryLen int64 = 3
    bytesTwo := make([]byte, secondTryLen)
    byteCount = byteCount - (len(b) - int(secondTryLen))

    a, c := file.Seek(-(secondTryLen + 1), 1)
    log.Println("seeking...", a, c)
    _, err := file.Read(bytesTwo)
    if err != nil {
      log.Fatal(err)
    }
    ReadBinaryToUInt32(bytesTwo, transactionIndex)
    return bytesTwo, nil
  }

//LookForMagic handles instance when encounter string of zeros in searching for Magic Number
func LookForMagic(file *os.File) ([]byte) {
  iter := make([]byte, 1)
  zero := make([]byte, 1)
  for bytes.Equal(iter, zero) {
    log.Print("found zero")
    _, err := file.Read(iter)
    if err != nil {
      log.Fatal(err)
    }
  }
  byteCount = byteCount + 1
  ret := append(iter[:], ReadNextBytes(file, 3)[:] ...)
  byteCount = byteCount + 3
  return ret
}

//ReadUInt8ByteArray reads a bytestream into an array of their values
func ReadUInt8ByteArray(b []byte, passedVariable *[]uint8) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, passedVariable)
  return err
}

//ReadUInt8ByteArrayLength32 reads a bytestream into a string hash with 32 characters
func ReadUInt8ByteArrayLength32(b []byte, passedVariable *string) {
  *passedVariable = hex.EncodeToString(b)
}

//ReadBinaryToUInt8 reads a binary bytestream into an unsigned integer byte
func ReadBinaryToUInt8(b []byte, passedVariable *uint8) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, passedVariable)
  return err
}

//ReadBinaryToUInt16 reads a binary bytestream into an unsigned integer of 2 bytes
func ReadBinaryToUInt16(b []byte, passedVariable *uint16) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, passedVariable)
  return err
}

//ReadBinaryToUInt32 reads a binary bytestream into an unsigned integer of 4 bytes
func ReadBinaryToUInt32(b []byte, passedVariable *uint32) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, passedVariable)
  return err
}

//ReadBinaryToUInt64 reads a binary bytestream into an unsigned integer of 8 bytes
func ReadBinaryToUInt64(b []byte, passedVariable *uint64) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, passedVariable)
  return err
}


//ReadVariableLengthInteger reads a variable length integer as described by the bitcoin protocol into an unsigned 8 byte integer
func ReadVariableLengthInteger(file *os.File) (uint64, []byte, error) {

  var ret uint64
  var eight uint8
  var byteret []byte

  bytes := make([]byte, 1)
  _, err := file.Read(bytes)
  if err != nil {
    return ret, nil, err
  }
  err = ReadBinaryToUInt8(bytes, &eight)
  if err != nil {
    return ret, nil, err
  }
  if eight < 0xFD {       // If it's less than 0xFD use this value as the unsigned integer
    byteCount++
    ret = uint64(eight)
    byteret = bytes
  } else {
      var sixteen uint16
      bytes = make([]byte, 3)
      _, err = file.Read(bytes)
      if err != nil {
        return ret, nil, err
      }
      err = ReadBinaryToUInt16(bytes, &sixteen)
      if err != nil {
        return ret, nil, err
      }
      if sixteen < 0xFFFF {
        byteCount += 3
        ret = uint64(sixteen)
        byteret = bytes
      } else {
          var thirtytwo uint32
          bytes = make([]byte, 5)
          _, err = file.Read(bytes)
          if err != nil {
            return ret, nil, err
          }
          err = ReadBinaryToUInt32(bytes, &thirtytwo)
          if err != nil {
            return ret, nil, err
          }
          if thirtytwo < 0xFFFFFFFF {
            byteCount += 5
            ret = uint64(thirtytwo)
            byteret = bytes
          } else {      // never expect to actually encounter a 64bit integer in the block-chain stream; it's outside of any reasonable expected value
              var sixtyfour uint64
              bytes = make([]byte, 9)
              _, err = file.Read(bytes)
              if err != nil {
                return ret, nil, err
              }
              err = ReadBinaryToUInt64(bytes, &sixtyfour)
              if err != nil {
                return ret, nil, err
              }
              byteCount += 9
              ret = uint64(sixtyfour)
              byteret = bytes
            }
          }
      }
  return ret, byteret, nil
}

//ResetBlockHeadPointer points the byte-reader to the next block in the chain
func ResetBlockHeadPointer(blockLength uint32, file *os.File) ([]byte, error) {
  if byteCount <= int(blockLength) {
    bytes := make([]byte, int(blockLength) - byteCount)
    ReadNextBytes(file, int(blockLength) - byteCount)
    _, err := file.Read(bytes)
    if err != nil {
        log.Fatal(err)
    }
    return bytes, nil
  }
  return nil, errors.New("used more bytes than listed in blocklength")
}

//GetByteCount returns the global byteCount variable in filefunctions class
func GetByteCount() (int) {
  return byteCount
}

//SetByteCount sets the global byteCount variable in filefunctions class
func SetByteCount(newVal int) {
  byteCount = newVal
}

//IncrementByteCount increments by incrementVal the global byteCount variable in filefunctions class
func IncrementByteCount(incrementVal int) {
  byteCount += incrementVal
}
