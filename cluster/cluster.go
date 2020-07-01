package cluster

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"hoster/host"
	"os"
	"sort"
	"sync"
)

type Cluster struct {
	hostNames []string
	host      []*host.Host
}

// TODO: change type Cluster to export hosts?
func New(hostNames []string) Cluster {
	// TODO: add debug len of hostsSlice and rename hostsSlice to hostdb os different name?
	log.Debugf("New() - hosts len: %d\n", len(hostNames))
	hostsSlice := make([]*host.Host, 0, len(hostNames))
	return Cluster{hostNames, hostsSlice}
}

func (c *Cluster) InitHCA() error {
	log.Debugf("InitHCA\n")

	if err := c.MSTBuilder(); err != nil {
		log.Debugf("PrintHostHCATable->MSTBuilder() Failure: %s\n", err.Error())
		return err
	}

	if err := c.HCABuilder(); err != nil {
		log.Debugf("PrintHostHCATable->HCABuilder() Failure: %s\n", err.Error())
		return err
	}

	return nil
}

func (c *Cluster) Init() error {
	// append hosts
	log.Debugf("Init()\n")
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}
	for hostnameIndex := range c.hostNames {
		wg.Add(1)
		go func(c *Cluster, hostnameIndex int) error {
			defer wg.Done()
			h, err := host.New(c.hostNames[hostnameIndex])
			if err != nil {
				// TODO: logic bail if persistent is required - check global var if > 1 and make log
				//return
				log.Debugf("Host: %s failed to create h not appending due to: %s\n", c.hostNames[hostnameIndex], err.Error())
				//GError++ TODO: this add DATARACE need to be synced or removed
				return err
			}
			mutex.Lock()
			c.host = append(c.host, h)
			mutex.Unlock()
			return nil
		}(c, hostnameIndex)
	}
	wg.Wait()
	return nil
}

func (c *Cluster) InitNonSsh() error {
	log.Debugf("InitNonSsh()\n")
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}
	for hostnameIndex := range c.hostNames {
		wg.Add(1)
		go func(c *Cluster, hostnameIndex int) {
			defer wg.Done()
			h := host.NewNonSsh(c.hostNames[hostnameIndex])
			//if err != nil {
			//	// TODO: logic bail if persistent is required - check global var if > 1 and make log
			//	//return
			//	log.Errorf("Host: %s failed to create host not appending due to: %s\n", c.hostNames[hostnameIndex], err.Error())
			//	//GError++ TODO: this add DATARACE need to be synced or removed
			//	return err
			//}
			mutex.Lock()
			c.host = append(c.host, h)
			mutex.Unlock()
		}(c, hostnameIndex)
	}
	wg.Wait()
	return nil
}

func (c *Cluster) PrintHostHCATable() {
	//if err := MSTBuilder(c.host); err != nil {
	//	log.Debugf("PrintHostHCATable->MSTBuilder() Failure: %s\n", err.Error())
	//	return
	//}
	//
	//if err := HCABuilder(c.host); err != nil {
	//	log.Debugf("PrintHostHCATable->HCABuilder() Failure: %s\n", err.Error())
	//	return
	//}

	if len(c.host) < 1 {
		fmt.Printf("c.Cluster is 0. Try to use -C for config file or -H hostname\n")
		return
	}

	//TODO: change sort might be global or method?
	sort.Slice(c.host, func(i, j int) bool {
		return c.host[i].Name() < c.host[j].Name()
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "InterfaceName", "HardwareName", "FWVer", "DriverVersion", "MSTDevice", "AddrLen", "MacAddr", "IP", "IfIndex", "PciSlotID"})
	table.SetAutoMergeCells(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.SetRowLine(true)
	table.SetAutoFormatHeaders(false)

	data := getHostHCATableData(c)
	table.AppendBulk(data)
	fmt.Println("\nExtended Report:")
	table.Render()
}

func (c *Cluster) HostRestart() {
	c.hostRestart()
}

func FWRestart(Hosts []*host.Host, hostnames []string) {
	Hosts = NewHosts(hostnames, Hosts)
	// sort by name TODO: later flag if to sort or not
	// TODO: add ? LEN() and SWAP() might be also required
	sort.Slice(Hosts, func(i, j int) bool {
		return Hosts[i].Name() < Hosts[j].Name()
	})
	if err := MSTBuilder(Hosts); err != nil {
		log.Debugf("PrintHostHCATable->MSTBuilder() Failure: %s\n", err.Error())
		return
	}

	if err := HCABuilder(Hosts); err != nil {
		log.Debugf("PrintHostHCATable->HCABuilder() Failure: %s\n", err.Error())
		return
	}

	fwRestart(Hosts)
}

//func PrintHostHCATable(Hosts []*host.Host) {
//	if err := MSTBuilder(Hosts); err != nil {
//		log.Debugf("PrintHostHCATable->MSTBuilder() Failure: %s\n", err.Error())
//		return
//	}
//
//	if err := HCABuilder(Hosts); err != nil {
//		log.Debugf("PrintHostHCATable->HCABuilder() Failure: %s\n", err.Error())
//		return
//	}
//
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader([]string{"Name", "InterfaceName", "HardwareName", "FWVer", "DriverVersion", "MSTDevice", "AddrLen", "MacAddr", "IP", "IfIndex", "PciSlotID"})
//	table.SetAutoMergeCells(true)
//	table.SetAlignment(tablewriter.ALIGN_CENTER)
//	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
//	table.SetRowLine(true)
//	table.SetAutoFormatHeaders(false)
//
//	data := getHostHCATableData(Hosts)
//	table.AppendBulk(data)
//	fmt.Println("\nExtended Report:")
//	table.Render()
//}

func getHostHCATableData(cluster *Cluster) [][]string {
	var data [][]string
	for hostIndex := range cluster.host {
		for hcaIndex := range cluster.host[hostIndex].Hca() {
			data = append(data, []string{
				cluster.host[hostIndex].Name(),
				cluster.host[hostIndex].Hca()[hcaIndex].InterfaceName(),
				cluster.host[hostIndex].Hca()[hcaIndex].HardwareName(),
				cluster.host[hostIndex].Hca()[hcaIndex].FwVer(),
				cluster.host[hostIndex].Hca()[hcaIndex].DriverVersion(),
				cluster.host[hostIndex].Hca()[hcaIndex].MST.Name(),
				cluster.host[hostIndex].Hca()[hcaIndex].AddrLen(),
				cluster.host[hostIndex].Hca()[hcaIndex].MacAddr(),
				"Future",
				cluster.host[hostIndex].Hca()[hcaIndex].IfIndex(),
				cluster.host[hostIndex].Hca()[hcaIndex].PciSlotID(),
			})
		}
	}
	return data
}

func HCABuilder(Hosts []*host.Host) error {
	log.Debugf("HCABuilder()\n")
	var wg sync.WaitGroup
	for HostsIndex := range Hosts {
		wg.Add(1)
		go func(Hosts *[]*host.Host, HostsIndex int) error {
			defer wg.Done()
			h := (*Hosts)[HostsIndex]

			// TODO: need to decide how to make newmst() from where and how to handle error, e.g. host repos or lazy by PrintHostHCATable

			infinibandDevices, err := h.GetInfinibandDevices()
			if err != nil {
				log.Debugf("getInfinibandDevices fail: ", err)
				return err
			}
			// TODO: hardwareName infinibandDevices need to be renamed and might be build differently
			for infinibandDevicesIndex := range infinibandDevices {
				hardwareName := infinibandDevices[infinibandDevicesIndex]

				if err := h.NewHCA(hardwareName); err != nil {
					log.Debugf("Host: %s failed to run: %s\n", h.Name(), err.Error())
					return err
				}
			}
			return nil
		}(&Hosts, HostsIndex)
	}
	wg.Wait()

	return nil
}

func (c *Cluster) MSTBuilder() error {
	log.Debugf("MSTBuilder()\n")
	var wg sync.WaitGroup
	for hostIndex := range c.host {
		wg.Add(1)
		go func(c *Cluster, hostIndex int) error {
			defer wg.Done()
			h := c.host[hostIndex]

			// TODO: need to decide how to make newmst() from where and how to handle error, e.g. host repos or lazy by PrintHostHCATable
			err := h.NewMST()
			if err != nil {
				log.Debugf("MSTDevices fail: %s\n", err)
				return err
			}
			return nil
		}(c, hostIndex)
	}
	wg.Wait()
	return nil
}

func (c *Cluster) HCABuilder() error {
	log.Debugf("Cluster HCABuilder()\n")
	var wg sync.WaitGroup
	for hostIndex := range c.host {
		wg.Add(1)
		go func(c *Cluster, hostIndex int) error {
			defer wg.Done()
			h := c.host[hostIndex]

			// TODO: need to decide how to make newmst() from where and how to handle error, e.g. host repos or lazy by PrintHostHCATable

			infinibandDevices, err := h.GetInfinibandDevices()
			if err != nil {
				log.Debugf("GetInfinibandDevices fail: %s\n", err)
				return err
			}
			// TODO: hardwareName infinibandDevices need to be renamed and might be build differently
			for infinibandDevicesIndex := range infinibandDevices {
				hardwareName := infinibandDevices[infinibandDevicesIndex]

				if err := h.NewHCA(hardwareName); err != nil {
					log.Debugf("Host: %s failed to run: %s\n", h.Name(), err)
					return err
				}
			}
			return nil
		}(c, hostIndex)
	}
	wg.Wait()

	return nil
}

func MSTBuilder(Hosts []*host.Host) error {
	log.Debugf("MSTBuilder()\n")
	var wg sync.WaitGroup
	for HostsIndex := range Hosts {
		wg.Add(1)
		go func(Hosts *[]*host.Host, HostsIndex int) error {
			defer wg.Done()
			h := (*Hosts)[HostsIndex]

			// TODO: need to decide how to make newmst() from where and how to handle error, e.g. host repos or lazy by PrintHostHCATable
			err := h.NewMST()
			if err != nil {
				log.Debugf("MSTDevices fail: %s\n", err)
				return err
			}
			return nil
		}(&Hosts, HostsIndex)
	}
	wg.Wait()

	return nil
}

func MSTInstaller(Hosts []*host.Host) {
	log.Debugf("MSTInstaller()\n")
	var wg sync.WaitGroup
	errorCounter := 0
	verbose := false
	for HostsIndex := range Hosts {
		wg.Add(1)
		go func(Hosts *[]*host.Host, HostsIndex int, errorCounter *int, verbose bool) {
			defer wg.Done()
			h := (*Hosts)[HostsIndex]

			log.Debugf("Host: %s Starting - MSTInstaller()\n", h.Name())

			output, err := h.SshRunCmd("/some_path/release/mft/last_stable/install.sh")
			if err != nil {
				*errorCounter++
				log.Debugf("Host: %s MSTInstaller() - failure Fatal %s\n", h.Name(), err)
			}

			if verbose {
				log.Debugf("Host: %s MSTInstaller() - output %s\n", h.Name(), output)
				log.Debugf("Host: %s DONE MSTInstaller()\n", h.Name())
			}
		}(&Hosts, HostsIndex, &errorCounter, verbose)
	}
	wg.Wait()
	log.Infof("MSTInstaller is DONE len(hosts) is: %d errorCounter is: %d \n", len(Hosts), errorCounter)
}

func (c *Cluster) MSTStart() {
	log.Debugf("MSTStart()\n")
	var wg sync.WaitGroup
	errorCounter := 0
	verbose := false
	for HostsIndex := range c.host {
		wg.Add(1)
		go func(c *Cluster, HostsIndex int, errorCounter *int, verbose bool) {
			defer wg.Done()
			h := c.host[HostsIndex]

			log.Debugf("Host: %s Starting - MSTStart()\n", h.Name())

			output, err := h.SshRunCmd("mst start")
			if err != nil {
				*errorCounter++
				log.Debugf("Host: %s MSTStart() - failure Fatal %s\n", h.Name(), err)
			}

			if verbose {
				log.Debugf("Host: %s MSTStart() - output %s\n", h.Name(), output)
				log.Debugf("Host: %s DONE MSTStart()\n", h.Name())
			}
		}(c, HostsIndex, &errorCounter, verbose)
	}
	wg.Wait()
	log.Infof("MSTStart is DONE len(hosts) is: %d errorCounter is: %d \n", len(c.host), errorCounter)
}

func (c *Cluster) MSTStop() {
	log.Debugf("MSTStop()\n")
	var wg sync.WaitGroup
	errorCounter := 0
	verbose := false
	for HostsIndex := range c.host {
		wg.Add(1)
		go func(c *Cluster, HostsIndex int, errorCounter *int, verbose bool) {
			defer wg.Done()
			h := c.host[HostsIndex]

			log.Debugf("Host: %s Starting - MSTStop()\n", h.Name())

			output, err := h.SshRunCmd("mst stop")
			if err != nil {
				*errorCounter++
				log.Debugf("Host: %s MSTStop() - failure Fatal %s\n", h.Name(), err)
			}

			if verbose {
				log.Debugf("Host: %s MSTStop() - output %s\n", h.Name(), output)
				log.Debugf("Host: %s DONE MSTStop()\n", h.Name())
			}
		}(c, HostsIndex, &errorCounter, verbose)
	}
	wg.Wait()
	log.Infof("MSTStop is DONE len(hosts) is: %d errorCounter is: %d \n", len(c.host), errorCounter)
}

func (c *Cluster) MSTInstaller() {
	log.Debugf("MSTInstaller()\n")
	var wg sync.WaitGroup
	errorCounter := 0
	verbose := false
	command := fmt.Sprintf("/some_path/release/mft/last_stable/install.sh")

	for HostsIndex := range c.host {
		wg.Add(1)
		go func(c *Cluster, HostsIndex int, errorCounter *int, verbose bool) {
			defer wg.Done()
			h := c.host[HostsIndex]

			log.Debugf("Host: %s Starting - MSTInstaller() command: %s\n", h.Name(), command)


			output, err := h.SshRunCmd(command)
			if err != nil {
				*errorCounter++
				log.Debugf("Host: %s MSTInstaller() - failure Fatal %s\n", h.Name(), err)
			}

			if verbose {
				log.Debugf("Host: %s MSTInstaller() - output %s\n", h.Name(), output)
				log.Debugf("Host: %s DONE v()\n", h.Name())
			}
		}(c, HostsIndex, &errorCounter, verbose)
	}
	wg.Wait()
	log.Infof("MSTInstaller is DONE len(hosts) is: %d errorCounter is: %d \n", len(c.host), errorCounter)
}

func (c *Cluster) OFEDStart() {
	log.Debugf("OFEDStart()\n")
	var wg sync.WaitGroup
	errorCounter := 0
	verbose := false
	command := fmt.Sprintf("/etc/init.d/openibd start")

	for HostsIndex := range c.host {
		wg.Add(1)
		go func(c *Cluster, HostsIndex int, errorCounter *int, verbose bool) {
			defer wg.Done()
			h := c.host[HostsIndex]

			log.Debugf("Host: %s Starting - OFEDStart()\n", h.Name())

			output, err := h.SshRunCmd(command)
			if err != nil {
				*errorCounter++
				log.Debugf("Host: %s OFEDStart() - failure Fatal %s\n", h.Name(), err)
			}

			if verbose {
				log.Debugf("Host: %s OFEDStart() - output %s\n", h.Name(), output)
				log.Debugf("Host: %s DONE OFEDStart()\n", h.Name())
			}
		}(c, HostsIndex, &errorCounter, verbose)
	}
	wg.Wait()
	log.Infof("OFEDStart is DONE len(hosts) is: %d errorCounter is: %d \n", len(c.host), errorCounter)
}

func (c *Cluster) OFEDInstall(build, addFlag string, addKernelSupport, withNvmeF bool) {
	log.Debugf("OFEDInstall()\n")
	var wg sync.WaitGroup
	errorCounter := 0
	verbose := false
	cmd := cmdBuilderOFEDInstall(build, addKernelSupport, withNvmeF, addFlag)

	for HostsIndex := range c.host {
		wg.Add(1)
		go func(c *Cluster, HostsIndex int, errorCounter *int, verbose bool) {
			defer wg.Done()
			h := c.host[HostsIndex]

			log.Debugf("Host: %s Starting - OFEDInstall() cmd: %s\n", h.Name(), cmd)

			output, err := h.SshRunCmd(cmd)
			if err != nil {
				*errorCounter++
				log.Debugf("Host: %s OFEDInstall() - failure Fatal %s\n", h.Name(), err)
			}
			//
			if verbose {
				log.Debugf("Host: %s OFEDInstall() - output %s\n", h.Name(), output)
				log.Debugf("Host: %s DONE OFEDInstall()\n", h.Name())
			}
		}(c, HostsIndex, &errorCounter, verbose)
	}
	wg.Wait()
	log.Infof("OFEDInstall is DONE len(hosts) is: %d errorCounter is: %d \n", len(c.host), errorCounter)
}

func cmdBuilderOFEDInstall(build string, addKernelSupport bool, withNvmeF bool, addFlag string) string {
	cmd := fmt.Sprintf("build=%s /some_path/release/MLNX_OFED/mlnx_ofed_install", build)
	if addKernelSupport {
		cmd = fmt.Sprintf("%s %s", cmd, "--add-kernel-support")
	}
	if withNvmeF {
		cmd = fmt.Sprintf("%s %s", cmd, "--with-nvmf")
	}
	if addFlag != "" {
		cmd = fmt.Sprintf("%s %s", cmd, addFlag)
	}

	return cmd
}


func OFEDInstaller(Hosts []*host.Host, build string) {
	log.Debugf("OFEDInstaller()\n")
	var wg sync.WaitGroup
	errorCounter := 0
	verbose := false
	command := fmt.Sprintf("build=%s /some_path/release/MLNX_OFED/mlnx_ofed_install --add-kernel-support --with-nvmf", build)

	for HostsIndex := range Hosts {
		wg.Add(1)
		go func(Hosts *[]*host.Host, HostsIndex int, errorCounter *int, verbose bool, command string) {
			defer wg.Done()
			h := (*Hosts)[HostsIndex]

			log.Debugf("Host: %s Starting - OFEDInstaller()\n", h.Name())

			output, err := h.SshRunCmd(command)
			if err != nil {
				*errorCounter++
				log.Debugf("Host: %s OFEDInstaller() - failure Fatal %s\n", h.Name(), err)
			}

			if verbose {
				log.Debugf("Host: %s OFEDInstaller() - output %s\n", h.Name(), output)
				log.Debugf("Host: %s DONE OFEDInstaller()\n", h.Name())
			}
		}(&Hosts, HostsIndex, &errorCounter, verbose, command)
	}
	wg.Wait()
	log.Infof("OFEDInstaller is DONE len(hosts) is: %d errorCounter is: %d \n", len(Hosts), errorCounter)
}

func fwRestart(Hosts []*host.Host) {
	log.Debugf("fwRestart()\n")
	var wg sync.WaitGroup
	errorCounter := 0
	verbose := false
	for HostsIndex := range Hosts {
		wg.Add(1)
		go func(Hosts *[]*host.Host, HostsIndex int, errorCounter *int, verbose bool) {
			defer wg.Done()
			h := (*Hosts)[HostsIndex]
			h.FwRestart()

			if verbose {
				log.Debugf("Host: %s DONE fwRestart()\n", h.Name())
			}
		}(&Hosts, HostsIndex, &errorCounter, verbose)

	}
	wg.Wait()
	log.Debugf("fwRestart is DONE len(hosts) is: %d errorCounter is: %d \n", len(Hosts), errorCounter)
}

// TODO: this need be done not locally. It should be via os.exec not via ssh remotely
func hostRestart(Hosts []*host.Host) {
	log.Debugf("hostRestart()\n")
	var wg sync.WaitGroup
	var errorCounter = 0
	var mutex = &sync.Mutex{}
	for HostsIndex := range Hosts {
		wg.Add(1)
		go func(Hosts *[]*host.Host, HostsIndex int, errorCounter *int) {
			defer wg.Done()
			h := (*Hosts)[HostsIndex]

			if err := h.HostRESTART(); err != nil {
				log.Debugf("Host: %s failed to hostRestart: %s\n", h.Name(), err.Error())
				mutex.Lock()
				*errorCounter++
				mutex.Unlock()
			}
		}(&Hosts, HostsIndex, &errorCounter)

	}
	wg.Wait()
	log.Debugf("hostRestart is DONE len(hosts) is: %d errorCounter is: %d \n", len(Hosts), errorCounter)
}

func (c *Cluster) hostRestart() {
	log.Debugf("hostRestart()\n")
	var wg sync.WaitGroup
	var errorCounter = 0
	var mutex = &sync.Mutex{}
	for HostsIndex := range c.host {
		wg.Add(1)
		go func(c *Cluster, HostsIndex int, errorCounter *int) {
			defer wg.Done()
			h := c.host[HostsIndex]

			if err := h.HostRESTART(); err != nil {
				log.Warnf("Host: %s failed to hostRestart: %s\n", h.Name(), err)
				mutex.Lock()
				*errorCounter++
				mutex.Unlock()
			}
		}(c, HostsIndex, &errorCounter)

	}
	wg.Wait()
	log.Infof("hostRestart is DONE len(c.Cluster) is: %d errorCounter is: %d \n", len(c.host), errorCounter)
}
func (c *Cluster) OFEDStop() {
	log.Debugf("OFEDStart()\n")
	var wg sync.WaitGroup
	errorCounter := 0
	verbose := false
	command := fmt.Sprintf("/etc/init.d/openibd stop")

	for HostsIndex := range c.host {
		wg.Add(1)
		go func(c *Cluster, HostsIndex int, errorCounter *int, verbose bool) {
			defer wg.Done()
			h := c.host[HostsIndex]

			log.Debugf("Host: %s Starting - OFEDStop() command: %s\n", h.Name(), command)

			output, err := h.SshRunCmd(command)
			if err != nil {
				*errorCounter++
				log.Debugf("Host: %s OFEDStop() - failure Fatal %s\n", h.Name(), err)
			}

			if verbose {
				log.Debugf("Host: %s OFEDStop() - output %s\n", h.Name(), output)
				log.Debugf("Host: %s DONE OFEDStop()\n", h.Name())
			}
		}(c, HostsIndex, &errorCounter, verbose)
	}
	wg.Wait()
	log.Infof("OFEDStop is DONE len(hosts) is: %d errorCounter is: %d \n", len(c.host), errorCounter)

}
func (c *Cluster) OFEDRestart() {
	log.Debugf("OFEDRestart()\n")
	var wg sync.WaitGroup
	errorCounter := 0
	verbose := false
	command := fmt.Sprintf("/etc/init.d/openibd restart --forrce")

	for HostsIndex := range c.host {
		wg.Add(1)
		go func(c *Cluster, HostsIndex int, errorCounter *int, verbose bool) {
			defer wg.Done()
			h := c.host[HostsIndex]

			log.Debugf("Host: %s Starting - OFEDRestart() command: %s\n", h.Name(), command)

			output, err := h.SshRunCmd(command)
			if err != nil {
				*errorCounter++
				log.Debugf("Host: %s OFEDRestart() - failure Fatal %s\n", h.Name(), err)
			}

			if verbose {
				log.Debugf("Host: %s OFEDRestart() - output %s\n", h.Name(), output)
				log.Debugf("Host: %s DONE OFEDRestart()\n", h.Name())
			}
		}(c, HostsIndex, &errorCounter, verbose)
	}
	wg.Wait()
	log.Infof("OFEDStop is DONE len(hosts) is: %d errorCounter is: %d \n", len(c.host), errorCounter)

}

func OFEDIBDRestart(Hosts []*host.Host) {
	log.Debugf("OPEDIBDRestartAsync()\n")
	var wg sync.WaitGroup
	errorCounter := 0
	verbose := false
	command := fmt.Sprintf("/etc/init.d/openibd restart --force")

	for HostsIndex := range Hosts {
		wg.Add(1)
		go func(Hosts *[]*host.Host, HostsIndex int, errorCounter *int, verbose bool, command string) {
			defer wg.Done()
			h := (*Hosts)[HostsIndex]

			log.Debugf("Host: %s Starting - OPEDIBDRestartAsync()\n", h.Name())

			output, err := h.SshRunCmd(command)
			if err != nil {
				*errorCounter++
				log.Debugf("Host: %s OPEDIBDRestartAsync() - failure Fatal %s\n", h.Name(), err)
			}

			if verbose {
				log.Debugf("Host: %s OPEDIBDRestartAsync() - output %s\n", h.Name(), output)
				log.Debugf("Host: %s DONE OPEDIBDRestartAsync()\n", h.Name())
			}
		}(&Hosts, HostsIndex, &errorCounter, verbose, command)
	}
	wg.Wait()
	log.Debugf("OPEDIBDRestartAsync is DONE len(hosts) is: %d errorCounter is: %d \n", len(Hosts), errorCounter)
}

// temp TODO: make method of hosts and Cluster will crate Cluster by struct type or lazy when required
func NewHostsNonSsh(hostnames []string, Hosts []*host.Host) []*host.Host {
	// append hosts
	for hostname := range hostnames {
		h := host.NewNonSsh(hostnames[hostname])
		Hosts = append(Hosts, h)
	}
	return Hosts
}

func NewHosts(hostnames []string, Hosts []*host.Host) []*host.Host {
	// append hosts
	log.Debugf("Init()\n")
	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}
	for hostnameIndex := range hostnames {
		wg.Add(1)
		go func(Hosts *[]*host.Host, hostnameIndex int, hostnames []string) {
			defer wg.Done()
			h, err := host.New(hostnames[hostnameIndex])
			if err != nil {
				// TODO: logic bail if persistent is required - check global var if > 1 and make log
				//return
				log.Debugf("Host: %s failed to create host not appending due to: %s\n", hostnames[hostnameIndex], err.Error())
				//GError++ TODO: this add DATARACE need to be synced or removed
				return
			}
			mutex.Lock()
			*Hosts = append(*Hosts, h)
			mutex.Unlock()
		}(&Hosts, hostnameIndex, hostnames)
	}
	wg.Wait()
	return Hosts
}
