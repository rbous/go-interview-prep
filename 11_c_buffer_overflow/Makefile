CC = gcc
CFLAGS = -Wall -Wextra -g -fsanitize=address
LDFLAGS = -fsanitize=address

.PHONY: test clean

test: test_runner
	./test_runner

test_runner: version_parser.c test_parser.c
	$(CC) $(CFLAGS) -o $@ $^ $(LDFLAGS)

clean:
	rm -f test_runner
