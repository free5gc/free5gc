package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"gofree5gc/lib/path_util"
	"gofree5gc/lib/pfcp/pfcpType"
	"gofree5gc/src/smf/smf_context"
	"io/ioutil"
)

var topo_path = path_util.Gofree5gcPath("gofree5gc/smf_upf_topo.json")

// func Read_topo() (upf_root *smf_context.UPF) {

// }

func main() {

	var upf_root *smf_context.UPF
	links := make(map[string][]string)
	upf_nodes := make(map[string]*smf_context.UPF)

	content, err := ioutil.ReadFile(topo_path)

	if err != nil {
		fmt.Println("Failed to open file from path: ", topo_path)
		return
	}

	json := string(content)

	result := gjson.Get(json, "switches")

	if !result.Exists() {
		fmt.Println("There is no UPF in the Topology!")
		return
	}
	result.ForEach(func(key, value gjson.Result) bool {

		switch_type := gjson.Get(value.String(), "opts.switchType").String()

		upf_name := gjson.Get(value.String(), "opts.hostname").String()
		//create new upf
		NodeID := new(pfcpType.NodeID)
		upf := smf_context.AddUPF(NodeID)
		upf_nodes[upf_name] = upf

		if switch_type == "legacyRouter" {
			upf_root = upf
		}

		return true // keep iterating
	})

	result := gjson.Get(json, "links")

	if !result.Exists() {
		fmt.Println("There is no data link in the Topology!")
		return
	}

	result.ForEach(func(key, value gjson.Result) bool {
		src := gjson.Get(value.String(), "src").String()
		dest := gjson.Get(value.String(), "dest").String()

		links[src] = append(links[src], dest)
		return true // keep iterating
	})

}
