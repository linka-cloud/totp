## totp completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	totp completion fish | source

To load completions for every new session, execute once:

	totp completion fish > ~/.config/fish/completions/totp.fish

You will need to start a new shell for this setup to take effect.


```
totp completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string   The path to the TOTP accounts configuration [$TOTP_CONFIG]
```

### SEE ALSO

* [totp completion](totp_completion.md)	 - Generate the autocompletion script for the specified shell

