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

1. [Create a backup](#creating-a-daylio-backup) of your Daylio data.

   If you're using Daylio on an iPhone with iCloud enabled, save the file to
   your "Downloads" directory to have the exporter automatically find it.

1. [Download](https://github.com/carlosonunez/daylio-to-day-one/releases) the
   latest release for your platform.

2. Run it! `./exporter-$VERSION-$OS-$ARCH` or
   `./exporter-$VERSION-$OS-$ARCH [PATH_TO_BACKUP]` if you are not using iCloud
   or saved the file outside of the "Downloads" directory.

## Quirks

These were quirks I made to support my particular use case along with
environment variables you can set to disable them.

Create a file called `.env` and set the environment variables in there for
the flags to take effect.

| Quirk                                                                                      | Flag Environment Variable |
| :----                                                                                      | :------                   |
| Set Home Location when "Home" Activity detected and `HOME_ADDRESS_JSON` detected in dotenv | `NO_AUTO_HOME_LOCATION`   |
| Score alone time when "No", "A Little Bit", and "Yes" activities detected                  | `NO_ALONE_TIME_SCORING`   |

## Creating a Daylio Backup

Creating a Daylio backup is very easy.

First, tap on the **(...) More** tab then on "Backup and Restore" later in the
page.

<img
src="https://github.com/carlosonunez/daylio-to-day-one/raw/main/static/daylio-1.png" 
width=40%>

Once there, tap on **Advanced Options**

<img
src="https://github.com/carlosonunez/daylio-to-day-one/raw/main/static/daylio-2.png" 
width=40%>

Finally, tap on the big "Export" button and save the created file somewhere
convenient (like the Downloads directory).

<img src="https://github.com/carlosonunez/daylio-to-day-one/raw/main/static/daylio-3.png" 
width=40%>
