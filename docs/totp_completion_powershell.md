## totp completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	totp completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
totp completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string   The path to the TOTP accounts configuration [$TOTP_CONFIG] (default "~/.config/totp/data")
```

### SEE ALSO

* [totp completion](totp_completion.md)	 - Generate the autocompletion script for the specified shell

