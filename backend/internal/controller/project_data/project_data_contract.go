package project_data

import (
	"context"

	v1 "engineServ/api/project_data/v1"
	projectDataLogic "engineServ/internal/logic/api/project_data"
	uploadLogic "engineServ/internal/logic/upload"
)

func (c *ControllerV1) ContractList(ctx context.Context, req *v1.ContractListReq) (*v1.ContractListRes, error) {
	items, total, err := c.data.ListContracts(ctx, req.ProjectId, req.Current, req.Size)
	if err != nil {
		return nil, err
	}
	records := make([]*v1.ContractItem, 0, len(items))
	for _, item := range items {
		attachments := make([]*v1.ContractAttachment, 0, len(item.Attachments))
		for _, attachment := range item.Attachments {
			attachments = append(attachments, &v1.ContractAttachment{
				Path: attachment.Path,
				Name: attachment.Name,
				Url:  uploadLogic.AbsoluteURL(ctx, attachment.Path),
			})
		}
		records = append(records, &v1.ContractItem{
			Id: item.Id, ProjectId: item.ProjectId, Name: item.Name, ContractNo: item.ContractNo,
			Amount: item.Amount, SignedAt: item.SignedAt, Remark: item.Remark,
			Attachments: attachments, CreatedAt: item.CreatedAt,
		})
	}
	return &v1.ContractListRes{Current: req.Current, Size: req.Size, Total: total, Records: records}, nil
}

func (c *ControllerV1) ContractCreate(ctx context.Context, req *v1.ContractCreateReq) (*v1.ContractCreateRes, error) {
	id, err := c.data.CreateContract(ctx, contractInput(req.ContractForm))
	if err != nil {
		return nil, err
	}
	return &v1.ContractCreateRes{Id: id}, nil
}

func (c *ControllerV1) ContractUpdate(ctx context.Context, req *v1.ContractUpdateReq) (*v1.ContractUpdateRes, error) {
	if err := c.data.UpdateContract(ctx, req.Id, contractInput(req.ContractForm)); err != nil {
		return nil, err
	}
	return &v1.ContractUpdateRes{Id: req.Id}, nil
}

func (c *ControllerV1) ContractDelete(ctx context.Context, req *v1.ContractDeleteReq) (*v1.ContractDeleteRes, error) {
	if err := c.data.DeleteContract(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.ContractDeleteRes{}, nil
}

func contractInput(form v1.ContractForm) *projectDataLogic.ContractInput {
	attachments := make([]*projectDataLogic.ContractAttachment, 0, len(form.Attachments))
	for _, attachment := range form.Attachments {
		if attachment != nil {
			attachments = append(attachments, &projectDataLogic.ContractAttachment{Path: attachment.Path, Name: attachment.Name})
		}
	}
	return &projectDataLogic.ContractInput{
		ProjectId: form.ProjectId, Name: form.Name, ContractNo: form.ContractNo, Amount: form.Amount,
		SignedAt: form.SignedAt, Remark: form.Remark, Attachments: attachments,
	}
}
