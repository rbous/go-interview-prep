CC = gcc
CFLAGS = -Wall -Wextra -g -fsanitize=address
LDFLAGS = -fsanitize=address

.PHONY: test clean

test: test_runner
	./test_runner

test_runner: update_queue.c test_queue.c
	$(CC) $(CFLAGS) -o $@ $^ $(LDFLAGS)

clean:
	rm -f test_runner
