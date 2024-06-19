package service

import (
	"crypto/x509"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"faraway/internal/domain"
)

const (
	testKey           = "3082025c020100028181009b19532980488df83bf5e12229595a65d73bf11532ae06ce179533f5035f27f6e0be71bc8d76d28661c7961a51dd9175befafea415c544be911e58f67a78fbe0447fe2a4471306619601fa5622510a9eedd239c4c115e682b4a25cd981e8e5b2bc384f12bdfbe69076da25d06399a87355f284f1f334e0a616d4dfeb5cf080150203010001028180321fd11c8c74e64cdf33eb7a5adaa1b86002e33af2920368ff7e1cb8864a6e63fee60d63de64144d91b42af27e9a98b3f0f0b4f2da86525d341116b731858000cdd6d1401f6c94a798f484c630efd41491b59cb2d576e816b4a52e639a9da645df4266bc4dc9aa1cf65e33d8340ddd942d5344f87f60042e144b709a125fd281024100c77bbf1fef8e0ab020573b0eedf4da00702f8646621af566cb5f30651564cf9eff627249d6b6c458827f19b5169c95a7472f44101f64520c70590074cb0a38ad024100c70a6f56b19514fa0ee0f8e027c3ee6d33f9d7b3556f4ba4f7abdb8f575f79f80418a791e0abc6cf5f80f92c771f6162ac61da94e7009eacb1897efddbccca0902405ae82fdf33e23d48aa54565ba56151ffa5206346abeab12ed93b55e89ae9481ca3318ff7ca5b9bfae1ed5e1fc260356af7ebb84ec89f852c99fe5550e43e92390240368a2f740bf913e4694b5026ebfe8e48b22355edb80d6526f10ed07cf8ae1ad7d11788633ab317291fbc518ad3a16fa80020582ad119a46121ccd15572732d31024100ad6ca95fdcc70304a7108c8032279b76f6aed4b4d18486aa0474cdeb5f6b441a1789b68f22f71e6615e2b0784499f36b10114c74b232408637d5ecbc2f268ba9"
	defaultID         = "127.0.0.1"
	defaultDifficulty = 2
)

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(serviceSuite))
}

type serviceSuite struct {
	suite.Suite

	p *Puzzler

	challenge string
	proof     string
}

func (s *serviceSuite) SetupSuite() {
	k, err := hex.DecodeString(testKey)
	s.Require().NoError(err)

	pk, err := x509.ParsePKCS1PrivateKey(k)
	s.Require().NoError(err)
	s.Require().NotNil(pk)

	ctrl := gomock.NewController(s.T())
	kr := NewMockKeyRepository(ctrl)
	kr.EXPECT().GetPrivateKey().Return(pk, nil).AnyTimes()

	s.p = NewPuzzler(kr, defaultDifficulty)

	challenge, difficulty, err := s.p.GenerateChallenge(defaultID)
	s.Require().NoError(err)
	s.Require().NotEmpty(challenge)
	s.Require().NotEmpty(difficulty)

	proof, err := SolveChallenge(challenge, difficulty)
	s.Require().NoError(err)
	s.Require().NotEmpty(proof)

	s.challenge = challenge
	s.proof = proof
}

func (s *serviceSuite) TestVerifyProof() {
	tests := []struct {
		name      string
		id        string
		challenge string
		proof     string
		expErr    error
	}{
		{
			name:      "positive",
			id:        defaultID,
			challenge: s.challenge,
			proof:     s.proof,
		},
		{
			name:      "invalid id",
			id:        "invalid",
			challenge: s.challenge,
			proof:     s.proof,
			expErr:    domain.InvalidChallenge,
		},
		{
			name:      "invalid challenge",
			id:        defaultID,
			challenge: hex.EncodeToString([]byte("challenge")),
			proof:     s.proof,
			expErr:    domain.InvalidChallenge,
		},
		{
			name:      "invalid proof",
			id:        defaultID,
			challenge: s.challenge,
			proof:     hex.EncodeToString([]byte("proof")),
			expErr:    domain.InvalidProof,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Require().ErrorIs(s.p.VerifyProof(tt.id, tt.challenge, tt.proof), tt.expErr)
		})
	}
}
