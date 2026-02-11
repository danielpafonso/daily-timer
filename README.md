# Daily Timer

A terminal timer for Daily Standups meetings

<p align="center">
  <img src="/docs/hero.png" alt="Hero screenshot"/>
</p>

## Command-Line Options

| Flag | Description                                 | Default       |
| ---- | ------------------------------------------- | ------------- |
| `-c` | Path to configuration file                  | `config.json` |
| `-m` | Statistics storage mode (`sqlite` or `csv`) | `sqlite`      |
| `-v` | Display version and exit                    | -             |

## Configurations

| Field             |  Type   | Description                                                                   |
| ----------------- | :-----: | ----------------------------------------------------------------------------- |
| time              | integer | Start time value for timer mode or soft limit in stopwatch mode               |
| warning           | integer | Time value where the warning color will be displayed                          |
| participants      |  array  | List of participants for the daily                                            |
| randomOrder       | boolean | When `true`, randomize the participant order                                  |
| stopwatch         | boolean | Timer mode: `true` for stopwatch, `false` for countdown timer                 |
| addTemp           | boolean | When `true`, temporary participants are included in statistics                |
| ignoreZeroTime    | boolean | When `true`, participants with zero elapsed time are excluded from statistics |
| stats.display     | boolean | Toggle display of past dailies statistics on participants list                |
| stats.lastDailies | integer | Number of past dailies to include in statistics calculations                  |

Example:

```json
{
  "time": 20,
  "warning": 10,
  "participants": ["John", "Marcus", "Abigal"],
  "randomOrder": true,
  "stopwatch": true,
  "addTemp": false,
  "ignoreZeroTime": true,
  "stats": {
    "display": true,
    "lastDailies": 30
  }
}
```

## Storage Formats

- SQLite (Default)
  - Stores data in `daily.db` file
  - More efficient for large datasets
  - Supports complex queries
  - Recommended for most use cases

- CSV
  - Stores data in `daily.csv` file
  - Human-readable format
  - Easy to import into spreadsheets
  - Useful for manual data analysis

Switch between formats using the `-m` flag. Data is not automatically migrated between formats.
