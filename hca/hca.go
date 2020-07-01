package hca

import (
	"hoster/mst"
	"strings"
)

const (
	firstDevice string = ".0"
	secondDevice string = ".1"
)

type HCA struct {
	hardwareName  string
	fwVer         string
	interfaceName string
	addrLen       string
	driverVersion string
	macAddr       string
	ifIndex       string
	pciSlotID     string
	driver        string
	*mst.MST
}

func New(hardwareName, fwVer, interfaceName, addrLen, driverVersion, macAddr, ifIndex, pciSlotID, driver string, MST *mst.MST) *HCA {
	return &HCA{
		hardwareName,
		fwVer,
		interfaceName,
		addrLen,
		driverVersion,
		macAddr,
		ifIndex,
		pciSlotID,
		driver,
		MST,
	}
}

func (H *HCA) Driver() string { return H.driver }

func (H *HCA) PciSlotID() string { return H.pciSlotID }

func (H *HCA) IfIndex() string { return H.ifIndex }

func (H *HCA) MacAddr() string { return H.macAddr }

func (H *HCA) DriverVersion() string { return H.driverVersion }

func (H *HCA) AddrLen() string { return H.addrLen }

func (H *HCA) InterfaceName() string { return H.interfaceName }

func (H *HCA) FwVer() string { return H.fwVer }

func (H *HCA) HardwareName() string { return H.hardwareName }

// return true of false if HCA.FUNC is first == ".0
func (H *HCA) IsFirstFunc() bool { return strings.Contains(H.PciSlotID(), firstDevice) }
