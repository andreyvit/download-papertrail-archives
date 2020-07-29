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


## Changelog

* 1.0.0 (2020-07-29) — initial release
