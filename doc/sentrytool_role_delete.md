## sentrytool role delete

delete Sentry roles

### Synopsis


delete Sentry roles.

```
sentrytool role delete
```

### Options

```
      --force          force deletion
  -m, --match string   regexp matching role
```

### Options inherited from parent commands

```
  -C, --component string   sentry client component
      --config string      config file (default is $HOME/.sentrytool.yaml)
  -H, --host string        hostname for Sentry server (default "localhost")
  -J, --jstack             show Java stack on for errors
  -P, --port string        port for Sentry server (default "8038")
  -U, --username string    user name (default "akolb")
  -v, --verbose            verbose mode
```

### SEE ALSO
* [sentrytool role](sentrytool_role.md)	 - Sentry roles manipulation

###### Auto generated by spf13/cobra on 14-Dec-2016
