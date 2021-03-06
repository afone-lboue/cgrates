package engine

import (
	"testing"

	"github.com/cgrates/cgrates/cache2go"
	"github.com/cgrates/cgrates/utils"
)

func init() {
	aliasService = NewAliasHandler(accountingStorage)
}
func TestAliasesGetAlias(t *testing.T) {
	alias := Alias{}
	err := aliasService.Call("AliasesV1.GetAlias", &Alias{
		Direction: "*out",
		Tenant:    "cgrates.org",
		Category:  "call",
		Account:   "dan",
		Subject:   "dan",
		Context:   "*rating",
	}, &alias)
	if err != nil ||
		len(alias.Values) != 2 ||
		len(alias.Values[0].Pairs) != 2 {
		t.Error("Error getting alias: ", err, alias, alias.Values)
	}
}

func TestAliasesGetMatchingAlias(t *testing.T) {
	var response string
	err := aliasService.Call("AliasesV1.GetMatchingAlias", &AttrMatchingAlias{
		Direction:   "*out",
		Tenant:      "cgrates.org",
		Category:    "call",
		Account:     "dan",
		Subject:     "dan",
		Context:     "*rating",
		Destination: "444",
		Target:      "Subject",
		Original:    "rif",
	}, &response)
	if err != nil || response != "rif1" {
		t.Error("Error getting alias: ", err, response)
	}
}

func TestAliasesSetters(t *testing.T) {
	var out string
	if err := aliasService.Call("AliasesV1.SetAlias", &AttrAddAlias{
		Alias: &Alias{
			Direction: "*out",
			Tenant:    "cgrates.org",
			Category:  "call",
			Account:   "set",
			Subject:   "set",
			Context:   "*rating",
			Values: AliasValues{&AliasValue{
				DestinationId: utils.ANY,
				Pairs:         AliasPairs{"Account": map[string]string{"1234": "1235"}},
				Weight:        10,
			}},
		},
		Overwrite: true,
	}, &out); err != nil || out != utils.OK {
		t.Error("Error setting alias: ", err, out)
	}
	r := &Alias{}
	if err := aliasService.Call("AliasesV1.GetAlias", &Alias{
		Direction: "*out",
		Tenant:    "cgrates.org",
		Category:  "call",
		Account:   "set",
		Subject:   "set",
		Context:   "*rating",
	}, r); err != nil || len(r.Values) != 1 || len(r.Values[0].Pairs) != 1 {
		t.Errorf("Error getting alias: %+v", r)
	}

	if err := aliasService.Call("AliasesV1.SetAlias", &AttrAddAlias{
		Alias: &Alias{
			Direction: "*out",
			Tenant:    "cgrates.org",
			Category:  "call",
			Account:   "set",
			Subject:   "set",
			Context:   "*rating",
			Values: AliasValues{&AliasValue{
				DestinationId: utils.ANY,
				Pairs:         AliasPairs{"Subject": map[string]string{"1234": "1235"}},
				Weight:        10,
			}},
		},
		Overwrite: false,
	}, &out); err != nil || out != utils.OK {
		t.Error("Error updateing alias: ", err, out)
	}
	if err := aliasService.Call("AliasesV1.GetAlias", &Alias{
		Direction: "*out",
		Tenant:    "cgrates.org",
		Category:  "call",
		Account:   "set",
		Subject:   "set",
		Context:   "*rating",
	}, r); err != nil ||
		len(r.Values) != 1 ||
		len(r.Values[0].Pairs) != 2 ||
		r.Values[0].Pairs["Subject"]["1234"] != "1235" ||
		r.Values[0].Pairs["Account"]["1234"] != "1235" {
		t.Errorf("Error getting alias: %+v", r.Values[0])
	}
	if err := aliasService.Call("AliasesV1.SetAlias", &AttrAddAlias{
		Alias: &Alias{
			Direction: "*out",
			Tenant:    "cgrates.org",
			Category:  "call",
			Account:   "set",
			Subject:   "set",
			Context:   "*rating",
			Values: AliasValues{&AliasValue{
				DestinationId: utils.ANY,
				Pairs:         AliasPairs{"Subject": map[string]string{"1111": "2222"}},
				Weight:        10,
			}},
		},
		Overwrite: false,
	}, &out); err != nil || out != utils.OK {
		t.Error("Error updateing alias: ", err, out)
	}
	if err := aliasService.Call("AliasesV1.GetAlias", &Alias{
		Direction: "*out",
		Tenant:    "cgrates.org",
		Category:  "call",
		Account:   "set",
		Subject:   "set",
		Context:   "*rating",
	}, r); err != nil || len(r.Values) != 1 || len(r.Values[0].Pairs) != 2 || r.Values[0].Pairs["Subject"]["1111"] != "2222" {
		t.Errorf("Error getting alias: %+v", r.Values[0].Pairs["Subject"])
	}
	if err := aliasService.Call("AliasesV1.SetAlias", &AttrAddAlias{
		Alias: &Alias{
			Direction: "*out",
			Tenant:    "cgrates.org",
			Category:  "call",
			Account:   "set",
			Subject:   "set",
			Context:   "*rating",
			Values: AliasValues{&AliasValue{
				DestinationId: "NAT",
				Pairs:         AliasPairs{"Subject": map[string]string{"3333": "4444"}},
				Weight:        10,
			}},
		},
		Overwrite: false,
	}, &out); err != nil || out != utils.OK {
		t.Error("Error updateing alias: ", err, out)
	}
	if err := aliasService.Call("AliasesV1.GetAlias", &Alias{
		Direction: "*out",
		Tenant:    "cgrates.org",
		Category:  "call",
		Account:   "set",
		Subject:   "set",
		Context:   "*rating",
	}, r); err != nil ||
		len(r.Values) != 2 ||
		len(r.Values[1].Pairs) != 1 ||
		r.Values[1].Pairs["Subject"]["3333"] != "4444" ||
		len(r.Values[0].Pairs) != 2 ||
		r.Values[0].Pairs["Subject"]["1111"] != "2222" ||
		r.Values[0].Pairs["Subject"]["1234"] != "1235" {
		t.Errorf("Error getting alias: %+v", r.Values[0])
	}
}

func TestAliasesLoadAlias(t *testing.T) {
	var response string
	cd := &CallDescriptor{
		Direction:   "*out",
		Tenant:      "cgrates.org",
		Category:    "call",
		Account:     "rif",
		Subject:     "rif",
		Destination: "444",
		ExtraFields: map[string]string{
			"Cli":   "0723",
			"Other": "stuff",
		},
	}
	err := LoadAlias(
		&AttrMatchingAlias{
			Direction:   "*out",
			Tenant:      "cgrates.org",
			Category:    "call",
			Account:     "dan",
			Subject:     "dan",
			Context:     "*rating",
			Destination: "444",
		}, cd, "ExtraFields")
	if err != nil || cd == nil {
		t.Error("Error getting alias: ", err, response)
	}
	if cd.Subject != "rif1" ||
		cd.ExtraFields["Cli"] != "0724" {
		t.Errorf("Aliases failed to change interface: %+v", cd)
	}
}

func TestAliasesCache(t *testing.T) {
	key := "*out:cgrates.org:call:remo:remo:*rating"
	_, err := accountingStorage.GetAlias(key, false)
	if err != nil {
		t.Error("Error getting alias: ", err)
	}
	a, found := cache2go.Get(utils.ALIASES_PREFIX + key)
	if !found || a == nil {
		//log.Printf("Test: %+v", cache2go.GetEntriesKeys(utils.REVERSE_ALIASES_PREFIX))
		t.Error("Error getting alias from cache: ", err, a)
	}
	rKey1 := "minuAccount*rating"
	_, err = accountingStorage.GetReverseAlias(rKey1, false)
	if err != nil {
		t.Error("Error getting reverse alias: ", err)
	}
	ra1, found := cache2go.Get(utils.REVERSE_ALIASES_PREFIX + rKey1)
	if !found || len(ra1.([]string)) != 2 {
		t.Error("Error getting reverse alias 1: ", ra1)
	}
	rKey2 := "minuSubject*rating"
	_, err = accountingStorage.GetReverseAlias(rKey2, false)
	if err != nil {
		t.Error("Error getting reverse alias: ", err)
	}
	ra2, found := cache2go.Get(utils.REVERSE_ALIASES_PREFIX + rKey2)
	if !found || len(ra2.([]string)) != 2 {
		t.Error("Error getting reverse alias 2: ", ra2)
	}
	accountingStorage.RemoveAlias(key)
	a, found = cache2go.Get(utils.ALIASES_PREFIX + key)
	if found {
		t.Error("Error getting alias from cache: ", found)
	}
	_, err = accountingStorage.GetReverseAlias(rKey1, false)
	if err != nil {
		t.Error("Error getting reverse alias: ", err)
	}
	ra1, found = cache2go.Get(utils.REVERSE_ALIASES_PREFIX + rKey1)
	if !found || len(ra1.([]string)) != 1 {
		t.Error("Error getting reverse alias 1: ", ra1)
	}
	_, err = accountingStorage.GetReverseAlias(rKey2, false)
	if err != nil {
		t.Error("Error getting reverse alias: ", err)
	}
	ra2, found = cache2go.Get(utils.REVERSE_ALIASES_PREFIX + rKey2)
	if !found || len(ra2.([]string)) != 1 {
		t.Error("Error getting reverse alias 2: ", ra2)
	}
}
