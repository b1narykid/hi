# Hi!

A dead simple chat server and client without any external dependencies.

## Bugs

There are probably bugs.

## Channels

Channels are weird. Anyone can send messages to any channel. A channel
exists as long as there is someone in that channel. If there isn't, the
channel will be closed. What "being in a channel" means is basically
that you subscribe to it, i.e. you will receive messages posted in that
channel. It doesn't affect the way you post there, though. Channels
are conversations. You can scream something at a group of people, or
choose to join that group and interact with them.

## Security

Nothing is stored in a database. Everything is kept in memory, there
is no history. That means that all knowledge and all communication is
ephemeral, like at a party where you join and leave conversations as
you please. Please use HTTPS when deploying this chat, preferably through
caddy or any other web server.

## Deployment

This is where it gets tricky. If you know how to (cross-)compile Go programs and
put them on your server and run it, deploying `hi` is as simple as compiling,
e.g. with `GOOS=<os, probably linux> go build`, and copying the resulting binary
and the public directory to a server, then letting it run. Letting it run could
be more or less involved, depending on your hosting provider.
I'll be glad to help.

## Running

The compiled program can be run without arguments (runs on port 8080) or with
the port set (using the `-p` option). It assumes that the public directory is in
the directory where the program is started.

## Code style

### Indentation

Tabs for indentation, no alignment.

### Naming

Use descriptive names for constants, types, functions and methods.  If possible,
abbreviate variable names to the lowercase first letter of its type.

## Credits

Chat server was originally based on [hellerve/hi](https://github.com/hellerve/hi),
but was rewritten from scratch.  Client is still based on the upstream.

## License

To the extent possible under law, the author(s) have dedicated all copyright and
related and neighboring rights to this software to the public domain worldwide.
This software is distributed without any warranty.

You should have received a copy of the CC0 Public Domain Dedication along with
this software. If not, see [CC0](https://creativecommons.org/publicdomain/zero/1.0/).

--------------------------------------------------------------------------------

Have fun!
