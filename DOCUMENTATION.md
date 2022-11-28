# imscli
CLI tool to login and interact with the IMS API.

## Usage

The results of the command will be written to *stdout*, allowing it to be redirected to a file or an env var.

Any other output like verbose output or errors will be sent to *stderr* to not interfere with the token.

The command will return 0 in case of success or 1 in case of an error.

## Subcommands
### Authorize

imscli login will negotiate an ***access token*** with IMS following the specified ***flow***.

#### imscli authorize user (standard Authorization Code Grant Flow)

This command will launch a browser and execute the normal OAuth2 flow done by users when log into IMS to use a service.

#### imscli authz service (an IMS specific flow similar to the Client Credentials Grant Flow).

The imscli client will exchange client credentials and an additional service token to obtain the access token.

It is used to access an "Application", an Adobe API exposed through Adobe I/O Gateway.

#### imscli authz jwt (JWT Bearer Flow).

This command will build a JWT will all specified claims, sign it with a private key and exchange it for an access token.

It is used for "Adobe I/O" integrations.

### Completion

Generate a script to enable autocompletion for imscli. 

To configure bash, add the following line to your .bashrc (or alternative config file):

    eval "$(imscli autocompletion bash)"

### Profile

Provided a user's access token, gathers the user profile.

### Organizations

Provided a user's access token, gathers the user organizations.

### Exchange

Performs the "cluster authorization token exchange flow", exchanging a valid access token for another access token for the
same user in a different IMS Org.

### Validate

Validates a token using the IMS API.

### Invalidate

Invalidates a token using the IMS API.


## Configuration

Usage is defined by what the CLI libraries [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) support.

There are three sources of parameters, listed from higher to lower priority.

The three source of parameters can be combined, priority will be taken in account in case of overlap.

#### CLI flags

Direct CLI parameters, with extensive documentation executing each subcommand with --help.
```
imscli login user --scopes AdobeID,openid,session
```

#### Environment variables

Each parameter can be provided using the flag name and the IMS_ suffix.
```
IMS_SCOPES="AdobeID,openid,session" imscli login user
```

#### Configuration files

There are two ways to specify config files:

- Default configuration file, useful for defaults when used directly by users, not thought for automation.
    - The file should be stored in ~/.config/
    - The file should be named `imscli.extension`, where extension has to match the file format.
    - Supported formats among others include YAML and JSON.
- Directly specified configuration file with the -f flag.
```
user@host$ cat ~/.config/imscli.yaml
scopes:
  - AdobeID
  - openid
  - session

user@host$ imscli login user
```

