package main

type Options struct {
	All        Option
	Borderless Option
	Count      Option
	SubRange   Option
	Fields     Fields
}

type Fields struct {
	Prefix           Option
	NetworkAddress   Option
	FullRange        Option
	UsableRange      Option
	BroadcastAddress Option
	SubnetMask       Option
	MaskBits         Option
	TotalCount       Option
	UsableCount      Option
}

type Option struct {
	Active   bool
	UsageStr string
}
