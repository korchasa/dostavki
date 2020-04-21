package locator

import (
	"encoding/json"
	"fmt"
	"github.com/korchasa/kogdaeda/types"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ZakazClient struct {
	cl *http.Client
}

func New(cl *http.Client) *ZakazClient {
	return &ZakazClient{
		cl: cl,
	}
}

func (z *ZakazClient) Schedule(lat, long float64) ([][]types.Schedule, error) {
	windows := make([]types.StoreWindow, 0, 8)
	stores, err := z.getStoresIDs(lat, long)
	if err != nil {
		return nil, fmt.Errorf("can't get stores from zakaz.ua: %v", err)
	}
	for _, s := range stores {
		sw, err := z.getWindows(s.ID, lat, long)
		if err != nil {
			return nil, fmt.Errorf("can't get delivery windows for store %d: %v", s.ID, err)
		}
		for _, w := range sw {
			w.Store = s
			windows = append(windows, w)
		}
	}
	return fold(windows), nil
}

func fold(windows []types.StoreWindow) [][]types.Schedule {
	sort.Slice(windows, func(i, j int) bool {
		return windows[i].Start.Before(windows[j].Start)
	})
	byDate := make([][]types.Schedule, 0, 8)
	for _, w := range windows {
		dateExists := false
		for di, d := range byDate {
			if d[0].Start.Format("20060102") == w.Start.Format("20060102") {
				dateExists = true
				windowExists := false
				for fi, f := range d {
					if w.Start.Equal(f.Start) && w.End.Equal(f.End) {
						windowExists = true
						byDate[di][fi].Chains[w.Store.Chain] = true
					}
				}
				if !windowExists {
					byDate[di] = append(byDate[di], types.Schedule{
						Start: w.Start,
						End: w.End,
						Chains: map[string]bool{w.Store.Chain: true},
					})
				}
			}
		}
		if !dateExists {
			byDate = append(byDate, []types.Schedule{{
				Start:  w.Start,
				End:    w.End,
				Chains: map[string]bool{w.Store.Chain: true},
			}})
		}
	}

	return byDate
}

func (z *ZakazClient) getStoresIDs(lat, long float64) ([]types.Store, error) {
	var data []struct {
		Address struct {
			Building string `json:"building"`
			City     string `json:"city"`
			Street   string `json:"street"`
		} `json:"address"`
		City          string   `json:"city"`
		Coords        string   `json:"coords"`
		DeliveryTypes []string `json:"delivery_types"`
		Email         string   `json:"email"`
		ID            string   `json:"id"`
		Name          string   `json:"name"`
		OpeningHours  struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"opening_hours"`
		PaymentTypes []string `json:"payment_types"`
		Phones       []string `json:"phones"`
		RegionID     string   `json:"region_id"`
		RetailChain  string   `json:"retail_chain"`
	}

	u := fmt.Sprintf("https://stores-api.zakaz.ua/stores/?coords=%f,%f", lat, long)
	r, err := z.cl.Get(u)
	if err != nil {
		return nil, fmt.Errorf("can't get windows from zakaz.ua: %v", err)
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("can't decode zakaz.ua JSON: %v", err)
	}

	stores := make([]types.Store, 0, 2)
	for _, s := range data {
		id, err := strconv.Atoi(s.ID)
		if err != nil {
			return nil, fmt.Errorf("can't convert store id: %v", err)
		}
		stores = append(stores, types.Store{
			ID:    uint(id),
			Chain: s.RetailChain,
		})
	}
	return stores, nil
}

func (z *ZakazClient) getWindows(storeID uint, lat, long float64) (windows []types.StoreWindow, err error) {
	type zakBodySpec struct {
		Date  string `json:"date"`
		Items []struct {
			Currency        string  `json:"currency"`
			Date            string  `json:"date"`
			EndOrderingTime float64 `json:"end_ordering_time"`
			ID              string  `json:"id"`
			IsOpen          bool    `json:"is_open"`
			Price           int     `json:"price"`
			TimeRange       string  `json:"time_range"`
		} `json:"items"`
	}
	u := fmt.Sprintf("https://stores-api.zakaz.ua/stores/%d/delivery_schedule/plan/?coords=%f,%f", storeID, lat, long)
	r, err := z.cl.Get(u)
	if err != nil {
		return windows, fmt.Errorf("can't get windows from zakaz.ua: %v", err)
	}
	defer r.Body.Close()

	zd := make([]zakBodySpec, 0, 12)
	err = json.NewDecoder(r.Body).Decode(&zd)
	if err != nil {
		return windows, fmt.Errorf("can't decode zakaz.ua JSON: %v", err)
	}

	for _, d := range zd {
		for _, w := range d.Items {
			if !w.IsOpen {
				continue
			}
			from, to, err := parseTimeString(d.Date, w.TimeRange)
			if err !=  nil {
				return windows, err
			}
			windows = append(windows, types.StoreWindow{
				Start: from,
				End:   to,
				Price: uint(w.Price),
			})
		}
	}
	return
}

func parseTimeString(d string, s string) (start time.Time, end time.Time, err error) {
	loc, err := time.LoadLocation("Europe/Kiev")
	if err != nil {
		return start, end, fmt.Errorf("can't load Kiev location: %v", err)
	}

	times := strings.Split(s, " - ")
	if len(times) != 2 {
		err = fmt.Errorf("bad start-end string from zakaz: %s", s)
		return
	}
	format := "2006-01-02T15:04"
	start, err = time.ParseInLocation(format, fmt.Sprintf("%sT%s", d, times[0]), loc)
	if err != nil {
		err = fmt.Errorf("can't parse start time string: %v", err)
		return
	}
	end, err = time.ParseInLocation(format, fmt.Sprintf("%sT%s", d, times[1]), loc)
	if err != nil {
		err = fmt.Errorf("can't parse end time string: %v", err)
		return
	}
	return start, end, nil
}
