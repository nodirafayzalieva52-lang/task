package errs

import "errors"

var (
	ErrDatabaseAction = errors.New("ошибка выполнения операции в базе данных")
	ErrEventBusAction = errors.New("ошибка отправки события в шину")
)