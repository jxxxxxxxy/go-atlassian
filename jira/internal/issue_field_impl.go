package internal

import (
	"context"
	"fmt"
	model "github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/ctreminiom/go-atlassian/service"
	"github.com/ctreminiom/go-atlassian/service/jira"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func NewIssueFieldService(client service.Client, version string, configuration *IssueFieldConfigService, context *IssueFieldContextService) (*IssueFieldService, error) {

	if version == "" {
		return nil, model.ErrNoVersionProvided
	}

	return &IssueFieldService{
		internalClient: &internalIssueFieldServiceImpl{c: client, version: version},
		Configuration:  configuration,
		Context:        context,
	}, nil
}

type IssueFieldService struct {
	internalClient jira.Field
	Configuration  *IssueFieldConfigService
	Context        *IssueFieldContextService
}

func (i *IssueFieldService) Gets(ctx context.Context) ([]*model.IssueFieldScheme, *model.ResponseScheme, error) {
	return i.internalClient.Gets(ctx)
}

func (i *IssueFieldService) Create(ctx context.Context, payload *model.CustomFieldScheme) (*model.IssueFieldScheme, *model.ResponseScheme, error) {
	return i.internalClient.Create(ctx, payload)
}

func (i *IssueFieldService) Search(ctx context.Context, options *model.FieldSearchOptionsScheme, startAt, maxResults int) (*model.FieldSearchPageScheme, *model.ResponseScheme, error) {
	return i.internalClient.Search(ctx, options, startAt, maxResults)
}

func (i *IssueFieldService) Delete(ctx context.Context, fieldId string) (*model.TaskScheme, *model.ResponseScheme, error) {
	return i.internalClient.Delete(ctx, fieldId)
}

type internalIssueFieldServiceImpl struct {
	c       service.Client
	version string
}

func (i *internalIssueFieldServiceImpl) Gets(ctx context.Context) ([]*model.IssueFieldScheme, *model.ResponseScheme, error) {

	endpoint := fmt.Sprintf("rest/api/%v/field", i.version)

	request, err := i.c.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	var fields []*model.IssueFieldScheme
	response, err := i.c.Call(request, &fields)
	if err != nil {
		return nil, response, err
	}

	return fields, response, nil
}

func (i *internalIssueFieldServiceImpl) Create(ctx context.Context, payload *model.CustomFieldScheme) (*model.IssueFieldScheme, *model.ResponseScheme, error) {

	reader, err := i.c.TransformStructToReader(payload)
	if err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("rest/api/%v/field", i.version)

	request, err := i.c.NewRequest(ctx, http.MethodPost, endpoint, reader)
	if err != nil {
		return nil, nil, err
	}

	field := new(model.IssueFieldScheme)
	response, err := i.c.Call(request, field)
	if err != nil {
		return nil, response, err
	}

	return field, response, nil
}

func (i *internalIssueFieldServiceImpl) Search(ctx context.Context, options *model.FieldSearchOptionsScheme, startAt, maxResults int) (*model.FieldSearchPageScheme, *model.ResponseScheme, error) {

	params := url.Values{}
	params.Add("startAt", strconv.Itoa(startAt))
	params.Add("maxResults", strconv.Itoa(maxResults))

	if options != nil {

		if len(options.Expand) != 0 {
			params.Add("expand", strings.Join(options.Expand, ","))
		}

		if len(options.Types) != 0 {
			params.Add("type", strings.Join(options.Types, ","))
		}

		if len(options.IDs) != 0 {
			params.Add("id", strings.Join(options.IDs, ","))
		}

		if len(options.OrderBy) != 0 {
			params.Add("orderBy", options.OrderBy)
		}

		if len(options.Query) != 0 {
			params.Add("query", options.Query)
		}
	}

	endpoint := fmt.Sprintf("rest/api/%v/field/search?%v", i.version, params.Encode())

	request, err := i.c.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	page := new(model.FieldSearchPageScheme)
	response, err := i.c.Call(request, page)
	if err != nil {
		return nil, response, err
	}

	return page, response, nil
}

func (i *internalIssueFieldServiceImpl) Delete(ctx context.Context, fieldId string) (*model.TaskScheme, *model.ResponseScheme, error) {

	if fieldId == "" {
		return nil, nil, model.ErrNoFieldIDError
	}

	endpoint := fmt.Sprintf("rest/api/%v/field/%v", i.version, fieldId)

	request, err := i.c.NewRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	task := new(model.TaskScheme)
	response, err := i.c.Call(request, task)
	if err != nil {
		return nil, response, err
	}

	return task, response, nil
}
