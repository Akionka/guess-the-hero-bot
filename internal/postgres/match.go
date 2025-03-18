package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/akionka/akionkabot/internal/data"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewMatchRepository(db *pgxpool.Pool, logger *slog.Logger) *MatchRepository {
	return &MatchRepository{
		db:     db,
		logger: logger,
	}
}

func (r *MatchRepository) GetMatch(ctx context.Context, id data.MatchID) (*data.Match, error) {
	const sql = `
		SELECT
			m.match_id, m.winning_team, m.actual_rank, m.avg_mmr, m.started_at,
			mp.player_steam_id, mp.team, mp.position,
			h.hero_id, h.short_name AS h_short_name, h.display_name AS h_display_name,
			mpi."order", i.item_id, i.short_name AS i_short_name, i.display_name AS i_display_name

		FROM		matches				m
		LEFT JOIN	match_players		mp	ON	m.match_id = mp.match_id
		INNER JOIN	heroes				h	ON	mp.hero_id = h.hero_id
		LEFT JOIN	match_player_items	mpi	ON	mp.player_steam_id = mpi.player_steam_id
											AND	m.match_id = mpi.match_id
		INNER JOIN	items				i	ON	mpi.item_id = i.item_id

		WHERE m.match_id = $1
		ORDER BY mp.team, mp.position, mpi."order"
	`

	r.logger.DebugContext(ctx, "getting match", slog.Int64("match_id", int64(id)))
	rows, err := r.db.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}
	defer rows.Close()

	match, err := r.scanMatchRows(rows)

	return match, err
}

func (r MatchRepository) scanMatchRows(rows pgx.Rows) (*data.Match, error) {
	type (
		matchRow struct {
			MatchID     data.MatchID `db:"match_id"`
			WinningTeam data.Team    `db:"winning_team"`
			ActualRank  data.Rank    `db:"actual_rank"`
			AvgMMR      *int         `db:"avg_mmr"`
			StartedAt   time.Time    `db:"started_at"`

			PlayerSteamID *data.SteamID  `db:"player_steam_id"`
			PlayerTeam    *data.Team     `db:"team"`
			PlayerPos     *data.Position `db:"position"`

			HeroID          *data.HeroID `db:"hero_id"`
			HeroShortName   *string      `db:"h_short_name"`
			HeroDisplayName *string      `db:"h_display_name"`

			ItemID          *data.ItemID `db:"item_id"`
			ItemOrder       *int         `db:"order"`
			ItemShortName   *string      `db:"i_short_name"`
			ItemDisplayName *string      `db:"i_display_name"`
		}
		playerItemKey struct {
			SteamID data.SteamID
			Order   int
		}
	)
	var (
		match       *data.Match
		players     = make(map[data.SteamID]int, maxPlayerCount)
		playerItems = make(map[playerItemKey]int, maxPlayerCount*maxItemCount)
	)

	mRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[matchRow])
	if err != nil {
		return nil, fmt.Errorf("failed to collect match rows: %w", err)
	}

	for _, row := range mRows {
		if match == nil {
			match = &data.Match{
				ID:          row.MatchID,
				WinningTeam: row.WinningTeam,
				ActualRank:  row.ActualRank,
				AvgMMR:      row.AvgMMR,
				StartedAt:   row.StartedAt,
				Players:     make([]data.Player, 0, maxPlayerCount),
			}
		}

		if row.PlayerSteamID != nil {
			playerIdx, ok := players[*row.PlayerSteamID]
			if !ok {
				match.Players = append(match.Players, data.Player{
					SteamAccountID: *row.PlayerSteamID,
					Hero: data.Hero{
						ID:          *row.HeroID,
						DisplayName: *row.HeroDisplayName,
						ShortName:   *row.HeroShortName,
					},
					Team:     *row.PlayerTeam,
					Position: *row.PlayerPos,
					Items:    make([]data.Item, 0, maxItemCount),
				})
				playerIdx = len(match.Players) - 1
				players[*row.PlayerSteamID] = playerIdx
			}

			itemKey := playerItemKey{*row.PlayerSteamID, *row.ItemOrder}
			itemIdx, ok := playerItems[itemKey]
			if !ok && row.ItemID != nil {
				match.Players[playerIdx].Items = append(match.Players[playerIdx].Items, data.Item{
					ID:          *row.ItemID,
					DisplayName: *row.ItemDisplayName,
					ShortName:   *row.ItemShortName,
				})
				itemIdx = len(match.Players[playerIdx].Items) - 1
				playerItems[itemKey] = itemIdx
			}
			_ = itemIdx
		}
	}

	return match, nil
}

// SaveMatch implements service.MatchRepository.
func (r *MatchRepository) SaveMatch(ctx context.Context, match *data.Match) (data.MatchID, error) {
	const upsertMatch = `
		INSERT INTO matches (match_id, winning_team, actual_rank, avg_mmr, started_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (match_id)
			DO UPDATE SET
				winning_team = EXCLUDED.winning_team,
				actual_rank = EXCLUDED.actual_rank,
				avg_mmr = EXCLUDED.avg_mmr,
				started_at = EXCLUDED.started_at
	`
	const upsertPlayer = `
		INSERT INTO match_players (player_steam_id, match_id, hero_id, team, position)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (match_id, player_steam_id)
			DO UPDATE SET
				team = EXCLUDED.team,
				position = EXCLUDED.position,
				hero_id = EXCLUDED.hero_id
	`
	const upsertPlayerItem = `
		INSERT INTO match_player_items (player_steam_id, match_id, item_id, "order")
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (player_steam_id, match_id, "order")
			DO UPDATE SET
				item_id = EXCLUDED.item_id
	`

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				r.logger.ErrorContext(ctx, "failed to rollback transaction", slog.String("error", err.Error()))
			}
		} else {
			if err := tx.Commit(ctx); err != nil {
				r.logger.ErrorContext(ctx, "failed to commit transaction", slog.String("error", err.Error()))
			}
		}
	}()

	r.logger.DebugContext(ctx, "saving match", slog.Int64("match_id", int64(match.ID)))
	if _, err := tx.Exec(ctx, upsertMatch,
		match.ID, match.WinningTeam, match.ActualRank, match.AvgMMR, match.StartedAt,
	); err != nil {
		return 0, fmt.Errorf("failed to save match: %w", err)
	}
	for _, player := range match.Players {
		if _, err := tx.Exec(ctx, upsertPlayer,
			player.SteamAccountID, match.ID, player.Hero.ID, player.Team, player.Position,
		); err != nil {
			return 0, fmt.Errorf("failed to save player: %w", err)
		}
		for i, item := range player.Items {
			if _, err := tx.Exec(ctx, upsertPlayerItem,
				player.SteamAccountID, match.ID, item.ID, i,
			); err != nil {
				return 0, fmt.Errorf("failed to save player item: %w", err)
			}
		}
	}

	return match.ID, nil
}
