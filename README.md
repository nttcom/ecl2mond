# ecl2mond - ECL2.0 Monitoring Custom meter agent
ECL2.0 Custom meter agent (a.k.a ecl2mond) is a support tool for NTT Communications' Enterprise Cloud 2.0 Monitoring service.

## Purpose

This agent enables you to accumulate and visualize usage data which you want automatically in ECL2.0 Monitoring service. The agent runs on a Linux machine and collects usage data of some kinds of computing resources. Moreover, the agent sends the usage data periodically to ECL2.0 Monitoring service by using Monitoring API.

## Target

Linux servers running on or connect with ECL2.0 environment

## List of collectable usage data

No. |   OS  | Resource Category |                  Meter Name                   |               Description               |     Unit     | Meter Type | Data Source
----|-------|-------------------|-----------------------------------------------|-----------------------------------------|--------------|------------|-------------
 1  | Linux | CPU               | cpu.user.percents                             | CPU Usage (user mode)                   | Percent      | delta      | /proc/stat
 2  | Linux | CPU               | cpu.nice.percents                             | CPU Usage (user mode with low priority) | Percent      | delta      | /proc/stat
 3  | Linux | CPU               | cpu.system.percents                           | CPU Usage (system mode)                 | Percent      | delta      | /proc/stat
 4  | Linux | CPU               | cpu.idle.percents                             | CPU Usage (idle task)                   | Percent      | delta      | /proc/stat
 5  | Linux | CPU               | cpu.iowait.percents                           | CPU Usage (I/O wait)                    | Percent      | delta      | /proc/stat
 6  | Linux | CPU               | cpu.irq.percents                              | CPU Usage (hard irq)                    | Percent      | delta      | /proc/stat
 7  | Linux | CPU               | cpu.softirq.percents                          | CPU Usage (soft irq)                    | Percent      | delta      | /proc/stat
 8  | Linux | CPU               | cpu.steal.percents                            | CPU Usage (stolen time)                 | Percent      | delta      | /proc/stat
 9  | Linux | CPU               | cpu.guest.percents                            | CPU Usage (vCPU for guest OS)           | Percent      | delta      | /proc/stat
 10 | Linux | Disk              | disk.{device name}.reads.completed.count      | Disk Read I/O (completed successfully)  | Count        | delta      | /proc/diskstats
 11 | Linux | Disk              | disk.{device name}.reads.merged.count         | Disk Read I/O (merged)                  | Count        | delta      | /proc/diskstats
 12 | Linux | Disk              | disk.{device name}.reads.sectors.count        | Disk Read Sectors                       | Count        | delta      | /proc/diskstats
 13 | Linux | Disk              | disk.{device name}.reads.milliseconds         | Disk Read milliseconds                  | millisecond  | delta      | /proc/diskstats
 14 | Linux | Disk              | disk.{device name}.writes.completed.count     | Disk Write I/O (completed successfully) | Count        | delta      | /proc/diskstats
 15 | Linux | Disk              | disk.{device name}.writes.merged.count        | Disk Write I/O (merged)                 | Count        | delta      | /proc/diskstats
 16 | Linux | Disk              | disk.{device name}.writes.sectors.count       | Disk Write Sectors                      | Count        | delta      | /proc/diskstats
 17 | Linux | Disk              | disk.{device name}.writes.millisecond         | Disk Write milliseconds                 | millisecond  | delta      | /proc/diskstats
 18 | Linux | Disk              | disk.{device name}.currently.ios.count        | Current Disk I/O                        | Count        | delta      | /proc/diskstats
 19 | Linux | Disk              | disk.{device name}.ios.milliseconds           | Disk I/O milliseconds                   | millisecond  | delta      | /proc/diskstats
 20 | Linux | Disk              | disk.{device name}.weighted.ios.milliseconds  | Weighted Disk I/O milliseconds          | millisecond  | delta      | /proc/diskstats
 21 | Linux | Network           | network.{nwif name}.receive.bytes             | IF Received bytes                       | Byte         | delta      | /proc/net/dev
 22 | Linux | Network           | network.{nwif name}.receive.packets.count     | IF Received packets                     | Count        | delta      | /proc/net/dev
 23 | Linux | Network           | network.{nwif name}.receive.errs.count        | IF Received error packets               | Count        | delta      | /proc/net/dev
 24 | Linux | Network           | network.{nwif name}.receive.drop.count        | IF Received and dropped packets         | Count        | delta      | /proc/net/dev
 25 | Linux | Network           | network.{nwif name}.receive.fifo.count        | IF Received fifo error packets          | Count        | delta      | /proc/net/dev
 26 | Linux | Network           | network.{nwif name}.receive.frame.count       | IF Received frame error packets         | Count        | delta      | /proc/net/dev
 27 | Linux | Network           | network.{nwif name}.receive.compressed.count  | IF Received compressed packets          | Count        | delta      | /proc/net/dev
 28 | Linux | Network           | network.{nwif name}.receive.multicast.count   | IF Received multicast packets           | Count        | delta      | /proc/net/dev
 29 | Linux | Network           | network.{nwif name}.transmit.bytes            | IF Transmitted bytes                    | Byte         | delta      | /proc/net/dev
 30 | Linux | Network           | network.{nwif name}.transmit.packets.count    | IF Transmitted packets                  | Count        | delta      | /proc/net/dev
 31 | Linux | Network           | network.{nwif name}.transmit.errs.count       | IF Transmitted error packets            | Count        | delta      | /proc/net/dev
 32 | Linux | Network           | network.{nwif name}.transmit.drop.count       | IF Transmitted and dropped packets      | Count        | delta      | /proc/net/dev
 33 | Linux | Network           | network.{nwif name}.transmit.fifo.count       | IF Transmitted fifo error packets       | Count        | delta      | /proc/net/dev
 34 | Linux | Network           | network.{nwif name}.transmit.coll.count       | IF Transmitted collision error packets  | Count        | delta      | /proc/net/dev
 35 | Linux | Network           | network.{nwif name}.transmit.carrier.count    | IF Transmitted carrier loss packets     | Count        | delta      | /proc/net/dev
 36 | Linux | Network           | network.{nwif name}.transmit.compressed.count | IF Transmitted compressed packets       | Count        | delta      | /proc/net/dev
 37 | Linux | Load Average      | loadavg.1.count                               | Load average (1 min)                    | Count        | gauge      | /proc/loadavg
 38 | Linux | Load Average      | loadavg.5.count                               | Load average (5 min)                    | Count        | gauge      | /proc/loadavg
 39 | Linux | Load Average      | loadavg.15.count                              | Load average (15 min)                   | Count        | gauge      | /proc/loadavg
 40 | Linux | Memory            | memory.memtotal.kilobytes                     | Total memory                            | kilobytes    | gauge      | /proc/meminfo (MemTotal)
 41 | Linux | Memory            | memory.memfree.kilobytes                      | Memory left unused                      | kilobytes    | gauge      | /proc/meminfo (MemFree)
 42 | Linux | Memory            | memory.buffers.kilobytes                      | Memory used for file buffers            | kilobytes    | gauge      | /proc/meminfo (Buffers)
 43 | Linux | Memory            | memory.cached.kilobytes                       | Memory used as cache                    | kilobytes    | gauge      | /proc/meminfo (Cached)
 44 | Linux | Memory            | memory.swapcached.kilobytes                   | Memory used as swap cache               | kilobytes    | gauge      | /proc/meminfo (SwapCached)
 45 | Linux | Memory            | memory.active.kilobytes                       | Memory used actively                    | kilobytes    | gauge      | /proc/meminfo (Active)
 46 | Linux | Memory            | memory.inactive.kilobytes                     | Memory used inactively                  | kilobytes    | gauge      | /proc/meminfo (Inactive)
 47 | Linux | Memory            | memory.swaptotal.kilobytes                    | Total amount of swap available          | kilobytes    | gauge      | /proc/meminfo (SwapTotal)
 48 | Linux | Memory            | memory.swapfree.kilobytes                     | Total amount of swap free               | kilobytes    | gauge      | /proc/meminfo (SwapFree)
 49 | Linux | Memory            | memory.usedtotal.kilobytes                    | Total used Memory                       | kilobytes    | gauge      | /proc/meminfo (MemTotal - (MemFree + Buffers + Cached))

# How to build packages

## Build on CentOS

### Preparation

Edit your `~/.bash_profile`

```
export GOROOT=/usr/local/go
export GOPATH=$HOME
export GHQ_ROOT=$GOPATH/src
PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

Install golang

```
# Download
wget https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz

# Extract
sudo tar zxvf go1.8.linux-amd64.tar.gz -C /usr/local/

# Modify the file owner
sudo chown $USER:$USER /usr/local/go/pkg/
```

Install docker

```
# Add epel repository
sudo yum install epel-release

# Install docker-io
sudo yum install docker-io --enablerepo=epel

# Create docker group
sudo groupadd docker
sudo gpasswd -a $USER docker
sudo usermod -g docker $USER
sudo service docker start

# Logout and relogin to enable the group
exit
```

Install ghq

```
go get github.com/motemen/ghq
```

Clone this repository

```
# clone to $HOME/src/github.com/nttcom/ecl2mond/ by using ghq
ghq get git@github.com:nttcom/ecl2mond.git
```

### Generate the packages

Execute the following commands at `$HOME/src/github.com/nttcom/ecl2mond/` directory and generate the packages.
(The first run will take a long time.)

```
make deps

make package-all-docker

# generated packages are here
ls $HOME/src/github.com/nttcom/ecl2mond/_packaging/output
```

# How to install

## Supported Platform

* Linux
    * CentOS7.1
    * Ubuntu16.04LTS

Above is the list of tested platforms. `ecl2mond` should work on other linux platforms.


## Install procedure

### Add nttcom package repository

#### RHEL /  CentOS

Create `/etc/yum.repos.d/ecl2mond.repo` and edit as follows.

```
[ecl2mond]
name=Enterprise Cloud Custom Meter Agent
baseurl=http://repo.ecl.ntt.com/packages/rpms
enabled=1
gpgcheck=1
```

Download and import the public key for signature verification.

```
curl -O http://repo.ecl.ntt.com/packages/GPG-KEY-ecl; rpm --import GPG-KEY-ecl
```

Install by using `yum`, a package management tool.

```
yum install ecl2mond
```

#### Debian / Ubuntu

Create `/etc/apt/sources.list.d/ecl2mond.list` and edit as follows.

```
deb http://repo.ecl.ntt.com/packages/debs ./
```

Download and import the public key for signature verification.

```
curl -O http://repo.ecl.ntt.com/packages/GPG-KEY-ecl; apt-key add GPG-KEY-ecl
```

Install by using `apt`, a package management tool.

```
apt-get update; apt-cache search ecl2mond
apt-get install ecl2mond
```

(Additional)
If the `apt-transport-https` package is not installed, you need to install it.

```
apt-get update; apt-cache search apt-transport-https
apt-get install apt-transport-https
```

# Configuration and default

The `ecl2mond` configuration file is generated at `/etc/ecl2mond/ecl2mond.conf` .
Please refer to the following description and modify the parameters.

## Details of the configuration parameters

Parameter name | Required | Data type    | Default value | Value range                                                                                                                          | Description
---------------|----------|--------------|---------------|--------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------
monitoringUrl  | yes      | String       | N/A           |                                                                                                                                      | Monitoring API Endpoint <br> (Set your region's one.)
interval       |          | Integer      | 5             | 1 - 3599                                                                                                                             | Interval of data collection and sending (minute)
authUrl        | yes      | String       | N/A           |                                                                                                                                      |
authInterval   |          | Integer      | 60            | 5 - 3599                                                                                                                             | Interval for keystone token refreshing (minute) <br> (Basically it should be remained the default.)
resourceId     | yes      | String       | N/A           |                                                                                                                                      | Resource ID of the Monitoring target resource
tenantId       | yes      | String       | N/A           |                                                                                                                                      | Your tenant ID
userName       | yes      | String       | N/A           |                                                                                                                                      | Your API key
password       | yes      | String       | N/A           |                                                                                                                                      | Your API secret key
meters         | yes      | String array | N/A           | Please refer to [Monitoring service description](https://ecl.ntt.com/documents/service-descriptions/monitoring/monitoring.html#id9). | Meter names you want to monitor
logLevel       |          | String       | Info          | "TRACE", "DEBUG", "INFO", "WARNING", "ERROR", "FATAL"                                                                                | Log level

## Sample configuration

```
monitoringUrl: Enter the ECL2.0 Monitoring API Endpoint URL
# ex. "https://monitoring-jp1-ecl.api.ntt.com/"
monitoringUrl = "https://monitoring-jp1-ecl.api.ntt.com/"

# interval: Specify the time interval to collect and send the value by the agent
# [ 1 - 3599 ] minutes, default: 5
interval = 5

# authUrl: Enter the ECL2.0 Keystone API Endpoint URL
# ex. "https://keystone-jp1-ecl.api.ntt.com/"
authUrl = "https://keystone-jp1-ecl.api.ntt.com/"

# authInterval: Specify the time interval to refresh auth token
# [ 5 - 3599 ] min, default: 60
authInterval = 60

# resourceId: Enter the target resource id for custom meter creation
# ex. nova_12ab-12cd56gh9-ab34ef78i-34cd
resourceId = "nova_6708b574-434d-4476-b8c1-2a128ab5edc1"

# tenantId: Enter your tenant id
tenantId = "<your tenant id>"

# userName: Enter your API key
userName = "<your API key>"

# password: Enter your API secret key
password = "<your API secret key>"

meters = [
  "loadavg.1.count",
]
```

# Usage

```
Usage: ecl2mond [options]

options:
  -m, --mode {run(default)|run-once|dry-run}
      set mode
  -c, --config
      set config file
  -h, --help
      show help message
```

For operation check

```
# dry-run (run with data collection only)
ecl2mond -m dry-run

# run-once (collect data and send it to Monitoring only once)
ecl2mond -m run-once
```

Start `ecl2mond` in the background

```
service ecl2mond start
```

Stop `ecl2mond`

```
service ecl2mond stop
```

Check the status of `ecl2mond`

```
service ecl2mond status
```

# Logs

When `ecl2mond` is started in the background, a log is output to `/var/log/ecl2mond.log`.

# Documentation

Please find more usage documentation [here](https://ecl.ntt.com/documents/service-descriptions/monitoring/monitoring.html) or [here](https://ecl.ntt.com/documents/tutorials/rsts/Monitoring/custom_meter/index.html).

# Support

ECL2.0 users can raise requests via NTT Communication's ticket portal.

# Contact

* ecl2-mon-app-cl@ntt.com

# Contributing

Please contribute [Github Flow](https://guides.github.com/introduction/flow/ "Github Flow")  -  Create a branch, add commits, and [open a pull request](https://github.com/nttcom/ecl2mond/compare/ "Pull Request").

# License

* Apache 2.0
