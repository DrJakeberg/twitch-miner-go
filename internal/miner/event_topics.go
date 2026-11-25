package miner

import (
	"context"
	"time"

	"github.com/Guliveer/twitch-miner-go/internal/model"
)

// eventTopicGrace is how long a channel keeps its event topics after it went
// offline: raids are usually started at the very end of a stream, and
// prediction results can arrive after the stream has ended.
const eventTopicGrace = 10 * time.Minute

// eventTopicState is the event-topic subscription of one channel.
type eventTopicState struct {
	topics []*model.PubSubTopic
	// offlineSince is when the channel was first seen offline while it still
	// held event topics; zero while it is live.
	offlineSince time.Time
}

// streamerEventTopics returns the event topics (raid, predictions, moments,
// community goals) enabled for a streamer. Twitch caps an account's PubSub
// topics, so these are only held while the channel is live; the video
// playback topic, which reports stream up and down, is held at all times.
func (m *Miner) streamerEventTopics(s *model.Streamer) []*model.PubSubTopic {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	if s.Settings == nil {
		return nil
	}

	var topics []*model.PubSubTopic
	if s.Settings.FollowRaid {
		topics = append(topics, model.NewStreamerTopic(model.PubSubTopicRaid, s))
	}
	if s.Settings.MakePredictions {
		topics = append(topics, model.NewStreamerTopic(model.PubSubTopicPredictions, s))
	}
	if s.Settings.ClaimMoments {
		topics = append(topics, model.NewStreamerTopic(model.PubSubTopicCommunityMoments, s))
	}
	if s.Settings.CommunityGoalsEnabled {
		topics = append(topics, model.NewStreamerTopic(model.PubSubTopicCommunityGoals, s))
	}
	return topics
}

// syncEventTopicsToOnlineState reconciles a channel's event topics with its
// online state as last determined by an online check.
func (m *Miner) syncEventTopicsToOnlineState(ctx context.Context, s *model.Streamer) {
	s.Mu.RLock()
	live := s.IsOnline
	s.Mu.RUnlock()

	m.syncEventTopics(ctx, s, live)
}

// syncEventTopics reconciles a channel's event topics with whether it is live.
// It is idempotent and meant to run after every online check.
func (m *Miner) syncEventTopics(ctx context.Context, s *model.Streamer, live bool) {
	m.reconcileEventTopics(ctx, s, live, time.Now())
}

// reconcileEventTopics subscribes the enabled event topics while the channel
// is live, keeps them for eventTopicGrace after it went offline, and drops
// them afterwards. Topics the pool could not subscribe are retried on the
// next call.
func (m *Miner) reconcileEventTopics(ctx context.Context, s *model.Streamer, live bool, now time.Time) {
	s.Mu.RLock()
	channelID := s.ChannelID
	username := s.Username
	s.Mu.RUnlock()
	if channelID == "" {
		return
	}

	m.eventTopicsMu.Lock()
	defer m.eventTopicsMu.Unlock()

	state := m.eventTopics[channelID]
	if state == nil && !live {
		return
	}
	// A channel removed concurrently must not get its topics back.
	if m.getStreamerByChannelID(channelID) == nil {
		return
	}
	if state == nil {
		state = &eventTopicState{}
		m.eventTopics[channelID] = state
	}

	wanted := state.topics
	switch {
	case live:
		state.offlineSince = time.Time{}
		wanted = m.streamerEventTopics(s)
	case state.offlineSince.IsZero():
		state.offlineSince = now
	case now.Sub(state.offlineSince) >= eventTopicGrace:
		wanted = nil
	}

	toAdd := withoutTopics(wanted, state.topics)
	toRemove := withoutTopics(state.topics, wanted)
	if len(toRemove) > 0 {
		if err := m.unsubscribeTopics(toRemove); err != nil {
			m.log.Warn("Failed to unsubscribe event topics",
				"streamer", username, "error", err)
		}
		state.topics = withoutTopics(state.topics, toRemove)
		m.log.Debug("Unsubscribed event topics", "streamer", username, "count", len(toRemove))
	}
	if len(toAdd) > 0 {
		failed := m.subscribeTopics(ctx, toAdd)
		subscribed := withoutTopics(toAdd, failed)
		state.topics = append(state.topics, subscribed...)
		m.log.Debug("Subscribed event topics",
			"streamer", username, "count", len(subscribed), "failed", len(failed))
	}

	if len(state.topics) == 0 {
		delete(m.eventTopics, channelID)
	}
}

// forgetEventTopics drops the event-topic state of a removed channel. Its
// topics are unsubscribed together with the channel's other topics.
func (m *Miner) forgetEventTopics(channelID string) {
	m.eventTopicsMu.Lock()
	delete(m.eventTopics, channelID)
	m.eventTopicsMu.Unlock()
}

// withoutTopics returns the topics that are not in excluded, compared by their
// topic string.
func withoutTopics(topics, excluded []*model.PubSubTopic) []*model.PubSubTopic {
	excludedSet := make(map[string]bool, len(excluded))
	for _, topic := range excluded {
		excludedSet[topic.String()] = true
	}

	var out []*model.PubSubTopic
	for _, topic := range topics {
		if !excludedSet[topic.String()] {
			out = append(out, topic)
		}
	}
	return out
}
