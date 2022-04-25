## Implementation strategy

I added ttl to the struct value of the map that holds the order ID and history. So, I created a cron job using go/cron library and running under go routine to delete the order that it's expiredAt time is less than current time in unix()