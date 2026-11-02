package healthfacilities

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

type serviceStub struct {
	createFn func(context.Context, SaveRequest) (*Facility, error)
	updateFn func(context.Context, string, SaveRequest) (*Facility, error)
	getFn    func(context.Context, string) (*Facility, error)
	listFn   func(context.Context, ListFilter) ([]Facility, error)
}

func (m *serviceStub) Create(c context.Context, r SaveRequest) (*Facility, error) {
	if m.createFn != nil {
		return m.createFn(c, r)
	}
	return &Facility{ID: "created"}, nil
}
func (m *serviceStub) Update(c context.Context, id string, r SaveRequest) (*Facility, error) {
	if m.updateFn != nil {
		return m.updateFn(c, id, r)
	}
	return &Facility{ID: id}, nil
}
func (m *serviceStub) GetByID(c context.Context, id string) (*Facility, error) {
	if m.getFn != nil {
		return m.getFn(c, id)
	}
	return &Facility{ID: id}, nil
}
func (m *serviceStub) List(c context.Context, f ListFilter) ([]Facility, error) {
	if m.listFn != nil {
		return m.listFn(c, f)
	}
	return []Facility{}, nil
}
func testRouter(s Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(s)
	r.POST("/health-facilities", h.Create)
	r.PATCH("/health-facilities/:id", h.Update)
	r.GET("/public/health-facilities", h.List)
	r.GET("/public/health-facilities/nearby", h.List)
	r.GET("/public/health-facilities/:id", h.GetByID)
	return r
}
func send(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

const goodBody = `{"name":"Unidade Teste","unit_type":"hospital","latitude":-22.2105,"longitude":-49.6564,"city":"Garça","state":"SP"}`

func TestHandlerCreate(t *testing.T) {
	cases := []struct {
		name, body string
		err        error
		want       int
	}{
		{"sucesso", goodBody, nil, 201},
		{"json ausente", "", nil, 400},
		{"json invalido", "{", nil, 400},
		{"campo obrigatorio ausente", `{"name":"Unidade"}`, nil, 400},
		{"cnes duplicado", goodBody, ErrDuplicateCNES, 409},
		{"dados invalidos", goodBody, ErrInvalidInput, 400},
		{"erro interno", goodBody, errors.New("falha"), 500},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := &serviceStub{createFn: func(context.Context, SaveRequest) (*Facility, error) {
				if tc.err != nil {
					return nil, tc.err
				}
				return &Facility{ID: "created"}, nil
			}}
			w := send(testRouter(s), "POST", "/health-facilities", tc.body)
			if w.Code != tc.want {
				t.Fatalf("status=%d want=%d body=%s", w.Code, tc.want, w.Body.String())
			}
		})
	}
}
func TestHandlerUpdate(t *testing.T) {
	id := uuid.NewString()
	s := &serviceStub{updateFn: func(_ context.Context, got string, _ SaveRequest) (*Facility, error) {
		if got != id {
			t.Fatalf("id=%s", got)
		}
		return &Facility{ID: got}, nil
	}}
	w := send(testRouter(s), "PATCH", "/health-facilities/"+id, goodBody)
	if w.Code != 200 {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	w = send(testRouter(s), "PATCH", "/health-facilities/"+id, "")
	if w.Code != 400 {
		t.Fatal(w.Code)
	}
}
func TestHandlerGetByID(t *testing.T) {
	id := uuid.NewString()
	s := &serviceStub{getFn: func(_ context.Context, got string) (*Facility, error) {
		if got != id {
			t.Fatalf("id=%s", got)
		}
		return nil, ErrNotFound
	}}
	w := send(testRouter(s), "GET", "/public/health-facilities/"+id, "")
	if w.Code != 404 {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}
func TestHandlerListAndNearby(t *testing.T) {
	s := &serviceStub{listFn: func(_ context.Context, f ListFilter) ([]Facility, error) {
		if f.Latitude != nil {
			if *f.Latitude != -22.2105 || *f.Longitude != -49.6564 || f.RadiusMeters != 10000 {
				t.Fatalf("filtro nearby: %+v", f)
			}
		}
		return []Facility{{ID: "one"}}, nil
	}}
	r := testRouter(s)
	for _, path := range []string{"/public/health-facilities", "/public/health-facilities/nearby?latitude=-22.2105&longitude=-49.6564"} {
		w := send(r, "GET", path, "")
		if w.Code != 200 {
			t.Fatalf("%s status=%d", path, w.Code)
		}
		var response struct {
			Total int `json:"total"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || response.Total != 1 {
			t.Fatalf("response=%s err=%v", w.Body.String(), err)
		}
	}
}
func TestHandlerInvalidQuery(t *testing.T) {
	r := testRouter(&serviceStub{})
	paths := []string{
		"/public/health-facilities?limit=abc",
		"/public/health-facilities?limit=0",
		"/public/health-facilities/nearby?longitude=-49",
		"/public/health-facilities/nearby?latitude=x&longitude=-49",
		"/public/health-facilities/nearby?latitude=-22&longitude=-49&radius_meters=abc",
	}
	for _, path := range paths {
		w := send(r, "GET", path, "")
		if w.Code != 400 {
			t.Fatalf("%s: status=%d body=%s", path, w.Code, w.Body.String())
		}
	}
}
