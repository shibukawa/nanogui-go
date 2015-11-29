#!/bin/sh
pushd util
python gen_icon_symbols.py
popd
go-bindata -pkg nanogui -o fonts.go fonts
