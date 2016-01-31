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
func ResetBlockHeadPointer(blockLength uint32, file *os.File) (error) {
  if byteCount <= int(blockLength) {
    ReadNextBytes(file, int(blockLength) - byteCount)
    return nil
  }
  return errors.New("used more bytes than listed in blocklength")
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
