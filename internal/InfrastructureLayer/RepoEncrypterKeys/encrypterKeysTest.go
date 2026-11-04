package RepoEncrypterKeys

type TestEncrypterKeys struct {
	d []byte
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

func (t TestEncrypterKeys) UpdateNewKey(bytes []byte) {

}

func (t TestEncrypterKeys) UpdateOldKey() {

}
