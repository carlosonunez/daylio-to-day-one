# Daylio to Day One

This is a tool that migrates your Daylio entries to Day One from an iCloud
or Dropbox backup.

## Motivation

I loved using Daylio every day for three years, but there wasn't an easy way to
create journal entries from my MacBook. I considered using BlueStacks to run
Daylio on Android, but Android emulation doesn't work yet on Apple Silicon.

I discovered Day One while looking for solutions. It supports editing on my Mac
and Apple Watch, supports audio and video recordings, allows multiple journals
(which also replaces my work journal on GitHub) and synchronizes everything via
iCloud.

As much as I loved using Daylio, Day One is just better. Sorry!

## How To Use

1. [Download](https://github.com/carlosonunez/daylio-to-day-one/releases) the
   latest release for your platform

2. Run it! `./exporter-$VERSION-$OS-$ARCH`

This will automatically try to find your Daylio backup from a few known
locations. If that doesn't work, you can also provide the path to a backup file
like this:

```sh
./exporter-$VERSION-$OS-$ARCH [PATH_TO_BACKUP]
```


## Quirks

These were quirks I made to support my particular use case along with
environment variables you can set to disable them.

Create a file called `.env` and set the environment variables in there for
the flags to take effect.

| Quirk                                                                                      | Flag Environment Variable |
| :----                                                                                      | :------                   |
| Set Home Location when "Home" Activity detected and `HOME_ADDRESS_JSON` detected in dotenv | `NO_AUTO_HOME_LOCATION`   |
| Score alone time when "No", "A Little Bit", and "Yes" activities detected                  | `NO_ALONE_TIME_SCORING`   |
