package main

import (
    "fmt"
    "log"
    "os"
    "flag"
    "strconv"
    "github.com/tgebhart/goparsebtc/blockchainbuilder"
    //"github.com/tgebhart/goparsebtc/blockvalidation"
    "github.com/tgebhart/goparsebtc/block"
)

//CHECKEVERY determines how many blocks go unchecked before we check the next block using blockchain.info
var CHECKEVERY = 200

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

  var blockCounter = 0

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
    Block := block.Block{}
    err = nil
    for err == nil {
      fmt.Println("++++++++++++++++++++++++++++++++++++ BLOCK ", blockCounter, " +++++++++++++++++++++++++++++++++++++++++++")
      err = blockchainbuilder.ParseIndividualBlock(&Block, file)
      if err != nil {
        log.Println("error in parseIndividualBlock ", err)
      }
      //if blockCounter % CHECKEVERY == 0 {
        //fmt.Println("?? Checking Block ??")
        //err = blockvalidation.BlockChainInfoValidation(&Block)
      //}
      //if err != nil {
        //log.Fatal("error in blockchain.info validation")
      //}
      blockCounter++
    }
    fmt.Println("Closing file...")
    defer file.Close()
  }
}
