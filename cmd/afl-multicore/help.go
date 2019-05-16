package main

const (
	// TODO: use the actual afl version for help text
	helpText = `afl-fuzz 2.52b by <lcamtuf@google.com>

afl-fuzz [ options ] -- /path/to/fuzzed_app [ ... ]

Required parameters:

	-i dir        - input directory with test cases
	-o dir        - output directory for fuzzer findings

Execution control settings:

	-f file       - location read by the fuzzed program (stdin)
	-t msec       - timeout for each run (auto-scaled, 50-1000 ms)
	-m megs       - memory limit for child process (50 MB)
	-Q            - use binary-only instrumentation (QEMU mode)

Fuzzing behavior settings:

	-d            - quick & dirty mode (skips deterministic steps)
	-n            - fuzz without instrumentation (dumb mode)
	-x dir        - optional fuzzer dictionary (see README)

Other stuff:

	-T text       - text banner to show on the screen
	-M / -S id    - distributed mode (see parallel_fuzzing.txt)
	-C            - crash exploration mode (the peruvian rabbit thing)

For additional tips, please consult docs/README.

Multicore settings:

	-session name - prefix for all fuzz workers' name (fuzzer)
	-workers n    - number of workers to spin up (# CPUs)	
	-nocheck      - don't check if binary works before spinning up
`
)
