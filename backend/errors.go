package main

import "errors"

var (
	ErrDBOperation   = errors.New("DB를 읽는 중 오류가 발생했습니다")
	ErrFileOperation = errors.New("파일을 읽거나 쓰지 못했습니다")
	ErrKeyNotFound   = errors.New("키가 잘못되었거나, 해당 키의 파일을 찾지 못했습니다")
)
