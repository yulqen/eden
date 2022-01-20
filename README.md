# eden

`eden` is a command line tool for creating and manipulating daily log notes. It started life as a series of different bash script that did various tasks, but have been brought together as one tool.

Enjoy using this awesome tool!

## Functionality

### Adding an entry to the journal

A *journal* is a database, or registry, of *entries*. An entry is like a line in a daily journal, that is date/time-stamped to the moment of creation. In its most basic form, an entry is just a line of text. Associated with the timestamp, that in itself could be considered useful. You can record what you are doing at any time of the day (when you are in the terminal, anyway), and then output a list of each entry for a particular day. That is the basic functionality we are going to go for first. Things like *tags* and stuff come later.

To add a new entry to the journal:

`eden add "ENTRY"`

The program will confirm that the entry was saved and that's it.
