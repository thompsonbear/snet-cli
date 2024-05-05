# snet-cli

Subnet from the command line

## Usage

#### Address with Mask Bits (CIDR Notation)

```Bash
snet <options> <host-address>/<mask-bits>
snet 192.168.20.10/23
```

#### Address with Subnet Mask

```Bash
snet <options> <host-address> <subnet-mask>
snet 192.168.20.10 255.255.254.0
```

#### Output with No Field Options (Prints subnet table with default fields)

```Bash
Prefix           NetworkAddress  UsableRange                  BroadcastAddress  SubnetMask
192.168.20.0/23  192.168.20.0    192.168.20.1-192.168.21.254  192.168.21.255    255.255.254.0
```

#### Single Field Specified (Prints single value)

```Bash
snet -na 192.168.20.10/23

Output:

192.168.20.0
```

#### Multiple Fields Specified (Prints table)

```Bash
snet -na -ba 192.168.20.10/23

Output:

NetworkAddress  BroadcastAddress
192.168.20.0    192.168.21.255
```

## General Options

**-a, -all**

> Display All possible networks within the specified subnet. (Always displays in a table)

**-bl, -borderless**

> Display output without borders.

**-h, -help**

> Display help message

**-s, -sub-range**

> Display IP ranges with sub-range notation

## Field Options

**-b, -mask-bits**

> Display Subnet Mask Bits Field

**-ba, -broadcast-address**

> Display Broadcast Address Field

**-c, -count**

> Both -ct and -cu options

**-fr, -full-range**

> Display Full IP Range Field

**-m, -subnet-mask**

> Display Subnet Mask Field

**-p, -prefix**

> Display CIDR Prefix Field

**-tc, -total-count**

> Display Total Hosts Count Field

**-uc, -usable-count**

> Display Usable Hosts Count Field

**-ur, -usable-range**

> Display Usable IP Range Field

 
## Installation

#### 0. [Install Go](https://go.dev/dl) (Prerequisite)

#### 1. Build from source

```Bash
git clone https://github.com/thompsonbear/snet-cli
cd ./snet-cli/snet
go build .
```

#### 2. Install in path environment variable for your respective OS
