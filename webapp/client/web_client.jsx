// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import Client from './client.jsx';

import TeamStore from 'stores/team_store.jsx';
import BrowserStore from 'stores/browser_store.jsx';

import * as GlobalActions from 'actions/global_actions.jsx';
import {reconnect} from 'actions/websocket_actions.jsx';

import request from 'superagent';

const HTTP_UNAUTHORIZED = 401;

class WebClientClass extends Client {
    constructor() {
        super();
        this.enableLogErrorsToConsole(true);
        this.hasInternetConnection = true;
        TeamStore.addChangeListener(this.onTeamStoreChanged.bind(this));
    }

    onTeamStoreChanged() {
        this.setTeamId(TeamStore.getCurrentId());
    }

    // Deprecated. This method will be removed in 3.7 as it is part of the Segment Analytics feature.
    deprecatedTrack(category, action, label, property, value) {
        if (window.mm_config.SegmentDeveloperKey != null && window.mm_config.SegmentDeveloperKey !== '') {
            if (global.window && global.window.analytics) {
                global.window.analytics.track(action, {category, label, property, value});
            }
        }
    }

    // Deprecated. This method will be removed in 3.7 as it is part of the Segment Analytics feature.
    deprecatedTrackPage() {
        if (window.mm_config.SegmentDeveloperKey != null && window.mm_config.SegmentDeveloperKey !== '') {
            if (global.window && global.window.analytics) {
                global.window.analytics.page();
            }
        }
    }

    trackEvent(category, event, props) {
        if (window.mm_config.SegmentDeveloperKey != null && window.mm_config.SegmentDeveloperKey !== '') {
            // Segment is in use for analytics, so diagnostics is disabled, making this function a no-op.
            return;
        }

        if (global.window && global.window.analytics) {
            const properties = Object.assign({category, type: event}, props);
            const options = {
                context: {
                    ip: '0.0.0.0'
                },
                page: {
                    path: '',
                    referrer: '',
                    search: '',
                    title: '',
                    url: ''
                }
            };
            global.window.analytics.track('event', properties, options);
        }
    }

    handleError(err, res) {
        if (res.body.id === 'api.context.mfa_required.app_error') {
            window.location.reload();
            return;
        }

        if (err.status === HTTP_UNAUTHORIZED && res.req.url !== '/api/v3/users/login') {
            GlobalActions.emitUserLoggedOutEvent('/login');
        }

        if (err.status == null) {
            this.hasInternetConnection = false;
        }
    }

    handleSuccess = (res) => { // eslint-disable-line no-unused-vars
        if (res && !this.hasInternetConnection) {
            reconnect();
            this.hasInternetConnection = true;
        }
    }

    // not sure why but super.login doesn't work if using an () => arrow functions.
    // I think this might be a webpack issue.
    webLogin(loginId, password, token, success, error) {
        this.login(
            loginId,
            password,
            token,
            (data) => {
                this.deprecatedTrack('api', 'api_users_login_success', '', 'login_id', loginId);
                BrowserStore.signalLogin();

                if (success) {
                    success(data);
                }
            },
            (err) => {
                this.deprecatedTrack('api', 'api_users_login_fail', '', 'login_id', loginId);
                if (error) {
                    error(err);
                }
            }
        );
    }

    webLoginByLdap(loginId, password, token, success, error) {
        this.loginByLdap(
            loginId,
            password,
            token,
            (data) => {
                this.deprecatedTrack('api', 'api_users_login_success', '', 'login_id', loginId);
                BrowserStore.signalLogin();

                if (success) {
                    success(data);
                }
            },
            (err) => {
                this.deprecatedTrack('api', 'api_users_login_fail', '', 'login_id', loginId);
                if (error) {
                    error(err);
                }
            }
        );
    }

    getYoutubeVideoInfo(googleKey, videoId, success, error) {
        request.get('https://www.googleapis.com/youtube/v3/videos').
        query({part: 'snippet', id: videoId, key: googleKey}).
        end((err, res) => {
            if (err) {
                return error(err);
            }

            if (!res.body) {
                console.error('Missing response body for getYoutubeVideoInfo'); // eslint-disable-line no-console
            }

            return success(res.body);
        });
    }
}

var WebClient = new WebClientClass();
export default WebClient;
