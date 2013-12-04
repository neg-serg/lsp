package main

const (
	modeIRWXU = 00700
	modeIRUSR = 00400
	modeIWUSR = 00200
	modeIXUSR = 00100

	modeIRWXG = 00070
	modeIRGRP = 00040
	modeIWGRP = 00020
	modeIXGRP = 00010

	modeIRWXO = 00007
	modeIROTH = 00004
	modeIWOTH = 00002
	modeIXOTH = 00001

	modeIRUGO = (modeIRUSR | modeIRGRP | modeIROTH)
	modeIWUGO = (modeIWUSR | modeIWGRP | modeIWOTH)
	modeIXUGO = (modeIXUSR | modeIXGRP | modeIXOTH)
)
