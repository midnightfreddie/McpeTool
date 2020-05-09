The mcpetool command line executable offers some db put/get/delete options from command line, but it can also start a local web server with the API to the Bedrock world.

```
NAME:
   mcpetool - Reads and writes a Minecraft Bedrock Edition world directory.

USAGE:
   mcpetool [global options] command [command options] [arguments...]

VERSION:
   0.3.2

AUTHOR:
   Jim Nelson <jim@jimnelson.us>

COMMANDS:
   leveldat  Get or put level.dat data
   db        List, get, put, or delete leveldb keys
   api, www  Open world, start API at http://127.0.0.1:8080 . Control-c to exit.
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --path FILEPATH, -p FILEPATH  FILEPATH of world (default: ".") [$MCPETOOL_WORLD]
   --help, -h                    show help (default: false)
   --version, -v                 print the version (default: false)

COPYRIGHT:
   (c) 2018, 2020 Jim Nelson
```