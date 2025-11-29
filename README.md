# Goredacted

Go CLI programm to redact and anonymize sensitive information.

## Getting Started

1. Download the latest version of Goredacted from the Release Page
2. Make it executable with `chmod +x goredacted`
3. Move to $PATH: `mv goredacted /usr/local/bin/goredacted`

Once done initiliaze `goredacted` with (use `-s` for custom seed path)

```bash
goredacted init
```

> [!IMPORTANT]
> Each time `goredacted init` is run it creates and overwrites the seed file.
> Since the seed is used for consistent data masking using `init` multiple times can lead to inconsistent placeholders.

## Usage

Goredacted has multiple ways to run the main data masking logic.

### Parse

This mode is the simplest way to run Goredacted and is mainly used for bash scripts or as a utility.

```bash
echo "private@sensitive.domain.com" | goredacted parse
```

**Output:**

````bash
john@test.example.org
```

The `parse` subcommand uses `stdin` to receive and parse the data and outputs into `stdout` the result.

### Clipboard Hook

In the **Clipboard Hook** mode Goredacted runs in a loop where it gets the clipboard content parses it and writes it back into the clipboard.

```bash
goredacted start
````

### HTTP Listener

In the **HTTP Listener** mode Goredacted runs behind a http endpoint and responds to `GET` request to the `/parse` endpoint.
Internally it fetches the clipboard content parses it, processes it and sends it back as json inside of the `output` field.

```bash
goredacted start -l 127.0.0.1:5678
```

## Configuration

### Flags

| Flag | Example                                    | Note                                                                                                                                      |
| :--- | :----------------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------- |
| `-d` | `goredacted start -d`                      | Detach and spawn a background service; Only available for the `start` subcommand                                                          |
| `-l` | `goredacted start -l 127.0.0.1:5678`       | Specify listen address for http server; Only available for the `start` subcommand                                                         |
| `-s` | `goredacted parse -s /path/to/seed.file`   | Path of seed file to be used for data masking, needed to provide consistent masking (default: `/home/user/.config/grdct/tcdrq.grdct`)     |
| `-c` | `goredacted parse -s /path/to/config.json` | Path of config json file for configuring exludes and includes (default: `/home/user/.config/grdct/config.json`; [see more](#config-file)) |

### Parsers

Available parsers are:

| Parser ID | Note                                                              |
| :-------- | :---------------------------------------------------------------- |
| `domain`  | Parses domains with valid tlds whilst preserving domain structure |
| `ip`      | Parses IPv4 and IPv6 addresses                                    |
| `email`   | Parses the local email part (`local.part@domain.com`)             |

### Config File

> [!NOTE]
> Use `-s` flag to modify default config path

**Example:**

```json
{
	"include": [
		{
			"parser": "domain",
			"regex": "\\b(?:[a-zA-Z0-9-]{1,63}\\.)*domain\\.com\\b"
		},
		{
			"parser": "domain",
			"regex": "\\b(?:[a-zA-Z0-9-]{1,63}\\.)*example\\.org\\b"
		}
	],
	"exclude": [
		{
			"parser": "ip",
			"regex": "\\b\\d\\.\\d\\.\\d\\.\\d\\b"
		}
	]
}
```

This example modifies the `domain` parser to only parse `*.domain.com` and `*.example.org`.
Aswell as excluding any `X.X.X.X` ip from the masking process (example: `1.1.1.1` will not be parsed).

## Integrations

### Chrome Extension

When running Goredacted as **http listener** you can use the following chrome extension.

First you need to **download** the `integration/chrome-extension` folder,
for that we'll be using `git clone https://github.com/codeshelldev/goredacted` to clone the entire repository to disk.

To be able to add this chrome extension you will need to head to `chrome://extension` or wherever your extensions are managed and **enable** the `Developer Mode`.
Now you can **load a folder as an extension**.

Load the newly downloaded `integration/chrome-extension` in your chromium based browser and go to the **extensions settings**.

![Options](./imgs/integrations/chrome-extension/options.png)

Here you can **add disabled and enabled domains** and also the **http endpoint** that the extension should request.
Once done any domain that ends with a disabled or enabled will be treated respectively.
