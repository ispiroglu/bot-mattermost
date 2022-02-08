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
	botUserID string
}

const (
	botUserName    = "bilisim"
	botDisplayName = "bilisim"

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
	/*
		{
			x, err := p.MattermostPlugin.API.GetUserByUsername("bilisim")
			if err != nil {
				p.MattermostPlugin.API.LogError(err.Error())
			}
			err = p.MattermostPlugin.API.PermanentDeleteBot(x.Id)
			if err != nil {
				p.MattermostPlugin.API.LogError(err.Error())
			}
		} */

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

func (p *Plugin) ExecuteDayOffCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	input := strings.TrimSpace(strings.TrimPrefix(args.Command, "/"+triggerDayOff))

	if input == "" || input == "help" {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         hintDayOff,
		}, nil
	}
	dayOffRequest := getDayOffRequest(input) // TODO: dayOff formatı kontrol edilecek
	dayOffRequest.userID = args.UserId
	resp := &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         "İzin talebiniz başarı ile sisteme girilmiştir.",
	}

	/*
		- Talep gönderici
		- Talebin onaylanıp onaylanmadığını kontrol eden routine
	*/
	// ch := make(chan dayOff)

	p.sendOffReqToAdmin(dayOffRequest)
	// go p.lookingForReaction(postID, dayOffRequest)

	return resp, nil
}

// for ok := true; ok; ok = EXPR { }
/*
func (p *Plugin) lookingForReaction(postID string, off dayOff) {
	var reactions []model.Reaction
	for ok := true; ok; {
		p.API.LogInfo("Hala reaction bekliyorum. ")
		reactions, appErr := p.API.GetReactions(postID)
		if appErr != nil {
			p.API.LogError(appErr.Error())
			return
		}
		if len(reactions) > 0 {
			ok = false
		}
	}

	dm, _ := p.MattermostPlugin.API.GetDirectChannel(p.botUserID, off.userID)
	post := model.Post{
		ChannelId: dm.Id,
		UserId:    p.botUserID,
	}

	for _, v := range reactions {
		if v.EmojiName == "while_check_mark" {
			post.Message = "Izniniz onaylanmisitr. Iyi tatiller dileriz."
		} else {
			post.Message = "Maalesef izniniz onaylanmamistir."
		}
	}

	_, appErr := p.API.CreatePost(&post)
	if appErr != nil {
		p.API.LogError(appErr.Error())
	}
}
*/
func (p *Plugin) sendOffReqToAdmin(dayOffRequest dayOff) string {
	receiver, err := p.API.GetUserByUsername("evren")
	if err != nil {
		p.API.LogError(err.Error())
	}
	dm, _ := p.MattermostPlugin.API.GetDirectChannel(p.botUserID, receiver.Id)
	post := model.Post{
		ChannelId: dm.Id,
		UserId:    p.botUserID,
		Message:   "hello",
		Props:     dayOffPostInit(dayOffRequest),
	}

	var appErr *model.AppError
	_, appErr = p.API.CreatePost(&post)
	if appErr != nil {
		p.MattermostPlugin.API.LogError(appErr.Error())
	}
	return post.Id
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
