package main

type writer interface {
	Write([]byte) (int, error)
	WriteByte(byte) error
	WriteString(string) (int, error)
}
