package sui_types

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fardream/go-bcs/bcs"
	"github.com/thorli9527/sui-wallet-sdk/crypto"
	"github.com/thorli9527/sui-wallet-sdk/lib"
	"golang.org/x/crypto/blake2b"
	"hash"
	"reflect"
)

type DefaultHash struct {
	hash.Hash
}

func NewDefaultHash() DefaultHash {
	digest, err := blake2b.New256([]byte{})
	if err != nil {
		panic(err)
	}
	return DefaultHash{
		digest,
	}
}

type Signable interface {
	Write(digest hash.Hash)
}

type BcsSignableKind interface {
	TransactionData | Object
	//Not Implement | Committee | CheckpointSummary | CheckpointContents | TransactionEffects| TransactionEvents
	//| SenderSignedData  | Accumulator | Foo
}
type BcsSignable[K BcsSignableKind] struct {
	Data K
}

func (s BcsSignable[K]) Write(digest hash.Hash) {
	rv := reflect.ValueOf(s.Data)
	name := rv.Type().Name()
	_, err := digest.Write([]byte(fmt.Sprintf("%s::", name)))
	if err != nil {
		panic("Hasher should not fail")
	}
	bcsData, err := bcs.Marshal(s.Data)
	if err != nil {
		panic("Message serialization should not fail")
	}
	_, err = digest.Write(bcsData)
	if err != nil {
		panic("Hasher should not fail")
	}
}

func internalHash(digest hash.Hash, signable Signable) []byte {
	signable.Write(digest)
	return digest.Sum([]byte{})
}

func UseDefaultHash(signable Signable) []byte {
	return internalHash(NewDefaultHash(), signable)
}

type Signature struct {
	*Ed25519SuiSignature
	*Secp256k1SuiSignature
	*Secp256r1SuiSignature
}

func (s Signature) MarshalJSON() ([]byte, error) {
	switch {
	case s.Ed25519SuiSignature != nil:
		return json.Marshal(s.Ed25519SuiSignature.Signature[:])
	case s.Secp256k1SuiSignature != nil:
		return json.Marshal(s.Secp256k1SuiSignature.Signature[:])
	case s.Secp256r1SuiSignature != nil:
		return json.Marshal(s.Secp256r1SuiSignature.Signature[:])
	default:
		return nil, errors.New("nil signature")
	}

}

func (s *Signature) UnmarshalJSON(data []byte) error {
	var signature []byte
	err := json.Unmarshal(data, &signature)
	if err != nil {
		return err
	}
	switch signature[0] {
	case 0:
		if len(signature) != ed25519.PublicKeySize+ed25519.SignatureSize+1 {
			return errors.New("invalid ed25519 signature")
		}
		var signatureBytes [ed25519.PublicKeySize + ed25519.SignatureSize + 1]byte
		copy(signatureBytes[:], signature)
		s.Ed25519SuiSignature = &Ed25519SuiSignature{
			Signature: signatureBytes,
		}
	default:
		return errors.New("unsupport signature")
	}
	return nil
}

func NewSignatureSecure[T IntentValue](value IntentMessage[T], secret crypto.Signer[Signature]) (Signature, error) {
	message, err := bcs.Marshal(value)
	if err != nil {
		return Signature{}, err
	}
	hash := blake2b.Sum256(message)
	return secret.Sign(hash[:]), nil
}

type SignatureScheme struct {
	ED25519   *lib.EmptyEnum
	Secp256k1 *lib.EmptyEnum
	Secp256r1 *lib.EmptyEnum
	MultiSig  *lib.EmptyEnum
	BLS12381  *lib.EmptyEnum
}

func (s *SignatureScheme) Flag() byte {
	switch {
	case s.ED25519 != nil:
		return 0
	case s.Secp256k1 != nil:
		return 1
	case s.Secp256r1 != nil:
		return 2
	case s.MultiSig != nil:
		return 3
	case s.BLS12381 != nil:
		return 4
	default:
		return 0
	}
}

func NewSignatureScheme(flag byte) (SignatureScheme, error) {
	switch flag {
	case 0:
		return SignatureScheme{
			ED25519: &lib.EmptyEnum{},
		}, nil
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		fallthrough
	default:
		return SignatureScheme{}, errors.New("unsupported scheme")
	}
}

type Secp256k1SuiSignature struct {
	Signature []byte //secp256k1.pubKey + Secp256k1Signature + 1
}

type Secp256r1SuiSignature struct {
	Signature []byte //secp256k1.pubKey + Secp256k1Signature + 1
}

type Ed25519SuiSignature struct {
	Signature [ed25519.PublicKeySize + ed25519.SignatureSize + 1]byte
}

func NewSuiKeyPair(scheme SignatureScheme, seed []byte) SuiKeyPair {
	switch scheme.Flag() {
	case 0:
		return SuiKeyPair{
			Ed25519: crypto.NewEd25519KeyPair(ed25519.NewKeyFromSeed(seed[:])),
		}
	default:
		return SuiKeyPair{}
	}
}

type SuiKeyPair struct {
	Ed25519 *crypto.Ed25519KeyPair
	//Secp256k1 *Secp256k1KeyPair
	//Secp256r1 *Secp256r1KeyPair
	SignatureScheme
}

func (s *SuiKeyPair) PublicKey() []byte {
	switch s.Flag() {
	case 0:
		return s.Ed25519.PublicKey()
	default:
		return []byte{}
	}
}

func (s *SuiKeyPair) PrivateKey() []byte {
	switch s.Flag() {
	case 0:
		return s.Ed25519.PrivateKey()
	default:
		return []byte{}
	}
}

func (s *SuiKeyPair) Sign(msg []byte) Signature {
	switch s.Flag() {
	case 0:
		return Signature{
			Ed25519SuiSignature: NewEd25519SuiSignature(s.Ed25519, msg),
		}
	default:
		return Signature{}
	}
}

func NewEd25519SuiSignature(keyPair crypto.KeyPair, message []byte) *Ed25519SuiSignature {
	sig := keyPair.Sign(message)

	var signatureBytes [ed25519.PublicKeySize + ed25519.SignatureSize + 1]byte
	signatureBuffer := bytes.NewBuffer([]byte{})
	scheme := SignatureScheme{ED25519: &lib.EmptyEnum{}}
	signatureBuffer.WriteByte(scheme.Flag())
	signatureBuffer.Write(sig)
	signatureBuffer.Write(keyPair.PublicKey())
	copy(signatureBytes[:], signatureBuffer.Bytes())
	return &Ed25519SuiSignature{
		Signature: signatureBytes,
	}
}
