// Copyright (c) 2016-2020 Tigera, Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package calico

import (
	"context"
	"fmt"
	"github.com/projectcalico/calicoctl/calicoctl/commands/clientmgr"
	"github.com/projectcalico/libcalico-go/lib/ipam"
	"math"
	"strings"
)

type IPPResult struct {
	Name     string
	CIDR     string
	Total    float64
	Inuse    float64
	Free     float64
	FreeRate float64
	Selector string
	IPIP     string
}

func showBlockUtilization(ctx context.Context, ipamClient ipam.Interface, showBlocks bool) ([]IPPResult, error) {
	usage, err := ipamClient.GetUtilization(ctx, ipam.GetUtilizationArgs{})
	if err != nil {
		return nil, nil
	}
	//table := tablewriter.NewWriter(os.Stdout)
	//table.SetHeader([]string{"GROUPING", "CIDR", "IPS TOTAL", "IPS IN USE", "IPS FREE"})
	var ippList []IPPResult
	genRow := func(kind, cidr string, inUse, capacity float64) []string {
		return []string{
			kind,
			cidr,
			fmt.Sprintf("%.5g", capacity),
			// Note: the '+capacity/2' bits here give us rounding to the nearest
			// integer, instead of rounding down, and so ensure that the two percentages
			// add up to 100.
			fmt.Sprintf("%.5g (%.f%%)", inUse, 100*inUse/capacity),
			fmt.Sprintf("%.5g (%.f%%)", capacity-inUse, 100*(capacity-inUse)/capacity),
		}
	}
	for _, poolUse := range usage {
		var blockRows [][]string
		var poolInUse float64
		if strings.Contains(poolUse.Name, "orphaned allocation blocks") && poolUse.CIDR.String() == "0.0.0.0/0" {
			continue
		}
		for _, blockUse := range poolUse.Blocks {
			blockRows = append(blockRows, genRow("Block", blockUse.CIDR.String(), float64(blockUse.Capacity-blockUse.Available), float64(blockUse.Capacity)))
			poolInUse += float64(blockUse.Capacity - blockUse.Available)
		}
		ones, bits := poolUse.CIDR.Mask.Size()
		poolCapacity := math.Pow(2, float64(bits-ones))
		if ones > 0 {
			// Only show the IP Pool row for a real IP Pool and not for the orphaned
			// block case.
			//table.Append(genRow("IP Pool", poolUse.CIDR.String(), poolInUse, poolCapacity))
		}
		if showBlocks {
			//table.AppendBulk(blockRows)
		}
		name := strings.Replace(poolUse.Name, "-", "_", -1)
		name = strings.Replace(name, " ", "_", -1)
		re := IPPResult{Name: name, CIDR: poolUse.CIDR.String(), Total: poolCapacity, Inuse: poolInUse, Free: poolCapacity - poolInUse, FreeRate: (poolCapacity - poolInUse) / poolCapacity, Selector: poolUse.Selector, IPIP: poolUse.IPIP}
		ippList = append(ippList, re)
	}
	//table.Render()

	return ippList, nil
}

const DefaultConfigPath = "/etc/calico/calicoctl.cfg"

// IPAM takes keyword with an IP address then calls the subcommands.
func Show() ([]IPPResult, error) {
	ctx := context.Background()

	client, err := clientmgr.NewClient(DefaultConfigPath)
	if err != nil {
		return nil, err
	}
	ipamClient := client.IPAM()
	return showBlockUtilization(ctx, ipamClient, false)
}
