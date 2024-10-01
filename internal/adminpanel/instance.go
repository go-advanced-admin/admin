package adminpanel

import "fmt"

type Instance struct {
	InstanceID  interface{}
	Data        interface{}
	Model       *Model
	Permissions Permissions
}

func (i *Instance) GetLink() string {
	return fmt.Sprintf("%s/%v", i.Model.GetLink(), i.InstanceID)
}

func (i *Instance) GetFullLink() string {
	return i.Model.App.Panel.Config.GetLink(i.GetLink())
}

func (i *Instance) GetEditLink() string {
	return fmt.Sprintf("%s/edit", i.GetLink())
}

func (i *Instance) GetFullEditLink() string {
	return i.Model.App.Panel.Config.GetLink(i.GetEditLink())
}
