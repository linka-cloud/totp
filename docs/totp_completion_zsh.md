## totp completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions for every new session, execute once:

#### Linux:

	totp completion zsh > "${fpath[1]}/_totp"

#### macOS:

	totp completion zsh > /usr/local/share/zsh/site-functions/_totp

You will need to start a new shell for this setup to take effect.


```
totp completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string   The path to the TOTP accounts configuration [$TOTP_CONFIG]
```

### SEE ALSO

* [totp completion](totp_completion.md)	 - Generate the autocompletion script for the specified shell

