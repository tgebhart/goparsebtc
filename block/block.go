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
  ByteHeader ByteHeader
  Transactions []Transaction
  ByteTransactions []ByteTransaction
}

//Header holds the interpreted Header fields read from the byte stream
type Header struct {
  FormatVersion uint32
  ByteFormatVersion []byte
  PreviousBlockHash string
  MerkleRoot string
  TimeStamp uint32
  TargetValue uint32
  Nonce uint32
}

//Transaction holds the interpreted Transaction fields read from the byte stream
type Transaction struct {
  TransactionHash string
  TransactionVersionNumber uint32
  InputCount uint64
  Inputs []Input
  OutputCount uint64
  Outputs []Output
  TransactionLockTime uint32
}

//Input holds the interpreted Input fields read from the byte stream
type Input struct {
  TransactionHash string
  TransactionIndex uint32
  InputScriptLength uint64
  InputScript string
  SequenceNumber uint32
}

//Output holds the interpreted Output fields read from the byte stream
type Output struct {
  OutputValue uint64
  ChallengeScriptLength uint64
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

//ByteHeader holds the fields relevant to the block header in byte arrays
type ByteHeader struct {
  FormatVersion []byte
  PreviousBlockHash []byte
  MerkleRoot []byte
  TimeStamp []byte
  TargetValue []byte
  Nonce []byte
}

//ByteTransaction holds fields relevant to each transaction in byte arrays
type ByteTransaction struct {
  TransactionLength []byte
  TransactionVersionNumber []byte
  InputCount []byte
  OutputCount []byte
  Inputs []ByteInput
  Outputs []ByteOutput
  TransactionLockTime []byte
}

//ByteInput holds fields related to the block input in byte arrays
type ByteInput struct {
  TransactionHash []byte
  TransactionIndex []byte
  InputScriptLength []byte
  InputScriptBytes []byte
  SequenceNumber []byte
}

//ByteOutput holds fields related to the block output in byte arrays
type ByteOutput struct {
  OutputValue []byte
  ChallengeScriptLength []byte
  ChallengeScript []byte
}
