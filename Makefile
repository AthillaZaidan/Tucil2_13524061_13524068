build:
	cd src && go build -o ../bin/OctoVox .

run:
	./bin/OctoVox

clean:
	rm -f bin/OctoVox
