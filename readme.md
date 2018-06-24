# Eye

Eye is a program that sends you an email when the free
space in the filesystem is less than the minimum.

## Build

Build Eye with make:

    make

## Package

Currently, only Debian packaging is supported:

    make package-for-debian

This will leave the package in the current directory.

## Install

### Debian

If you want to install the Debian package, just use ```dpkg```:

     sudo dpkg -i eye_1.0-1.deb

Installing will configure Eye as a ```systemd``` service, enable it, and start
it up.

### From sources

Installing from sources is done with ```make```:

    make install

You'll most likely need need ```sudo```:

    sudo make install

This will:
* Copy the built binary.
* Copy the configuration file. It doesn't override it if it already exists.
* Configure the service with systemd and start it.

Check ```Makefile``` for default files and directories.

## Uninstall

### Debian

Just use ```dpkg```

    dpkg -r eye

### From sources

Uninstall with ```make```:

    make uninstall

You'll most likely need ```sudo```:

    sudo make uninstall

## Operation  

Use systemd to start and stop:

    systemctl start eye.service
    systemctl stop eye.service

### Configuration

The configuration is stored in ```/etc/eye/eye.cfg```.

Example configuration:

    Threshold = 4294967296
    SensePeriod = "5m"

    [Email]

    Host = "smtp.gmail.com"
    Port = 587
    Username = "username@gmail.com"
    Password = "password"
    From = "username@gmail.com"
    To = "username@gmail.com"
