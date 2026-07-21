package coinquant

import (
	"context"
	"fmt"
	"net/http"
)

// ListLeaderboards returns community leaderboard entries.
func (c *Client) ListLeaderboards(ctx context.Context, opts LeaderboardOptions) (*LeaderboardsData, error) {
	return doJSON[LeaderboardsData](ctx, c, http.MethodGet, "/v1/community/leaderboards", nil, opts)
}

// ListCommunityBacktests returns public/shared backtests.
func (c *Client) ListCommunityBacktests(ctx context.Context, opts CommunityBacktestListOptions) (*CommunityBacktestsList, error) {
	return doJSON[CommunityBacktestsList](ctx, c, http.MethodGet, "/v1/community/backtests", nil, opts)
}

// GetCommunityBacktest returns the public-facing view of a shared backtest.
func (c *Client) GetCommunityBacktest(ctx context.Context, backtestID string, include string) (*CommunityBacktest, error) {
	return doJSON[CommunityBacktest](ctx, c, http.MethodGet, fmt.Sprintf("/v1/community/backtests/%s", backtestID), nil, CommunityBacktestDetailOptions{Include: include})
}

// GetCommunityMe returns the current user's community-facing stats.
func (c *Client) GetCommunityMe(ctx context.Context, include string) (*CommunityProfile, error) {
	return doJSON[CommunityProfile](ctx, c, http.MethodGet, "/v1/community/me", nil, CommunityMeOptions{Include: include})
}

// ListCommunityActivities returns community activity for the current user.
func (c *Client) ListCommunityActivities(ctx context.Context, opts CommunityActivityOptions) (*CommunityActivitiesData, error) {
	return doJSON[CommunityActivitiesData](ctx, c, http.MethodGet, "/v1/community/me/activities", nil, opts)
}
