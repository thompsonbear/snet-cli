package main

import (
	"flag"
	"fmt"
	"net/netip"
	"reflect"
	"strconv"
	"strings"

	"github.com/thompsonbear/netmath"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func getHostRange(na netip.Addr, ba netip.Addr, short bool, usable bool) string {
	if(na.Is6() && short){
		return "Short IPv6 Not Supported"
	}

	if ba.Prev().Less(na.Next()) && usable {
		return "None"
	}

	var first netip.Addr
	var last netip.Addr
	if usable {
		first = na.Next()
		last = ba.Prev()
	} else {
		first = na
		last = ba
	}

	if(!short){
		return first.String() + "-" + last.String()
	}

	firstBytes := first.AsSlice()
	lastBytes := last.AsSlice()

	var hostRange string

	for i := 0; i < len(firstBytes); i++ {
		if(firstBytes[i] == lastBytes[i]){
			hostRange += strconv.Itoa(int(firstBytes[i]))
		} else if (firstBytes[i] < lastBytes[i]){
			hostRange += strconv.Itoa(int(firstBytes[i])) + "-" + strconv.Itoa(int(lastBytes[i]))
		} else {
			return "None"
		}

		if(i < len(firstBytes) - 1) {
			hostRange += "."
		}
	}

	return hostRange
}

type Fields struct {
	Prefix bool
	Network bool
	FullRange bool
	UsableRange bool
	Broadcast bool
	SubnetMask bool
	MaskBits bool
	CountTotal bool
	CountUsable bool
}

type Options struct{
	All bool
	Borderless bool
	Count bool
	ShortRange bool
	Fields Fields
}

func printSubnetTable(s netmath.Subnet, opts Options) {
	colHeaders := []string{}

	fv := reflect.ValueOf(opts.Fields)
	ft := fv.Type()

	for i := 0; i < fv.NumField(); i++ {
		if fv.Field(i).Interface() == true {
			colHeaders = append(colHeaders, ft.Field(i).Name)
		}
	}
	
	t := table.New().
	BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("12"))).
	StyleFunc(func(row, col int) lipgloss.Style {
		switch {
			case row == 0:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true).PaddingRight(1).PaddingLeft(1)
			case row % 2 != 0:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("15")).PaddingRight(1).PaddingLeft(1)
			default: 
				return lipgloss.NewStyle().Foreground(lipgloss.Color("7")).PaddingRight(1).PaddingLeft(1)
		}
	}).Headers(colHeaders...)

	if(opts.Borderless){
		t.BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderRow(false).
		BorderColumn(false).
		BorderHeader(false)
	}

	var subnets []netmath.Subnet
	if opts.All {
		subnets = s.ListAll()
	} else {
		subnets = make([]netmath.Subnet, 0, 1)
		subnets = append(subnets, s)
	}

	for _, subnet := range subnets {
		cols := []string{}
		for i := 0; i < fv.NumField(); i++ {
			if fv.Field(i).Interface() == true {
				cols = append(cols, getSingleField(subnet, ft.Field(i).Name, opts.ShortRange))
			}
		}
		t.Row(cols...)
	}
	fmt.Println(t)
}

func getSingleField(s netmath.Subnet, field string, short_range bool) string {
	switch field {
	case "Prefix":
		na, _ := s.Network()
		return fmt.Sprintf("%v/%v", na.String(), s.Bits())
	case "Network":
		na, _ := s.Network()
		return na.String()
	case "FullRange":
		na, _ := s.Network()
		ba, _ := s.Broadcast()
		return getHostRange(na, ba, short_range, false)
	case "UsableRange":
		na, _ := s.Network()
		ba, _ := s.Broadcast()
		return getHostRange(na, ba, short_range, true)
	case "Broadcast":
		ba, _ := s.Broadcast()
		return ba.String()
	case "SubnetMask":
		m, _ := s.Mask()
		return m.String()
	case "MaskBits":
		return strconv.Itoa(s.Bits())
	case "CountTotal":
		ct, _ := s.Count()
		return strconv.FormatFloat(ct, 'g', -1, 64)
	case "CountUsable":
		ct, _ := s.Count()
		return strconv.FormatFloat(ct-2, 'g', -1, 64)
	default: // Same as prefix
		na, _ := s.Network()
		return fmt.Sprintf("%v/%v", na, s.Bits())
	}
}

func countTrueFields(fields Fields) (int, string) {
	fv := reflect.ValueOf(fields)
	ft := fv.Type()

	var count int
	var last string

	for i := 0; i < fv.NumField(); i++ {
		if fv.Field(i).Interface() == true {
			count += 1
			last = ft.Field(i).Name
		}
	}

	return count, last
}

func main() {

	var help bool
	opts := Options{}

	flag.BoolVar(&help, "h", false, "Display help message")
	flag.BoolVar(&help, "help", false, "")
	
	// Base Options
	flag.BoolVar(&opts.All, "a", false, "Display ALL possible networks within the specified subnet.")
	flag.BoolVar(&opts.All, "all", false, "")

	flag.BoolVar(&opts.Borderless, "bl", false, "Display output without borders.")
	flag.BoolVar(&opts.Borderless, "borderless", false, "")

	flag.BoolVar(&opts.ShortRange, "s", false, "Display ip ranges shorthand/abbreviated notation.")
	flag.BoolVar(&opts.ShortRange, "short", false, "")

	flag.BoolVar(&opts.Count, "c", false, "Display ip ranges with host count.")
	flag.BoolVar(&opts.Count, "count", false, "")
	
	// Fields
	flag.BoolVar(&opts.Fields.Prefix, "p", false, "CIDR Prefix Field")
	flag.BoolVar(&opts.Fields.Prefix, "prefix", false, "")

	flag.BoolVar(&opts.Fields.Network, "na", false, "Network Address Field")
	flag.BoolVar(&opts.Fields.Network, "network", false, "")
	flag.BoolVar(&opts.Fields.Network, "network-address", false, "")

	flag.BoolVar(&opts.Fields.FullRange, "fr", false, "Full IP Range Field")
	flag.BoolVar(&opts.Fields.FullRange, "full-range", false, "")

	flag.BoolVar(&opts.Fields.UsableRange, "ur", false, "Usable IP Range Field")
	flag.BoolVar(&opts.Fields.UsableRange, "usable-range", false, "")

	flag.BoolVar(&opts.Fields.Broadcast, "ba", false, "Broadcast Address Field")
	flag.BoolVar(&opts.Fields.Broadcast, "broadcast", false, "")
	flag.BoolVar(&opts.Fields.Broadcast, "broadcast-address", false, "")

	flag.BoolVar(&opts.Fields.SubnetMask, "m", false, "Subnet Mask Field")
	flag.BoolVar(&opts.Fields.SubnetMask, "mask", false, "")
	flag.BoolVar(&opts.Fields.SubnetMask, "subnet-mask", false, "")

	flag.BoolVar(&opts.Fields.MaskBits, "b", false, "Subnet Mask Field")
	flag.BoolVar(&opts.Fields.MaskBits, "bits", false, "")
	flag.BoolVar(&opts.Fields.MaskBits, "mask-bits", false, "")

	flag.BoolVar(&opts.Fields.CountTotal, "ct", false, "Count Total Hosts Field")
	flag.BoolVar(&opts.Fields.CountTotal, "count-total", false, "")

	flag.BoolVar(&opts.Fields.CountUsable, "cu", false, "Count Usable Hosts Field")
	flag.BoolVar(&opts.Fields.CountUsable, "count-usable", false, "")
	
	
	flag.Parse()
	args := flag.Args()
		
	if(help || len(args) == 0){
		cmdstyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).PaddingLeft(2)
		orstyle := lipgloss.NewStyle().PaddingLeft(4)
		fmt.Println("Usage:")
		fmt.Println(cmdstyle.Render("snet <options> <host-address> <subnet-mask>"))
		fmt.Println(orstyle.Render("or"))
		fmt.Println(cmdstyle.Render("snet <options> <host-address>/<mask-bits>"))
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		return
	}

	var input []string

	for i:= 0; i < len(args); i++ {
		if !strings.HasPrefix(args[i], "-") {
			input = append(input, args[i])
		}
	}
	var s netmath.Subnet
	var parseErr error
	if len(input) == 1  {
		s, parseErr = netmath.ParseCIDR(input[0])
	} else if len(input) > 1  {
		s, parseErr = netmath.Parse(input[0], input[1])
	}
	if parseErr != nil {
		fmt.Println(parseErr)
		return
	}

	// Enable both count options if count option is specified
	if opts.Count {
		opts.Fields.CountTotal = true
		opts.Fields.CountUsable = true
	}

	fieldCount, lastField := countTrueFields(opts.Fields)

	if fieldCount > 1 || fieldCount == 0 {
		// Set default options if no Fields are specified
		if fieldCount == 0 {
			opts.Fields.Prefix = true
			opts.Fields.Network = true
			opts.Fields.UsableRange = true
			opts.Fields.Broadcast = true
			opts.Fields.SubnetMask = true
		}

		printSubnetTable(s, opts)
	} else if fieldCount == 1 {
		fmt.Println(getSingleField(s, lastField, opts.ShortRange))
	}
	
}
