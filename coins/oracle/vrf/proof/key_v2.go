/*
*
The MIT License (MIT)

Copyright (c) 2018 SmartContract ChainLink, Ltd.
*/
package proof

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/bxlkm/go-wallet-sdk/crypto/vrf/secp256k1"
	"github.com/bxlkm/go-wallet-sdk/crypto/vrf/utils"
	"math/big"

	"go.dedis.ch/kyber/v3"
)

var suite = secp256k1.NewBlakeKeccackSecp256k1()

type Raw []byte

func (raw Raw) Key() (KeyV2, error) {
	rawKeyInt := new(big.Int).SetBytes(raw)
	k := secp256k1.IntToScalar(rawKeyInt)
	key, err := keyFromScalar(k)
	if err != nil {
		return KeyV2{}, err
	}
	return key, nil
}

func (raw Raw) String() string {
	return "<VRF Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

var _ fmt.GoStringer = &KeyV2{}

type KeyV2 struct {
	k         *kyber.Scalar
	PublicKey secp256k1.PublicKey
}

func NewV2() (KeyV2, error) {
	k := suite.Scalar().Pick(suite.RandomStream())
	return keyFromScalar(k)
}

func (key KeyV2) ID() string {
	return hexutil.Encode(key.PublicKey[:])
}

func (key KeyV2) Raw() Raw {
	return secp256k1.ToInt(*key.k).Bytes()
}

// GenerateProofWithNonce allows external nonce generation for testing purposes
//
// As with signatures, using nonces which are in any way predictable to an
// adversary will leak your secret key! Most people should use GenerateProof
// instead.
func (key KeyV2) GenerateProofWithNonce(seed, nonce *big.Int) (Proof, error) {
	secretKey := secp256k1.ScalarToHash(*key.k).Big()
	if !(secp256k1.RepresentsScalar(secretKey) && seed.BitLen() <= 256) {
		return Proof{}, fmt.Errorf("badly-formatted key or seed")
	}
	skAsScalar := secp256k1.IntToScalar(secretKey)
	publicKey := Secp256k1Curve.Point().Mul(skAsScalar, nil)
	h, err := HashToCurve(publicKey, seed, func(*big.Int) {})
	if err != nil {
		return Proof{}, fmt.Errorf("vrf.makeProof#HashToCurve, %v", err)
	}
	gamma := Secp256k1Curve.Point().Mul(skAsScalar, h)
	sm := secp256k1.IntToScalar(nonce)
	u := Secp256k1Curve.Point().Mul(sm, Generator)
	uWitness := secp256k1.EthereumAddress(u)
	v := Secp256k1Curve.Point().Mul(sm, h)
	c, err := ScalarFromCurvePoints(h, publicKey, gamma, uWitness, v)
	if err != nil {
		return Proof{}, fmt.Errorf("vrf.makeProof#ScalarFromCurvePoints, %v", err)
	}
	// (m - c*secretKey) % GroupOrder
	s := utils.Mod(utils.Sub(nonce, utils.Mul(c, secretKey)), secp256k1.GroupOrder)
	if e := checkCGammaNotEqualToSHash(c, gamma, s, h); e != nil {
		return Proof{}, e
	}
	outputHash := utils.MustHash(string(append(RandomOutputHashPrefix,
		secp256k1.LongMarshal(gamma)...)))
	rv := Proof{
		PublicKey: publicKey,
		Gamma:     gamma,
		C:         c,
		S:         s,
		Seed:      seed,
		Output:    outputHash.Big(),
	}
	valid, err := rv.VerifyVRFProof()
	if !valid || err != nil {
		return Proof{}, errors.New("constructed invalid proof")
	}
	return rv, nil
}

// GenerateProof returns gamma, plus proof that gamma was constructed from seed
// as mandated from the given secretKey, with public key secretKey*Generator
//
// secretKey and seed must be less than secp256k1 group order. (Without this
// constraint on the seed, the samples and the possible public keys would
// deviate very slightly from uniform distribution.)
func (key KeyV2) GenerateProof(seed *big.Int) (Proof, error) {
	for {
		nonce, err := rand.Int(rand.Reader, secp256k1.GroupOrder)
		if err != nil {
			return Proof{}, err
		}
		proof, err := key.GenerateProofWithNonce(seed, nonce)
		switch {
		case errors.Is(err, ErrCGammaEqualsSHash):
			// This is cryptographically impossible, but if it were ever to happen, we
			// should try again with a different nonce.
			continue
		case err != nil: // Any other error indicates failure
			return Proof{}, err
		default:
			return proof, err // err should be nil
		}
	}
}

func (key KeyV2) String() string {
	return fmt.Sprintf("VRFKeyV2{PublicKey: %s}", key.PublicKey)
}

func (key KeyV2) GoString() string {
	return key.String()
}

func keyFromScalar(k kyber.Scalar) (KeyV2, error) {
	rawPublicKey, err := secp256k1.ScalarToPublicPoint(k).MarshalBinary()
	if err != nil {
		return KeyV2{}, fmt.Errorf("could not marshal public key,%v", err)
	}
	if len(rawPublicKey) != secp256k1.CompressedPublicKeyLength {
		return KeyV2{}, fmt.Errorf("public key %x has wrong length", rawPublicKey)
	}
	var publicKey secp256k1.PublicKey
	if l := copy(publicKey[:], rawPublicKey); l != secp256k1.CompressedPublicKeyLength {
		return KeyV2{}, fmt.Errorf("failed to copy correct length in serialized public key")
	}
	return KeyV2{
		k:         &k,
		PublicKey: publicKey,
	}, nil
}
