package blockchainreader

import (
    "errors"
    "fmt"
    "os"
    "github.com/tgebhart/goparsebtc/block"
  //  "github.com/tgebhart/goparsebtc/filefunctions"
    "github.com/tgebhart/goparsebtc/blockchainbuilder"
    "github.com/tgebhart/goparsebtc/blockvalidation"
    "encoding/csv"
    "strconv"
    //"github.com/aws/aws-sdk-go/aws"
    //"github.com/aws/aws-sdk-go/aws/session"
    //"github.com/aws/aws-sdk-go/service/dynamodb"
)

//Blockchain holds the BlockMap object
type Blockchain struct {
  BlockMap map[string]block.DBlock
}
//NewBlockchain constructs a Blockchain instance
func NewBlockchain() *Blockchain {
  var b Blockchain
  b.BlockMap = make(map[string]block.DBlock)
  return &b
}
//ReadChain holds read csv file structure
type ReadChain struct {
  ReadBlocks []ReadBlock
}
//NewReadChain constructs new read chain
func NewReadChain() *ReadChain {
  var c ReadChain
  return &c
}
//ReadBlock holds row information from csv file
type ReadBlock struct {
  BlockHash string
  FileEndpoint string
  ByteOffset int
  BlockLength int
  RawBlockNumber int
  TimeStamp int
}

//ErrCompareHashes is thrown when the block hashes of the two sources do not match
var ErrCompareHashes = errors.New("Error comparing read file and dat file block hashes")
//ErrHashExists is thrown when trying to add a block to blockchain and its key already exists
var ErrHashExists = errors.New("Error: hash in blockchain already exists")

//ReadReferenceFile reads csv data from reference file and fills ReadChain and ReadBlock
func ReadReferenceFile(r *ReadChain, location string) (error) {

  fcsv, err := os.Open(location)
  if err != nil {
    return err
  }

  defer fcsv.Close()

  reader := csv.NewReader(fcsv)

  rawCSV, err := reader.ReadAll()
  if err != nil {
    return err
  }

  var tempBlock ReadBlock
  var tempChain []ReadBlock

  for _, each := range rawCSV {
    tempBlock.BlockHash = each[0]
    tempBlock.FileEndpoint = each[1]
    tempBlock.ByteOffset, _ = strconv.Atoi(each[2])
    tempBlock.BlockLength, _ = strconv.Atoi(each[3])
    tempBlock.RawBlockNumber, _ = strconv.Atoi(each[4])
    //tempBlock.TimeStamp, _ = strconv.Atoi(each[5])
    tempChain = append(tempChain, tempBlock)
  }

  r.ReadBlocks = tempChain
  return nil
}

//LoadChain populates the Blockchain hashmap by reading in all files designated
//in the readChain struct. Requires location of .dat files
func LoadChain(chain *Blockchain, readchain *ReadChain, datLocation string) (error) {

  var dBlock block.DBlock
  var fBlock block.Block
  var datEndpoint string
  var file *os.File
  var err error

  for i := 0; i < len(readchain.ReadBlocks) - 2; i++ {
    b := readchain.ReadBlocks[i]
    fmt.Println(b.FileEndpoint)

    nextEndpoint := b.FileEndpoint

    fmt.Println("Read b: ", b.BlockHash)

    if nextEndpoint == "" {
      err := blockvalidation.BridgeWithBlockchainInfo(&dBlock, b.BlockHash)
      if err != nil {
        return err
      }
      fmt.Println("dBlock: ", dBlock)
    } else {
      fmt.Println("compare", nextEndpoint, datEndpoint)
      if nextEndpoint != datEndpoint {
        file, err = os.Open(datLocation + nextEndpoint)
        if err != nil {
            return err
        }
        defer file.Close()
        datEndpoint = nextEndpoint
      }

      err := readBlock(&fBlock, readchain.ReadBlocks[i+1].ByteOffset, readchain.ReadBlocks[i+1].BlockLength, file)
      if err != nil {
        fmt.Println("LoadChain: ")
        return err
      }

      fmt.Println(fBlock.Header, fBlock.BlockHash)

      for fBlock.BlockHash != b.BlockHash {
        return ErrCompareHashes
      }

      err = MapBlockToDBlock(&fBlock, &dBlock)
    }

    err := putBlock(chain, dBlock)
    if err != nil {
      return err
    }

  }

  return nil
}
/*
//UploadChain uploads the full chain to AWS DynamoDB
func UploadFullChain(chain *block.Blockchain) (error) {

  for hashkey, dblock := range chain.BlockMap {

  }
}
*/


func readBlock(b *block.Block, startByte int, length int, file *os.File) (error) {

  file.Seek(int64(startByte - length - 4), 0)
  err := blockchainbuilder.ParseBlockOnly(b, file)
  if err != nil {
    fmt.Println("readBlock")
    return err
  }
  return nil
}

func putBlock(chain *Blockchain, d block.DBlock) (error) {

  _, exists := chain.BlockMap[d.BlockHash]
  if exists {
    return ErrHashExists
  }

  chain.BlockMap[d.BlockHash] = d

  return nil
}


//MapBlockToDBlock maps a block read from the .dat files to a database block
func MapBlockToDBlock(b *block.Block, d *block.DBlock) (error) {

  d.MagicNumber = int(b.MagicNumber)
  d.BlockLength = int(b.BlockLength)
  d.BlockHash = b.BlockHash
  d.FormatVersion = int(b.Header.FormatVersion)
  d.PreviousBlockHash = b.Header.PreviousBlockHash
  d.MerkleRoot = b.Header.MerkleRoot
  d.TimeStamp = int(b.Header.TimeStamp)
  d.TargetValue = int(b.Header.TargetValue)
  d.Nonce = int(b.Header.Nonce)
  d.TransactionCount = int(b.TransactionCount)

  for t := 0; t < d.TransactionCount - 1; t++ {
    d.Transactions[t].TransactionIndex = 0
    d.Transactions[t].Time = d.TimeStamp
    d.Transactions[t].TransactionHash = b.Transactions[t].TransactionHash
    d.Transactions[t].TransactionVersionNumber = int(b.Transactions[t].TransactionVersionNumber)
    d.Transactions[t].InputCount = int(b.Transactions[t].InputCount)
    for i := 0; i < d.Transactions[t].InputCount - 1; i++ {
      d.Transactions[t].Inputs[i].TransactionHash = b.Transactions[t].Inputs[i].TransactionHash
      d.Transactions[t].Inputs[i].TransactionIndex = int(b.Transactions[t].Inputs[i].TransactionIndex)
      d.Transactions[t].Inputs[i].InputScriptLength = int(b.Transactions[t].Inputs[i].InputScriptLength)
      d.Transactions[t].Inputs[i].InputScript = b.Transactions[t].Inputs[i].InputScript
      d.Transactions[t].Inputs[i].SequenceNumber = int(b.Transactions[t].Inputs[i].SequenceNumber)
    }
    d.Transactions[t].OutputCount = int(b.Transactions[t].OutputCount)
    for o := 0; o < d.Transactions[t].OutputCount - 1; o++ {
      d.Transactions[t].Outputs[o].OutputValue = int(b.Transactions[t].Outputs[o].OutputValue)
      d.Transactions[t].Outputs[o].ChallengeScriptLength = int(b.Transactions[t].Outputs[o].ChallengeScriptLength)
      d.Transactions[t].Outputs[o].ChallengeScript = b.Transactions[t].Outputs[o].ChallengeScript
      d.Transactions[t].Outputs[o].KeyType = b.Transactions[t].Outputs[o].KeyType
      d.Transactions[t].Outputs[o].NumAddresses = len(b.Transactions[t].Outputs[o].Addresses)
      for a := 0; a < d.Transactions[t].Outputs[o].NumAddresses; a++ {
        d.Transactions[t].Outputs[o].Addresses[a].Address = b.Transactions[t].Outputs[o].Addresses[a].Address
      }
    }
    d.Transactions[t].TransactionLockTime = int(b.Transactions[t].TransactionLockTime)
  }
  return nil
}
