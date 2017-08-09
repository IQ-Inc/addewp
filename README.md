## addewp

The command line utility adds a file into IAR's EWP file

Interactive usage: `$ ./addewp.exe`

Quick usage with flags: `$ ./addewp.exe -ewp MyEwpFile.ewp -file main.cpp`

Flags may be omitted for interactivity. For example, the following invocation
prompts for an EWP file, since the new file is provided as an argument:

```
$ ./addewp.exe -file main.cpp
```