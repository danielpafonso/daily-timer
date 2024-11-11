# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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