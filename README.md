# TCP Scanner

TCP port scanner, spews SYN packets asynchronously, scanning entire Internet in under 5 minutes.


## Usage 

```bash
go install github.com/Guaderxx/scanner

# tcp mode
scanner -i 0.0.0.0/0 -p 1-65535 -m t

# syn
sudo scanner -i 0.0.0.0/0 -p 1-65535
```
