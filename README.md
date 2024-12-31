# Daily Timer

A terminal timer for Daily Standups meetings

<p align="center">
  <img src="/docs/hero.png" alt="Hero screenshot"/>
</p>

## Configurations

| Field             |  Type   | Description                                                                                         |
| ----------------- | :-----: | --------------------------------------------------------------------------------------------------- |
| time              | integer | Start time value for timer mode or soft limit in stopwatch mode                                     |
| warning           | integer | Time value where the warning color will be showed                                                   |
| participants      |  array  | List of participants for the daily                                                                  |
| randomOrder       | boolean | If `true` randomize the participant order                                                           |
| stopwatch         | boolean | Select time mode, if `true` the script will work as a stopwatch else will run as a count down timer |
| addTemp           | boolean | Flag that indicates if temporary participants should be written to statistics file                  |
| stats.display     | boolean | Toggle display of past dailies statistics on participants list                                      |
| stats.lastDailies | integer | Number of past dailies used to compute statistics                                                   |

Example:

```json
{
  "time": 20,
  "warning": 10,
  "participants": ["John", "Marcus", "Abigal"],
  "randomOrder": true,
  "stopwatch": true,
  "addTemp": false,
  "stats": {
    "display": true,
    "lastDailies": 30
  }
}
```
