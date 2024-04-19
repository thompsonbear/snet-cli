# snet cli
Subnet from the command line

## Installation

#### Build from source

```Bash
git clone https://github.com/thompsonbear/snet-cli
cd snet-cli
go build snet.go
```

## Usage

#### Address with Mask Bits

```Bash
snet <options> <host-address>/<mask-bits>
```

```Bash
snet -b 192.168.20.10/23
```

#### Address with Subnet Mask

```Bash
snet <options> <host-address> <subnet-mask>
```

```Bash
snet -b 192.168.20.10 255.255.254.0
```

#### Sample Output

```Bash
Prefix            Network       Useable Hosts        Broadcast       Mask
192.168.20.10/23  192.168.20.0  192.168.20-21.1-254  192.168.21.255  255.255.254.0
```

## Options

**-a, -all**

- Display ALL possible networks within the specified subnet.

**-b, -borderless**

- Display with a BORDERLESS table.

**-c, -count**

- Display ip ranges with host count.

**-h, -help**

- Display help message

**-s, -simple**

- Display ip ranges without sub-range notation.

