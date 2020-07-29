# download-papertrail-archives

A Golang script to download all missing archives from Papertrail.


## Usage

Build:

    go install

Run:

    download-papertrail-archives -token <TOKEN> -o <DIR>

where:

* `<TOKEN>` is your Papertrail HTTP API token from [the profile page](https://papertrailapp.com/account/profile),
* `<DIR>` is the directory to save downloaded logs to.


## Options

* `-token <TOKEN>` sets the Papertrail HTTP API token

* `-o <DIR>` sets output directory (default `.`)

* `-since <YYYY-MM-DD>` only downloads logs on or after this date

* `-until <YYYY-MM-DD>` only downloads logs on or before this date

* `-timeout <DURATION>` sets timeout for HTTP operations (default `30s`)

* `-q` enables quiet operation (don't print any progress information)


## Changelog

* 1.0.0 (2020-07-29) — initial release
