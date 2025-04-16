# looperutil - a utility for interacting with Roland/Boss RC loopers

## NOTE: I've only tested this with the BOSS RC-5. Patches welcome for other loopers. This is also mostly guesswork, unless there's a spec for the file format knocking around nobody's telling me about. Caveat fuckaboutor.

This utility mostly operates on the MEMORYX.RC0 and SYSTEMX.RC0 file you find on your looper's storage. These are basic XML files. For the RC-5, you need to power on the looper with USB conencted, then press 'SETUP', scroll to 'STORAGE', and scroll down once. The looper screen will display "Connecting..." and you should be able to mount your filesystem. The RC0 files are then in //ROLAND/DATA

### General Usage

```
# Common flags, mostly just trhe location of the memory XML file to work on.
# This is generally at /ROLAND/DATA/MEMORY1.RC0 from your looper's mount point.

`looperutil --memory=/media/doc/BOSS\ RC-5/ROLAND/DATA/MEMORY1.RC0 ls`
```

```
# Read-only stuff.

# List all Memory slots. Default for --memory is ./MEMORY1.RC0
`looperutil --memory=/path/to/MEMORY1.RC0 ls`

# Show Information about a Memory slot (slot 1 in this case).
`looperutil ls 1`

# Just spit out a <NAME> tag and contents for a given slot name.
`looperutil xmlname "My Loop"`

# Get a specific attribute for a loop
`looperutil get 01 Track1.Pan`
`looperutil get 01 Name`

```

```
# Read-Write Stuff

# These commands spit a new XML file to stdout, unless you specify `--xml-output=` to another file.

# Rename a memory slot
`looperutil set 01 Name myloop`

#  Set other various bits (the argument is the same as the 'ls' output keys)
`looperutil set 01 Track1.One 1` # Set one-time mode for loop playback.
`looperutil set 01 Master.Tempo 1400`  # This is x10 for some reason.

# Convert and install an audio file to a slot. The looper only likes 44.1khz stereo wav files.
# This uses ffmpeg in the background, and will overwrite /ROLAND/WAVE/001_01/001.wav for slot 1 (and so on).
# use --looper-data-dir to point to your mount point, and --clear-existing-wav if you want to nuke anything that's already in that slot.
`looperutil --looper-data-dir=/media/doc/BOSS\ RC-5/ installwav 01 ~/Music/myrandomfile.mp3`

```
