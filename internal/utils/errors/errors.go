package errors

import "errors"

var (
	NotUnique               = errors.New("Запись с указанными данными уже существует")
	WrongLoginOrPasswordErr = errors.New("Неверный логин или пароль")
	NotEnoughCoinErr        = errors.New("У вас недостаточно стредств")
	NoUserErr               = errors.New("Пользователь не найден")
	NoMerchErr              = errors.New("Мерч не найден")
)
