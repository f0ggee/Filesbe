package Controller

//func Check(r *http.Request) error {
//	store := Store()
//	session, err := store.Get(r, "token6")
//	if err != nil {
//		slog.Error("Error check", "Err", err)
//
//		return err
//
//	}
//
//	key, _ := hex.DecodeString(os.Getenv("KEYFORJWT"))
//
//	jwts, _ := session.Values["JWT"].(string)
//
//	_, err = jwt.ParseWithClaims(jwts, &Dto.MyCustomCookie{}, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
//		}
//		return key, nil
//	})
//	if err != nil {
//		rtToken, _ := session.Values["RT"].(string)
//
//		RF, err := jwt.ParseWithClaims(rtToken, &Dto.MyCustomCookie{}, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
//			}
//			return key, nil // Возвращаем секретный ключ для проверки
//		})
//		if err != nil {
//			slog.Error("Error parse token JWT", "Err", err)
//			return err
//		}
//
//		if !RF.Valid {
//			slog.Error("Token is not valid")
//
//		}
//	}
//
//	return nil
//
//}
