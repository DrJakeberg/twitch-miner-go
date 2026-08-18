package miner

import (
	"context"
	"testing"

	"github.com/Guliveer/twitch-miner-go/internal/model"
)

// pointsEarnedMessage builds the PubSub message Twitch sends when points land.
func pointsEarnedMessage(reasonCode string, amount int) *model.Message {
	return &model.Message{
		Type: model.MsgTypePointsEarned,
		Data: map[string]any{
			"balance": map[string]any{"balance": 1000},
			"point_gain": map[string]any{
				"total_points": amount,
				"reason_code":  reasonCode,
			},
		},
	}
}

func onlineStreamerWithBroadcast(username, broadcastID string) *model.Streamer {
	s := model.NewStreamer(username)
	s.Settings = model.DefaultStreamerSettings()
	s.IsOnline = true
	s.Stream.BroadcastID = broadcastID
	s.Stream.IsWatchStreakMissing = true
	return s
}

// The reason code arrives as "WATCH_STREAK" but is stored under the mapped
// event name, which is how the flag stopped being cleared in the first place.
func TestWatchStreakPayoutClearsTheFlag(t *testing.T) {
	m, _ := newTestMiner(t)
	s := onlineStreamerWithBroadcast("krissi", "b-1")

	m.handleCommunityPoints(context.Background(), pointsEarnedMessage("WATCH_STREAK", 350), s)

	if s.Stream.IsWatchStreakMissing {
		t.Fatal("a WATCH_STREAK payout must release the channel from the rotation")
	}
	if got := s.History[string(model.EventGainForWatchStreak)]; got == nil || got.Amount != 350 {
		t.Errorf("expected the payout recorded under the mapped event name, got %+v", got)
	}
}

func TestOrdinaryWatchPayoutLeavesTheFlag(t *testing.T) {
	m, _ := newTestMiner(t)
	s := onlineStreamerWithBroadcast("krissi", "b-1")

	m.handleCommunityPoints(context.Background(), pointsEarnedMessage("WATCH", 10), s)

	if !s.Stream.IsWatchStreakMissing {
		t.Fatal("plain watch points are not a streak and must not release the channel")
	}
}

func TestRememberWatchStreakIsSafeWithoutAStore(t *testing.T) {
	m, _ := newTestMiner(t)
	m.streaks = nil

	// Must not panic; persistence is optional.
	m.rememberWatchStreak("krissi", "b-1")

	if m.WatchStreakEarned("krissi", "b-1") {
		t.Error("without a store nothing can be reported as earned")
	}
}
