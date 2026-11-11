package project_list

import (
	"context"
	"strings"

	v1 "engineServ/api/project_list/v1"
)

func (c *ControllerV1) Options(ctx context.Context, req *v1.OptionsReq) (res *v1.OptionsRes, err error) {
	customers, managers, err := c.project.Options(ctx)
	if err != nil {
		return nil, err
	}

	res = &v1.OptionsRes{
		Customers: make([]*v1.OptionItem, 0, len(customers)),
		Managers:  make([]*v1.OptionItem, 0, len(managers)),
	}

	for _, customer := range customers {
		label := customer.Name
		if phone := strings.TrimSpace(customer.Phone); phone != "" {
			label = label + " · " + phone
		}

		res.Customers = append(res.Customers, &v1.OptionItem{Value: customer.Id, Label: label})
	}

	for _, user := range managers {
		label := strings.TrimSpace(user.Nickname)
		if label == "" {
			label = user.Username
		}

		res.Managers = append(res.Managers, &v1.OptionItem{Value: user.Id, Label: label})
	}

	return res, nil
}
