# gome-schedule Release Log

### Merged & not released
- Add: badger DB memory settings for small memory footprints
- Add: load config from environment for easier docker runtime support
- Add: debug logging to troubleshoot schedules
- Add: new state schedule supports "daily", "weekday", "weekend" to group days in a schedule
- Change: dockerfile for compose
- Fix: database release lock on db error
- Fix: initialize database make dir full recursive path

### 1.0.0 - untagged
* initial release

### 1.0.1
* Updated schedules to map
* fixed bug in set power state from post to put
* workaround database closures in for each schedule loop