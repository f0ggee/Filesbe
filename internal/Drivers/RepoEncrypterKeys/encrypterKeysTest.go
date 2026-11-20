package RepoEncrypterKeys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

type TestEncrypterKeys struct {
	d []byte
}

func (t TestEncrypterKeys) GetServerKey() []byte {
	return t.d
}

func (t TestEncrypterKeys) GetMasterKey() []byte {
	return t.d
}

func NewTestEncrypterKeys(d []byte) *TestEncrypterKeys {
	return &TestEncrypterKeys{d: d}
}

func (t TestEncrypterKeys) GetKey() []byte {
	return t.d
}

func (t TestEncrypterKeys) GetOldKey() []byte {
	return t.d
}

func (t TestEncrypterKeys) UpdateNewKey(newKey []byte) {

}

func (t TestEncrypterKeys) UpdateOldKey([]byte) {

}

type TestGetFileData struct {
	old  []byte
	_new []byte
}

func (t *TestGetFileData) GetServerKey() []byte {
	return nil
}

func (t *TestGetFileData) GetMasterKey() []byte {
	return nil
}

func NewTestGetFileData() *TestGetFileData {
	return &TestGetFileData{}
}

func (t *TestGetFileData) GetKey() []byte {
	return t._new
}

func (t *TestGetFileData) GetOldKey() []byte {
	return t.old
}

func (t *TestGetFileData) UpdateNewKey(newKey []byte) {
	t._new = newKey

}

func (t *TestGetFileData) UpdateOldKey(key []byte) {
	t.old = key
}

type TestPrepareData struct {
	masterKey  []byte
	serverKeys []byte
	newKey     []byte
	oldKey     []byte
}

func NewTestPrepareData(masterKey []byte, serverKey []byte) *TestPrepareData {
	return &TestPrepareData{
		masterKey:  masterKey,
		serverKeys: serverKey,
		newKey:     nil,
		oldKey:     nil,
	}
}

func (t *TestPrepareData) GetKey() []byte {
	return t.newKey
}

func (t *TestPrepareData) GetOldKey() []byte {
	return t.oldKey
}

func (t *TestPrepareData) UpdateNewKey(newKey []byte) {
	t.newKey = newKey
}

func (t *TestPrepareData) UpdateOldKey(otherKey []byte) {
	if otherKey != nil {
		t.oldKey = otherKey
		return
	}
	t.oldKey = t.newKey
}

func (t *TestPrepareData) GetMasterKey() []byte {
	return t.masterKey
}

func (t *TestPrepareData) GetServerKey() []byte {
	return t.serverKeys
}

type NormalKeyTest struct {
	k    []byte
	oldK []byte
}

func NewNormalKeyTest() *NormalKeyTest {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	oldKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return &NormalKeyTest{
		k:    x509.MarshalPKCS1PrivateKey(key),
		oldK: x509.MarshalPKCS1PrivateKey(oldKey),
	}
}

func (n *NormalKeyTest) GetKey() []byte {
	return n.k
}

func (n *NormalKeyTest) GetOldKey() []byte {
	return n.oldK
}

func (n *NormalKeyTest) UpdateNewKey(newKey []byte) {
	n.k = newKey
}

func (n *NormalKeyTest) UpdateOldKey(otherKey []byte) {
	if otherKey != nil {
		n.oldK = otherKey
		return
	}
	n.oldK = n.k
}

func (n *NormalKeyTest) GetMasterKey() []byte {
	return nil
}

func (n *NormalKeyTest) GetServerKey() []byte {
	return nil
}
