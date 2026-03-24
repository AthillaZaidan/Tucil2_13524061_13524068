SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command

build:
	cd src; go build -o ../bin/OctoVox.exe .; $$gopath = go env GOPATH; $$dllDir = Get-ChildItem "$$gopath\pkg\mod\github.com\g3n\engine@*\audio\windows\bin" -Directory -ErrorAction SilentlyContinue | Select-Object -First 1 -ExpandProperty FullName; if ($$dllDir) { Copy-Item "$$dllDir\*.dll" ../bin -Force -ErrorAction SilentlyContinue }

run:
	./bin/OctoVox.exe

clean:
	Remove-Item -Force -ErrorAction SilentlyContinue bin/OctoVox.exe, bin/*.dll
