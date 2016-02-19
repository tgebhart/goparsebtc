package main

import (
    "fmt"
    "log"
    "os"
    "io"
    "flag"
    "strconv"
    "github.com/tgebhart/goparsebtc/blockchainbuilder"
    "github.com/tgebhart/goparsebtc/block"
    "github.com/tgebhart/goparsebtc/filefunctions"
    "github.com/tgebhart/goparsebtc/blockvalidation"
)

//CHECKEVERY determines how many blocks go unchecked before we check the next block using blockchain.info
var CHECKEVERY = 200

/******************************MAIN********************************************/

func main() {
  path := "/Users/tgebhart/Library/Application Support/Bitcoin/blocks/"
  flag.Parse()
    s := flag.Arg(0)
    f := flag.Arg(1)
    dumpLocation := flag.Arg(2)
    start, err := strconv.Atoi(s)
    if err != nil {
      fmt.Println(err)
      os.Exit(2)
    }
    finish, err := strconv.Atoi(f)
    if err != nil {
      fmt.Println(err)
      os.Exit(2)
    }

  var blockCounter = 0
  var bytesRead = 0

  chain :=  blockchainbuilder.NewBlockchain()
  var key string

  for j := start; j <= finish; j++ {
    path = "/Users/tgebhart/Library/Application Support/Bitcoin/blocks/"
    e := strconv.Itoa(j)
    tempString := e
    for k := len(e); k < 5; k++ {
      tempString = "0" + tempString
    }
    pathEndpoint := "blk" + tempString + ".dat"
    path = path + pathEndpoint

    file, err := os.Open(path)
    if err != nil {
        log.Fatal("Error while opening file", err)
    }
    fmt.Printf("%s opened\n", path)
    defer file.Close()
    err = nil
    for err == nil {
      fmt.Println("++++++++++++++++++++++++++++++++++++ BLOCK ", blockCounter, " +++++++++++++++++++++++++++++++++++++++++++")
      Block := block.Block{}
      err = chain.ParseIndividualBlockSuppressOutput(&Block, file)
      if err != nil {
        if err == io.EOF { //reached end of file
          fmt.Println("EOF, opening next file")
          break
        }
        if err == blockvalidation.ErrMultiSig {
          fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ \n MultiSigErr \n @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
          err = nil
          //log.Fatal(err)
        }
        if err == blockchainbuilder.ErrBadFormatVersion {
          fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ \n Found bad format version \n @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
          err = nil
          chain.PrepareSkipBlock(&Block, pathEndpoint, blockCounter, bytesRead, file)
        }
        if err == blockchainbuilder.ErrBadMagic {
          log.Fatal(err)
        }
      }
      fmt.Println(blockvalidation.ReverseEndian(Block.BlockHash))
      if Block.HashBlock.CompressedBlockHash != "" {
        key = Block.HashBlock.CompressedBlockHash
      }
      Block.HashBlock.FileEndpoint = pathEndpoint
      Block.HashBlock.RawBlockNumber = blockCounter
      Block.HashBlock.ByteOffset = bytesRead

      //Add HashBlock to Blockchain hashmap
      if Block.HashBlock.CompressedBlockHash == "" {
        fmt.Println("ZERO HASH BLOCK AT", Block.HashBlock.PreviousCompressedBlockHash)
      }
      chain.BlockMap[Block.HashBlock.CompressedBlockHash] = Block.HashBlock

      //if blockCounter % CHECKEVERY == 0 {
        //fmt.Println("?? Checking Block ??")
        //err = blockvalidation.BlockChainInfoValidation(&Block)
      //}
      //if err != nil {
        //log.Fatal("error in blockchain.info validation")
      //}
      bytesRead += filefunctions.GetByteCount()
      blockCounter++
    }
    fmt.Println("Bytes read: ", bytesRead)
    fmt.Println("Closing file...", path)
    file.Close()
  }

  fmt.Println("About to call main write")
  err = blockchainbuilder.WriteMainChainToFile(chain, key, dumpLocation)
  if err != nil {
    log.Fatal(err)
  }
}
