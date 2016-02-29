package blockvalidation

import (
  "github.com/tgebhart/goparsebtc/block"
  "github.com/tgebhart/goparsebtc/btchashing"
  "net/http"
  //"bytes"
  "fmt"
  "io/ioutil"
  "encoding/json"
  "errors"
  //"log"
  "time"
)

//MaxReasonableTransactionIndex holds the upper bound for transaction index
const MaxReasonableTransactionIndex uint32 = 10000

//SatoshiConst is Satoshi's transaction index for the genesis block
const SatoshiConst uint32 = 4294967295

//ErrMultiSig is thrown when we cannot read multisig output script
var ErrMultiSig = errors.New("unable to parse multisig")
//ErrReplacementKey is thrown when blockchain.info validation cannot find previous block key
var ErrReplacementKey = errors.New("could not find replacement key")
//ErrZeroOutputScript is thrown when zero length output script is present
var ErrZeroOutputScript = errors.New("block may have zero length outputs script")

//public key types
const (
  UncompressedPublicKey = "UNCOMPRESSED_PUBLIC_KEY"
  StealthKey = "STEALTH"
  CompressedPublicKey = "COMPRESSED_PUBLIC_KEY"
  TruncatedCompressedKey = "TRUNCATED_COMPRESSED_KEY"
  ScriptHashKey = "SCRIPT_HASH"
  RipeMD160Key = "RIPEMD160"
  MultiSigKey = "MULTISIG"
  NullKey = block.NullHash
)


//challenge script op codes
const (
	OP0 			=  0x00
	OPPUSHDATA1 	=  0x4c
	OPPUSHDATA2 	=  0x4d
	OPPUSHDATA4 	=  0x4e
	OP1NEGATE 		=  0x4f
	OPRESERVED 	=  0x50
	OP1 			=  0x51
	OP2 			=  0x52
	OP3 			=  0x53
	OP4 			=  0x54
	OP5 			=  0x55
	OP6 			=  0x56
	OP7 			=  0x57
	OP8 			=  0x58
	OP9 			=  0x59
	OP10 			=  0x5a
	OP11 			=  0x5b
	OP12 			=  0x5c
	OP13 			=  0x5d
	OP14 			=  0x5e
	OP15 			=  0x5f
	OP16 			=  0x60
	OPNOP 			=  0x61
	OPVER 			=  0x62
	OPIF 			=  0x63
	OPNOTIF 		=  0x64
	OPVERIF 		=  0x65
	OPVERNOTIF 	=  0x66
	OPELSE 		=  0x67
	OPENDIF 		=  0x68
	OPVERIFY 		=  0x69
	OPRETURN 		=  0x6a
	OPTOALTSTACK 	=  0x6b
	OPFROMALTSTACK =  0x6c
	OP2DROP 		=  0x6d
	OP2DUP 		=  0x6e
	OP3DUP 		=  0x6f
	OP2OVER 		=  0x70
	OP2ROT 		=  0x71
	OP2SWAP 		=  0x72
	OPIFDUP 		=  0x73
	OPDEPTH 		=  0x74
	OPDROP 		=  0x75
	OPDUP 			=  0x76
	OPNIP 			=  0x77
	OPOVER 		=  0x78
	OPPICK 		=  0x79
	OPROLL 		=  0x7a
	OPROT 			=  0x7b
	OPSWAP 		=  0x7c
	OPTUCK 		=  0x7d
	OPCAT 			=  0x7e	// Currently disabled
	OPSUBSTR 		=  0x7f	// Currently disabled
	OPLEFT 		=  0x80	// Currently disabled
	OPRIGHT 		=  0x81	// Currently disabled
	OPSIZE 		=  0x82	// Currently disabled
	OPINVERT 		=  0x83	// Currently disabled
	OPAND 			=  0x84	// Currently disabled
	OPOR 			=  0x85	// Currently disabled
	OPXOR 			=  0x86	// Currently disabled
	OPEQUAL 		=  0x87
	OPEQUALVERIFY 	=  0x88
	OPRESERVED1 	=  0x89
	OPRESERVED2 	=  0x8a
	OP1ADD 		=  0x8b
	OP1SUB 		=  0x8c
	OP2MUL 		=  0x8d	// Currently disabled
	OP2DIV 		=  0x8e	// Currently disabled
	OPNEGATE 		=  0x8f
	OPABS 			=  0x90
	OPNOT 			=  0x91
	OP0NOTEQUAL 	=  0x92
	OPADD 			=  0x93
	OPSUB 			=  0x94
	OPMUL 			=  0x95	// Currently disabled
	OPDIV 			=  0x96	// Currently disabled
	OPMOD 			=  0x97	// Currently disabled
	OPLSHIFT 		=  0x98	// Currently disabled
	OPRSHIFT 		=  0x99	// Currently disabled
	OPBOOLAND 		=  0x9a
	OPBOOLOR 		=  0x9b
	OPNUMEQUAL 	=  0x9c
	OPNUMEQUALVERIFY =  0x9d
	OPNUMNOTEQUAL 	=  0x9e
	OPLESSTHAN 	=  0x9f
	OPGREATERTHAN 	=  0xa0
	OPLESSTHANOREQUAL =  0xa1
	OPGREATERTHANOREQUAL =  0xa2
	OPMIN 			=  0xa3
	OPMAX 			=  0xa4
	OPWITHIN 		=  0xa5
	OPRIPEMD160 	=  0xa6
	OPSHA1 		=  0xa7
	OPSHA256		=  0xa8
	OPHASH160 		=  0xa9
	OPHASH256 		=  0xaa
	OPCODESEPARATOR =  0xab
	OPCHECKSIG 	=  0xac
	OPCHECKSIGVERIFY =  0xad
	OPCHECKMULTISIG =  0xae
	OPCHECKMULTISIGVERIFY = 0xaf
	OPNOP1 		=  0xb0
	OPNOP2 		=  0xb1
	OPNOP3 		=  0xb2
	OPNOP4 		=  0xb3
	OPNOP5 		=  0xb4
	OPNOP6 		=  0xb5
	OPNOP7 		=  0xb6
	OPNOP8 		=  0xb7
	OPNOP9 		=  0xb8
	OPNOP10 		=  0xb9
	OPSMALLINTEGER =  0xfa
	OPPUBKEYS 		=  0xfb
	OPPUBKEYHASH 	=  0xfd
	OPPUBKEY 		=  0xfe
	OPINVALIDOPCODE =  0xff
)


//BLOCKCHAININFOENDPOINT is the API endpoint for json information from blockchain.info
var BLOCKCHAININFOENDPOINT = "https://blockchain.info/rawblock/"

//APICode is passed with blockchain.info call for extended API requests
var APICode = "865c2783-0c23-45b7-a808-29dbca5435df"

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
  if formatVersion == 1 || formatVersion == 2 || formatVersion == 3 || formatVersion == 4 {  //format version should still be 1 for now
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

//ConvertUnixEpochToDate converts the integer timestamp to a time.Time object to output
func ConvertUnixEpochToDate(timeStamp uint32) (time.Time) {
  stamp64 := int64(timeStamp)
  ret := time.Unix(stamp64, 0)
  return ret
}

//ValidateTransactionVersion checks transaction version. Should be equal to 1 currently
func ValidateTransactionVersion(transactionVersion uint32) (bool) {
  if transactionVersion == 1 {  //current transaction version
    return true
  }
  return false
}

//ValidateTransactionIndex checks the transaction index to make sure number is reasonable
func ValidateTransactionIndex(transactionIndex uint32) (bool) {
  if transactionIndex < MaxReasonableTransactionIndex || transactionIndex == SatoshiConst {
    return true
  }
  return false
}

//ValidateSequenceNumber checks to make sure sequence number is below the maximum integer value
func ValidateSequenceNumber(b []byte) (bool) {
  if b[0] == 255 && b[1] == 255 && b[2] == 255 && b[3] != 255 {  //current largest sequence number
    return false
  }
  return true
}

//ValidateOutputValue checks to see if the parsed output value is withing a reasonable range. If not, could be a read error.
func ValidateOutputValue(outputValue uint64) (bool) {
  if outputValue < 1501439850948224747 { //stupid-ass error number returned while debugging
    return true
  }
  return false
}

//ValidateTransactionLockTime checks transaction lock time is equal to 0
func ValidateTransactionLockTime(transactionLockTime uint32) (bool) {
  if transactionLockTime <= SatoshiConst && transactionLockTime != 16777216 {
    return true
  }
  return false
}

// ParseOutputScript iterates an output script and validates interior op_codes. Returns keytype
func ParseOutputScript(output *block.Output) (string, error) {
  var multiSigFormat int
  var keytype string

  if output.ChallengeScript != "" {
    lastInstruction := output.ChallengeScriptBytes[output.ChallengeScriptLength - 1]
    if output.ChallengeScriptLength == 67 && output.ChallengeScriptBytes[0] == 65 && output.ChallengeScriptBytes[66] == OPCHECKSIG {
      output.Addresses[0].PublicKeyBytes = output.ChallengeScriptBytes[1:output.ChallengeScriptLength-1]
      keytype = UncompressedPublicKey
    }
    if output.ChallengeScriptLength == 40 && output.ChallengeScriptBytes[0] == OPRETURN {
      output.Addresses[0].PublicKeyBytes = output.ChallengeScriptBytes[1:]
      output.KeyType = StealthKey
    } else if output.ChallengeScriptLength == 66 && output.ChallengeScriptBytes[65] == OPCHECKSIG {
      output.Addresses[0].PublicKeyBytes = output.ChallengeScriptBytes[:]
      keytype = UncompressedPublicKey
    } else if output.ChallengeScriptLength == 35 && output.ChallengeScriptBytes[34] == OPCHECKSIG {
      output.Addresses[0].PublicKeyBytes = output.ChallengeScriptBytes[1:]
      keytype = CompressedPublicKey
    } else if output.ChallengeScriptLength == 33 && output.ChallengeScriptBytes[0] == 0x20 {
      output.Addresses[0].PublicKeyBytes = output.ChallengeScriptBytes[1:]
      keytype = TruncatedCompressedKey
    } else if output.ChallengeScriptLength == 23 && output.ChallengeScriptBytes[0] == OPHASH160 && output.ChallengeScriptBytes[1] == 20 && output.ChallengeScriptBytes[22] == OPEQUAL {
      output.Addresses[0].PublicKeyBytes = output.ChallengeScriptBytes[2:output.ChallengeScriptLength-1]
      keytype = ScriptHashKey
    } else if output.ChallengeScriptLength >= 25 && output.ChallengeScriptBytes[0] == OPDUP && output.ChallengeScriptBytes[1] == OPHASH160 && output.ChallengeScriptBytes[2] == 20 {
      output.Addresses[0].PublicKeyBytes = output.ChallengeScriptBytes[3:23]
      keytype = RipeMD160Key
    } else if output.ChallengeScriptLength == 5 && output.ChallengeScriptBytes[0] == OPDUP && output.ChallengeScriptBytes[1] == OPHASH160 && output.ChallengeScriptBytes[2] == OP0 && output.ChallengeScriptBytes[3] == OPEQUALVERIFY && output.ChallengeScriptBytes[4] == OPCHECKSIG {
      fmt.Println("WARNING : Encountered unusual but expected output script. ")
      keytype = NullKey
    } else if lastInstruction == OPCHECKMULTISIG && output.ChallengeScriptLength > 25 { //effin multisig
      scanIndex := 0
      scanbegin := output.ChallengeScriptBytes[scanIndex]
      scanend := output.ChallengeScriptBytes[output.ChallengeScriptLength - 2]
      expectedPrefix := false
      expectedSuffix := false
      switch scanbegin {
      case OP0:
        expectedPrefix = true
        break
      case OP1:
        expectedPrefix = true
        break
      case OP2:
        expectedPrefix = true
        break
      case OP3:
        expectedPrefix = true
        break
      case OP4:
        expectedPrefix = true
        break
      case OP5:
        expectedPrefix = true
        break
      default:
        //unexpected
        break
      }

      switch scanend {
      case OP1:
        expectedSuffix = true
        break
      case OP2:
        expectedSuffix = true
        break
      case OP3:
        expectedSuffix = true
        break
      case OP4:
        expectedSuffix = true
        break
      case OP5:
        expectedSuffix = true
        break
      default:
        //unexpected
        break
      }

      if expectedPrefix && expectedSuffix {
        scanIndex++
        scanbegin = output.ChallengeScriptBytes[scanIndex]
        var keyIndex uint8
        for keyIndex < 5 && scanbegin < scanend {
          if scanbegin == 0x21 {
            output.KeyType = MultiSigKey
            scanIndex++
            scanbegin = output.ChallengeScriptBytes[scanIndex]
            output.Addresses[keyIndex].PublicKeyBytes = output.ChallengeScriptBytes[scanIndex:]
            scanbegin += 0x21
            bitMask := 1<<keyIndex
            multiSigFormat|=bitMask
            keyIndex++
          } else if scanbegin == 0x41 {
            output.KeyType = MultiSigKey
            scanIndex++
            scanbegin = output.ChallengeScriptBytes[scanIndex]
            output.Addresses[keyIndex].PublicKeyBytes = output.ChallengeScriptBytes[scanIndex:]
            scanbegin += 0x41
            keyIndex++
          } else {
            break
          }
        }
      }
      if output.Addresses[0].PublicKeyBytes == nil {
        fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&& error multisig &&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")
        return "", ErrMultiSig
      }
    } else { //scan for pattern OP_DUP, OP_HASH160, 0x14, 20 bytes, 0x88, 0xac
      if output.ChallengeScriptLength > 25 {
        endIndex := output.ChallengeScriptLength - 25
        for i := 0; i < int(endIndex); i++ {
          scan := output.ChallengeScriptBytes[i:]
          if scan[0] == OPDUP && scan[1] == OPHASH160 && scan[2] == 20 && scan[23] == OPEQUALVERIFY && scan[24] == OPCHECKSIG {
            output.Addresses[0].PublicKeyBytes = scan[3:]
            output.KeyType = RipeMD160Key
            //fmt.Println("WARNING: Unusual output script in scan")

          }
        }
      }
    }
    //if output.Addresses[0].PublicKey == "" {
    //  fmt.Println("FAILED TO LOCATE PUBLIC KEY")
    //}
  } else {
    output.KeyType = NullKey
    return output.KeyType ,ErrZeroOutputScript
  }

  if output.Addresses[0].PublicKey == "" {
    if output.ChallengeScriptLength == 0 {
      output.Addresses[0].PublicKey = NullKey
    } else {
      output.Addresses[0].PublicKey = NullKey
    }
    output.KeyType = RipeMD160Key
    //fmt.Println("WARNING : Failed to decode public key in output script ")
  }

  switch keytype {
  case RipeMD160Key:
    btchashing.BitcoinRipeMD160ToAddress(output.Addresses[0].PublicKeyBytes, &output.Addresses[0])
    output.KeyType = keytype
    return output.KeyType, nil
  case ScriptHashKey:
    btchashing.BitcoinRipeMD160ToAddress(output.Addresses[0].PublicKeyBytes, &output.Addresses[0])
    output.KeyType = keytype
    return output.KeyType, nil
  case StealthKey:
    btchashing.BitcoinRipeMD160ToAddress(output.Addresses[0].PublicKeyBytes, &output.Addresses[0])
    output.KeyType = keytype
  case UncompressedPublicKey:
    btchashing.BitcoinPublicKeyToAddress(output.Addresses[0].PublicKeyBytes, &output.Addresses[0])
    output.KeyType = keytype
    return output.KeyType, nil
  case CompressedPublicKey:
    btchashing.BitcoinCompressedPublicKeyToAddress(output.Addresses[0].PublicKeyBytes, &output.Addresses[0])
    output.KeyType = keytype
    return output.KeyType, nil
  case TruncatedCompressedKey:
    tempkey := make([]byte, 1)
    tempkey[0] = 0x2
    key := append(tempkey[:], output.Addresses[0].PublicKey[:] ...)
    btchashing.BitcoinCompressedPublicKeyToAddress(key, &output.Addresses[0])
    output.KeyType = keytype
    return output.KeyType, nil
  case MultiSigKey:
    var i uint32
    for i = 0; i < block.MaxMultiSig; i++ {
      key := output.Addresses[i].PublicKey
      if key == "" {
        break
      }
      mask := 1<<i
      if multiSigFormat & mask != 0 {
        btchashing.BitcoinCompressedPublicKeyToAddress([]byte(output.Addresses[i].PublicKey), &output.Addresses[i])
      } else {
         btchashing.BitcoinPublicKeyToAddress([]byte(output.Addresses[i].PublicKey), &output.Addresses[i])
      }
    }
    output.KeyType = keytype
  }
  return keytype, nil
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
  ResponseBlock := block.ResponseBlock{}
  blockHash := ReverseEndian(Block.BlockHash)
  fmt.Println("block hash", blockHash)
  resp, err := http.Get(BLOCKCHAININFOENDPOINT + blockHash)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  json.Unmarshal(body, &ResponseBlock)

  if blockHash == ResponseBlock.Hash {
    fmt.Println("Height: ", ResponseBlock.Height)
    return nil
  }
  return errors.New("Hashes do not match")
}

//GetReplacementKey returns blockchain.info's reported previous block hash given the parameter block hash
func GetReplacementKey(hash string) (string, string, error) {
  //narcolepsy()
  blockHash := hash
  fmt.Println("Looking for previous block hash ...")
  resp, err := http.Get(BLOCKCHAININFOENDPOINT + blockHash + "?" + APICode)
  if err != nil {
    return "", "", err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(err.Error())
  }
  tx, err := getTxs(body)
  fmt.Println("prevblock" , tx.Prevblock)
  if tx.Prevblock != "" {
    return tx.Prevblock, tx.Hash, nil
  }

  return "", "", ErrReplacementKey
}

func getTxs(body []byte) (*block.ResponseBlock, error) {
  var r = new(block.ResponseBlock)
  err := json.Unmarshal(body, &r)
  if err != nil {
    fmt.Println("couldn't unmarshal", body, err)
  }
  return r, nil
}



//BridgeWithBlockchainInfo bridges data that could not be parsed with block from blockchain.info
func BridgeWithBlockchainInfo(dBlock *block.DBlock, hash string) (error) {
  var r = new(block.ResponseBlock)
  resp, err := http.Get(BLOCKCHAININFOENDPOINT + hash + "?" + APICode)
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  json.Unmarshal(body, &r)

  fmt.Println("Hashes: ", r.Hash , hash)

  if hash == r.Hash {
    mapResponseToBlock(r, dBlock)
    return nil
  }
  return errors.New("Hashes do not match")
}


func mapResponseToBlock(r *block.ResponseBlock, d *block.DBlock) {

  d.BlockLength = r.Size
  d.BlockHash = r.Hash
  d.TransactionCount = r.Ntx
  d.FormatVersion = r.Ver
  d.PreviousBlockHash = r.Prevblock
  d.MerkleRoot = r.Mrklroot
  d.TimeStamp = r.Time
  d.TargetValue = 0
  d.Nonce = r.Nonce

  d.Transactions = make([]block.DTransaction , d.TransactionCount)

  for t := 0; t < d.TransactionCount - 1; t++ {
    d.Transactions[t].TransactionHash = r.Tx[t].Hash
    d.Transactions[t].TransactionVersionNumber = r.Tx[t].Ver
    d.Transactions[t].InputCount = r.Tx[t].Vinsz
    d.Transactions[t].TransactionIndex = r.Tx[t].Txindex
    d.Transactions[t].Time = r.Tx[t].Time

    d.Transactions[t].Inputs = make([]block.DInput , d.Transactions[t].InputCount)

    for i := 0; i < d.Transactions[t].InputCount - 1; i++ {
      d.Transactions[t].Inputs[i].InputScript = r.Tx[t].Inputs[i].Script
      d.Transactions[t].Inputs[i].TransactionIndex = d.Transactions[t].TransactionIndex
      d.Transactions[t].Inputs[i].TransactionHash = ""
      d.Transactions[t].Inputs[i].InputScriptLength = len(d.Transactions[t].Inputs[i].InputScript)
      d.Transactions[t].Inputs[i].SequenceNumber = r.Tx[t].Inputs[i].Sequence
    }

    d.Transactions[t].OutputCount = r.Tx[t].Voutsz
    d.Transactions[t].Outputs = make([]block.DOutput , d.Transactions[t].OutputCount)

    for o := 0; o < d.Transactions[t].OutputCount - 1; o++ {
      d.Transactions[t].Outputs[o].OutputValue = r.Tx[t].Out[o].Value
      d.Transactions[t].Outputs[o].TransactionIndex = r.Tx[t].Out[o].Txindex
      d.Transactions[t].Outputs[o].ChallengeScript = r.Tx[t].Out[o].Script
      d.Transactions[t].Outputs[o].ChallengeScriptLength = len(d.Transactions[t].Outputs[o].ChallengeScript)
      d.Transactions[t].Outputs[o].KeyType = ""
      d.Transactions[t].Outputs[o].NumAddresses = 1

      d.Transactions[t].Outputs[o].Addresses = make([]block.DAddress , 1)
      d.Transactions[t].Outputs[o].Addresses[0].Address = r.Tx[t].Out[o].Addr
    }

    d.Transactions[t].TransactionLockTime = r.Tx[t].Locktime

  }

}
