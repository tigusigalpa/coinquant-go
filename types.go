package coinquant

import "time"

// Meta is the standard API response envelope metadata.
type Meta struct {
	RequestID string  `json:"request_id"`
	Cursor    *string `json:"cursor,omitempty"`
	HasMore   bool    `json:"has_more,omitempty"`
}

// APIResponse is the wrapper shape returned by every JSON endpoint.
type APIResponse[T any] struct {
	Data  T       `json:"data"`
	Meta  *Meta   `json:"meta,omitempty"`
	Error *string `json:"error,omitempty"`
}

// PaginatedOptions contains shared pagination parameters.
type PaginatedOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
	Sort   string `json:"sort,omitempty"`
}

// HealthResponse is returned by GET /health.
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// LatestMessagePreview is the preview object on chats.
type LatestMessagePreview struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ArtifactCounts holds report/version counts on chats.
type ArtifactCounts struct {
	Reports          int `json:"reports"`
	StrategyVersions int `json:"strategy_versions"`
}

// Chat is a chat session.
type Chat struct {
	ID                   string                `json:"id"`
	Title                string                `json:"title"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
	Archived             bool                  `json:"archived"`
	Starred              bool                  `json:"starred"`
	Tags                 []string              `json:"tags"`
	LatestMessagePreview *LatestMessagePreview `json:"latest_message_preview,omitempty"`
	ArtifactCounts       *ArtifactCounts       `json:"artifact_counts,omitempty"`
	MessageCount         int                   `json:"message_count,omitempty"`
	StrategyID           *string               `json:"strategy_id,omitempty"`
	StrategyVersionID    *string               `json:"strategy_version_id,omitempty"`
	Metadata             map[string]any        `json:"metadata,omitempty"`
}

// ChatsList is the data container for GET /v1/chats.
type ChatsList struct {
	Chats []Chat `json:"chats"`
}

// ChatsListOptions are query parameters for GET /v1/chats.
type ChatsListOptions struct {
	Limit       int    `url:"limit,omitempty"`
	Cursor      string `url:"cursor,omitempty"`
	Sort        string `url:"sort,omitempty"`
	Archived    *bool  `url:"archived,omitempty"`
	HasStrategy *bool  `url:"has_strategy,omitempty"`
	HasReport   *bool  `url:"has_report,omitempty"`
	Tags        string `url:"tags,omitempty"`
}

// CreateChatRequest is the request body for POST /v1/chats.
type CreateChatRequest struct {
	Title          string         `json:"title,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
	InitialMessage string         `json:"initial_message,omitempty"`
}

// UpdateChatRequest is the request body for PATCH /v1/chats/{chat_id}.
type UpdateChatRequest struct {
	Title    string   `json:"title,omitempty"`
	Archived *bool    `json:"archived,omitempty"`
	Starred  *bool    `json:"starred,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

// ChatDeleteResponse is returned by DELETE /v1/chats/{chat_id}.
type ChatDeleteResponse struct {
	Deleted bool   `json:"deleted"`
	ID      string `json:"id"`
}

// Message is a chat message.
type Message struct {
	ID                string         `json:"id"`
	ChatID            string         `json:"chat_id"`
	Role              string         `json:"role"`
	Content           string         `json:"content"`
	CreatedAt         time.Time      `json:"created_at"`
	StrategyID        *string        `json:"strategy_id,omitempty"`
	StrategyVersionID *string        `json:"strategy_version_id,omitempty"`
	Artifacts         []Artifact     `json:"artifacts"`
	Metadata          map[string]any `json:"metadata,omitempty"`
}

// Artifact appears on messages.
type Artifact struct {
	Type              string  `json:"type"`
	StrategyID        *string `json:"strategy_id,omitempty"`
	StrategyVersionID *string `json:"strategy_version_id,omitempty"`
}

// MessagesList is the data container for GET /v1/chats/{chat_id}/messages.
type MessagesList struct {
	Messages []Message `json:"messages"`
}

// AppendMessageRequest is the request body for POST /v1/chats/{chat_id}/messages.
type AppendMessageRequest struct {
	Content  string         `json:"content"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// StreamingChatRequest is the request body for POST /v1/chats/{chat_id}/messages:stream.
type StreamingChatRequest struct {
	Content string         `json:"content"`
	Options map[string]any `json:"options,omitempty"`
}

// StreamingPromptRequest is the request body for POST /v1/prompts/stream.
type StreamingPromptRequest struct {
	Message string         `json:"message"`
	ChatID  *string        `json:"chat_id,omitempty"`
	Options map[string]any `json:"options,omitempty"`
}

// StrategyCondition is a single entry/exit condition within a strategy schema.
type StrategyCondition struct {
	ID               string         `json:"id"`
	Type             string         `json:"type"`
	Action           string         `json:"action"`
	Operator         string         `json:"operator"`
	Series1          string         `json:"series_1"`
	Series2          string         `json:"series_2,omitempty"`
	Series1Params    map[string]any `json:"series_1_params,omitempty"`
	Series2Params    map[string]any `json:"series_2_params,omitempty"`
	Series1Timeframe string         `json:"series_1_timeframe,omitempty"`
	Series2Timeframe string         `json:"series_2_timeframe,omitempty"`
}

// StrategySchema is the executable strategy payload.
type StrategySchema struct {
	Instrument        string              `json:"instrument"`
	Timeframe         string              `json:"timeframe"`
	StrategyName      string              `json:"strategy_name,omitempty"`
	InitialCapital    float64             `json:"initial_capital,omitempty"`
	PositionSizeType  string              `json:"position_size_type,omitempty"`
	PositionSizeValue float64             `json:"position_size_value,omitempty"`
	Conditions        []StrategyCondition `json:"conditions"`
}

// StrategyLatestBacktest is the lightweight backtest summary on strategy/version objects.
type StrategyLatestBacktest struct {
	BacktestID  string    `json:"backtest_id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	TotalReturn float64   `json:"total_return,omitempty"`
	SharpeRatio float64   `json:"sharpe_ratio,omitempty"`
	Status      string    `json:"status,omitempty"`
}

// StrategyVersionSummary is a short version descriptor.
type StrategyVersionSummary struct {
	ID             string                  `json:"id"`
	VersionNumber  int                     `json:"version_number"`
	State          string                  `json:"state"`
	IsPublished    bool                    `json:"is_published"`
	PublishedAt    *time.Time              `json:"published_at,omitempty"`
	LatestBacktest *StrategyLatestBacktest `json:"latest_backtest,omitempty"`
	Schema         *StrategySchema         `json:"schema,omitempty"`
}

// Strategy is a strategy container.
type Strategy struct {
	ID            string                  `json:"id"`
	ChatID        *string                 `json:"chat_id,omitempty"`
	Name          string                  `json:"name"`
	Description   string                  `json:"description,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
	State         string                  `json:"state"`
	VersionCount  int                     `json:"version_count"`
	LatestVersion *StrategyVersionSummary `json:"latest_version,omitempty"`
}

// StrategyVersion is a full strategy version snapshot.
type StrategyVersion struct {
	ID             string                  `json:"id"`
	StrategyID     string                  `json:"strategy_id"`
	VersionNumber  int                     `json:"version_number"`
	State          string                  `json:"state"`
	IsPublished    bool                    `json:"is_published"`
	PublishedAt    *time.Time              `json:"published_at,omitempty"`
	Title          string                  `json:"title,omitempty"`
	Description    string                  `json:"description,omitempty"`
	Schema         *StrategySchema         `json:"schema,omitempty"`
	LatestBacktest *StrategyLatestBacktest `json:"latest_backtest,omitempty"`
	CreatedAt      time.Time               `json:"created_at"`
}

// StrategiesList is the data container for GET /v1/strategies.
type StrategiesList struct {
	Strategies []Strategy `json:"strategies"`
}

// StrategyListOptions are query parameters for GET /v1/strategies.
type StrategyListOptions struct {
	Limit       int    `url:"limit,omitempty"`
	Cursor      string `url:"cursor,omitempty"`
	State       string `url:"state,omitempty"`
	IsPublished *bool  `url:"is_published,omitempty"`
	Q           string `url:"q,omitempty"`
	Sort        string `url:"sort,omitempty"`
}

// CreateStrategyRequest is the request body for POST /v1/strategies.
type CreateStrategyRequest struct {
	ChatID      string `json:"chat_id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateStrategyRequest is the request body for PATCH /v1/strategies/{strategy_id}.
type UpdateStrategyRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateStrategyVersionRequest is the request body for PATCH /v1/strategies/{strategy_id}/versions/{version_id}.
type UpdateStrategyVersionRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	IsPublished *bool  `json:"is_published,omitempty"`
}

// StrategyVersionsList is the data container for GET /v1/strategies/{strategy_id}/versions.
type StrategyVersionsList struct {
	Versions []StrategyVersion `json:"versions"`
}

// Backtest is a backtest resource.
type Backtest struct {
	BacktestID            string           `json:"backtest_id"`
	StrategyID            string           `json:"strategy_id"`
	StrategyVersionID     string           `json:"strategy_version_id"`
	Status                string           `json:"status"`
	CreatedAt             time.Time        `json:"created_at"`
	CompletedAt           *time.Time       `json:"completed_at,omitempty"`
	SourceBacktestID      *string          `json:"source_backtest_id,omitempty"`
	LatestBacktestSummary *BacktestSummary `json:"latest_backtest_summary,omitempty"`
}

// BacktestSummary contains lightweight backtest metrics.
type BacktestSummary struct {
	TotalReturn float64 `json:"total_return"`
	SharpeRatio float64 `json:"sharpe_ratio"`
	MaxDrawdown float64 `json:"max_drawdown"`
	TotalTrades int     `json:"total_trades"`
}

// BacktestsList is the data container for GET /v1/backtests.
type BacktestsList struct {
	Backtests []Backtest `json:"backtests"`
}

// BacktestListOptions are query parameters for GET /v1/backtests.
type BacktestListOptions struct {
	Limit             int    `url:"limit,omitempty"`
	Cursor            string `url:"cursor,omitempty"`
	Sort              string `url:"sort,omitempty"`
	StrategyID        string `url:"strategy_id,omitempty"`
	StrategyVersionID string `url:"strategy_version_id,omitempty"`
	Status            string `url:"status,omitempty"`
}

// CreateBacktestRequest is the request body for POST /v1/backtests.
type CreateBacktestRequest struct {
	StrategyVersionID string `json:"strategy_version_id"`
}

// BacktestResults contains the full results of a completed backtest.
type BacktestResults struct {
	BacktestID  string             `json:"backtest_id"`
	Status      string             `json:"status"`
	Metrics     map[string]float64 `json:"metrics"`
	EquityCurve []map[string]any   `json:"equity_curve,omitempty"`
	Trades      []map[string]any   `json:"trades,omitempty"`
	Drawdown    []map[string]any   `json:"drawdown,omitempty"`
}

// CompareBacktestsRequest is the request body for POST /v1/backtests/compare.
type CompareBacktestsRequest struct {
	BacktestIDs []string `json:"backtest_ids"`
}

// CompareBacktestsResponse is returned by POST /v1/backtests/compare.
type CompareBacktestsResponse struct {
	Backtests           []map[string]any `json:"backtests"`
	WinnerByTotalReturn string           `json:"winner_by_total_return,omitempty"`
}

// DuplicateBacktestResponse is returned by POST /v1/backtests/{backtest_id}/duplicate.
type DuplicateBacktestResponse struct {
	BacktestID       string    `json:"backtest_id"`
	SourceBacktestID string    `json:"source_backtest_id"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

// BacktestPollResult is the aggregate output of CreateBacktestAndWait.
type BacktestPollResult struct {
	Detail     *Backtest
	Results    *BacktestResults
	SummaryCSV string
	TradesCSV  string
}

// Report is a read-only report artifact.
type Report struct {
	ID        string         `json:"id"`
	ChatID    *string        `json:"chat_id,omitempty"`
	MessageID *string        `json:"message_id,omitempty"`
	Title     string         `json:"title"`
	Content   string         `json:"content,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// ReportsList is the data container for GET /v1/reports.
type ReportsList struct {
	Reports []Report `json:"reports"`
}

// ReportListOptions are query parameters for GET /v1/reports.
type ReportListOptions struct {
	Limit  int    `url:"limit,omitempty"`
	Cursor string `url:"cursor,omitempty"`
	ChatID string `url:"chat_id,omitempty"`
	Sort   string `url:"sort,omitempty"`
}

// TemplateCategory describes a template category.
type TemplateCategory struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// Template is a reusable strategy template.
type Template struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Category    *TemplateCategory `json:"category,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

// TemplatesList is the data container for GET /v1/templates.
type TemplatesList struct {
	Templates []Template `json:"templates"`
}

// TemplateListOptions are query parameters for GET /v1/templates.
type TemplateListOptions struct {
	Limit    int    `url:"limit,omitempty"`
	Cursor   string `url:"cursor,omitempty"`
	Category string `url:"category,omitempty"`
	Q        string `url:"q,omitempty"`
	GroupBy  string `url:"group_by,omitempty"`
}

// TemplateDetail is the full template definition.
type TemplateDetail struct {
	ID              string            `json:"id"`
	Title           string            `json:"title"`
	Slug            string            `json:"slug"`
	Prompt          string            `json:"prompt"`
	Description     string            `json:"description,omitempty"`
	LongDescription string            `json:"long_description,omitempty"`
	IsFeatured      bool              `json:"is_featured"`
	IsPro           bool              `json:"is_pro"`
	IsPopular       bool              `json:"is_popular"`
	StrategyType    string            `json:"strategy_type,omitempty"`
	Category        *TemplateCategory `json:"category,omitempty"`
}

// CreditLot describes a credit lot.
type CreditLot struct {
	LotID            string    `json:"lot_id"`
	Bucket           string    `json:"bucket"`
	RemainingCredits int       `json:"remaining_credits"`
	TotalCredits     int       `json:"total_credits"`
	ExpiresAt        time.Time `json:"expires_at"`
}

// CreditBalance is returned by GET /v1/credits.
type CreditBalance struct {
	AvailableCreditsTotal        int         `json:"available_credits_total"`
	AvailableCreditsSubscription int         `json:"available_credits_subscription"`
	AvailableCreditsTopup        int         `json:"available_credits_topup"`
	PendingCreditsSubscription   int         `json:"pending_credits_subscription"`
	Lots                         []CreditLot `json:"lots"`
	TickMaxDays                  int         `json:"tick_max_days"`
	ServerTime                   time.Time   `json:"server_time"`
}

// CreditUsageTotals holds aggregated credit usage.
type CreditUsageTotals struct {
	PromptsCredits int `json:"prompts_credits"`
	PromptsCount   int `json:"prompts_count"`
	CandleCredits  int `json:"candle_credits"`
	CandleCount    int `json:"candle_count"`
	TickCredits    int `json:"tick_credits"`
	TickMinutes    int `json:"tick_minutes"`
}

// CreditUsage is returned by GET /v1/credits/usage.
type CreditUsage struct {
	FromDate string            `json:"from_date"`
	ToDate   string            `json:"to_date"`
	Totals   CreditUsageTotals `json:"totals"`
	Series   []map[string]any  `json:"series"`
}

// CreditUsageOptions are query parameters for GET /v1/credits/usage.
type CreditUsageOptions struct {
	From    string `url:"from,omitempty"`
	To      string `url:"to,omitempty"`
	GroupBy string `url:"group_by,omitempty"`
}

// CreditTransaction is a single ledger entry.
type CreditTransaction struct {
	ID           string    `json:"id"`
	Direction    string    `json:"direction"`
	Amount       int       `json:"amount"`
	EventType    string    `json:"event_type"`
	ReferenceID  *string   `json:"reference_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	BalanceAfter int       `json:"balance_after"`
}

// CreditTransactionsList is the data container for GET /v1/credits/transactions.
type CreditTransactionsList struct {
	Transactions []CreditTransaction `json:"transactions"`
}

// CreditTransactionOptions are query parameters for GET /v1/credits/transactions.
type CreditTransactionOptions struct {
	Limit       int    `url:"limit,omitempty"`
	Cursor      string `url:"cursor,omitempty"`
	From        string `url:"from,omitempty"`
	To          string `url:"to,omitempty"`
	Direction   string `url:"direction,omitempty"`
	EventType   string `url:"event_type,omitempty"`
	ReferenceID string `url:"reference_id,omitempty"`
}

// TickEstimate is returned by GET /v1/credits/estimates/tick.
type TickEstimate struct {
	DataDays         int `json:"data_days"`
	TimeframeMinutes int `json:"timeframe_minutes"`
	EstimatedCredits int `json:"estimated_credits"`
	TickMaxDays      int `json:"tick_max_days"`
}

// TickEstimateOptions are query parameters for GET /v1/credits/estimates/tick.
type TickEstimateOptions struct {
	DataDays         int `url:"data_days"`
	TimeframeMinutes int `url:"timeframe_minutes"`
}

// LeaderboardEntry is a single leaderboard row.
type LeaderboardEntry struct {
	Rank        int    `json:"rank"`
	Username    string `json:"username"`
	TotalPoints int    `json:"total_points"`
}

// LeaderboardsData is the data container for GET /v1/community/leaderboards.
type LeaderboardsData struct {
	Entries     []LeaderboardEntry `json:"entries"`
	CurrentUser *LeaderboardEntry  `json:"current_user,omitempty"`
	Season      int                `json:"season"`
}

// LeaderboardOptions are query parameters for GET /v1/community/leaderboards.
type LeaderboardOptions struct {
	Season  int    `url:"season,omitempty"`
	Limit   int    `url:"limit,omitempty"`
	Cursor  string `url:"cursor,omitempty"`
	Include string `url:"include,omitempty"`
}

// CommunityBacktest is a public/shared backtest.
type CommunityBacktest struct {
	BacktestID        string           `json:"backtest_id"`
	StrategyName      string           `json:"strategy_name"`
	CreatedByUsername string           `json:"created_by_username"`
	PublishedAt       time.Time        `json:"published_at"`
	Metrics           map[string]any   `json:"metrics,omitempty"`
	FavouritesCount   int              `json:"favourites_count"`
	IsFavouritedByMe  bool             `json:"is_favourited_by_me"`
	Orders            []map[string]any `json:"orders,omitempty"`
	AssetClass        string           `json:"asset_class,omitempty"`
	SqsScore          int              `json:"sqs_score,omitempty"`
	Timeframe         string           `json:"timeframe,omitempty"`
	StartDate         *time.Time       `json:"start_date,omitempty"`
	EndDate           *time.Time       `json:"end_date,omitempty"`
}

// CommunityBacktestsList is the data container for GET /v1/community/backtests.
type CommunityBacktestsList struct {
	Backtests []CommunityBacktest `json:"backtests"`
}

// CommunityBacktestListOptions are query parameters for GET /v1/community/backtests.
type CommunityBacktestListOptions struct {
	Season     int    `url:"season,omitempty"`
	StrategyID string `url:"strategy_id,omitempty"`
	Q          string `url:"q,omitempty"`
	Sort       string `url:"sort,omitempty"`
	Limit      int    `url:"limit,omitempty"`
	Cursor     string `url:"cursor,omitempty"`
	Include    string `url:"include,omitempty"`
}

// CommunityProfile is returned by GET /v1/community/me.
type CommunityProfile struct {
	UserID             string `json:"user_id"`
	Username           string `json:"username"`
	TotalPoints        int    `json:"total_points"`
	Rank               int    `json:"rank"`
	PublishedBacktests int    `json:"published_backtests"`
	Followers          int    `json:"followers"`
}

// CommunityActivity is a community activity record.
type CommunityActivity struct {
	ID        string    `json:"id"`
	Task      string    `json:"task"`
	Status    string    `json:"status"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// CommunityActivitiesData is the data container for GET /v1/community/me/activities.
type CommunityActivitiesData struct {
	Activities []CommunityActivity `json:"activities"`
	Totals     map[string]any      `json:"totals,omitempty"`
}

// CommunityActivityOptions are query parameters for GET /v1/community/me/activities.
type CommunityActivityOptions struct {
	Limit   int    `url:"limit,omitempty"`
	Cursor  string `url:"cursor,omitempty"`
	Task    string `url:"task,omitempty"`
	Status  string `url:"status,omitempty"`
	Include string `url:"include,omitempty"`
}

// CommunityBacktestDetailOptions are query parameters for GET /v1/community/backtests/{backtest_id}.
type CommunityBacktestDetailOptions struct {
	Include string `url:"include,omitempty"`
}

// CommunityMeOptions are query parameters for GET /v1/community/me.
type CommunityMeOptions struct {
	Include string `url:"include,omitempty"`
}
