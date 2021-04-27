# ip-logger
Log peoples IP addresses if they visit a link from you


# Project plan

Create an IP logger similar to that one that I used recently but honested on my own url. 


### 1
Concepts:

lnk.certainlyawesome.com/a1KdlKoW

a1KdlKoW is a random value, that identifies the person that Im issuing it to. 

Should redirect the user to the target after a very short amount of time, but uses javascript in order to collect screen res, browser type, while the server gets IP. 


### 2 

URLs should become invalid after either a set amount of time, or number of visits. So that we get more accurate targettings. Although setting it so its unlimited might be fun to make graphs of how it gets shared. 

### 3

Tool should be interactable through command line, so I can easily issue a command to create a new url and immediately paste it. 

E.g 

`iplog https://google.co.nz` 

With sane defaults for number of visits/time (which will have to be worked out with use.)

Use UNIX sockets with an active agent. 

Some more examples

`iplog -a 4 https://google.co.nz` // Only allows 4 accesses before db entry is removed

`iplog https://google.co.nz Andy` // Create a URL for a specific person

`iplog ls` //Show all entries, shows brief num visits/if visited/ expiry info

`iplog ls Andy` // For specific person SELECT * FROM x where?

`iplog ls -la` // Shows all specific information

`iplog rm Andy`

### 4 

If an arbitrary url is requested 404? Maybe?