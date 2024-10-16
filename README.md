# ADR-ER

> a friendly little thing for managing architectural decision records

this is a ~~super~~ pretty simple cli app for working with ADRs using [the Nygard format](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions.html).

ADRs are stored in a local directory that can/should be tracked along with the code:

under the current working directory, the application will look for a viable
directory named adr, .adr, or architectural-decision-records.

"viable": the contents of the candidate directories are scanned. a "viable" one has only markdown files that fit the
ADR naming convention. eg: `0003-security-audit.md`, `0007-team-expansion.md`, etc. Subdirectories are allowed in the ADR dir, but
its immediate files must be only ADRs.

This mechanism is to prevent the app from writing markdown files into source-code directories

## Install

### From Releases 

this repo publishes its release archives here https://github.com/therealkevinard/adr-er/releases.    
download the relevant archive for your platform and extract the contained binary to your path. 

#### From Source
`go install github.com/therealkevinard/adr-er@latest` 


## Usage 

### Creating an ADR

Run `adr-er create` to make a new ADR. this opens a tui form to fill in the deets.   
Only `Title` is strictly required, but all fields are recommended.

When you're all done, the ADR file will be created with an incremented sequence number.

Templating is a one-way job, but the file can be edited all you want as text once it's created.

![demo-create.gif](doc/demo/demo-create.gif)

### Viewing exising ADRs

Use `adr-er view` to open a handy little navigator for existing ADR files.  
The app has simple keyboard navigation and supports filtering the list. for tall files, the viewer is scrollable - you
just have to tab/arrow over to the viewer to scroll (otherwise, you're scrolling the file list, yknow?)

![demo-view.gif](doc/demo/demo-view.gif)