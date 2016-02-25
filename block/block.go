package block


//NullHash serves as default error hash when searching for RipeMD in output scripts
const NullHash string = "0000000000000000000000000000000000000000"

//MaxMultiSig holds the approximated maximum number of signatures in an output
const MaxMultiSig uint32 = 10

//Block holds fields for each new block
type Block struct {
  MagicNumber uint32
  BlockLength uint32
  Header Header
  BlockHash string
  TransactionCount uint64
  Transactions []Transaction
  HashBlock HashBlock
}

//HashBlock holds a compressed version of a block to hash to our blockchain
type HashBlock struct {
  FileEndpoint string
  CompressedBlockHash string
  BlockHash string
  PreviousCompressedBlockHash string
  PreviousBlockHash string
  TimeStamp uint32
  ByteOffset int
  LengthRead int
  ParsedBlockLength uint32
  RawBlockNumber int
}

//Header holds the interpreted Header fields read from the byte stream
type Header struct {
  FormatVersion uint32
  ByteFormatVersion []byte
  PreviousBlockHash string
  BytePreviousBlockHash []byte
  MerkleRoot string
  ByteMerkleRoot []byte
  TimeStamp uint32
  ByteTimeStamp []byte
  TargetValue uint32
  ByteTargetValue []byte
  Nonce uint32
  ByteNonce []byte
}

//Transaction holds the interpreted Transaction fields read from the byte stream
type Transaction struct {
  TransactionHash string
  TransactionVersionNumber uint32
  ByteTransactionVersionNumber []byte
  InputCount uint64
  ByteInputCount []byte
  Inputs []Input
  OutputCount uint64
  ByteOutputCount []byte
  Outputs []Output
  TransactionLockTime uint32
  ByteTransactionLockTime []byte
}

//Input holds the interpreted Input fields read from the byte stream
type Input struct {
  TransactionHash string
  ByteTransactionHash []byte
  TransactionIndex uint32
  ByteTransactionIndex []byte
  InputScriptLength uint64
  ByteInputScriptLength []byte
  InputScript string
  ByteInputScript []byte
  SequenceNumber uint32
  ByteSequenceNumber []byte
}

//Output holds the interpreted Output fields read from the byte stream
type Output struct {
  OutputValue uint64
  ByteOutputValue []byte
  ChallengeScriptLength uint64
  ByteChallengeScriptLength []byte
  ChallengeScript string
  ChallengeScriptBytes []byte
  KeyType string
  Addresses [MaxMultiSig]Address
}

//Address holds information for a given bitcoin address
type Address struct {
  Address string
  PublicKey string
  PublicKeyBytes []byte
  RipeMD160 string
  Transactions []Transaction
}

//ResponseBlock holds the blockchain.info response json when querying a block through
//blockchain.info API. It should be noted that compound names like Prevblock are
//represented as prev_block with underscores. However, underscores are to be avoided
//in Go.
type ResponseBlock struct {
  Hash string `json:"hash"`
  Ver int `json:"ver"`
  Prevblock string `json:"prev_block"`
  Mrklroot string `json:"mrkl_root"`
  Time int `json:"time"`
  Bits int `json:"bits"`
  Fee int `json:"fee"`
  Nonce int `json:"nonce"`
  Ntx int `json:"n_tx"`
  Size int `json:"size"`
  Blockindex int `json:"block_index"`
  Mainchain bool `json:"main_chain"`
  Height int `json:"height"`
  Tx []ResponseTransaction `json:"tx"`
}

//ResponseTransaction holds the blockchain.info response json for a given transaction in a block
type ResponseTransaction struct {
  Locktime int `json:"lock_time"`
  Ver int `json:"ver"`
  Size int `json:"size"`
  Inputs []ResponseInput `json:"inputs"`
  Time int `json:"time"`
  Txindex int `json:"tx_index"`
  Vinsz int `json:"vin_sz"`
  Hash string `json:"hash"`
  Voutsz int `json:"vout_sz"`
  Relayedby string `json:"relayed_by"`
  Out []ResponseOutput `json:"out"`
}

//ResponseInput holds the blockchain.info response json for a block's input.
type ResponseInput struct {
  Sequence int `json:"sequence"`
  Script string `json:"script"`
}

//ResponseOutput holds the blockchain.info response json for a block's output.
//Note that responsetype is returned with key "type" by blockchain.info, but this
//is a Go reserved word
type ResponseOutput struct {
  Spent bool `json:"spent"`
  Txindex int `json:"tx_index"`
  Responsetype int `json:"type"`
  Addr string `json:"addr"`
  Value int `json:"value"`
  N int `json:"n"`
  Script string `json:"script"`
}


//DBlock holds database structure for block
type DBlock struct {
  MagicNumber int
  BlockLength int
  BlockHash string
  FormatVersion int
  PreviousBlockHash string
  MerkleRoot string
  TimeStamp int
  TargetValue int
  Nonce int
  TransactionCount int
  Transactions []DTransaction
}

//DTransaction holds structure for input object in database
type DTransaction struct {
  TransactionHash string
  TransactionVersionNumber int
  InputCount int
  TransactionIndex int
  Inputs []DInput
  OutputCount int
  Outputs []DOutput
  TransactionLockTime int
  Time int
}

//DInput holds structure for input object in database
type DInput struct {
  TransactionHash string
  TransactionIndex int
  InputScriptLength int
  InputScript string
  SequenceNumber int
}

//DOutput holds structure for output object in database
type DOutput struct {
  OutputValue int
  ChallengeScriptLength int
  ChallengeScript string
  KeyType string
  TransactionIndex int
  NumAddresses int
  Addresses []DAddress
}

//DAddress holds structure for address object in database
type DAddress struct {
  Address string
  Transactions []DTransaction
}
