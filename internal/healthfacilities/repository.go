package healthfacilities

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(context.Context, SaveRequest) (*Facility, error)
	Update(context.Context, string, SaveRequest) (*Facility, error)
	GetByID(context.Context, string) (*Facility, error)
	List(context.Context, ListFilter) ([]Facility, error)
}
type PostgresRepository struct{ db *pgxpool.Pool }

func NewPostgresRepository(db *pgxpool.Pool) Repository { return &PostgresRepository{db: db} }

const fields = `id,name,cnes,unit_type,ST_Y(location::geometry),ST_X(location::geometry),address,city,state,phone,is_24_hours,treats_envenomation,has_antivenom,last_verified_at,is_active,created_at,updated_at`

func scanFacility(row pgx.Row) (*Facility, error) {
	var f Facility
	err := row.Scan(&f.ID, &f.Name, &f.CNES, &f.UnitType, &f.Latitude, &f.Longitude, &f.Address, &f.City, &f.State, &f.Phone, &f.Is24Hours, &f.TreatsEnvenomation, &f.HasAntivenom, &f.LastVerifiedAt, &f.IsActive, &f.CreatedAt, &f.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar unidade: %w", err)
	}
	return &f, nil
}
func mapError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "health_facilities_cnes_key" {
		return ErrDuplicateCNES
	}
	return err
}
func active(r SaveRequest) bool {
	if r.IsActive == nil {
		return true
	}
	return *r.IsActive
}
func (p *PostgresRepository) Create(ctx context.Context, r SaveRequest) (*Facility, error) {
	q := `INSERT INTO health_facilities (name,cnes,unit_type,location,address,city,state,phone,is_24_hours,treats_envenomation,has_antivenom,last_verified_at,is_active)
 VALUES ($1,$2,$3,ST_SetSRID(ST_MakePoint($4,$5),4326)::geography,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING ` + fields
	f, err := scanFacility(p.db.QueryRow(ctx, q, r.Name, r.CNES, r.UnitType, *r.Longitude, *r.Latitude, r.Address, r.City, r.State, r.Phone, r.Is24Hours, r.TreatsEnvenomation, r.HasAntivenom, r.LastVerifiedAt, active(r)))
	return f, mapError(err)
}
func (p *PostgresRepository) Update(ctx context.Context, id string, r SaveRequest) (*Facility, error) {
	q := `UPDATE health_facilities SET name=$2,cnes=$3,unit_type=$4,location=ST_SetSRID(ST_MakePoint($5,$6),4326)::geography,address=$7,city=$8,state=$9,phone=$10,is_24_hours=$11,treats_envenomation=$12,has_antivenom=$13,last_verified_at=$14,is_active=COALESCE($15,is_active),updated_at=NOW() WHERE id=$1 RETURNING ` + fields
	f, err := scanFacility(p.db.QueryRow(ctx, q, id, r.Name, r.CNES, r.UnitType, *r.Longitude, *r.Latitude, r.Address, r.City, r.State, r.Phone, r.Is24Hours, r.TreatsEnvenomation, r.HasAntivenom, r.LastVerifiedAt, r.IsActive))
	return f, mapError(err)
}
func (p *PostgresRepository) GetByID(ctx context.Context, id string) (*Facility, error) {
	return scanFacility(p.db.QueryRow(ctx, `SELECT `+fields+` FROM health_facilities WHERE id=$1 AND is_active=TRUE`, id))
}
func (p *PostgresRepository) List(ctx context.Context, f ListFilter) ([]Facility, error) {
	q := `SELECT ` + fields
	args := []any{f.State, f.City, f.Limit}
	if f.Latitude != nil {
		q += `,ST_Distance(location,ST_SetSRID(ST_MakePoint($4,$5),4326)::geography) AS distance_meters`
		args = append(args, *f.Longitude, *f.Latitude, f.RadiusMeters)
	}
	q += ` FROM health_facilities WHERE is_active=TRUE AND ($1='' OR state=$1) AND ($2='' OR LOWER(city)=LOWER($2))`
	if f.Latitude != nil {
		q += ` AND ST_DWithin(location,ST_SetSRID(ST_MakePoint($4,$5),4326)::geography,$6) ORDER BY distance_meters ASC, name ASC LIMIT $3`
	} else {
		q += ` ORDER BY name ASC LIMIT $3`
	}
	rows, err := p.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar unidades: %w", err)
	}
	defer rows.Close()
	result := make([]Facility, 0)
	for rows.Next() {
		var x Facility
		dest := []any{&x.ID, &x.Name, &x.CNES, &x.UnitType, &x.Latitude, &x.Longitude, &x.Address, &x.City, &x.State, &x.Phone, &x.Is24Hours, &x.TreatsEnvenomation, &x.HasAntivenom, &x.LastVerifiedAt, &x.IsActive, &x.CreatedAt, &x.UpdatedAt}
		if f.Latitude != nil {
			dest = append(dest, &x.DistanceMeters)
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("erro ao ler unidade: %w", err)
		}
		result = append(result, x)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
