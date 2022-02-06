package main

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v6/plugin"

	"github.com/mattermost/mattermost-server/v6/model"
)

type Plugin struct {
	plugin.MattermostPlugin
}

func (p *Plugin) OnActivate() error {
	if p.MattermostPlugin.API.GetConfig().ServiceSettings.SiteURL == nil {
		p.MattermostPlugin.API.LogError("SiteURL must be set. Some features depend on it")
	}
	if err := p.MattermostPlugin.API.RegisterCommand(createHelloCommand()); err != nil {
		return errors.Wrapf(err, "failed to register command")
	}
	if err := p.MattermostPlugin.API.RegisterCommand(createDayOff()); err != nil {
		return errors.Wrapf(err, "\n\n\n\n\nfailed to register command -- Day Off")
	}
	return nil
}

func createDayOff() *model.Command {
	c := model.NewAutocompleteData(triggerDayOff, "", hintDayOff)

	return &model.Command{
		Trigger:          triggerDayOff,
		AutoComplete:     true,
		AutoCompleteDesc: "autocompletefordayoff",
		AutocompleteData: c,
	}
}

func createHelloCommand() *model.Command {
	c := model.NewAutocompleteData(triggerDayOff, "", "don't forget to say hello!")

	return &model.Command{
		Trigger:          triggerHello,
		AutoComplete:     true,
		AutoCompleteDesc: "autocompleteforhello",
		AutocompleteData: c,
	}
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {

} //
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	if strings.Contains(args.Command, "/hello") {
		return p.ExecuteHelloCommand(c, args)
	} else if strings.Contains(args.Command, "/izin") {
		return p.ExecuteDayOffCommand(c, args)
	}
	return nil, nil
}
func (p *Plugin) ExecuteHelloCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	// siteURL := p.GetSiteURL()
	input := strings.TrimSpace(strings.TrimPrefix(args.Command, "/hello"))

	if input == "" || input == "help" {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Just type 'hello' ",
		}, nil
	}
	resp := &model.CommandResponse{
		ResponseType: model.CommandResponseTypeInChannel,
		Text:         "Hi from mattermost plugin, you did it!",
	}
	return resp, nil
}

func (p *Plugin) ExecuteDayOffCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	// siteURL := p.GetSiteURL()
	input := strings.TrimSpace(strings.TrimPrefix(args.Command, "/"+triggerDayOff))

	if input == "" || input == "help" {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         hintDayOff,
		}, nil
	}
	dayOffRequest := getDayOffRequest(input)

	resp := &model.CommandResponse{
		ResponseType: model.CommandResponseTypeInChannel,
		Text:         dayOffRequest.toString() + "\nİzin talebiniz başarı ile sisteme girilmiştir.",
	}
	return resp, nil
}

func (p *Plugin) GetSiteURL() string {
	siteURL := ""
	ptr := p.MattermostPlugin.API.GetConfig().ServiceSettings.SiteURL
	if ptr != nil {
		siteURL = *ptr
	}
	return siteURL
}

func main() {
	plugin.ClientMain(&Plugin{})
}
