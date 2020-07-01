package host

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"hoster/hca"
	"hoster/mst"
	"os/exec"
	"strings"
	"time"
)

type Host struct {
	name   string
	client *ssh.Client
	hca    []hca.HCA
	mst    []mst.MST
}

func New(name string) (*Host, error) {
	log.Debugf("New: %s\n", name)
	host, err := newHost(name)
	if err != nil {
		log.Warnf("New Host: %s fail: %s\n", name, err)
		return nil, err
	}

	return host, nil
}

var sshPort = "22"
func newHost(name string) (*Host, error) {
	nameAndsshPort := fmt.Sprintf("%s:%s", name, sshPort)
	client, err := newsshClient(nameAndsshPort)
	if err != nil {
		log.Debugf("Host: %s Failed to dial: %s\n", name, err)
		return nil, err
	}

	return &Host{name, client, []hca.HCA{}, []mst.MST{}}, nil
}

func NewNonSsh(name string) *Host {
	log.Debugf("NewNonSsh: %s\n", name)
	return newNonSsh(name)
}

func newNonSsh(name string) *Host {
	return &Host{name, nil, nil, nil}
}

func newsshClient(name string) (*ssh.Client, error) {
	//var timeout = time.Second * 1 TODO: add timeout as var func
	config := &ssh.ClientConfig{ // need to be pointed flag
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("pass"), // need to be pointed to var
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(time.Second * 2),
	}
	client, err := ssh.Dial("tcp", name, config)
	if err != nil {
		log.Debugf("Host: %s Failed to dial: %s\n", name, err)
		//return "Failed to dial: ", err
		return nil, err
	}
	return client, nil
}

func (h *Host) Hca() []hca.HCA {
	return h.hca
}

func (h *Host) Name() string {
	return h.name
}

const (
	infinibandSysFsPath string = "/sys/class/infiniband/"
	mstPath 			string = "/dev/mst/"
	pciSlotName         string = "PCI_SLOT_NAME"
	read                string = "cat"
	list                string = "ls"
)

func (h *Host) newHCA(hardwareName string) (*hca.HCA, error) {

	var fwVer string
	var interfaceName string
	var addrLen string
	var driverVersion string
	var macAddr string
	var ifIndex string
	var uevent string
	var pciSlotID string
	var driver string

		//fmt.Printf("getInfinibandDevices fail: %v\n", byte(infinibandSysFsPath))
	infinibandPath := fmt.Sprintf("%s%s/", infinibandSysFsPath, hardwareName) // e.g. /sys/class/infiniband/mlx5_0/

	fwVer, err := h.SshRunCmd(fmt.Sprintf("%s %sfw_ver", read, infinibandPath))
	if err != nil {
		log.Debugf("getInfinibandDevices fail: %s\n", err)
		return nil, err
	}

	interfaceName, err = h.SshRunCmd(fmt.Sprintf("%s %sdevice/net/", list, infinibandPath))
	if err != nil {
		log.Debugf("getInfinibandDevices fail: %s\n", err)
		return nil, err
	}

	addrLen, err = h.SshRunCmd(fmt.Sprintf("%s %sdevice/net/%s/addr_len", read, infinibandPath, interfaceName))
	if err != nil {
		log.Debugf("getInfinibandDevices fail: %s\n", err)
		return nil, err
	}

	driverVersion, err = h.SshRunCmd(fmt.Sprintf("%s %sdevice/driver/module/version", read, infinibandPath))
	if err != nil {
		log.Debugf("getInfinibandDevices fail: %s\n", err)
		return nil, err
	}

	macAddr, err = h.SshRunCmd(fmt.Sprintf("%s %sdevice/net/%s/address", read, infinibandPath, interfaceName))
	if err != nil {
		log.Debugf("getInfinibandDevices fail: %s\n", err)
		return nil, err
	}

	ifIndex, err = h.SshRunCmd(fmt.Sprintf("%s %sdevice/net/%s/ifindex", read, infinibandPath, interfaceName))
	if err != nil {
		log.Debugf("getInfinibandDevices fail: %s\n", err)
		return nil, err
	}

	uevent, err = h.SshRunCmd(fmt.Sprintf("%s %sdevice/uevent", read, infinibandPath))
	if err != nil {
		log.Debugf("getInfinibandDevices fail: %s\n", err)
		return nil, err
	}

	ueventMap := getueventMap(uevent)

	pciSlotID = ueventMap[pciSlotName]
	driver = ueventMap["DRIVER"]

	MST, err := h.getMstDeviceByID(pciSlotID)
	if err != nil {
		log.Debugf("getInfinibandDevices fail: %s\n", err)
		return nil, err
	}

	return hca.New(hardwareName, fwVer, interfaceName, addrLen, driverVersion, macAddr, ifIndex, pciSlotID, driver, MST), nil

}

func (h *Host) NewHCA(hardwareName string) error {

	HCA, err := h.newHCA(hardwareName)
	if err != nil {
		log.Debugf("NewHCA fail: %s\n", err)
		return err
	}

	h.hca = append(h.hca, *HCA)

	return nil
}

func (h *Host) SshRunCmd(cmd string) (string, error) {
	log.Debugf("Host: %s cmd: \"%s\"\n", h.name, cmd)
	output, err := h.sshRunCmd(cmd)
	if err != nil {
		log.Debugf("Host: %s cmd: %s - failure: %s\n", h.name, cmd, err)
		return "", err
	}
	return output, nil
}

// TODO: decide on error handle and session client handle
func (h *Host) sshRunCmd(cmd string) (string, error) {
	session, err := h.client.NewSession()
	if err != nil {
		log.Debugf("Host: %s failed to create session: %s\n", h.name, err.Error())
		//h.state = false
		return "", err
	}

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(cmd); err != nil {
		log.Debugf("Host: %s failed to run: %s\n", h.name, err.Error())
		//h.res = false
		//return
		return "", err
	}
	//s := fmt.Sprint(session.Stdout)

	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Println("Recovered in f", r)
	//	}
	//}()

	//fmt.Printf("Host: %v stdout: %v\n", h.name, session.Stdout)
	// TODO: decide if session.Close() is required and if different approach is better
	defer session.Close()

	s := fmt.Sprint(session.Stdout)
	s = strings.TrimSuffix(s, "\n")

	return s, nil
}

func ck(err error) {
	if err != nil {
		log.Fatalf("ck() failure Fatal %s\n", err)
	}
}

func (h *Host) HostRESTART() error {
	var sleepTime = time.Duration(6000)
	log.Debugf("HostRESTART()\n")

	if err := h.hostOFF(); err != nil {
		log.Debugf("Host: %s failed to HostRESTART: %s\n", h.name, err)
		return err
	}

	log.Debugf("Host: %s Sleep %d ms\n", h.name, sleepTime)
	time.Sleep(sleepTime)

	if err := h.hostON(); err != nil {
		log.Debugf("Host: %s failed to HostRESTART: %s\n", h.name, err)
		return err
	}

	return nil
}

func (h *Host) hostOFF() error {
	verbose := false
	log.Debugf("Host: %s Starting - hostOFF()\n", h.name)

	output, err := exec.Command("/some_path/bin", "-c", "/some_path/some.conf", h.name, "OFF").Output()
	if err != nil {
		log.Debugf("Host: %s hostOFF() - failure Fatal %s\n", h.name, err)
		return err
	}

	if verbose {
		log.Debugf("Host: %s hostOFF() output %s\n", h.name, output)
		log.Debugf("Host: %s DONE hostOFF()\n", h.name)
	}

	return nil
}

func (h *Host) hostON() error {
	log.Debugf("Host: %s Starting - hostON()\n", h.name)

	verbose := false

	output, err := exec.Command("/some_path/bin", "-c", "/some_path/some.conf", h.name, "ON").Output()
	if err != nil {
		log.Debugf("Host: %s hostON() - failure Fatal %s\n", h.name, err)
		return err
	}

	if verbose {
		log.Debugf("Host: %s hostON() output %s\n", h.name, output)
		log.Debugf("Host: %s DONE hostON()\n", h.name)
	}

	return nil
}

func getueventMap(uevent string) map[string]string {
	ueventMap := map[string]string{}
	for _, stringField := range strings.Fields(uevent) {
		//fmt.Println(stringField)
		ueventField := strings.Split(stringField, "=")
		ueventMap[ueventField[0]] = ueventField[1]
	}
	//fmt.Printf("%#v\n", ueventMap)
	//fmt.Printf("%#v\n", ueventMap["PCI_SLOT_NAME"])
	//fmt.Printf("%#v\n", ueventMap["DRIVER"])
	return ueventMap
}

func (h *Host) GetInfinibandDevices() ([]string, error) {
	infinibandDevices, err := h.getInfinibandDevices()
	if err != nil {
		log.Debugf("GetInfinibandDevices fail: %s\n", err)
		return nil, err
	}
	return infinibandDevices, nil
}

func (h *Host) getInfinibandDevices() ([]string, error) {
	//session, err := h.client.NewSession()
	//if err != nil {
	//	log.Debugf("Failed to create session: %s\n", err)
	//	return nil, err
	//}
	//defer session.Close()
	//// Once a Session is created, you can execute a single command on
	//// the remote side using the Run method.
	//var b bytes.Buffer
	//session.Stdout = &b
	//if err := session.Run("ls /sys/class/infiniband/"); err != nil {
	//	log.Debugf("Failed to run: %s\n", err)
	//	return nil, err
	//}



	b, err := h.SshRunCmd(fmt.Sprintf("%s %s", list, infinibandSysFsPath))
	if err != nil {
		log.Debugf("getInfinibandDevices fail: %s\n", err)
		return nil, err
	}


	infinibandDevices := strings.Fields(b)
	return infinibandDevices, nil
}

// return []string of all /dev/mst/
func (h *Host) getMstDevices() ([]string, error) {
	log.Debugf("Host: %s - getMstDevices()\n", h.name)
	errCounter := 0
Retry:

	b, err := h.SshRunCmd(fmt.Sprintf("%s %s", list, mstPath))
	if err != nil {
		//log.Debugf("getInfinibandDevices fail: %s\n", err)
		//return nil, err
		if errCounter < 2 {
			log.Debugf("Host: %s - Attempt to correct - \"mst start\" will be done Retrying retry is: %d\n", h.name, errCounter)
			h.SshRunCmd("mst start")
			//fmt.Printf("getMstDevices sleep\n")
			//time.Sleep(1)
			errCounter++
			goto Retry
		} // TODO: add log verbose maybe in else note errCounter and brake
		return nil, err
	}

	mstDevices := strings.Fields(b)
	return mstDevices, nil
}

func (h *Host) NewMST() error {
	log.Debugf("Host: %s - NewMST()\n", h.name)
	err := h.newMST()
	if err != nil {
		log.Debugf("MSTDevices fail: %s\n", err)
		return err
	}
	return nil
}

// TODO: need to decide how to make newmst() from where and how to handle error, e.g. host repos or lazy by PrintHostHCATable
// TODO: /dev/mst change to const mstPath
func (h *Host) newMST() error {
	//log.Debugf("Host: %s - newMST()\n", h.name)
	mstDevices, err := h.getMstDevices() // return []string of all /dev/mst/
	if err != nil {
		log.Debugf("NewMST() Failed: %s\n", err)
		return err
	}

	for mstDevicesIndex := range mstDevices {
		var mstDevicesOutput string

		mstDevicesOutput, err = h.SshRunCmd(fmt.Sprintf("%s /dev/mst/%s", read, mstDevices[mstDevicesIndex]))
		if err != nil {
			log.Debugf("getInfinibandDevices fail: %s\n", err)
			return err
		}

		// TODO: Add logic find possible bugs and return error
		mstDevicesOutputSplited := strings.Split(mstDevicesOutput, "=")
		name := mstDevices[mstDevicesIndex]
		path := fmt.Sprintf("/dev/mst/%s", mstDevices[mstDevicesIndex])
		pciSlotID := strings.Split(mstDevicesOutputSplited[1], " ")[0]

		// for debug TODO: later change log verbose ?
		//fmt.Printf("name %#q\n", name)
		//fmt.Printf("path %#q\n", path)
		//fmt.Printf("pciSlotID %#q\n", pciSlotID)

		MST := mst.New(name, path, pciSlotID)

		h.mst = append(h.mst, MST)
	}
	return nil
}

// lookup mst by ID return error if not found
func (h *Host) getMstDeviceByID(pciSlotID string) (*mst.MST, error) {
	for mstIndex := range h.mst {
		if h.mst[mstIndex].PciSlotID() == pciSlotID {
			return &h.mst[mstIndex], nil
		}
	}
	return nil, errors.New("getMstDeviceByID() failure, pciSlotID not found")
}

func (h *Host) FwRestart() {
	log.Debugf("Host: %s Starting - FwRestart()\n", h.name)
	h.fwRestart()
}

func (h *Host) fwRestart() {
	verbose := false
	for hcaIndex := range h.hca {
		if h.hca[hcaIndex].IsFirstFunc() {

			command := fmt.Sprintf("mlxfwreset -d %s -l 3 r -y", h.hca[hcaIndex].Path())
			//log.Debugf("Host: %s Starting - fwRestart()\n", h.name)

			output, err := h.SshRunCmd(command)
			if err != nil {
				//errorCounter++
				log.Debugf("Host: %s fwRestart() - failure Fatal %s\n", h.name, err)
			}

			if verbose {
				log.Debugf("Host: %s command %s", h.name, command)
				log.Debugf("Host: %s fwRestart() output %s", h.name, output)
				log.Debugf("Host: %s DONE fwRestart()", h.name)
			}
		}
	}
}

func (h *Host) getFirstFuncOfHCA() {
	//var hca []HCA
	for hcaIndex := range h.hca {
		if h.hca[hcaIndex].IsFirstFunc() {
			h.hca[hcaIndex].Path()
			//hca = append(hca, h.hca[hcaIndex])
		}
	}
	//return hca
}
