# ADR-ER

> a friendly little thing for managing architectural decision records

this is a dead-simple cli app for working with ADRs using [the Nygard format](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions.html).

the tui leans heavily on [charm](https://charm.sh/). 

### CLI Help 
```text
NAME:
   adr-er - a friendly little thing for managing architectural decision records

USAGE:
   adr-er [global options] command [command options]

COMMANDS:
   create, new  create a new adr document
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dir value, -d value  root directory to store adr files.

      if empty: 
        the application will search for a viable directory according to some conventions.  
        directories in CWD named "architectural-decision-records", "adr", or ".adr" will be checked. 
        we will set --dir to the first in the the list that is  
        a) empty, or b) holds only adr files and optionally subdirectories.
      if provided:
        the application will not validate contents - we'll trust your judgement
        the special value "-" can be used to indicate stdout
   --help, -h  show help
```

### Creating an ADR
![demo-create.gif](doc/demo/demo-create.gif)