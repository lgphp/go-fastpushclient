package fastpushclient

import (
	"encoding/base64"
	"github.com/rogpeppe/fastuuid"
	"github.com/wumansgy/goEncrypt"
	"push-sdk-go/client/utils"
	"testing"
)

func Test_GetUUID(t *testing.T) {

	println(fastuuid.MustNewGenerator().Hex128())
}

func TestZIP(t *testing.T) {
	str := "fdfsakljkfkdskfsafkjsakjdfkasfjkasdkjfkasfkjsdakjdffadkdkksafkkkfkjdsjkfksajfklaskdlfaslk;fdlas;faslkfas;klfasfasfdasfasfa"
	zipBytes := utils.Gzip([]byte(str))
	zipstring := base64.StdEncoding.EncodeToString(zipBytes)
	println("zipstring", zipstring)

	unzipBytes := utils.UnGzip(zipBytes)

	println("uzip", string(unzipBytes))

}

func TestAes(t *testing.T) {
	str := "我是中国人"
	key := "NNL1V0oy3+I6XvQZnY+L9Q/G5U8rgeNS4SUELeqUtgIio16dqEiAjFtee5htl2Fl"
	iv := []byte("1234567812345678")
	keyBytes, _ := base64.StdEncoding.DecodeString(key)
	enckey := utils.GetMsgEncKey(keyBytes)
	println("enckey:", base64.StdEncoding.EncodeToString(enckey))
	re, _ := goEncrypt.AesCbcEncrypt([]byte(str), enckey, iv)
	//cipher := aesCbc.NewAesCipher(enckey, iv)
	//re := cipher.Encrypt([]byte(str))
	s := base64.StdEncoding.EncodeToString(re)

	println(s)
	//bytes, _ := base64.StdEncoding.DecodeString("2AybKNs4khfTc4HMjC73pg==")
	//decrypt := cipher.Decrypt(bytes)
	//
	//println("decrypt:" , string(decrypt))

}
