# afl-multicore

afl-multicore is a set of tools that makes it easy to manage multiple instances of AFL on a single machine

## afl-multicore

`afl-multicore` was built to be a drop-in replacement for the `afl-fuzz` command. If you take a standard afl command such as `afl-fuzz -i ~/fuzz_in -o ~/fuzz_out -- /my/binary`, and replace `afl-fuzz` with `afl-multicore`, it will spin up one instance for each CPU on your machine.

Running `afl-multikill` kills these instances.

`afl-multicore` has additional options that allow you to name your fuzzing session and decide how many instances are created. You can learn more by running `afl-multicore --help`:

```
--- Cut off standard AFL options from output

Multicore settings:

	-session name - prefix for all fuzz workers' name (fuzzer)
	-workers n    - number of workers to spin up (# CPUs)
	-nocheck      - don't check if binary works before spinning up
```

## afl-multikill

`afl-multikill` can stop sessions spawned by `afl-multicore`.

If you named your session, you can kill it by running `afl-multikill -session <session-name>`

`afl-multikill --help`:

```
afl-multikill [options]

Options:

	-session name - the session to kill (fuzzer)
```

## afl-multistats

`afl-multistats` gathers and aggregates stats from an AFl sync directory, and presents them either in a human-readable format, or in JSON. It also includes a HUD mode which is similar to AFL's retro-style UI.

`afl-multistats --help`:

```
afl-multistats [options] /path/to/sync/dir

Options:

	-format fmt   - the format to output stats in (human)
	-basic        - only output basic stats
	-hud          - display a persistent stats HUD
```
