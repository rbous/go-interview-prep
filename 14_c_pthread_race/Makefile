CC = gcc
CFLAGS = -Wall -Wextra -g -fsanitize=thread
LDFLAGS = -fsanitize=thread -pthread

.PHONY: test clean

test: test_runner
	./test_runner

test_runner: progress.c test_progress.c
	$(CC) $(CFLAGS) -o $@ $^ $(LDFLAGS)

clean:
	rm -f test_runner
