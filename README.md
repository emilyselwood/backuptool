# Backup tool

A simple command line backup tool.
Configure. Add a directory you care about. Set a cron schedule and forget about it.

## Features

* Automatic upload to google drive
* Simple to add new directories
* back up and restore cycle

## Design Goals

* Simple to use (easy to add new things to be backed up)
* Reliable (Set it and forget it)
* Minimal (Create and backups using google drive. Do this well)

## Building

`make clean dep all`

## Setup

Build the tool and put it in your path `~/bin/` is not a bad place to put it.

Run `backup init` this will prompt you for google api client ids and secrets. These can be obtained by following the wizard here:
<https://console.developers.google.com/start/api?id=drive>

Run `backup add <path to be backed up>` this will add a folder to be backed up. You can provide include and exclude regex rules that files must match

When you want to create a backup run `backup create` to generate a backup and store it in google drive

When you need to restore a backup download the file from google drive and run `backup restore <path to backup file>` and your files should be restored.

## TODO

* Pattern in file name
* Zip path fixing
* Restore from google drive directly
* Test restore from google drive
* Log last run and make create not do anything if not long enough
* Document adding to start up
* Document internals (config locations etc...)
* Blog post
* Go lang tutorial based on this project