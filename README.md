# Daylio to Day One

This is a tool that helps you export the text from your Daylio entries into a
DayOne-compatible format from a provided CSV.

> ⚠️  Audio and video captured in Daylio is not supported. I recommend backing
> that up to Dropbox, as they are not easily discoverable within iCloud and
> Daylio only works on iOS devices.

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

> ✅ You'll need to install Docker to use this. [Click here](https://get.docker.io)
> to learn how or run `brew install docker` if you're on a Mac.

1. Clone this repository.
2. Export your Daylio entries as a CSV and save them to this repository's
   directory in a file called `daylio.csv`.
3. Run the exporter: `docker-compose run --rm export-daylio-entries`.
   This will produce a JSON file called `dayone.json`.
4. Import the JSON into Daylio. The entries will get imported into your default
   journal.

## Quirks

These were quirks I made to support my particular use case along with
environment variables you can set to disable them.

Create a file called `.env` and set the environment variables in there for
the flags to take effect.

| Quirk                                                                                      | Flag Environment Variable |
| :----                                                                                      | :------                   |
| Set Home Location when "Home" Activity detected and `HOME_ADDRESS_JSON` detected in dotenv | `NO_AUTO_HOME_LOCATION`   |
| Score alone time when "No", "A Little Bit", and "Yes" activities detected                  | `NO_ALONE_TIME_SCORING`   |
