package main

import (
	"fmt"
	"net/netip"
	"reflect"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/thompsonbear/netmath"
)

func GetSubnetField(s netmath.Subnet, field string, sub_range bool) string {
	switch field {
	case "Prefix":
		na, _ := s.Network()
		return fmt.Sprintf("%v/%v", na.String(), s.Bits())
	case "NetworkAddress":
		na, _ := s.Network()
		return na.String()
	case "FullRange":
		na, _ := s.Network()
		ba, _ := s.Broadcast()
		return GetHostRange(na, ba, sub_range, false)
	case "UsableRange":
		na, _ := s.Network()
		ba, _ := s.Broadcast()
		return GetHostRange(na, ba, sub_range, true)
	case "BroadcastAddress":
		ba, _ := s.Broadcast()
		return ba.String()
	case "SubnetMask":
		m, _ := s.Mask()
		return m.String()
	case "MaskBits":
		return strconv.Itoa(s.Bits())
	case "TotalCount":
		ct, _ := s.Count()
		return strconv.FormatFloat(ct, 'g', -1, 64)
	case "UsableCount":
		ct, _ := s.Count()
		return strconv.FormatFloat(ct-2, 'g', -1, 64)
	default: // Same as prefix
		na, _ := s.Network()
		return fmt.Sprintf("%v/%v", na, s.Bits())
	}
}

func PrintSubnetTable(s netmath.Subnet, opts Options) {
	colHeaders := []string{}

	fv := reflect.ValueOf(opts.Fields)
	ft := fv.Type()

	for i := 0; i < fv.NumField(); i++ {
		if fv.Field(i).FieldByName("Active").Interface() == true {
			colHeaders = append(colHeaders, ft.Field(i).Name)
		}
	}

	t := table.New().
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("12"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true).PaddingRight(1).PaddingLeft(1)
			case row%2 != 0:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("15")).PaddingRight(1).PaddingLeft(1)
			default:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("7")).PaddingRight(1).PaddingLeft(1)
			}
		}).Headers(colHeaders...)

	if opts.Borderless.Active {
		t.BorderTop(false).
			BorderBottom(false).
			BorderLeft(false).
			BorderRight(false).
			BorderRow(false).
			BorderColumn(false).
			BorderHeader(false)
	}

	var subnets []netmath.Subnet
	if opts.All.Active {
		subnets = s.ListAll()
	} else {
		subnets = make([]netmath.Subnet, 0, 1)
		subnets = append(subnets, s)
	}

	for _, subnet := range subnets {
		cols := []string{}
		for i := 0; i < fv.NumField(); i++ {
			if fv.Field(i).FieldByName("Active").Interface() == true {
				cols = append(cols, GetSubnetField(subnet, ft.Field(i).Name, opts.SubRange.Active))
			}
		}
		t.Row(cols...)
	}
	fmt.Println(t)
}

func GetHostRange(na netip.Addr, ba netip.Addr, short bool, usable bool) string {
	if na.Is6() && short {
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

	if !short {
		return first.String() + "-" + last.String()
	}

	firstBytes := first.AsSlice()
	lastBytes := last.AsSlice()

	var hostRange string

	for i := 0; i < len(firstBytes); i++ {
		if firstBytes[i] == lastBytes[i] {
			hostRange += strconv.Itoa(int(firstBytes[i]))
		} else if firstBytes[i] < lastBytes[i] {
			hostRange += strconv.Itoa(int(firstBytes[i])) + "-" + strconv.Itoa(int(lastBytes[i]))
		} else {
			return "None"
		}

		if i < len(firstBytes)-1 {
			hostRange += "."
		}
	}

	return hostRange
}

func CountActiveFields(fields Fields) (int, string) {
	fv := reflect.ValueOf(fields)
	ft := fv.Type()

	var count int
	var last string

	for i := 0; i < fv.NumField(); i++ {
		if fv.Field(i).FieldByName("Active").Interface() == true {
			count += 1
			last = ft.Field(i).Name
		}
	}

	return count, last
}

func PrintHelp(opts Options) {
	cmdstyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).PaddingLeft(2)
	orstyle := lipgloss.NewStyle().PaddingLeft(4)

	fmt.Println("Usage:")
	fmt.Println(cmdstyle.Render("snet <options> <host-address> <subnet-mask>"))
	fmt.Println(orstyle.Render("or"))
	fmt.Println(cmdstyle.Render("snet <options> <host-address>/<mask-bits>"))

	fmt.Println("\nGeneral Options:")
	fmt.Printf("%v\n\t%v\n\n", cmdstyle.Render("-a, -all"), opts.All.UsageStr)
	fmt.Printf("%v\n\t%v\n\n", cmdstyle.Render("-bl, -borderless"), opts.Borderless.UsageStr)
	fmt.Printf("%v\n\t%v\n\n", cmdstyle.Render("-s, -short-range"), opts.SubRange.UsageStr)

	fmt.Println("Field Options (Specifying multiple field options will display a table):")
	fmt.Printf("%v\n\t%v\n\n", cmdstyle.Render("-b, -mask-bits"), opts.Fields.MaskBits.UsageStr)
	fmt.Printf("%v\n\t%v\n\n", cmdstyle.Render("-ba, -broadcast-address"), opts.Fields.BroadcastAddress.UsageStr)
	fmt.Printf("%v\n\t%v\n\n", cmdstyle.Render("-c, -count"), opts.Count.UsageStr)
	fmt.Printf("%v\n\t%v\n\n", cmdstyle.Render("-fr, -full-range"), opts.Fields.FullRange.UsageStr)
	fmt.Printf("%v\n\t%v\n\n", cmdstyle.Render("-m, -subnet-mask"), opts.Fields.SubnetMask.UsageStr)
	fmt.Printf("%v\n\t%v\n\n", cmdstyle.Render("-p, -prefix"), opts.Fields.Prefix.UsageStr)
	fmt.Printf("%v\n\t%v\n\n", cmdstyle.Render("-tc, -total-count"), opts.Fields.TotalCount.UsageStr)
	fmt.Printf("%v\n\t%v\n\n", cmdstyle.Render("-uc, -usable-count"), opts.Fields.UsableCount.UsageStr)
	fmt.Printf("%v\n\t%v\n\n", cmdstyle.Render("-ur, -usable-range"), opts.Fields.UsableRange.UsageStr)
}
