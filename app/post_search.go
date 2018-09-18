package app

import (
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
)

func (a *App) NewPostSearch(terms string, userId string, teamId string) (*model.PostSearchResults, *model.AppError) {
	esInterface := a.Elasticsearch
	if license := a.License(); esInterface != nil && *a.Config().ElasticsearchSettings.EnableSearching && license != nil && *license.Features.Elasticsearch {
		// We only allow the user to search in channels they are a member of.
		userChannels, err := a.GetChannelsForUser(teamId, userId, false)
		if err != nil {
			mlog.Error(fmt.Sprint(err))
			return nil, err
		}

		postIds, matches, err := a.Elasticsearch.NewSearchPosts(userChannels, terms)
		if err != nil {
			return nil, err
		}

		// Get the posts
		postList := model.NewPostList()
		if len(postIds) > 0 {
			if presult := <-a.Srv.Store.Post().GetPostsByIds(postIds); presult.Err != nil {
				return nil, presult.Err
			} else {
				for _, p := range presult.Data.([]*model.Post) {
					postList.AddPost(p)
				}
			}
		}

		for _, postId := range postIds {
			found := false

			for _, post := range postList.Posts {
				if postId == post.Id {
					found = true
					break
				}
			}

			if found {
				postList.AddOrder(postId)
			}
		}

		return model.MakePostSearchResults(postList, matches), nil
	}

	return nil, model.NewAppError("SearchPostsInTeam", "store.sql_post.search.disabled", nil, fmt.Sprintf("teamId=%v userId=%v", teamId, userId), http.StatusNotImplemented)
}
