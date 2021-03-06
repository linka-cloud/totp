## totp completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(totp completion bash)

To load completions for every new session, execute once:

#### Linux:

	totp completion bash > /etc/bash_completion.d/totp

#### macOS:

	totp completion bash > /usr/local/etc/bash_completion.d/totp

You will need to start a new shell for this setup to take effect.


```
totp completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string   The path to the TOTP accounts configuration [$TOTP_CONFIG]
```

### SEE ALSO

* [totp completion](totp_completion.md)	 - Generate the autocompletion script for the specified shell

