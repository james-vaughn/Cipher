## DNS Traffic Monitor

This project monitors DNS traffic that it comes across and triggers emails to be sent if
the volume of traffic exceeds a configures amount of packets within a given timespan.

To configure this app, add the following config under appsettings.config

```
{
  "interface" : "enp4s0",
  "dns" : {
    "cutoffMinutes" : 5,
    "triggerThreshold" : 100,
    "minutesBetweenTriggers" : 5
  },
  "email" : {
    "to" : "to@gmail.com",
    "from" : "from@gmail.com",
    "password" : "pass",
    "smtpServerHostname" : "smtp.gmail.com",
    "smtpServerPort" : 587
  }
}
 ```