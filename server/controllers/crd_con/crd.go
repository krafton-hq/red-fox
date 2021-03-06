package crd_con

import (
	"context"
	"strings"

	"github.com/krafton-hq/red-fox/apis/crds"
	"github.com/krafton-hq/red-fox/apis/idl_common"
	"github.com/krafton-hq/red-fox/server/controllers/utils"
	"github.com/krafton-hq/red-fox/server/pkg/domain_helper"
	"github.com/krafton-hq/red-fox/server/pkg/errors"
	"github.com/krafton-hq/red-fox/server/services/services"
)

type Controller struct {
	crds.UnimplementedCustomResourceDefinitionServerServer

	service services.ClusterService[*crds.CustomResourceDefinition]
}

func NewController(service services.ClusterService[*crds.CustomResourceDefinition]) *Controller {
	return &Controller{service: service}
}

func (c *Controller) GetCustomResourceDefinition(ctx context.Context, req *idl_common.SingleObjectReq) (*crds.GetCustomResourceDefinitionRes, error) {
	if req.Name == "" {
		return &crds.GetCustomResourceDefinitionRes{CommonRes: utils.CommonResNotEmpty("name")}, nil
	}

	crd, err := c.service.Get(ctx, req.Name)
	if err != nil {
		return &crds.GetCustomResourceDefinitionRes{CommonRes: utils.CommonResWithErrorTypes(err)}, nil
	}

	return &crds.GetCustomResourceDefinitionRes{
		CommonRes: &idl_common.CommonRes{Message: "Get Crd Success"},
		Crd:       crd,
	}, nil
}

func (c *Controller) ListCustomResourceDefinitions(ctx context.Context, req *idl_common.ListObjectReq) (*crds.ListCustomResourceDefinitionsRes, error) {
	if req.Namespace != "" {
		return &crds.ListCustomResourceDefinitionsRes{CommonRes: utils.CommonResEmpty("Namespace", "namespace")}, nil
	}

	crdList, err := c.service.List(ctx, req.LabelSelectors)
	if err != nil {
		return &crds.ListCustomResourceDefinitionsRes{CommonRes: utils.CommonResWithErrorTypes(err)}, nil
	}

	return &crds.ListCustomResourceDefinitionsRes{
		CommonRes: &idl_common.CommonRes{Message: "List Crds Success"},
		Crds:      crdList,
	}, nil
}

func (c *Controller) CreateCustomResourceDefinition(ctx context.Context, req *crds.CreateCustomResourceDefinitionReq) (*idl_common.CommonRes, error) {
	if err := domain_helper.ValidationMetadatable(req.Crd); err != nil {
		return utils.CommonResWithErrorTypes(err), nil
	}
	if err := domain_helper.ValidationQualifiedName(req.Crd); err != nil {
		return utils.CommonResWithErrorTypes(err), nil
	}
	if err := domain_helper.ValidationCrdSpec(req.Crd.Spec); err != nil {
		return utils.CommonResWithErrorTypes(err), nil
	}

	if domain_helper.IsSystemGvk(req.Crd.Spec.Gvk) {
		return utils.InvalidArguments(errors.NewInvalidArguments("User Cannot Modify System Gvk")), nil
	}

	err := c.service.Create(ctx, req.Crd)
	if err != nil {
		return utils.CommonResWithErrorTypes(err), nil
	}

	return &idl_common.CommonRes{Message: "Create Crd Success"}, nil
}

func (c *Controller) UpdateCustomResourceDefinition(ctx context.Context, req *crds.UpdateCustomResourceDefinitionReq) (*idl_common.CommonRes, error) {
	if err := domain_helper.ValidationMetadatable(req.Crd); err != nil {
		return utils.CommonResWithErrorTypes(err), nil
	}
	if err := domain_helper.ValidationQualifiedName(req.Crd); err != nil {
		return utils.CommonResWithErrorTypes(err), nil
	}
	if err := domain_helper.ValidationCrdSpec(req.Crd.Spec); err != nil {
		return utils.CommonResWithErrorTypes(err), nil
	}

	if domain_helper.IsSystemGvk(req.Crd.Spec.Gvk) {
		return utils.InvalidArguments(errors.NewInvalidArguments("User Cannot Modify System Gvk")), nil
	}

	err := c.service.Update(ctx, req.Crd)
	if err != nil {
		return utils.CommonResWithErrorTypes(err), nil
	}

	return &idl_common.CommonRes{Message: "Update Crd Success"}, nil
}

func (c *Controller) DeleteCustomResourceDefinition(ctx context.Context, req *idl_common.SingleObjectReq) (*idl_common.CommonRes, error) {
	if req.Name == "" {
		return utils.CommonResNotEmpty("name"), nil
	}

	lowerKind, group, err := domain_helper.ParseGvkName(req.Name)
	if err != nil {
		return utils.CommonResWithErrorTypes(err), nil
	}

	for _, spec := range domain_helper.GetSystemGvks() {
		if strings.EqualFold(spec.Kind, lowerKind) && spec.Group == group {
			return utils.InvalidArguments(errors.NewInvalidArguments("User Cannot Delete System Gvk")), nil
		}
	}

	err = c.service.Delete(ctx, req.Name)
	if err != nil {
		return utils.CommonResWithErrorTypes(err), nil
	}

	return &idl_common.CommonRes{Message: "Delete Crd Success"}, nil
}
