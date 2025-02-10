# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Add

- Add Flash mode to participant list
- Loop Jump, making participant list circular
- Reactability with windows size change

### Fix

- Make clean target failing when build folder is not present

### Remove

- Debug function, not used

## [1.1.0]

### Add

- Command line flag, `-m`, to change file interface

### Change

- Move SQLite and csv functions to dedicated GO plugins
- Project structure with the use of GO plugins
- Makefile to build and generate releases

## [1.0.1]

### Add

- Set minimal width to input participant widget, equal to timer size

### Change

- Exiting with `ctrl+c` wont write statistics to file

## [1.0.0]

### Add

- Add keybing/function to add Temp participant to current list
- Field `AddTemp` to configuration that controls if a temp participant is written to stats files
- Temp participant input widget
- Update keybings to have possibilities of opening Temp input widget, accept or cancel addition

### Fix

- (csv mode) Crash when participant is in stat file and not on configuration list

## [0.0.3]

### Fix

- Program crash by dividing by zero, when new elements were added to participants
- Stop resetting first and last timer participant when moving up or down respectively

## [0.0.2]

### Added

- Add command line flag, `-v`, to show binary version
- Change Help menu text, listing keys

### Change

- Split binaries into using Sqlite and other using csv, speeding dev/testing

## [0.0.1]

### Added

- Timer or stopwatch with changing colours, warning when near and past desired limit
- Read participant list and randomize order if option is set
- Select participant with cursor, enabling toggling the timer for each entry
- Create `sqlite` database for statistics, and create basic crud operations
- Add writing participants timers when application is closed
- Add toggling of participants as active/inactive, and disable database update of inactive users
- Inactive participants have their timer display as `--:--`
- Cursor will skip over inactive participants when timer is running
- Timer will only start on active participants

---

[unreleased]: https://github.com/danielpafonso/daily-timer/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/danielpafonso/daily-timer/releases/tag/v1.1.0
[1.0.1]: https://github.com/danielpafonso/daily-timer/releases/tag/v1.0.1
[1.0.0]: https://github.com/danielpafonso/daily-timer/releases/tag/v1.0.0
[0.0.3]: https://github.com/danielpafonso/daily-timer/releases/tag/v0.0.3
[0.0.2]: https://github.com/danielpafonso/daily-timer/releases/tag/v0.0.2
[0.0.1]: https://github.com/danielpafonso/daily-timer/releases/tag/v0.0.1
