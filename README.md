### looperutil - a utility for interacting with Roland/Boss RC loopers

## NOTE: I've only tested this with the RC-5. Patches welcome for other loopers. This is also mostly guesswork, unless there's a spec for the file format knocking around nobody's telling me about. Caveat fuckaboutor.

This utility mostly operates on the MEMORYX.RC0 and SYSTEMX.RC0 file you find on your looper's storage. These are basic XML files. For the RC-5, you need to power on the looper with USB conencted, then press 'SETUP', scroll to 'STORAGE', and scroll down once. The looper screen will display "Connecting..." and you should be able to mount your filesystem. The RC0 files are then in //ROLAND/DATA

# General Usage

# Read-only stuff.

# List all Memory slots. Default for --memory is ./MEMORY1.RC0
looperutil --memory=/path/to/MEMORY1.RC0 ls

# Show Information about a Memory slot (slot 1 in this case).
looperutil ls 1

# Read-Write Stuff

# Rename a memory slot
looperutil set name 1 myloop

