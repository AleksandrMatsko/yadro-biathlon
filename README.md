# yandro-biathlon

This project is a [test task](task/README.md) for internship to YADRO

## Build and run

To build executable run:
```shell
make build
```

This will create `biathlon-reporter` executable. Run it with proper command line arguments.

## Command line arguments

`--config`

> Required. Path to config file.

`--print-config`

> Optional. If present, prints provided config to stdout.

`--events`

> Required. Path to file with incoming events.

`--report`

> Optional. Path to file to save report. Default is `report.txt`.

# Inconsistencies and omissions

Unfortunately during task completion I have found some inconsistencies between task condintions and given examples. 
All mistamtchs listed below were found in task version located in [`task` directory](task/README.md)

## Firing lines

According to task description `firingLines` field in configuration file is the number of firing lines **per** lap. 
But in the [example events](task/sunny_5_skiers/events) given with task description `firingLines` is the total number of firing lines.

I've decided to use `firingLines` as total number of firing lines.

## Calculating times and speed

Unfortunately there is no clear explanation about calculating the times. So here are the rules I'm using:

### Total time

If the competitor has finished, it is time interval between:
- scheduled start time
- timestamp of ending the last lap

Otherwise: **NotFinished**

### Penalty lap time

Sum of time intervals between:
- entering penalty laps
- leaving penalty laps

### Average speed over penalty laps

It is `(Penalty lap time) / (penaltyLapsCount * penaltyLen)`

### Main lap time

It is calculated if the lap was completed.

For the first lap it is time interval between:
- scheduled start time
- timestamp of ending first lap

For other laps:
- timestamp of ending previous lap
- timestamp of ending current lap

### Average speed for each main lap

It is `(Main lap time) / lapLen`



