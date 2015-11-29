#!/bin/sh
pushd util
python gen_icon_symbols.py
popd

go-bindata -o font.go -pkg materialicons font
