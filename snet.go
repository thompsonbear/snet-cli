package main

import (
	"flag"
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	"github.com/thompsonbear/netmath"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func getHostRange(na netip.Addr, ba netip.Addr, short bool) string {
	if(!na.Is4() && short){
		return "Not Supported"
	}

	if ba.Prev().Less(na.Next()) {
		return "None"
	}

	if(!short){
		return na.Next().String() + "-" + ba.Prev().String()
	}
	firstHost := na.Next().AsSlice()
	lastHost := ba.Prev().AsSlice()

	var hostRange string

	for i := 0; i < len(firstHost); i++ {
		if(firstHost[i] == lastHost[i]){
			hostRange += strconv.Itoa(int(firstHost[i]))
		} else if (firstHost[i] < lastHost[i]){
			hostRange += strconv.Itoa(int(firstHost[i])) + "-" + strconv.Itoa(int(lastHost[i]))
		} else {
			return "None"
		}

		if(i < len(firstHost) - 1) {
			hostRange += "."
		}
	}

	return hostRange
}

type options struct{
	borderless bool 
	showCount bool
	shortRange bool
	listAll bool
}

func printSubnetTable(s netmath.Subnet, opts options) {
	colHeaders := []string{"Subnet", "Network", "Useable Hosts", "Broadcast", "Mask"}

	if(opts.showCount){
		colHeaders = append(colHeaders, "Total Count")
		colHeaders = append(colHeaders, "Useable Count")
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

	if(opts.borderless){
		t.BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderRow(false).
		BorderColumn(false).
		BorderHeader(false)
	}

	var subnets []netmath.Subnet
	if opts.listAll {
		subnets = s.ListAll()
	} else {
		subnets = make([]netmath.Subnet, 0, 1)
		subnets = append(subnets, s)
	}

	for _, subnet := range subnets {
		na, _ := subnet.Network()
		ba, _ := subnet.Broadcast()
		mask, _ := subnet.Mask()
		hostRange := getHostRange(na, ba, opts.shortRange)
		
		cols := []string{subnet.String(), na.String(), hostRange, ba.String(), mask.String()}

		if opts.showCount {
			c, err := subnet.Count()
			if err != nil {
				cols = append(cols, "Error")
				cols = append(cols, "Error")
			} else {
				cols = append(cols, strconv.FormatFloat(c, 'g', -1, 64)) // Total Hosts

				if c >= 2 {
					cols = append(cols, strconv.FormatFloat(c-2, 'g', -1, 64)) // Usable Hosts
				} else {
					cols = append(cols, "0")
				}
			}
		}
		t.Row(cols...)
	}
	fmt.Println(t)
}

func main() {

	var help bool
	opts := options{}
	
	flag.BoolVar(&opts.listAll, "a", false, "Display ALL possible networks within the specified subnet.")
	flag.BoolVar(&opts.listAll, "all", false, "")

	flag.BoolVar(&opts.borderless, "b", false, "Display with a BORDERLESS table.")
	flag.BoolVar(&opts.borderless, "borderless", false, "")

	flag.BoolVar(&opts.showCount, "c", false, "Display ip ranges with host count.")
	flag.BoolVar(&opts.showCount, "count", false, "")

	flag.BoolVar(&opts.shortRange, "s", false, "Display ip ranges shorthand/abbreviated notation.")
	flag.BoolVar(&opts.shortRange, "short", false, "")
	
	flag.BoolVar(&help, "h", false, "Display help message")
	flag.BoolVar(&help, "help", false, "")
	
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

	
	printSubnetTable(s, opts)
}
