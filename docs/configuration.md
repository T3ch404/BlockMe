# Configuration
Configuration is done through environment variables. The following subsections 
detail how to configure BlockMe to work in your environment.

## Database
Setting up a database enabled persistent IP storage to ensure that blocked IPs 
stay blocked.

**DB_TYPE** | `postgres` or `sqlite` | Default: `sqlite`  
> Choosing Postgres makes DB_HOST, DB_PORT, DB_NAME, DB_USER, & DB_PASS required 
> in order to connect to your postgres database. Sqlite database files are 
> stored the database subdirectory of your working directory (`/app/database/` 
> for Docker containers)

**DB_HOST** | IP address or FQDN | Default: N/A
> Routable address of your database server. The DB_HOST environment variable has 
> no default and will result in the program exiting with an error code if the 
> DB_TYPE has been set to postgres and no DB_HOST is specified.

**DB_PORT** | 1 - 65535 | Default: `sqlite`
> TCP port that your database server is listening on. This variable will default
> to the default port of the given DB_TYPE (Postgres → 5432, MariaDB → 3306).

**DB_NAME** | Existing database name | Default: `blockme`
> Specify a database name if you want to use a specific database name. If no 
> database exists with the given name, BlockMe will attempt to create it.

**DB_USER** | Existing database user | Default: None
> The DB_USER environment variable has no default and will result in the program 
> exiting with an error code if the DB_TYPE has been set to postgres and no 
> DB_USER is specified.

**DB_PASS** | Existing database user password | Default: None
> The DB_PASS environment variable has no default and will result in the program
> exiting with an error code if the DB_TYPE has been set to postgres and no
> DB_PASS is specified.

## Reset
The reset functionality allows user's to unblock themselves with a generated 
UUID.

**I_AM_A_LITTLE_BITCH** | `true` or `false` | Default: `false`
> Setting I_AM_A_LITTLE_BITCH to `true` enables the reset functionality.

**I_AM_A_LITTLE_BITCH_CRON** | Five field Cron schedule | Default: `1 * * * *`
> Defines how often the reset token is rotated and sent to the webhook.

**I_AM_A_LITTLE_BITCH_WEBHOOK_URL** | Existing webhook URL (Discord) | Default: None
> Defines where reset tokens will be sent to.

## Misc

**IGNORE_LIST** | Comma separated list of IPs | Default: ""
> Adding an IP address to the IGNORE_LIST environment variable prevents that IP 
> from being blocked. All other logic works the same, and BlockMe will still 
> return a "blocked" message but the address will not be added to the block list