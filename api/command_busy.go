// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	"github.com/mattermost/platform/model"
)

type BusyProvider struct {
}

const (
	CMD_BUSY = "busy"
)

func init() {
	RegisterCommandProvider(&BusyProvider{})
}

func (me *BusyProvider) GetTrigger() string {
	return CMD_BUSY
}

func (me *BusyProvider) GetCommand(c *Context) *model.Command {
	return &model.Command{
		Trigger:          CMD_BUSY,
		AutoComplete:     true,
		AutoCompleteDesc: c.T("api.command_busy.desc"),
		DisplayName:      c.T("api.command_busy.name"),
	}
}

func (me *BusyProvider) DoCommand(c *Context, channelId string, message string) *model.CommandResponse {
	rmsg := c.T("api.command_busy.success")
	if len(message) > 0 {
		rmsg = message + " " + rmsg
	}
	SetStatusBusyIfNeeded(c.Session.UserId, true)

	return &model.CommandResponse{ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: rmsg}
}
