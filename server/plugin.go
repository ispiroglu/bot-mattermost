package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v6/plugin"

	"github.com/mattermost/mattermost-server/v6/model"
)

type Plugin struct {
	plugin.MattermostPlugin
	botUserID string
}

const (
	botUserName    = "bilisim_hr"
	botDisplayName = "bilisim_hr"

	// MatterpollPostType = "custom_matterpoll"
)

func (p *Plugin) OnActivate() error {
	if p.MattermostPlugin.API.GetConfig().ServiceSettings.SiteURL == nil {
		p.MattermostPlugin.API.LogError("SiteURL must be set. Some features depend on it")
	}

	bot := &model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
	}
	botUserID, appErr := p.MattermostPlugin.API.CreateBot(bot)
	if appErr != nil {
		return errors.Wrapf(appErr, "failed to create bot.")
	}
	p.botUserID = botUserID.UserId

	if err := p.MattermostPlugin.API.RegisterCommand(createHelloCommand()); err != nil {
		return errors.Wrapf(err, "failed to register command - hello command")
	}
	if err := p.MattermostPlugin.API.RegisterCommand(createDayOff()); err != nil {
		return errors.Wrapf(err, "\n\n\n\n\nfailed to register command - Day Off")
	}
	return nil
}

func (p *Plugin) OnDeactivate() error {
	if err := p.MattermostPlugin.API.PermanentDeleteBot(p.botUserID); err != nil {
		return errors.Wrapf(err, "failed to delete bot.")
	}
	return nil
}

func (p *Plugin) OnCrash() error {
	return p.OnDeactivate()
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

// http://localhost:8065/hooks/ooce1nmw6i8tdqgxdu1tkqu3gw
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
func JSONMarshal(t interface{}) *bytes.Buffer {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(t)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("IN json", bf.String())
	return bf
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
	dayOffRequest := getDayOffRequest(input) // TODO: dayOff formatı kontrol edilecek
	fmt.Println(dayOffRequest.toString())
	resp := &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         "İzin talebiniz başarı ile sisteme girilmiştir.",
	}

	ch, _ := p.MattermostPlugin.API.GetDirectChannel(p.botUserID, args.UserId)
	post := model.Post{
		ChannelId: ch.Id,
		UserId:    p.botUserID,
		Message:   "selamlar",
	}
	if err := p.MattermostPlugin.API.SendEphemeralPost(args.UserId, &post); err != nil {
		p.MattermostPlugin.API.LogError("Could'n send ephemeral post inside bicicic")
	}

	copy := post.Clone()
	copy.StripActionIntegrations()

	fmt.Println(copy)
	hreps, err := http.Post("http://localhost:8065/hooks/ooce1nmw6i8tdqgxdu1tkqu3gw", "application/json",
		JSONMarshal(copy))
	if err != nil {
		p.MattermostPlugin.API.LogError(err.Error())
	}

	defer hreps.Body.Close()

	body, err := ioutil.ReadAll(hreps.Body)
	if err != nil {
		p.MattermostPlugin.API.LogError(err.Error())
	}
	p.MattermostPlugin.API.LogInfo("RESP BODY ---", string(body))

	/*defer func() { // Gonderilmeyebilir.

		if err := p.MattermostPlugin.API.SendEphemeralPost(args.UserId, &post); err != nil {
			p.MattermostPlugin.API.LogError("Could'n send ephemeral post inside bicicic")
		}
	}()*/ /*
		go func() {
			var reactions []*model.Reaction
			var err *model.AppError
			reactions, err = p.MattermostPlugin.API.GetReactions(post.Id)

			if err != nil {
				p.MattermostPlugin.API.LogError(err.Message, "Couldn't get plugins.")
			}

			for len(reactions) == 0 {
				reactions, err = p.MattermostPlugin.API.GetReactions(post.Id)
				if err != nil {
					p.MattermostPlugin.API.LogError(err.Message, "Couldn't get plugins.")
				}
			}

			if reactions[0].EmojiName == "while_check_mark" {
				post2 := model.Post{
					ChannelId: args.ChannelId,
					UserId:    p.botUserID,
					Message:   "dayOffRequest.toString()",
				}
				if err := p.MattermostPlugin.API.SendEphemeralPost(args.UserId, &post2); err != nil {
					p.MattermostPlugin.API.LogError("Could'n send ephemeral post inside bicicic")
				}
			}
		}()*/

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
