package pubsub

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"github.com/Guliveer/twitch-miner-go/internal/constants"
	"github.com/Guliveer/twitch-miner-go/internal/logger"
	"github.com/Guliveer/twitch-miner-go/internal/model"
)

// newTestPool returns a pool whose connections are in-memory test connections
// and whose log output can be read back through the returned function.
func newTestPool(t *testing.T) (*Pool, func() string) {
	t.Helper()

	logDir := t.TempDir()
	log, err := logger.Setup(logger.Config{
		Level:     slog.Level(100),
		FileLevel: slog.LevelDebug,
		LogDir:    logDir,
	})
	if err != nil {
		t.Fatalf("logger setup: %v", err)
	}
	t.Cleanup(func() { _ = log.Close() })

	pool := NewPool(&mockAuthProvider{token: "token"}, log, nil)
	pool.dial = func(_ context.Context, index int) (*Connection, error) {
		conn := newTestConnection(nil)
		conn.index = index
		return conn, nil
	}

	readLog := func() string {
		files, err := filepath.Glob(filepath.Join(logDir, "*.log"))
		if err != nil || len(files) != 1 {
			t.Fatalf("expected one log file in %s, got %v (err %v)", logDir, files, err)
		}
		data, err := os.ReadFile(files[0])
		if err != nil {
			t.Fatalf("reading log file: %v", err)
		}
		return string(data)
	}
	return pool, readLog
}

func testStreamer(username, channelID string) *model.Streamer {
	s := model.NewStreamer(username)
	s.ChannelID = channelID
	return s
}

// videoTopics returns one video-playback topic for each of n distinct channels.
func videoTopics(n int) []*model.PubSubTopic {
	topics := make([]*model.PubSubTopic, 0, n)
	for i := range n {
		s := testStreamer(fmt.Sprintf("channel%d", i), fmt.Sprintf("%d", 1000+i))
		topics = append(topics, model.NewStreamerTopic(model.PubSubTopicVideoPlayback, s))
	}
	return topics
}

func TestSubscribe_OverflowNamesChannelsAndTopicTypes(t *testing.T) {
	pool, readLog := newTestPool(t)
	pool.maxConns = 1

	if failed := pool.Subscribe(context.Background(), videoTopics(constants.MaxTopicsPerConn)); len(failed) != 0 {
		t.Fatalf("expected the first connection to take %d topics, %d failed", constants.MaxTopicsPerConn, len(failed))
	}

	late := testStreamer("latecomer", "9001")
	overflow := []*model.PubSubTopic{
		model.NewStreamerTopic(model.PubSubTopicRaid, late),
		model.NewStreamerTopic(model.PubSubTopicPredictions, late),
	}
	failed := pool.Subscribe(context.Background(), overflow)

	if len(failed) != len(overflow) {
		t.Fatalf("expected %d topics to be returned as not subscribed, got %d", len(overflow), len(failed))
	}
	out := readLog()
	if !strings.Contains(out, "PubSub topic limit reached") {
		t.Fatalf("expected a limit warning, log was:\n%s", out)
	}
	if !strings.Contains(out, "latecomer (raid, predictions-channel-v1)") {
		t.Errorf("expected the warning to name the channel and topic types, log was:\n%s", out)
	}
	if strings.Contains(out, "Failed to subscribe PubSub topics") {
		t.Errorf("did not expect a generic failure warning for a pure overflow, log was:\n%s", out)
	}
}

func TestSubscribe_DialFailureIsNotReportedAsLimit(t *testing.T) {
	pool, readLog := newTestPool(t)
	pool.dial = func(context.Context, int) (*Connection, error) {
		return nil, errors.New("dial tcp: lookup pubsub-edge.twitch.tv: i/o timeout")
	}

	topic := model.NewStreamerTopic(model.PubSubTopicRaid, testStreamer("somechannel", "42"))
	failed := pool.Subscribe(context.Background(), []*model.PubSubTopic{topic})

	if len(failed) != 1 {
		t.Fatalf("expected the topic to be returned as not subscribed, got %d", len(failed))
	}
	out := readLog()
	if strings.Contains(out, "limit reached") {
		t.Errorf("a dial failure must not be reported as the topic limit, log was:\n%s", out)
	}
	if !strings.Contains(out, "Failed to subscribe PubSub topics") || !strings.Contains(out, "i/o timeout") {
		t.Errorf("expected a failure warning with the dial error, log was:\n%s", out)
	}
	if !strings.Contains(out, "somechannel (raid)") {
		t.Errorf("expected the failure warning to name the channel, log was:\n%s", out)
	}
}

func TestSubscribe_SkipsTopicAlreadyOnAnotherConnection(t *testing.T) {
	pool, _ := newTestPool(t)

	topics := videoTopics(constants.MaxTopicsPerConn)
	if failed := pool.Subscribe(context.Background(), topics); len(failed) != 0 {
		t.Fatalf("unexpected failures: %d", len(failed))
	}
	extra := model.NewStreamerTopic(model.PubSubTopicRaid, testStreamer("other", "7"))
	if failed := pool.Subscribe(context.Background(), []*model.PubSubTopic{extra}); len(failed) != 0 {
		t.Fatalf("unexpected failures: %d", len(failed))
	}

	// The first connection is full, so without a pool-wide check the repeated
	// topic would land on the second connection a second time.
	if failed := pool.Subscribe(context.Background(), topics[:1]); len(failed) != 0 {
		t.Fatalf("unexpected failures: %d", len(failed))
	}

	if got, want := pool.TotalTopicCount(), constants.MaxTopicsPerConn+1; got != want {
		t.Fatalf("expected %d topics after re-subscribing an existing one, got %d", want, got)
	}
}

// listenRecorder is a PubSub test server that records the topics of every
// LISTEN request it receives.
func listenRecorder(t *testing.T) (*httptest.Server, <-chan string) {
	t.Helper()

	listens := make(chan string, 16)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer ws.CloseNow()
		for {
			var req Request
			if err := wsjson.Read(r.Context(), ws, &req); err != nil {
				return
			}
			if req.Type == TypeListen && req.Data != nil {
				for _, topic := range req.Data.Topics {
					listens <- topic
				}
			}
		}
	}))
	t.Cleanup(srv.Close)
	return srv, listens
}

func TestSubscribe_StartsConnectionCreatedWhileRunning(t *testing.T) {
	srv, listens := listenRecorder(t)
	pool, _ := newTestPool(t)
	pool.dial = func(ctx context.Context, index int) (*Connection, error) {
		ws, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(srv.URL, "http"), nil) //nolint:bodyclose // coder/websocket owns the handshake response body
		if err != nil {
			return nil, err
		}
		conn := newTestConnection(nil)
		conn.index = index
		conn.conn = ws
		return conn, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Go(func() { _ = pool.Run(ctx) })
	t.Cleanup(func() {
		cancel()
		wg.Wait()
		pool.Close()
	})

	deadline := time.Now().Add(5 * time.Second)
	for {
		pool.mu.Lock()
		running := pool.runGroup != nil
		pool.mu.Unlock()
		if running {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("pool did not start running")
		}
		time.Sleep(5 * time.Millisecond)
	}

	topic := model.NewStreamerTopic(model.PubSubTopicRaid, testStreamer("livechannel", "55"))
	if failed := pool.Subscribe(ctx, []*model.PubSubTopic{topic}); len(failed) != 0 {
		t.Fatalf("unexpected failures: %d", len(failed))
	}

	select {
	case got := <-listens:
		if got != topic.String() {
			t.Fatalf("expected LISTEN for %q, got %q", topic.String(), got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("a connection created after Run never sent its LISTEN")
	}
}
