# assetfinder

Find domains and subdomains potentially related to a given domain.


## Install

If you have Go installed and configured (i.e. with `$GOPATH/bin` in your `$PATH`):

```
go install github.com/0x1Jar/assetfinder@latest
```

Alternatively, you can [download a pre-compiled release for your platform](https://github.com/0x1Jar/assetfinder/releases).
To make it easier to execute you can put the binary in your `$PATH`.

## Usage

```
assetfinder [--subs-only] [--sources <source1,source2...>] <domain>
```

If `<domain>` is not provided, assetfinder will read domains from stdin (one domain per line).

**Flags:**
* `--subs-only`: Only include subdomains of the search domain.
* `--sources <source1,...,sourceN>`: Specify a comma-separated list of sources to use. If not provided, all available sources are used.
  Available sources: `crtsh`, `certspotter`, `hackertarget`, `threatcrowd`, `wayback`, `dns.bufferover.run`, `facebook`, `virustotal`, `findsubdomains`, `urlscan`.

## API Key Configuration

Some data sources require API keys or tokens to be set as environment variables.

*   **Facebook:**
    *   `FB_APP_ID`: Your Facebook App ID.
    *   `FB_APP_SECRET`: Your Facebook App Secret.
    *   (Obtain from [https://developers.facebook.com/](https://developers.facebook.com/))
    *   Be mindful of your app's rate limits.
*   **VirusTotal:**
    *   `VT_API_KEY`: Your VirusTotal API Key.
    *   (Obtain from [https://developers.virustotal.com/reference](https://developers.virustotal.com/reference))
*   **FindSubdomains (Spyse):**
    *   `SPYSE_API_TOKEN`: Your Spyse API Token.
    *   (Obtain from [https://spyse.com/apidocs](https://spyse.com/apidocs) - the free version provides limited access)

To set these environment variables, you can use the `export` command in your shell (for Linux/macOS):
```bash
export FB_APP_ID="your_app_id"
export FB_APP_SECRET="your_app_secret"
export VT_API_KEY="your_vt_api_key"
export SPYSE_API_TOKEN="your_spyse_api_token"
```
For Windows, you can use the `set` command in Command Prompt or `Set-Item` in PowerShell.
You may want to add these commands to your shell's startup file (e.g., `.bashrc`, `.zshrc`) to set them automatically.

## Sources

Please feel free to issue pull requests with new sources! :)

### Implemented
* crt.sh
* certspotter
* hackertarget
* threatcrowd
* wayback machine
* dns.bufferover.run
* facebook
    * Needs `FB_APP_ID` and `FB_APP_SECRET` environment variables set (https://developers.facebook.com/)
    * You need to be careful with your app's rate limits
* virustotal
    * Needs `VT_API_KEY` environment variable set (https://developers.virustotal.com/reference)
* findsubdomains
    * Needs `SPYSE_API_TOKEN` environment variable set (the free version always gives the first response page, and you also get "25 unlimited requests") — (https://spyse.com/apidocs)

### Sources to be implemented
* http://api.passivetotal.org/api/docs/
* https://community.riskiq.com/ (?)
* https://riddler.io/
* http://www.dnsdb.org/
* https://certdb.com/api-documentation

## TODO
* Implement more sources (see "Sources to be implemented" below).
