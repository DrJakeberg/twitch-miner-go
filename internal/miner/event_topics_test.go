package miner

import (
	"context"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/Guliveer/twitch-miner-go/internal/chat"
	"github.com/Guliveer/twitch-miner-go/internal/model"
)

// topicRecorder stands in for the PubSub pool and records which topics are
// currently subscribed.
type topicRecorder struct {
	mu               sync.Mutex
	subscribed       map[string]bool
	reject           map[string]bool
	subscribeCalls   int
	unsubscribeCalls int
}

func recordTopics(m *Miner) *topicRecorder {
	r := &topicRecorder{subscribed: make(map[string]bool), reject: make(map[string]bool)}
	m.subscribeTopics = func(_ context.Context, topics []*model.PubSubTopic) []*model.PubSubTopic {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.subscribeCalls++
		var failed []*model.PubSubTopic
		for _, topic := range topics {
			if r.reject[topic.String()] {
				failed = append(failed, topic)
				continue
			}
			r.subscribed[topic.String()] = true
		}
		return failed
	}
	m.unsubscribeTopics = func(topics []*model.PubSubTopic) error {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.unsubscribeCalls++
		for _, topic := range topics {
			delete(r.subscribed, topic.String())
		}
		return nil
	}
	return r
}

func (r *topicRecorder) topics() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, 0, len(r.subscribed))
	for topic := range r.subscribed {
		out = append(out, topic)
	}
	slices.Sort(out)
	return out
}

func (r *topicRecorder) calls() (subscribe, unsubscribe int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.subscribeCalls, r.unsubscribeCalls
}

// allEventTopics is what a channel with every event feature enabled holds
// while it is live, sorted.
var allEventTopics = []string{
	"community-moments-channel-v1.100",
	"community-points-channel-v1.100",
	"predictions-channel-v1.100",
	"raid.100",
}

// trackedChannel adds a streamer with every event feature enabled to m.
func trackedChannel(m *Miner) *model.Streamer {
	s := model.NewStreamer("livechannel")
	s.ChannelID = "100"
	s.Settings = &model.StreamerSettings{
		FollowRaid:            true,
		MakePredictions:       true,
		ClaimMoments:          true,
		CommunityGoalsEnabled: true,
	}
	m.streamers = append(m.streamers, s)
	return s
}

func expectTopics(t *testing.T, r *topicRecorder, want []string, when string) {
	t.Helper()
	if got := r.topics(); !slices.Equal(got, want) {
		t.Fatalf("%s: expected subscribed topics %v, got %v", when, want, got)
	}
}

func TestEventTopics_HeldWhileLiveAndForTheGracePeriod(t *testing.T) {
	m, _ := newTestMiner(t)
	r := recordTopics(m)
	s := trackedChannel(m)
	ctx := context.Background()
	start := time.Date(2026, 11, 25, 20, 0, 0, 0, time.UTC)

	m.reconcileEventTopics(ctx, s, true, start)
	expectTopics(t, r, allEventTopics, "after going live")

	m.reconcileEventTopics(ctx, s, true, start.Add(30*time.Second))
	if subscribe, _ := r.calls(); subscribe != 1 {
		t.Fatalf("expected a repeated live check to subscribe nothing, got %d subscribe calls", subscribe)
	}

	offline := start.Add(2 * time.Hour)
	m.reconcileEventTopics(ctx, s, false, offline)
	expectTopics(t, r, allEventTopics, "right after going offline")

	m.reconcileEventTopics(ctx, s, false, offline.Add(eventTopicGrace-time.Second))
	expectTopics(t, r, allEventTopics, "just before the grace period ends")

	m.reconcileEventTopics(ctx, s, false, offline.Add(eventTopicGrace))
	expectTopics(t, r, []string{}, "once the grace period has ended")

	m.eventTopicsMu.Lock()
	remaining := len(m.eventTopics)
	m.eventTopicsMu.Unlock()
	if remaining != 0 {
		t.Fatalf("expected no reconcile state once the topics are dropped, got %d entries", remaining)
	}
}

func TestEventTopics_ReturningWithinGraceRestartsIt(t *testing.T) {
	m, _ := newTestMiner(t)
	r := recordTopics(m)
	s := trackedChannel(m)
	ctx := context.Background()
	start := time.Date(2026, 11, 25, 20, 0, 0, 0, time.UTC)

	m.reconcileEventTopics(ctx, s, true, start)
	m.reconcileEventTopics(ctx, s, false, start.Add(time.Hour))
	m.reconcileEventTopics(ctx, s, true, start.Add(time.Hour+5*time.Minute))

	secondOffline := start.Add(time.Hour + 6*time.Minute)
	m.reconcileEventTopics(ctx, s, false, secondOffline)
	m.reconcileEventTopics(ctx, s, false, secondOffline.Add(eventTopicGrace-time.Second))
	expectTopics(t, r, allEventTopics, "within the grace period of the second offline")

	if subscribe, unsubscribe := r.calls(); subscribe != 1 || unsubscribe != 0 {
		t.Fatalf("expected a short break to cause no resubscription, got %d subscribe and %d unsubscribe calls",
			subscribe, unsubscribe)
	}

	m.reconcileEventTopics(ctx, s, false, secondOffline.Add(eventTopicGrace))
	expectTopics(t, r, []string{}, "after the grace period of the second offline")
}

func TestEventTopics_NeverSubscribedForAnOfflineChannel(t *testing.T) {
	m, _ := newTestMiner(t)
	r := recordTopics(m)
	s := trackedChannel(m)

	m.reconcileEventTopics(context.Background(), s, false, time.Now())

	if subscribe, unsubscribe := r.calls(); subscribe != 0 || unsubscribe != 0 {
		t.Fatalf("expected no PubSub changes for an offline channel, got %d subscribe and %d unsubscribe calls",
			subscribe, unsubscribe)
	}
}

func TestEventTopics_RejectedTopicsAreRetried(t *testing.T) {
	m, _ := newTestMiner(t)
	r := recordTopics(m)
	s := trackedChannel(m)
	ctx := context.Background()
	now := time.Now()

	r.reject["raid.100"] = true
	m.reconcileEventTopics(ctx, s, true, now)
	expectTopics(t, r, []string{
		"community-moments-channel-v1.100",
		"community-points-channel-v1.100",
		"predictions-channel-v1.100",
	}, "while the pool rejects the raid topic")

	r.mu.Lock()
	delete(r.reject, "raid.100")
	r.mu.Unlock()
	m.reconcileEventTopics(ctx, s, true, now.Add(time.Minute))
	expectTopics(t, r, allEventTopics, "on the next live check")
}

func TestEventTopics_FollowSettingsWhileLive(t *testing.T) {
	m, _ := newTestMiner(t)
	r := recordTopics(m)
	s := trackedChannel(m)
	ctx := context.Background()
	now := time.Now()

	m.reconcileEventTopics(ctx, s, true, now)
	s.Mu.Lock()
	s.Settings.ClaimMoments = false
	s.Mu.Unlock()
	m.reconcileEventTopics(ctx, s, true, now.Add(time.Minute))

	expectTopics(t, r, []string{
		"community-points-channel-v1.100",
		"predictions-channel-v1.100",
		"raid.100",
	}, "after disabling moments")
}

func TestEventTopics_RemovedChannelIsNotResubscribed(t *testing.T) {
	m, _ := newTestMiner(t)
	r := recordTopics(m)
	s := trackedChannel(m)
	ctx := context.Background()

	m.reconcileEventTopics(ctx, s, true, time.Now())
	m.streamers = nil
	m.forgetEventTopics(s.ChannelID)
	r.mu.Lock()
	clear(r.subscribed)
	r.mu.Unlock()

	// An online check that was already running when the channel was removed.
	m.reconcileEventTopics(ctx, s, true, time.Now())

	expectTopics(t, r, []string{}, "after the channel was removed")
}

func TestEventTopics_FollowStreamUpAndDown(t *testing.T) {
	m, _ := newTestMiner(t)
	m.chat = chat.NewManager("testaccount", "token", m.log)
	r := recordTopics(m)
	s := trackedChannel(m)
	ctx := context.Background()

	m.handleStreamUp(ctx, s)
	expectTopics(t, r, allEventTopics, "after stream-up")

	m.handleStreamDown(ctx, s)
	expectTopics(t, r, allEventTopics, "right after stream-down")

	m.eventTopicsMu.Lock()
	offlineSince := m.eventTopics[s.ChannelID].offlineSince
	m.eventTopicsMu.Unlock()
	if offlineSince.IsZero() {
		t.Fatal("expected stream-down to start the grace period")
	}
}

func TestCheckAllStreamersOnline_SubscribesEventTopicsOfLiveChannels(t *testing.T) {
	m, mock := newTestMiner(t)
	m.chat = chat.NewManager("testaccount", "token", m.log)
	r := recordTopics(m)
	trackedChannel(m)
	offline := model.NewStreamer("offlinechannel")
	offline.ChannelID = "200"
	offline.Settings = &model.StreamerSettings{FollowRaid: true, MakePredictions: true}
	m.streamers = append(m.streamers, offline)
	mock.online = map[string]bool{"livechannel": true}

	m.checkAllStreamersOnline(context.Background())

	expectTopics(t, r, allEventTopics, "after the initial online check")
}
