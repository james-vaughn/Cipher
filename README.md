## DNS Traffic Monitor

This project monitors DNS traffic that it comes across and triggers emails to be sent if
the volume of traffic exceeds a configures amount of packets within a given timespan.

To configure this app, add the following config under appsettings.config

```
{
   "dns" : {
     "interface" : "eth0",
     "cutoffDuration" : 5, //minutes
     "triggerThreshold" : 500, //packet count
     "durationBetweenTriggers" : 1 //minutes
   },
   "email" : {
     "username" : "user",
     "password" : "password",
     "to"       : "to@gmail.com",
     "smtpServer" : "smtp.gmail.com:587",
   }
 }
 ```