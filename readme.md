# Eye

Eye is a program that sends you an email when the free
space in the filesystem is less than the minimum.

## Configuration

The configuration is stored in ```/etc/eye/eye.cfg```.

Example configuration:

    Threshold = 4294967296

    [Email]

    Host = "smtp.gmail.com"
    Port = 587
    Username = "username@gmail.com"
    Password = "password"
    From = "username@gmail.com"
    To = "username@gmail.com"
