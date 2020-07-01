# hoster

hoster is a simple utility to manage multiple machines as you probably would do if had manage and worked on a single machine.

The util gathers useful repttive commands in concurrent way "non blocking apporach" the command will be blocked only when it needs or must be blocked e.g. waiting for previous command to be complete (waiting for completion). In this way, single operation which is constructed from many other commands will consume time close to the actual time that would of taken anyway by a single machine.

There are some nice warrpers that represents common data in tables. To achive the same "manually" would take bounch of commands and grep / slice the the relevant data and take care to conncrent (async approcach)
The cli has has diffrent context and probably match as much as possible to actual command you would do on a single machine.

There are no examples here but cli is intuitive and has good examples compile and run it. 

There a few context commands and flags as the hosts details and other as override the defualt path for hosts list and etc.
Multiple structural data documents are supported (tomal, json, yaml and other).

Add command is relatively easy. If the commands need to be more dynamic then it might be better to collect them from external location (file or db).
