# Prediction Strategies

The miner supports 8 betting strategies for channel point predictions. Set the strategy globally under `streamer_defaults.bet.strategy` or override it per streamer.

## Strategy overview

| Strategy | Votes on | Best for |
|----------|----------|----------|
| `SMART` | High-odds outcome if gap ≥ `percentage_gap`, otherwise majority | General use — the recommended default |
| `MOST_VOTED` | Whichever outcome has the most votes at bet time | Following the crowd |
| `HIGH_ODDS` | Outcome with the best odds (fewest points bet) | High risk, high reward |
| `PERCENTAGE` | Outcome with the highest vote percentage | Similar to MOST_VOTED but percentage-weighted |
| `SMART_MONEY` | Outcome where top predictors have placed bets | Following experienced bettors |
| `NUMBER_1` | Always outcome #1 (first in the list) | Blind fixed-outcome bets |
| `NUMBER_2` – `NUMBER_8` | Always outcome #N | Blind fixed-outcome bets |

---

## `SMART` (recommended)

A hybrid strategy that adapts based on how lopsided the odds are:

1. Calculate the odds ratio between the two outcomes.
2. If the gap is ≥ `percentage_gap` (default 20%), vote on the **high-odds** outcome (higher payout, fewer bettors).
3. Otherwise, vote with the **majority** (lower risk).

```yaml
bet:
  strategy: "SMART"
  percentage: 5        # Bet 5% of current points
  percentage_gap: 20   # Switch to high-odds when gap >= 20%
  max_points: 50000
```

**When to use:** Most situations. Balances risk and reward automatically.

---

## `MOST_VOTED`

Votes on whichever outcome has accumulated the most points or votes at the time of betting. Follows the crowd.

```yaml
bet:
  strategy: "MOST_VOTED"
  percentage: 5
```

**When to use:** When you trust collective wisdom and want to minimize losses.

---

## `HIGH_ODDS`

Votes on the outcome with the fewest points bet (highest potential payout multiplier). High variance — can yield large wins but also large losses.

```yaml
bet:
  strategy: "HIGH_ODDS"
  percentage: 3        # Keep percentage low — losses will happen
  max_points: 10000
```

**When to use:** When accumulating points slowly and willing to accept variance for bigger swings.

---

## `PERCENTAGE`

Votes on the outcome with the highest vote percentage among all participants. Similar to `MOST_VOTED` but uses the percentage share rather than absolute counts.

```yaml
bet:
  strategy: "PERCENTAGE"
  percentage: 5
```

---

## `SMART_MONEY`

Tracks how top predictors (users with consistent wins) have voted and follows their lead. Requires Twitch to expose predictor history data in the PubSub payload.

```yaml
bet:
  strategy: "SMART_MONEY"
  percentage: 5
  max_points: 30000
```

**When to use:** On channels where a small group of experienced bettors consistently predicts correctly.

---

## `NUMBER_1` through `NUMBER_8`

Always votes on outcome number N (1-indexed, in the order Twitch presents them). Ignores odds, votes, and any other signal.

```yaml
bet:
  strategy: "NUMBER_1"   # Always bet on the first outcome
```

**When to use:** When you know a specific outcome always wins for a particular streamer (e.g. the streamer always wins their own game).

---

## Betting amount calculation

The bet amount is calculated as:

```
amount = min(current_points × percentage / 100, max_points)
```

If `amount < minimum_points`, no bet is placed and a `BET_FILTERS` notification is emitted.

---

## Delay modes

The `delay` and `delay_mode` fields control when the bet is actually placed after the prediction opens.

| `delay_mode` | Behavior |
|-------------|----------|
| `FROM_START` | Wait `delay` seconds after the prediction opens |
| `FROM_END` | Place the bet `delay` seconds before the prediction closes |
| `PERCENTAGE` | Wait until `delay`% of the prediction window has elapsed |

```yaml
bet:
  delay: 6
  delay_mode: "FROM_END"   # Bet 6 seconds before it closes
```

`FROM_END` is the default and generally best — it allows odds to stabilize before committing.

---

## Stealth mode

When `stealth_mode: true`, the miner waits for other bettors to place their bets first before placing its own. This reduces the miner's influence on the displayed odds and makes the betting less detectable.

```yaml
bet:
  stealth_mode: true
  delay: 10
  delay_mode: "FROM_END"
```

---

## Filter conditions

Skip betting entirely if a condition is not met. Emits a `BET_FILTERS` notification when skipped.

```yaml
bet:
  filter_condition:
    by: "total_users"    # total_users | total_points
    where: "GTE"         # GT | GTE | LT | LTE | EQ
    value: 100           # Only bet when 100+ users have voted
```

Useful for skipping low-engagement predictions that are more likely to be scripted or manipulated.
