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
  var datEndpoint string
  var file *os.File
  var err error

  for i := 1; i < len(readchain.ReadBlocks) - 1; i++ {
    var fBlock block.Block
    b := readchain.ReadBlocks[i]
    fmt.Println(b.FileEndpoint)

    nextEndpoint := b.FileEndpoint

    fmt.Println("Read b: ", b.BlockHash)

    if nextEndpoint == "" || readchain.ReadBlocks[i-1].ByteOffset == 0 {
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

      err := readBlock(&fBlock, readchain.ReadBlocks[i-1].ByteOffset, readchain.ReadBlocks[i-1].BlockLength, file)
      if err != nil {
        if err == blockchainbuilder.ErrBadMagic {
          err = blockvalidation.BridgeWithBlockchainInfo(&dBlock, b.BlockHash)
          if err != nil {
            return err
          }

        } else {
          fmt.Println("LoadChain: ")
          return err
        }
      }

      if fBlock.BlockHash != "" {

        fmt.Println(fBlock.BlockHash, b.BlockHash)

        if fBlock.BlockHash != b.BlockHash {
          return ErrCompareHashes
        }

        err = MapBlockToDBlock(&fBlock, &dBlock)
        if err != nil {
          return err
        }

      }

    } //end else

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

  //bytes := make([]byte, length)

  /*_ , err := file.ReadAt(bytes, int64(startByte))
  if err != nil {
    fmt.Println("read in block error")
    return err
  }*/

  file.Seek(int64(startByte), 0)

  //err = blockchainbuilder.ParseBlockOnly(b, file)
  err := blockchainbuilder.ParseBlock(b, file)
  if err != nil {
    fmt.Println("parse bytes only")
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

  var dTxs []block.DTransaction

  for t := 0; t < d.TransactionCount - 1; t++ {

    var dIns []block.DInput
    var tx block.DTransaction

    tx.TransactionIndex = 0
    tx.Time = d.TimeStamp
    tx.TransactionHash = b.Transactions[t].TransactionHash
    tx.TransactionVersionNumber = int(b.Transactions[t].TransactionVersionNumber)
    tx.InputCount = int(b.Transactions[t].InputCount)

    for i := 0; i < tx.InputCount - 1; i++ {

      var in block.DInput

      in.TransactionHash = b.Transactions[t].Inputs[i].TransactionHash
      in.TransactionIndex = int(b.Transactions[t].Inputs[i].TransactionIndex)
      in.InputScriptLength = int(b.Transactions[t].Inputs[i].InputScriptLength)
      in.InputScript = b.Transactions[t].Inputs[i].InputScript
      in.SequenceNumber = int(b.Transactions[t].Inputs[i].SequenceNumber)

      dIns = append(dIns, in)
    }

    tx.OutputCount = int(b.Transactions[t].OutputCount)
    var dOuts []block.DOutput

    for o := 0; o < tx.OutputCount - 1; o++ {

      var out block.DOutput

      out.OutputValue = int(b.Transactions[t].Outputs[o].OutputValue)
      out.ChallengeScriptLength = int(b.Transactions[t].Outputs[o].ChallengeScriptLength)
      out.ChallengeScript = b.Transactions[t].Outputs[o].ChallengeScript
      out.KeyType = b.Transactions[t].Outputs[o].KeyType
      out.NumAddresses = len(b.Transactions[t].Outputs[o].Addresses)

      var dAdds []block.DAddress

      for a := 0; a < out.NumAddresses; a++ {

        var add block.DAddress

        add.Address = b.Transactions[t].Outputs[o].Addresses[a].Address

        dAdds = append(dAdds, add)
      }
      out.Addresses = dAdds
      dOuts = append(dOuts, out)
    }

    tx.TransactionLockTime = int(b.Transactions[t].TransactionLockTime)

    tx.Inputs = dIns
    tx.Outputs = dOuts
    dTxs = append(dTxs, tx)

  }
  d.Transactions = dTxs
  return nil
}
