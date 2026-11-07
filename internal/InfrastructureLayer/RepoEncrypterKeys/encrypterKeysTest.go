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

func (t TestEncrypterKeys) UpdateNewKey(newKey []byte) {

}

func (t TestEncrypterKeys) UpdateOldKey([]byte) {

}

type TestGetFileData struct {
	old  []byte
	_new []byte
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
