.PHONY: test

test:
	@if [ -d test/results ]; then rm -r test/results; fi
	@mkdir test/results
	protoc --typescript-definitions_out=test/results/ -I test/definitions/ test/definitions/*.proto
	diff test/results/ test/expected/

gen:
	@if [ -d test/expected ]; then rm -r test/expected; fi
	@mkdir test/expected
	protoc --typescript-definitions_out=test/expected/ -I test/definitions/ test/definitions/*.proto
