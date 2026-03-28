package KeysManager

func (u Updater) GetKey() []byte {
	defer u.Mu.RUnlock()

	u.Mu.RLock()

	return u.NewPrivateKey.Bytes()
}
