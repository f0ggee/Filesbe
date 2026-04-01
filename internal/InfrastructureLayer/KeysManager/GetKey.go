package KeysManager

func (u *Updater) GetKey() []byte {
	defer u.Mu.RUnlock()

	u.Mu.RLock()

	//e, err := x509.ParsePKCS1PrivateKey(u.NewPrivateKey.Bytes())
	//if err != nil {
	//	slog.Error("Error while unmarshalling NewSavingRsa", "Error", err.Error())
	//}
	return u.NewPrivateKey.Data()
}
