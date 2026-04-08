CC = gcc
CFLAGS = -Wall -Wextra -g -fsanitize=address
LDFLAGS = -fsanitize=address

.PHONY: test clean

test: test_runner
	./test_runner

test_runner: manifest.c test_manifest.c
	$(CC) $(CFLAGS) -o $@ $^ $(LDFLAGS)

clean:
	rm -f test_runner
