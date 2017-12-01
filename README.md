# Hi!

Chat server and client based on a dead simple implementation by [hellerve]
(https://github.com/hellerve/hi).

There __are__ bugs.
Room destruction is broken now,
I'll fix it as soon as I rewrite `Message` parsing.


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

`Server`, `Room`, `Client` and `Message` types are abbreviated as `s`, `r`, `c`
and `m` respectively.

Tabs for indentation, spaces for alignment.

## Commands

There are a few special commands that you can issue to interact with the server,
IRC-style. All of these commands need to be sent to a channel that the user
subscribes to.

```
/join <channelname>  # join a channel
/leave <channelname> # leave a channel (will send back an error if user is not subscribed to the channel)
/who                 # list all users in the current channel
/channels            # list all channels
/nick <username>     # change your nickname
/whoami              # prints your nickname@remote_addr
/whois <username>    # prints nickname@remote_addr for with that username
```

--------------------------------------------------------------------------------

Have fun!
