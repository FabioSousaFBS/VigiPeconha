package healthfacilities

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"math"
	"testing"
)

type repositoryStub struct {
	createFn func(context.Context, SaveRequest) (*Facility, error)
	updateFn func(context.Context, string, SaveRequest) (*Facility, error)
	getFn    func(context.Context, string) (*Facility, error)
	listFn   func(context.Context, ListFilter) ([]Facility, error)
	calls    int
}

func (m *repositoryStub) Create(c context.Context, r SaveRequest) (*Facility, error) {
	m.calls++
	if m.createFn != nil {
		return m.createFn(c, r)
	}
	return &Facility{Name: r.Name}, nil
}
func (m *repositoryStub) Update(c context.Context, id string, r SaveRequest) (*Facility, error) {
	m.calls++
	if m.updateFn != nil {
		return m.updateFn(c, id, r)
	}
	return &Facility{ID: id, Name: r.Name}, nil
}
func (m *repositoryStub) GetByID(c context.Context, id string) (*Facility, error) {
	m.calls++
	if m.getFn != nil {
		return m.getFn(c, id)
	}
	return &Facility{ID: id}, nil
}
func (m *repositoryStub) List(c context.Context, f ListFilter) ([]Facility, error) {
	m.calls++
	if m.listFn != nil {
		return m.listFn(c, f)
	}
	return []Facility{}, nil
}
func number(v float64) *float64  { return &v }
func stringPtr(v string) *string { return &v }
func validSave() SaveRequest {
	return SaveRequest{Name: "  Unidade Teste  ", UnitType: " HOSPITAL ", Latitude: number(-22.2105), Longitude: number(-49.6564), City: " Garça ", State: " sp "}
}

func TestCreateValidation(t *testing.T) {
	cases := []struct {
		name   string
		change func(*SaveRequest)
	}{
		{"sem nome", func(r *SaveRequest) { r.Name = " " }},
		{"sem cidade", func(r *SaveRequest) { r.City = " " }},
		{"estado invalido", func(r *SaveRequest) { r.State = "S" }},
		{"tipo invalido", func(r *SaveRequest) { r.UnitType = "clinic" }},
		{"latitude ausente", func(r *SaveRequest) { r.Latitude = nil }},
		{"longitude ausente", func(r *SaveRequest) { r.Longitude = nil }},
		{"latitude fora do intervalo", func(r *SaveRequest) { r.Latitude = number(91) }},
		{"longitude fora do intervalo", func(r *SaveRequest) { r.Longitude = number(-181) }},
		{"latitude NaN", func(r *SaveRequest) { r.Latitude = number(math.NaN()) }},
		{"longitude infinita", func(r *SaveRequest) { r.Longitude = number(math.Inf(1)) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := validSave()
			tc.change(&r)
			repo := &repositoryStub{}
			_, err := NewService(repo).Create(context.Background(), r)
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("erro = %v, esperado ErrInvalidInput", err)
			}
			if repo.calls != 0 {
				t.Fatal("repository nao deveria ser chamado")
			}
		})
	}
}
func TestCreateNormalizesAndReturns(t *testing.T) {
	repo := &repositoryStub{createFn: func(_ context.Context, r SaveRequest) (*Facility, error) {
		if r.Name != "Unidade Teste" || r.UnitType != "hospital" || r.City != "Garça" || r.State != "SP" {
			t.Fatalf("normalizacao incorreta: %+v", r)
		}
		if r.CNES != nil || r.Phone != nil || r.Address != nil {
			t.Fatalf("campos opcionais vazios deveriam ser nil: %+v", r)
		}
		return &Facility{ID: "created", Name: r.Name}, nil
	}}
	r := validSave()
	r.CNES = stringPtr(" ")
	r.Phone = stringPtr(" ")
	r.Address = stringPtr(" ")
	f, err := NewService(repo).Create(context.Background(), r)
	if err != nil || f.ID != "created" || repo.calls != 1 {
		t.Fatalf("resultado: %+v %v calls=%d", f, err, repo.calls)
	}
}
func TestCreatePropagatesRepositoryError(t *testing.T) {
	repo := &repositoryStub{createFn: func(context.Context, SaveRequest) (*Facility, error) { return nil, ErrDuplicateCNES }}
	_, err := NewService(repo).Create(context.Background(), validSave())
	if !errors.Is(err, ErrDuplicateCNES) {
		t.Fatal(err)
	}
}
func TestUpdateAndGetByID(t *testing.T) {
	id := uuid.NewString()
	repo := &repositoryStub{}
	svc := NewService(repo)
	if _, err := svc.Update(context.Background(), "bad", validSave()); !errors.Is(err, ErrInvalidID) {
		t.Fatal(err)
	}
	if _, err := svc.GetByID(context.Background(), "bad"); !errors.Is(err, ErrInvalidID) {
		t.Fatal(err)
	}
	if repo.calls != 0 {
		t.Fatal("UUID invalido chegou ao repository")
	}
	updated, err := svc.Update(context.Background(), id, validSave())
	if err != nil || updated.ID != id {
		t.Fatalf("update: %+v %v", updated, err)
	}
	found, err := svc.GetByID(context.Background(), id)
	if err != nil || found.ID != id {
		t.Fatalf("get: %+v %v", found, err)
	}
	if repo.calls != 2 {
		t.Fatalf("calls=%d", repo.calls)
	}
}
func TestListValidation(t *testing.T) {
	cases := []struct {
		name   string
		filter ListFilter
	}{
		{"estado invalido", ListFilter{State: "SPP"}},
		{"somente latitude", ListFilter{Latitude: number(-22)}},
		{"somente longitude", ListFilter{Longitude: number(-49)}},
		{"raio zero", ListFilter{Latitude: number(-22), Longitude: number(-49)}},
		{"raio acima do limite", ListFilter{Latitude: number(-22), Longitude: number(-49), RadiusMeters: 200001}},
		{"latitude invalida", ListFilter{Latitude: number(-91), Longitude: number(-49), RadiusMeters: 1000}},
		{"longitude invalida", ListFilter{Latitude: number(-22), Longitude: number(181), RadiusMeters: 1000}},
		{"raio NaN", ListFilter{Latitude: number(-22), Longitude: number(-49), RadiusMeters: math.NaN()}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &repositoryStub{}
			_, err := NewService(repo).List(context.Background(), tc.filter)
			if !errors.Is(err, ErrInvalidInput) || repo.calls != 0 {
				t.Fatalf("err=%v calls=%d", err, repo.calls)
			}
		})
	}
}
func TestListDefaultsCapsAndNormalizes(t *testing.T) {
	repo := &repositoryStub{listFn: func(_ context.Context, f ListFilter) ([]Facility, error) {
		if f.Limit != 50 || f.State != "SP" || f.City != "Garça" {
			t.Fatalf("filtro: %+v", f)
		}
		return []Facility{{ID: "1"}}, nil
	}}
	items, err := NewService(repo).List(context.Background(), ListFilter{State: " sp ", City: " Garça "})
	if err != nil || len(items) != 1 {
		t.Fatalf("items=%+v err=%v", items, err)
	}
	repo.listFn = func(_ context.Context, f ListFilter) ([]Facility, error) {
		if f.Limit != 100 {
			t.Fatalf("limite = %d", f.Limit)
		}
		return nil, nil
	}
	_, err = NewService(repo).List(context.Background(), ListFilter{Limit: 999})
	if err != nil {
		t.Fatal(err)
	}
}
func TestListPropagatesError(t *testing.T) {
	want := errors.New("database unavailable")
	repo := &repositoryStub{listFn: func(context.Context, ListFilter) ([]Facility, error) { return nil, want }}
	_, err := NewService(repo).List(context.Background(), ListFilter{})
	if !errors.Is(err, want) {
		t.Fatal(err)
	}
}
