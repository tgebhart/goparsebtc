package block

//Block holds fields for each new block
type Block struct {
  BMagicNumber uint32
  BBlockLength uint32
  BFormatVersion uint32
  BPreviousBlockHash string
  BMerkleRoot string
  BTimeStamp uint32
  BTargetValue uint32
  BNonce uint32
  BBlockHash string
  BTransactionCount uint64
  BTransactionVersionNumber uint32
  BInputCount uint64
  BTransactionHash string
  BTransactionIndex uint32
  BInputScriptLength uint64
  BInputScriptBytes []uint8
  BSequenceNumber uint32
  BOutputCount uint64
  BOutputValue uint64
  BChallengeScriptLength uint64
  BChallengeScriptBytes []uint8
  BTransactionLockTime uint32
  Header Header
  Transaction Transaction
}

//Header holds the fields relevant to the block header in byte arrays
type Header struct {
  HFormatVersion []byte
  HPreviousBlockHash []byte
  HMerkleRoot []byte
  HTimeStamp []byte
  HTargetValue []byte
  HNonce []byte
}

//Transaction holds fields relevant to each transaction in the block
type Transaction struct {
  TTransactionLength []byte
  TTransactionVersionNumber []byte
  TInputCount []byte
  TOutputCount []byte
  TTransactionLockTime []byte
  Input []Input
  Output []Output
}

//Input holds fields related to the block input
type Input struct {
  ITransactionHash []byte
  ITransactionIndex []byte
  IInputScriptLength []byte
  IInputScriptBytes []byte
  ISequenceNumber []byte
}

//Output holds fields related to the block output
type Output struct {
  OOutputValue []byte
  OChallengeScriptLength []byte
  OChallengeScriptBytes []byte
}
