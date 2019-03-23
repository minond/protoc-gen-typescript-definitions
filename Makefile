.PHONY: test

test:
	@if [ ! -d test/results ]; then mkdir test/results; fi
	protoc --typescript-definitions_out=test/results/ -I test/definitions/ test/definitions/*.proto
	diff test/results/* test/expected/*

testregen:
	protoc --typescript-definitions_out=test/expected/ -I test/definitions/ test/definitions/*.proto
