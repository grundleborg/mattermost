// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import Client from 'client/web_client.jsx';

// Deprecated. This method will be removed in 3.7 as it is part of the Segment Analytics feature.
export function deprecatedTrack(category, action, label, property, value) {
    Client.deprecatedTrack(category, action, label, property, value);
}

// Deprecated. This method will be removed in 3.7 as it is part of the Segment Analytics feature.
export function deprecatedTrackPage() {
    Client.deprecatedTrackPage();
}
