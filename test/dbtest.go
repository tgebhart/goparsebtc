package main

import (
    "github.com/tgebhart/neoism"
    "github.com/tgebhart/goparsebtc/blockchainreader"
    //"github.com/tgebhart/goparsebtc/blockchainbuilder"
    "github.com/tgebhart/goparsebtc/block"
    "github.com/tgebhart/goparsebtc/blockvalidation"
    "fmt"
    "os"
)

func main() {
  testBlockUpload("../part2.csv")
}


func testNodeUpload() {

  db, err := neoism.Connect("http://neo4j:stella6332@localhost:7474/db/data")
  if err != nil{
    fmt.Println("connect : ", err)
  }
  fmt.Println(db)
  n, err := db.CreateNode(neoism.Props{"name" : "Aragorn II Elessar"})
  if err != nil {
    fmt.Println("create node: ", err)
  }
  n.AddLabel("King")
  fmt.Println(n)

}

func testQueryUpload() {

  db, err := neoism.Connect("http://neo4j:stella6332@localhost:7474/db/data")
  if err != nil {
    fmt.Println("connect : ", err)
  }

  res := []struct{ N neoism.Node }{}

  cq := neoism.CypherQuery{
    Statement: "MERGE (n:Person {name: {name}}) RETURN n",
    Parameters: neoism.Props{"name" : "Aragorn II Elessar"},
    Result: &res,
  }

  db.Cypher(&cq)

 fmt.Println(res)
}

func testBlockUpload(reference string) {

  db, err := neoism.Connect("http://neo4j:stella6332@localhost:7474/db/data")
  if err != nil {
    fmt.Println("connect : " , err)
  }

  chain := blockchainreader.ReadChain{}
  err = blockchainreader.ReadReferenceFile(&chain, reference)
  if err != nil {
    fmt.Println("ReadReferenceFile: ", err)
  }

  b := block.Block{}
  d := block.DBlock{}
  file, err := os.Open("/Users/tgebhart/Library/Application Support/Bitcoin/blocks/" + chain.ReadBlocks[0].FileEndpoint)
  if err != nil {
      fmt.Println("open file: ", err)
  }

  fmt.Println(chain.ReadBlocks[1])
  err = blockchainreader.ScanBlock(&b, chain.ReadBlocks[0].ByteOffset, chain.ReadBlocks[0].BlockLength, file)
  if err != nil {
    err = blockvalidation.BridgeWithBlockchainInfo(&d, chain.ReadBlocks[0].BlockHash)
    if err != nil {
      fmt.Println("Bridge: ", err)
    }
  }

  if b.BlockHash != "" {
    err = blockchainreader.MapBlockToDBlock(&b, &d)
    if err != nil {
      fmt.Println("MapBlockToDBlock: ", err)
    }
  }

  res := []struct{ N neoism.Node }{}
  tres := []struct{ N neoism.Node }{}
  ires := []struct{ N neoism.Node }{}
  ores := []struct{ N neoism.Node }{}
  ares := []struct{ N neoism.Node }{}

  fmt.Println(b)

  cq := neoism.CypherQuery {
    Statement: "MERGE (n:Block {hash: {hash}}) RETURN n",
    Parameters: neoism.Props{"hash" : d.BlockHash},
    Result: &res,
  }

  db.Cypher(&cq)

  tq := neoism.CypherQuery {
    Statement: "MERGE (n:TimeStamp {time_stamp: {time_stamp}}) RETURN n",
    Parameters: neoism.Props{"time_stamp" : d.TimeStamp},
    Result: &tres,
  }

  db.Cypher(&tq)
  blocknode := res[0].N
  timenode := tres[0].N
  blocknode.Db = db
  timenode.Db = db

  err = blocknode.SetProperties(neoism.Props{"hash" : d.BlockHash, "merkle_root" : d.MerkleRoot,
  "block_length" : d.BlockLength, "format_version" : d.FormatVersion,
  "target_value" : d.TargetValue, "nonce" : d.Nonce, "tx_count" : d.TransactionCount})
  if err != nil {
    fmt.Println("SetProperties", err)
  }

  blocknode.Relate("mined on", timenode.Id(), neoism.Props{})

  for i := 0; i < d.TransactionCount; i++ {

    trq := neoism.CypherQuery {
      Statement: "MERGE (n:Transaction {tx_hash: {tx_hash}}) RETURN n",
      Parameters: neoism.Props{"tx_hash" : d.Transactions[i].TransactionHash},
      Result: &tres,
    }
    db.Cypher(&trq)

    tran := tres[0].N
    tran.Db = db

    err = tran.SetProperties(neoism.Props{"tx_hash" : d.Transactions[i].TransactionHash,
    "tx_version" : d.Transactions[i].TransactionVersionNumber, "input_count" : d.Transactions[i].InputCount,
    "tx_index" : d.Transactions[i].TransactionIndex, "output_count" : d.Transactions[i].OutputCount,
    "lock_time" : d.Transactions[i].TransactionLockTime})
    if err != nil {
      fmt.Println("transaction SetProperties", err)
    }

    tran.Relate("mined on", timenode.Id(), neoism.Props{})
    tran.Relate("in", blocknode.Id(), neoism.Props{})
    blocknode.Relate("contains", tran.Id(), neoism.Props{})

    for in := 0; in < d.Transactions[i].InputCount; in++ {

      inq := neoism.CypherQuery {
        Statement: "MERGE (n:Input {tx_index : {tx_index}}) RETURN n",
        Parameters: neoism.Props{"tx_index" : d.Transactions[i].Inputs[in].TransactionIndex},
        Result: &ires,
      }
      db.Cypher(&inq)

      input := ires[0].N
      input.Db = db

      err = input.SetProperties(neoism.Props{"tx_hash" : d.Transactions[i].Inputs[in].TransactionHash,
      "tx_index" : d.Transactions[i].Inputs[in].TransactionIndex, "input_script_length" : d.Transactions[i].Inputs[in].InputScriptLength,
      "input_script" : d.Transactions[i].Inputs[in].InputScript, "sequence_number" : d.Transactions[i].Inputs[in].SequenceNumber})
      if err != nil {
        fmt.Println("Error in input SetProperties", err)
      }

      tran.Relate("contains", input.Id(), neoism.Props{})
      input.Relate("in", tran.Id(), neoism.Props{})

    }

    for o := 0; o < d.Transactions[i].OutputCount; o++ {

      onq := neoism.CypherQuery {
        Statement: "MERGE (n:Output {tx_index : {tx_index}}) RETURN n",
        Parameters: neoism.Props{"tx_index" : d.Transactions[i].Outputs[o].TransactionIndex},
        Result: &ores,
      }
      db.Cypher(&onq)

      output := ores[0].N
      output.Db = db

      err = output.SetProperties(neoism.Props{"output_value" : d.Transactions[i].Outputs[o].OutputValue,
      "challenge_script_length" : d.Transactions[i].Outputs[o].ChallengeScriptLength,
      "challenge_script" : d.Transactions[i].Outputs[o].ChallengeScript, "key_type" : d.Transactions[i].Outputs[o].KeyType,
      "tx_index" : d.Transactions[i].Outputs[o].TransactionIndex, "num_address" : d.Transactions[i].Outputs[o].NumAddresses})
      if err != nil {
        fmt.Println("output SetProperties" , err)
      }

      tran.Relate("contains", output.Id(), neoism.Props{})
      output.Relate("in", tran.Id(), neoism.Props{})

      for a := 0; a < d.Transactions[i].Outputs[o].NumAddresses; a++ {

        adq := neoism.CypherQuery {
          Statement: "MERGE (n:Address {address : {address}}) RETURN n",
          Parameters: neoism.Props{"address" : d.Transactions[i].Outputs[o].Addresses[a].Address},
          Result: &ares,
        }
        db.Cypher(&adq)

        address := ares[0].N
        address.Db = db

        output.Relate("to address", address.Id(), neoism.Props{})
        address.Relate("in", output.Id(), neoism.Props{})

      }
    }

  }

}
