package main

const (
	ModeIRWXU = 00700
	ModeIRUSR = 00400
	ModeIWUSR = 00200
	ModeIXUSR = 00100

	ModeIRWXG = 00070
	ModeIRGRP = 00040
	ModeIWGRP = 00020
	ModeIXGRP = 00010

	ModeIRWXO = 00007
	ModeIROTH = 00004
	ModeIWOTH = 00002
	ModeIXOTH = 00001

	ModeIRUGO   = (ModeIRUSR | ModeIRGRP | ModeIROTH)
	ModeIWUGO   = (ModeIWUSR | ModeIWGRP | ModeIWOTH)
	ModeIXUGO   = (ModeIXUSR | ModeIXGRP | ModeIXOTH)
)
