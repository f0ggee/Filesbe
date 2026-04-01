package KeysManager

func (u *Updater) GetOldKey() []byte {

	return u.OldPrivateKey.Bytes()
}
