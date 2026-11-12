package RepoEncrypterKeys

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
