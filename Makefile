CC=go
PASS=meussegredos

vault:
	$(CC) build

run: vault
	./vault --k  meussegredos --tp plain.txt --cp cipher.txt --e
	./vault --k  meussegredos --cp cipher.txt --d

.PHONY: clean vault

clean:
	rm vault
