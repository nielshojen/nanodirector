# nanoDirector

nanoDirector is a rewrite of the amazing [MDMDirector](https://github.com/mamdirector/mdmdirector). I had to do a Proof of Concept for work, and wanted a component that could talk to [NanoMDM](https://github.com/NanoMDM/nanomdm) directly. I also wanted to use MySQL as a backend so we could keep all the nano components on the same DB, and had a need for some AccountCreation functionality.

NOTE that this is my first time working with go, so this is probably riddled with bad code and poor choices. Use at your peril!

## Usage

nanoDirector is a compiled binary and is configured using flags.

Requirements:

* Redis for the scheduled checkin queue
* MySQL database for storing device information
* (Recommended) Signing certificate for signing profiles
* (STRONGLY recommended) Load balancer/proxy to serve and terminate TLS for nanoDirector


### NanoMDM Setup

You must set the `-command-webhook-url` flag on NanoMDM to the URL of your nanoDirector instance (with the addition of `/webhook`).

```
-command-webhook-url=https://nanodirector.company.com/webhook
```

### Flags

- `-cert /path/to/certificate` - Path to the signing certificate or p12 file.
- `-clear-device-on-enroll` - Deletes device profiles and install applications when a device enrolls (default "false")
- `-db-host string` - **(Required)** Hostname or IP of the MySQL instance
- `-db-max-idle-connections int` - Maximum number of database connections in the idle connection pool (default -1, not set, uses the default for sql Go package)
- `-db-max-connections int` - Maximum number of database connections (default 100)
- `-db-name string` - **(Required)** Name of the database to connect to.
- `-db-password string` - **(Required)** Password of the DB user.
- `-db-port string` - The port of the MySQL instance (default 5432)
- `-db-sslmode` - The SSL Mode to use to connect to MySQL (default "false")
- `-db-username string` - **(Required)** Username used to connect to the MySQL instance.
- `-redis-host string` - Hostname of your Redis instance (default "localhost").
- `-redis-port string` - Port of your Redis instance (default 6379).
- `-redis-password string` - Password for your Redis instance (default is no password).
- `-debug` - Enable debug mode. Does things like shorten intervals for scheduled tasks. Only to be used during development.
- `-enrollment-profile` - Path to enrollment profile.
- `-enrollment-profile-signed` - Is the enrollment profile you are providing already signed (default: false)
- `-escrowurl` - HTTP(S) endpoint to escrow erase and unlock PINs to ([Crypt](https://github.com/grahamgilbert/crypt-server) and other compatible servers).
- `info-request-interval` - The amount of time in minutes to wait before requesting `DeviceInfo`, `ProfileList`, `SecurityInfo` etc. Defaults to 360.
- `-key-password string` - Password to decrypt the signing key or p12 file.
- `-loglevel string` - Log level. One of debug, info, warn, error (default "warn")
- `-logformat-format` - Log format. Either `logfmt` (the default) or `json`.
- `-nanomdmapikey string` - **(Required)** NanoMDM Server API Key.
- `-nanomdmurl string` - **(Required)** NanoMDM Server URL.
- `-once-in` - Number of minutes to wait before queuing an additional command for any device which already has commands queued. Defaults to 60. Ignored and overridden as 2 (minutes) if --debug is passed.
- `-password string` - **(Required)** Password used for basic authentication
- `-port string` - Port number to run nanoDirector on. (default "8000")
- `-prometheus` - Enable Prometheus metrics. (default false)
- `-push-new-build` - Re-push profiles if the device's build number changes. (default true)
- `-scep-cert-issuer` - The issuer of your SCEP certificate (default: "CN=NanoMDM,OU=NanoMDM SCEP CA,O=NanoMDM,C=US")
- `-scep-cert-min-validity` - The number of days at which the SCEP certificate has remaining before the enrollment profile is re-sent. (default: 180)
- `-sign` - Sign profiles prior to sending to NanoMDM. Requires `-cert` to be passed.
- `-signing-private-key string` - Path to the signing private key. Don't use with p12 file.
- `-admin-user-username string` - Username for the admin user creted at enrollment (default "administrator").
- `-admin-user-password string` - Password for the admin user creted at enrollmen (default "password").


## Todo

- Dynamically setting admin password and escrow to baackend (eg. Crypt)
- MySQL SSL mode

### Documentation

- Changes made to this fork needs to be fixed in the docs

### App

- Maybe a GUI

## Contributing

- File issues
- Open Pull Requests
