# imscli
CLI tool to obtain IMS tokens and interact with the IMS API.

## Usage

The results of the command will be written to *stdout*, allowing it to be redirected to a file or an env var.

Any other output like verbose output or errors will be sent to *stderr* to not interfere with the token.

The command will return 0 in case of success or 1 in case of an error.

## Subcommands
### Authorize

imscli authorize will negotiate an ***access token*** with IMS following the specified ***flow***.

#### imscli authorize user (standard Authorization Code Grant Flow)

This command will launch a browser and execute the normal OAuth2 flow done by users when logging into IMS to use a service.

#### imscli authorize service (an IMS specific flow similar to the Client Credentials Grant Flow).

The imscli client will exchange client credentials and an additional service token to obtain the access token.

It is used to access an "Application", an Adobe API exposed through Adobe I/O Gateway.

#### imscli authorize jwt (JWT Bearer Flow).

This command will build a JWT with all specified claims, sign it with a private key and exchange it for an access token.

It is used for "Adobe I/O" integrations.

#### imscli authorize pkce (Authorization Code Grant Flow with PKCE)

Like the user command, it uses the Authorization Code Grant Flow but with Proof Key for Code Exchange (PKCE). In IMS, PKCE is mandatory for public clients and recommended for private clients.

#### imscli authorize client (Client Credentials Grant Flow)

Exchanges client credentials (client ID + secret) and scopes directly for an access token, without user interaction.

### Profile

Provided a user's access token, gather the user profile.

### Organizations

Provided a user's access token, gather the user organizations.

### Exchange

Performs the "cluster authorization token exchange flow", exchanging a valid access token for another access token for the
same user in a different IMS Org.

### Validate

Validates a token using the IMS API.

### Invalidate

Invalidates a token using the IMS API.

### Decode

Decodes a JWT token locally, printing the header and payload without contacting IMS.

### Refresh

Refreshes an access token using a refresh token.

### Admin

Administrative operations using a service token:

- **admin profile**: Retrieve a user profile using a service token, client ID, guid and auth source.
- **admin organizations**: Retrieve organizations for a user using a service token.

### DCR (Dynamic Client Registration)

Register a new OAuth client using Dynamic Client Registration.

## Configuration

Usage is defined by what the CLI libraries [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) support.

There are three sources of parameters, listed from higher to lower priority.

The three sources of parameters can be combined, priority will be taken into account in case of overlap.

#### CLI flags

Direct CLI parameters, with extensive documentation executing each subcommand with --help.
```
imscli authz user --scopes AdobeID,openid,session
```

#### Environment variables

Each parameter can be provided using the flag name and the IMS_ prefix.
```
IMS_SCOPES="AdobeID,openid,session" imscli authorize user
```

#### Configuration files

There are two ways to specify config files:

- The default configuration file, useful for defaults when used directly by users. This is not practical for automation.
    - The file should be stored in ~/.config/
    - The file should be named `imscli.extension`, where the extension has to match the file format.
    - Supported formats among others include YAML and JSON.
- Directly specified configuration file with the -f flag.
```
user@host$ cat ~/.config/imscli.yaml
scopes:
  - AdobeID
  - openid
  - session

user@host$ imscli authorize user
```

