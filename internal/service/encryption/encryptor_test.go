package encryption

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/fs"
	"math/big"
	"os"
	"syscall"
	"testing"
)

func Test_encryptor_Encrypt(t *testing.T) {
	const (
		publicKey  = "-----BEGIN RSA PUBLIC KEY-----\nMIICCgKCAgEAtRVWix3Lt60p+PN/P34CQmEroHsVgUEBz9op8W0gAVpUS+/MV49W\nAG7ywJbyMRwfj1F5ZnMFpYT3BfqT5lYBXMA9Pf668Qe9n5iwu7DHgehrBfDBIXPU\njSAPeY+WZjaeYm8apaTZvf6/XOrcMrdIBplK4Zbn4UYE9gFMwzfAHu7Q0dF0hUVb\n96w5TVnRHG/L/3LFsjFRbo3K2rY9AtuS+prAhHDgrX7k18PaX62XEauF0cB7tIiK\nbKJjh2fKBj88snKw4GSrBLGyizEeNnbRLMcCy4UQhu6qOfYgM3VBJxYjLJhiIVqz\nMaku2fBQg24as1NRHFN20si5h+bCm2PpY+Ak6YVAlMoaRx3KXXhX1w42scNQKYMl\njdiZeZRMXehrva8pbY8aZUFS55xVqOdRWgkB2UNg+F3f3PkaW+xeJV3PHyRi5gOk\nVdxh9FjAtQdoHqtunOuUl7pFfMwYrqhtnwcACwZ2kYelLlXZlTst/tkD91m/iyqw\ntWBONS/v0fjmrmuqQt7JH//1czS0CZ6QoV2EzI4pRrSIlqB/CqAqChBKRXfLtSWU\nr2n8IlSJMiV+E4wciyth1DWFfe9rxvzw5/NKHuDJ4kO7oDUKR3NLKNBXzWaEBEnR\nuag15KLYCR/xFJW7qlUkabTQ2L9J2Wzxqi47hGB9Njk19/XpG4DIZlECAwEAAQ==\n-----END RSA PUBLIC KEY-----\n"
		privateKey = "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"
	)
	type fields struct {
		publicKey  func() *rsa.PublicKey
		privateKey func() *rsa.PrivateKey
	}
	type args struct {
		plaintext []byte
	}
	type want struct {
		err    error
		result []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "public key not specified",
			fields: fields{
				publicKey: func() *rsa.PublicKey {
					return nil
				},
			},
			want: want{
				result: nil,
				err:    errors.New("RSA public key not specified"),
			},
		},
		{
			name: "successfully case",
			fields: fields{
				publicKey: func() *rsa.PublicKey {
					block, _ := pem.Decode([]byte(publicKey))
					key, _ := x509.ParsePKCS1PublicKey(block.Bytes)

					return key
				},
				privateKey: func() *rsa.PrivateKey {
					block, _ := pem.Decode([]byte(privateKey))
					key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

					return key
				},
			},
			args: args{
				plaintext: []byte("Hello world!"),
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := &encryptor{publicKey: test.fields.publicKey()}
			ciphertext, err := e.Encrypt(test.args.plaintext)

			require.Equal(t, test.want.err, err)
			if test.want.err == nil {
				d := &decryptor{privateKey: test.fields.privateKey()}
				text, err := d.Decrypt(ciphertext)
				require.Nil(t, err)
				assert.Equal(t, test.args.plaintext, text)
			}
		})
	}
}

func TestNewEncryptor(t *testing.T) {
	n := new(big.Int)
	n.SetString("738755621979269410308665542380331961407983781253820926953091776470674864014650766740296558456758339993827646923379632754360036320659751921363963339890249820517331715954704747977022473905663477462395464413694667999222799925789715809208403968657624569018799351998401191841366640393028327573943443924793518086243377403268963213892628807375837803847687882106427806832463161222128712994736540209742501155757334891396489171195457140164186526669167671304818274052097901125926504052953136459605207336782181896617885783477731259209961728856829772542805771952506571968571059598371191960457902323745177303081688974466946069624272867582524145414837578673914854486137171076244329145679049812736146204450857868942000216609203556844766792394339845477398904856318727607085813887792472874061658186742475396206267002993955162166210390515692415836209822799907877174732819777359414939745422608980083765241452238875223247036775399190226547712203215251142335240434308119620336361341736083993740882904823115163966096024278256728516354546895465281504553749690337785517338238988341433024777828331795701917848946875009034289496519089778026580164647468114001055633003365511962085049165604197473085840775468258231190303598919734987501153050935990865403407001169", 10)
	type fields struct {
		data       []byte
		createFile bool
	}
	type args struct {
		keyFile string
	}
	type want struct {
		encryptor Encryptor
		err       error
	}
	tests := []struct {
		want   want
		args   args
		name   string
		fields fields
	}{
		{
			name: "key file not found",
			args: args{
				keyFile: "/tmp/grdn_key_test_encryptor",
			},
			fields: fields{
				createFile: false,
			},
			want: want{
				encryptor: nil,
				err: &fs.PathError{
					Op:   "open",
					Path: "/tmp/grdn_key_test_encryptor",
					Err:  syscall.ENOENT,
				},
			},
		},
		{
			name: "invalid key",
			args: args{
				keyFile: "/tmp/grdn_key_test_encryptor",
			},
			fields: fields{
				data:       []byte("Hello world!"),
				createFile: true,
			},
			want: want{
				encryptor: nil,
				err:       errors.New("failed to parse PEM block containing the key"),
			},
		},
		{
			name: "successfully case",
			args: args{
				keyFile: "/tmp/grdn_key_test_encryptor",
			},
			fields: fields{
				data:       []byte("-----BEGIN RSA PUBLIC KEY-----\nMIICCgKCAgEAtRVWix3Lt60p+PN/P34CQmEroHsVgUEBz9op8W0gAVpUS+/MV49W\nAG7ywJbyMRwfj1F5ZnMFpYT3BfqT5lYBXMA9Pf668Qe9n5iwu7DHgehrBfDBIXPU\njSAPeY+WZjaeYm8apaTZvf6/XOrcMrdIBplK4Zbn4UYE9gFMwzfAHu7Q0dF0hUVb\n96w5TVnRHG/L/3LFsjFRbo3K2rY9AtuS+prAhHDgrX7k18PaX62XEauF0cB7tIiK\nbKJjh2fKBj88snKw4GSrBLGyizEeNnbRLMcCy4UQhu6qOfYgM3VBJxYjLJhiIVqz\nMaku2fBQg24as1NRHFN20si5h+bCm2PpY+Ak6YVAlMoaRx3KXXhX1w42scNQKYMl\njdiZeZRMXehrva8pbY8aZUFS55xVqOdRWgkB2UNg+F3f3PkaW+xeJV3PHyRi5gOk\nVdxh9FjAtQdoHqtunOuUl7pFfMwYrqhtnwcACwZ2kYelLlXZlTst/tkD91m/iyqw\ntWBONS/v0fjmrmuqQt7JH//1czS0CZ6QoV2EzI4pRrSIlqB/CqAqChBKRXfLtSWU\nr2n8IlSJMiV+E4wciyth1DWFfe9rxvzw5/NKHuDJ4kO7oDUKR3NLKNBXzWaEBEnR\nuag15KLYCR/xFJW7qlUkabTQ2L9J2Wzxqi47hGB9Njk19/XpG4DIZlECAwEAAQ==\n-----END RSA PUBLIC KEY-----\n"),
				createFile: true,
			},
			want: want{
				encryptor: &encryptor{
					publicKey: &rsa.PublicKey{
						N: n,
						E: 65537,
					},
				},
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Remove(test.args.keyFile)
			if test.fields.createFile {
				f, err := os.Create(test.args.keyFile)
				require.Nil(t, err)
				_, err = f.Write(test.fields.data)
				require.Nil(t, err)
				err = f.Close()
				require.Nil(t, err)
			}
			got, err := NewEncryptor(test.args.keyFile)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.encryptor, got)
			if test.want.encryptor != nil {
				assert.Implements(t, (*Encryptor)(nil), got)
			}
		})
	}
}
