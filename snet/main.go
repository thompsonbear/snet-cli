package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/thompsonbear/netmath"
)

func main() {
	var help bool

	flag.BoolVar(&help, "h", false, "Display help message (shorthand)")
	flag.BoolVar(&help, "help", false, "Display help message")

	opts := Options{}

	// General Options
	opts.All.UsageStr = "Display All possible networks within the specified subnet."
	flag.BoolVar(&opts.All.Active, "a", false, fmt.Sprintf("%v (shorthand)", opts.All.UsageStr))
	flag.BoolVar(&opts.All.Active, "all", false, opts.All.UsageStr)

	opts.Borderless.UsageStr = "Display output without borders."
	flag.BoolVar(&opts.Borderless.Active, "bl", false, fmt.Sprintf("%v (shorthand)", opts.Borderless.UsageStr))
	flag.BoolVar(&opts.Borderless.Active, "borderless", false, opts.Borderless.UsageStr)

	opts.SubRange.UsageStr = "Display IP ranges with sub-range notation"
	flag.BoolVar(&opts.SubRange.Active, "s", false, fmt.Sprintf("%v (shorthand)", opts.SubRange.UsageStr))
	flag.BoolVar(&opts.SubRange.Active, "short-range", false, opts.SubRange.UsageStr)

	// Field Options
	opts.Fields.Prefix.UsageStr = "Display CIDR Prefix Field"
	flag.BoolVar(&opts.Fields.Prefix.Active, "p", false, fmt.Sprintf("%v (shorthand)", opts.Fields.Prefix.UsageStr))
	flag.BoolVar(&opts.Fields.Prefix.Active, "prefix", false, opts.Fields.Prefix.UsageStr)

	opts.Fields.NetworkAddress.UsageStr = "Display Network Address Field"
	flag.BoolVar(&opts.Fields.NetworkAddress.Active, "na", false, fmt.Sprintf("%v (shorthand)", opts.Fields.NetworkAddress.UsageStr))
	flag.BoolVar(&opts.Fields.NetworkAddress.Active, "network-address", false, opts.Fields.NetworkAddress.UsageStr)

	opts.Fields.FullRange.UsageStr = "Display Full IP Range Field"
	flag.BoolVar(&opts.Fields.FullRange.Active, "fr", false, fmt.Sprintf("%v (shorthand)", opts.Fields.FullRange.UsageStr))
	flag.BoolVar(&opts.Fields.FullRange.Active, "full-range", false, opts.Fields.FullRange.UsageStr)

	opts.Fields.UsableRange.UsageStr = "Display Usable IP Range Field"
	flag.BoolVar(&opts.Fields.UsableRange.Active, "ur", false, fmt.Sprintf("%v (shorthand)", opts.Fields.UsableRange.UsageStr))
	flag.BoolVar(&opts.Fields.UsableRange.Active, "usable-range", false, opts.Fields.UsableRange.UsageStr)

	opts.Fields.BroadcastAddress.UsageStr = "Display Broadcast Address Field"
	flag.BoolVar(&opts.Fields.BroadcastAddress.Active, "ba", false, fmt.Sprintf("%v (shorthand)", opts.Fields.BroadcastAddress.UsageStr))
	flag.BoolVar(&opts.Fields.BroadcastAddress.Active, "broadcast-address", false, opts.Fields.BroadcastAddress.UsageStr)

	opts.Fields.SubnetMask.UsageStr = "Display Subnet Mask Field"
	flag.BoolVar(&opts.Fields.SubnetMask.Active, "m", false, fmt.Sprintf("%v (shorthand)", opts.Fields.SubnetMask.UsageStr))
	flag.BoolVar(&opts.Fields.SubnetMask.Active, "mask", false, opts.Fields.SubnetMask.UsageStr)

	opts.Fields.MaskBits.UsageStr = "Display Subnet Mask Bits Field"
	flag.BoolVar(&opts.Fields.MaskBits.Active, "b", false, fmt.Sprintf("%v (shorthand)", opts.Fields.MaskBits.UsageStr))
	flag.BoolVar(&opts.Fields.MaskBits.Active, "bits", false, opts.Fields.MaskBits.UsageStr)

	opts.Fields.TotalCount.UsageStr = "Display Total Hosts Count Field"
	flag.BoolVar(&opts.Fields.TotalCount.Active, "tc", false, fmt.Sprintf("%v (shorthand)", opts.Fields.TotalCount.UsageStr))
	flag.BoolVar(&opts.Fields.TotalCount.Active, "total-count", false, opts.Fields.TotalCount.UsageStr)

	opts.Fields.UsableCount.UsageStr = "Display Usable Hosts Count Field"
	flag.BoolVar(&opts.Fields.UsableCount.Active, "uc", false, fmt.Sprintf("%v (shorthand)", opts.Fields.UsableCount.UsageStr))
	flag.BoolVar(&opts.Fields.UsableCount.Active, "usable-count", false, opts.Fields.UsableCount.UsageStr)

	opts.Count.UsageStr = "Both -uc and -tc options"
	flag.BoolVar(&opts.Count.Active, "c", false, fmt.Sprintf("%v (shorthand)", opts.Count.UsageStr))
	flag.BoolVar(&opts.Count.Active, "count", false, opts.Count.UsageStr)

	flag.Parse()
	args := flag.Args()

	// Display the help menu if help option is specified or no args
	if help || len(args) == 0 {
		PrintHelp(opts)
		return
	}

	var input []string

	for i := 0; i < len(args); i++ {
		if !strings.HasPrefix(args[i], "-") {
			input = append(input, args[i])
		}
	}

	// Parse user input
	var s netmath.Subnet
	var parseErr error
	if len(input) == 1 {
		s, parseErr = netmath.ParseCIDR(input[0])
	} else if len(input) > 1 {
		s, parseErr = netmath.Parse(input[0], input[1])
	}
	if parseErr != nil {
		fmt.Println(parseErr)
		return
	}

	// Enable both count options if count option is specified
	if opts.Count.Active {
		opts.Fields.TotalCount.Active = true
		opts.Fields.UsableCount.Active = true
	}

	fieldCount, lastField := CountActiveFields(opts.Fields)

	if fieldCount > 1 || fieldCount == 0 || opts.All.Active {
		// Set default options if no Fields are specified
		if fieldCount == 0 {
			opts.Fields.Prefix.Active = true
			opts.Fields.NetworkAddress.Active = true
			opts.Fields.UsableRange.Active = true
			opts.Fields.BroadcastAddress.Active = true
			opts.Fields.SubnetMask.Active = true
		}

		PrintSubnetTable(s, opts)
	} else if fieldCount == 1 {
		fmt.Println(GetSubnetField(s, lastField, opts.SubRange.Active))
	}

}
