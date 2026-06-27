package configeditor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
)

// RunTUI launches the interactive terminal config editor.
func RunTUI(configDir string) error {
	accounts := listAccountNames(configDir)

	if len(accounts) == 0 {
		fmt.Println("No config files found in " + configDir)
		fmt.Println("Create a new account:")
		return runCreateAccount(configDir)
	}

	accounts = append(accounts, "[+ Create new account]")

	var choice string
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Config Editor — Select account").
				Description(configDir).
				Options(huh.NewOptions(accounts...)...).
				Value(&choice),
		),
	).Run(); err != nil {
		return err
	}

	if choice == "[+ Create new account]" {
		return runCreateAccount(configDir)
	}
	return runEditAccount(configDir, choice)
}

func runCreateAccount(configDir string) error {
	var name string
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("New account name").
				Description("Enter the Twitch username (letters, numbers, _ and -)").
				Validate(func(s string) error {
					if !validName.MatchString(s) {
						return fmt.Errorf("use only letters, numbers, _ and -")
					}
					if _, err := os.Stat(filepath.Join(configDir, s+".yaml")); err == nil {
						return fmt.Errorf("account %q already exists", s)
					}
					return nil
				}).
				Value(&name),
		),
	).Run(); err != nil {
		return err
	}

	defaultCfg := map[string]any{
		"streamers": []any{map[string]any{"username": "placeholder"}},
		"features":  map[string]any{"enable_analytics": true},
	}
	s := &Server{configDir: configDir}
	if err := s.saveRaw(name, defaultCfg); err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	fmt.Printf("Account %q created.\n", name)
	return runEditAccount(configDir, name)
}

type editAccountFields struct {
	enabled            bool
	maxWatchStr        string
	proxy              string
	selectedPriorities []string
	claimDropsStartup  bool
	enableAnalytics    bool
	cwEnabled          bool
	cwInterval         string
	cwCategoriesStr    string
	twEnabled          bool
	twInterval         string
	twTeamsStr         string
	streamersStr       string
	action             string
}

func loadEditFields(cfg map[string]any) editAccountFields {
	features, _ := cfg["features"].(map[string]any)
	if features == nil {
		features = map[string]any{}
	}
	cw := subMap(cfg, "category_watcher")
	tw := subMap(cfg, "team_watcher")

	return editAccountFields{
		enabled:            boolVal(cfg, "enabled", true),
		maxWatchStr:        strconv.Itoa(intVal(cfg, "max_watch_streams", 2)),
		proxy:              stringVal(cfg, "proxy", ""),
		selectedPriorities: stringSliceVal(cfg, "priority", []string{"STREAK", "DROPS", "ORDER"}),
		claimDropsStartup:  boolVal(features, "claim_drops_startup", false),
		enableAnalytics:    boolVal(features, "enable_analytics", false),
		cwEnabled:          boolVal(cw, "enabled", false),
		cwInterval:         stringVal(cw, "poll_interval", "120s"),
		cwCategoriesStr:    strings.Join(parseCategorySlugs(cw), ", "),
		twEnabled:          boolVal(tw, "enabled", false),
		twInterval:         stringVal(tw, "poll_interval", "120s"),
		twTeamsStr:         strings.Join(parseTeamNames(tw), ", "),
		streamersStr:       strings.Join(parseStreamerUsernames(cfg), ", "),
	}
}

func subMap(cfg map[string]any, key string) map[string]any {
	m, _ := cfg[key].(map[string]any)
	if m == nil {
		return map[string]any{}
	}
	return m
}

func parseCategorySlugs(cw map[string]any) []string {
	raw, _ := cw["categories"].([]any)
	out := make([]string, 0, len(raw))
	for _, c := range raw {
		if cm, ok := c.(map[string]any); ok {
			if slug, ok := cm["slug"].(string); ok && slug != "" {
				out = append(out, slug)
			}
		}
	}
	return out
}

func parseTeamNames(tw map[string]any) []string {
	raw, _ := tw["teams"].([]any)
	out := make([]string, 0, len(raw))
	for _, t := range raw {
		if tm, ok := t.(map[string]any); ok {
			if n, ok := tm["name"].(string); ok && n != "" {
				out = append(out, n)
			}
		}
	}
	return out
}

func parseStreamerUsernames(cfg map[string]any) []string {
	raw, _ := cfg["streamers"].([]any)
	out := make([]string, 0, len(raw))
	for _, st := range raw {
		if sm, ok := st.(map[string]any); ok {
			if u, ok := sm["username"].(string); ok && u != "" {
				out = append(out, u)
			}
		}
	}
	return out
}

func applyEditFields(cfg map[string]any, f *editAccountFields) {
	maxWatchInt, _ := strconv.Atoi(f.maxWatchStr)
	if maxWatchInt < 1 {
		maxWatchInt = 2
	}
	if !f.enabled {
		cfg["enabled"] = false
	} else {
		delete(cfg, "enabled")
	}
	cfg["max_watch_streams"] = maxWatchInt
	if f.proxy != "" {
		cfg["proxy"] = f.proxy
	} else {
		delete(cfg, "proxy")
	}
	if len(f.selectedPriorities) > 0 {
		cfg["priority"] = f.selectedPriorities
	}
	applyFeaturesSection(cfg, f)
	applyCategoryWatcherSection(cfg, f)
	applyTeamWatcherSection(cfg, f)
	applyStreamersSection(cfg, f)
}

func applyFeaturesSection(cfg map[string]any, f *editAccountFields) {
	m := map[string]any{}
	if f.claimDropsStartup {
		m["claim_drops_startup"] = true
	}
	if f.enableAnalytics {
		m["enable_analytics"] = true
	}
	if len(m) > 0 {
		cfg["features"] = m
	} else {
		delete(cfg, "features")
	}
}

func applyCategoryWatcherSection(cfg map[string]any, f *editAccountFields) {
	m := map[string]any{}
	if f.cwEnabled {
		m["enabled"] = true
	}
	if f.cwInterval != "" && f.cwInterval != "120s" {
		m["poll_interval"] = f.cwInterval
	}
	if slugs := parseCSSV(f.cwCategoriesStr); len(slugs) > 0 {
		cats := make([]any, len(slugs))
		for i, slug := range slugs {
			cats[i] = map[string]any{"slug": slug}
		}
		m["categories"] = cats
	}
	if len(m) > 0 {
		cfg["category_watcher"] = m
	} else {
		delete(cfg, "category_watcher")
	}
}

func applyTeamWatcherSection(cfg map[string]any, f *editAccountFields) {
	m := map[string]any{}
	if f.twEnabled {
		m["enabled"] = true
	}
	if f.twInterval != "" && f.twInterval != "120s" {
		m["poll_interval"] = f.twInterval
	}
	if names := parseCSSV(f.twTeamsStr); len(names) > 0 {
		teams := make([]any, len(names))
		for i, n := range names {
			teams[i] = map[string]any{"name": n}
		}
		m["teams"] = teams
	}
	if len(m) > 0 {
		cfg["team_watcher"] = m
	} else {
		delete(cfg, "team_watcher")
	}
}

func applyStreamersSection(cfg map[string]any, f *editAccountFields) {
	list := parseCSSV(f.streamersStr)
	if len(list) > 0 {
		streamers := make([]any, len(list))
		for i, u := range list {
			streamers[i] = map[string]any{"username": u}
		}
		cfg["streamers"] = streamers
	} else {
		delete(cfg, "streamers")
	}
}

func runEditAccount(configDir string, name string) error {
	s := &Server{configDir: configDir}
	cfg, err := s.loadRaw(name)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	f := loadEditFields(cfg)

	priorityOptions := []string{"STREAK", "DROPS", "ORDER", "SUBSCRIBED", "POINTS_ASCENDING", "POINTS_DESCENDING"}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Editing: "+name).
				Description("Use Tab/Enter to navigate, Esc to cancel."),
			huh.NewConfirm().
				Title("Enabled").
				Description("Enable this account").
				Value(&f.enabled),
		).Title("General"),

		huh.NewGroup(
			huh.NewInput().
				Title("Max watch streams").
				Validate(func(s string) error {
					n, err := strconv.Atoi(s)
					if err != nil || n < 1 {
						return fmt.Errorf("must be a positive integer")
					}
					return nil
				}).
				Value(&f.maxWatchStr),
			huh.NewInput().
				Title("Proxy").
				Description("e.g. socks5://127.0.0.1:1080 (leave blank to disable)").
				Value(&f.proxy),
			huh.NewMultiSelect[string]().
				Title("Priority").
				Description("Select and order priorities (first wins)").
				Options(huh.NewOptions(priorityOptions...)...).
				Value(&f.selectedPriorities),
		).Title("General (continued)"),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Claim drops on startup").
				Value(&f.claimDropsStartup),
			huh.NewConfirm().
				Title("Enable analytics").
				Value(&f.enableAnalytics),
		).Title("Features"),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Category Watcher enabled").
				Value(&f.cwEnabled),
			huh.NewInput().
				Title("Category Watcher poll interval").
				Description("e.g. 120s, 5m, 1h").
				Validate(func(s string) error {
					if s != "" && !isValidDuration(s) {
						return fmt.Errorf("invalid duration (e.g. 120s, 5m, 1h30m)")
					}
					return nil
				}).
				Value(&f.cwInterval),
			huh.NewInput().
				Title("Categories (comma-separated slugs)").
				Description("e.g. just-chatting, science-and-technology").
				Value(&f.cwCategoriesStr),
		).Title("Category Watcher"),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Team Watcher enabled").
				Value(&f.twEnabled),
			huh.NewInput().
				Title("Team Watcher poll interval").
				Description("e.g. 120s, 5m").
				Validate(func(s string) error {
					if s != "" && !isValidDuration(s) {
						return fmt.Errorf("invalid duration (e.g. 120s, 5m, 1h30m)")
					}
					return nil
				}).
				Value(&f.twInterval),
			huh.NewInput().
				Title("Teams (comma-separated names)").
				Description("e.g. nrg, sentinels").
				Value(&f.twTeamsStr),
		).Title("Team Watcher"),

		huh.NewGroup(
			huh.NewInput().
				Title("Streamers (comma-separated usernames)").
				Description("e.g. streamer1, streamer2").
				Value(&f.streamersStr),
		).Title("Streamers"),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Save changes?").
				Options(
					huh.NewOption("Save and exit", "save"),
					huh.NewOption("Discard changes", "discard"),
					huh.NewOption("Delete this account", "delete"),
				).
				Value(&f.action),
		).Title("Confirm"),
	)

	if err := form.Run(); err != nil {
		return err
	}

	switch f.action {
	case "discard":
		fmt.Println("Changes discarded.")
		return nil
	case "delete":
		var confirm bool
		if err := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Delete account %q?", name)).
				Description("This will permanently remove the config file.").
				Value(&confirm),
		)).Run(); err != nil {
			return err
		}
		if confirm {
			if err := os.Remove(filepath.Join(configDir, name+".yaml")); err != nil {
				return fmt.Errorf("failed to delete: %w", err)
			}
			fmt.Printf("Account %q deleted.\n", name)
		}
		return nil
	}

	applyEditFields(cfg, &f)

	if errs := validateConfig(cfg); len(errs) > 0 {
		fmt.Println("Validation errors:")
		for _, e := range errs {
			fmt.Println("  - " + e)
		}
		return fmt.Errorf("config not saved due to validation errors")
	}

	if err := s.saveRaw(name, cleanConfig(cfg)); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}
	fmt.Printf("Account %q saved.\n", name)
	return nil
}

func listAccountNames(configDir string) []string {
	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".example") {
			continue
		}
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			names = append(names, strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml"))
		}
	}
	return names
}

func parseCSSV(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func boolVal(m map[string]any, key string, def bool) bool {
	v, ok := m[key]
	if !ok {
		return def
	}
	b, ok := v.(bool)
	if !ok {
		return def
	}
	return b
}

func intVal(m map[string]any, key string, def int) int {
	v, ok := m[key]
	if !ok {
		return def
	}
	switch n := v.(type) {
	case int:
		return n
	case float64:
		return int(n)
	}
	return def
}

func stringVal(m map[string]any, key string, def string) string {
	v, ok := m[key]
	if !ok {
		return def
	}
	s, ok := v.(string)
	if !ok {
		return def
	}
	return s
}

func stringSliceVal(m map[string]any, key string, def []string) []string {
	v, ok := m[key]
	if !ok {
		return def
	}
	raw, ok := v.([]any)
	if !ok {
		return def
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
