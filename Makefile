.PHONY: test binary

test:
	$(MAKE) -C tacyborg test
	$(MAKE) -C cyborgnerd test

binary:
	$(MAKE) -C tacyborg binary
	$(MAKE) -C cyborgnerd binary
