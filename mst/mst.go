package mst

type MST struct {
	name      string
	path      string
	pciSlotID string
}

func New(name, path, pciSlotID string) MST { return MST{name, path, pciSlotID} }

func (M *MST) PciSlotID() string { return M.pciSlotID }

func (M *MST) Path() string { return M.path }

func (M *MST) Name() string { return M.name }