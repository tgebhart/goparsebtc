package blockvalidation

import (
  "github.com/tgebhart/goparsebtc/block"
  "net/http"
  //"bytes"
  "fmt"
  "io/ioutil"
  "encoding/json"
  "errors"
  "time"
)

//ResponseBlock holds the blockchain.info response json when querying a block through
//blockchain.info API. It should be noted that compound names like prevblock are
//represented as prev_block with underscores. However, underscores are to be avoided
//in Go.
type ResponseBlock struct {
  hash string
  ver int
  prevblock string
  mrklroot string
  time int
  bits int
  fee int
  nonce int
  ntx int
  size int
  blockindex int
  mainchain bool
  height int
  tx ResponseTransaction
}

//ResponseTransaction holds the blockchain.info response json for a given transaction in a block
type ResponseTransaction struct {
  locktime int
  ver int
  size int
  inputs ResponseInput
  time int
  txindex int
  vinsz int
  hash string
  voutsz int
  relayedby string
  out ResponseOutput
}

//ResponseInput holds the blockchain.info response json for a block's input.
type ResponseInput struct {
  sequence int
  script string
}

//ResponseOutput holds the blockchain.info response json for a block's output.
//Note that responsetype is returned with key "type" by blockchain.info, but this
//is a Go reserved word
type ResponseOutput struct {
  spent bool
  txindex int
  responsetype int
  addr string
  value int
  n int
  script string
}


//BLOCKCHAININFOENDPOINT is the API endpoint for json information from blockchain.info
var BLOCKCHAININFOENDPOINT = "https://blockchain.info/rawblock/"

//REQUESTTYPE denotes the variable type when using http call
var REQUESTTYPE = "string"

//ValidateMagicNumber checks for correct magic number. Can take one of two values
func ValidateMagicNumber(magicNumber uint32) (bool) {
  if magicNumber == 3652501241 || magicNumber == 4190024921 {
    return true
  }
  return false
}

//ValidateBlockLength checks block length.  Should be less than maximum possible block length
func ValidateBlockLength(blockLength uint32) (bool) {
  if blockLength <= 4294967295 && blockLength > 0 { //2^32 -1 ~ 4GB or maximum possible block length
    return true
  }
  return false
}

//ValidateFormatVersion checks the block's format version (should be 1 for now)
func ValidateFormatVersion(formatVersion uint32) (bool) {
  if formatVersion == 1 {  //format version should still be 1 for now
    return true
  }
  return false
}

//ValidateTimeStamp checks for block timeStamp to be between timestamp of genesis block and maximum integer value
func ValidateTimeStamp(timeStamp uint32) (bool) {
  if timeStamp >= 1231006505 && timeStamp <= 4294967295 {  //genesis block UNIX epoch time && maximum value for unsigned integer
    return true
  }
  return false
}

//ValidateTransactionVersion checks transaction version. Should be equal to 1 currently
func ValidateTransactionVersion(transactionVersion uint32) (bool) {
  if transactionVersion == 1 {  //current transaction version
    return true
  }
  return false
}

//ValidateSequenceNumber checks to make sure sequence number is below the maximum integer value
func ValidateSequenceNumber(sequenceNumber uint32) (bool) {
  if sequenceNumber <= 4294967295 {  //current largest sequence number
    return true
  }
  return false
}

//ValidateTransactionLockTime checks transaction lock time is equal to 0
func ValidateTransactionLockTime(transactionLockTime uint32) (bool) {
  if transactionLockTime == 0 {
    return true
  }
  return false
}

//ReverseEndian switches the output of as 32 byte hash to Big-Endian from Little-Endian because blockchain.info is weird
func ReverseEndian(s string) (string) {
  var tempstring [64]string
  for i := 0; i < len(s) - 1; i+= 2 {
    tempstring[63 - i] = string(s[i]) + string(s[i+1])
  }
  var ret string
  for j := 0; j < len(tempstring); j++ {
    ret += tempstring[j]
  }
  return ret
}

func narcolepsy() {
  time.Sleep(100 * time.Millisecond)
}

//BlockChainInfoValidation calls blockchain.info and checks the block for near-real-time error-checking
func BlockChainInfoValidation(Block *block.Block) (error) {
  data := map[string]interface{}{}
  blockHash := ReverseEndian(Block.BBlockHash)
  fmt.Println("block hash", blockHash)
  resp, err := http.Get(BLOCKCHAININFOENDPOINT + blockHash)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  json.Unmarshal(body, &data)

  if blockHash == data["hash"]{
    return nil
  }
  return errors.New("Hashes do not match")
}
