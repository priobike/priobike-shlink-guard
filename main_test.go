package main

import (
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckRequestValid(t *testing.T) {
	req := httptest.NewRequest("GET", "/rest/v3/short-urls", nil)
	w := httptest.NewRecorder()
	if !checkRequest(w, req) {
		t.Errorf("Expected true, got false")
	}
}

func TestCheckRequestInvalidEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/rest/v3/short-url", nil)
	w := httptest.NewRecorder()
	if checkRequest(w, req) {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckGetRequestValid(t *testing.T) {
	req := httptest.NewRequest("GET", "/rest/v3/short-urls/123", nil)
	w := httptest.NewRecorder()
	if !checkGetRequest(w, req) {
		t.Errorf("Expected true, got false")
	}
}

func TestCheckGetRequestMissingShortlink(t *testing.T) {
	req := httptest.NewRequest("GET", "/rest/v3/short-urls", nil)
	w := httptest.NewRecorder()
	if checkGetRequest(w, req) {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckGetRequestEmptyShortlink(t *testing.T) {
	req := httptest.NewRequest("GET", "/rest/v3/short-urls/", nil)
	w := httptest.NewRecorder()
	if checkGetRequest(w, req) {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckBodyValid(t *testing.T) {
	req := httptest.NewRequest("POST", "/rest/v3/short-urls", nil)
	longUrl := "http://example.com/drth5"
	bodyJson := fmt.Sprintf(`{"longUrl": "%s"}`, longUrl)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(bodyJson))
	ok, _, longUrlReturn := checkBody(req)
	if !ok {
		t.Errorf("Expected true, got false")
	}
	if longUrl != longUrlReturn {
		t.Errorf("Expected %s, got %s", longUrl, longUrl)
	}
}

func TestCheckBodyInvalidContentType(t *testing.T) {
	req := httptest.NewRequest("POST", "/rest/v3/short-urls", nil)
	body := "aergfwrger"
	req.Header.Set("Content-Type", "text/plain")
	req.Body = io.NopCloser(strings.NewReader(body))
	ok, _, _ := checkBody(req)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckBodyMissingBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/rest/v3/short-urls", nil)
	req.Header.Set("Content-Type", "application/json")
	ok, _, _ := checkBody(req)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckBodyInvalidJson(t *testing.T) {
	req := httptest.NewRequest("POST", "/rest/v3/short-urls", nil)
	body := "aergfwrger"
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(body))
	ok, _, _ := checkBody(req)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckBodyMissingLongUrl(t *testing.T) {
	req := httptest.NewRequest("POST", "/rest/v3/short-urls", nil)
	body := `{longU: "http://example.com/drth5"}`
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(body))
	ok, _, _ := checkBody(req)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckBodyLongUrlNotString(t *testing.T) {
	req := httptest.NewRequest("POST", "/rest/v3/short-urls", nil)
	body := `{longUrl: 123}`
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(body))
	ok, _, _ := checkBody(req)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckLongUrlValid(t *testing.T) {
	longUrl := "http://example.com/import/eyJ0eXBlIjoiU2hvcnRjdXRMb2NhdGlvbiIsImlkIjoiWyM1ZTYzOV0iLCJuYW1lIjoiIiwid2F5cG9pbnQiOnsibGF0Ijo1My41NDE1NzAxMDc3NzY2LCJsb24iOjkuOTg0Mjc1NjA1Nzk0Njg2LCJhZGRyZXNzIjoiRWxicGhpbGhhcm1vbmllIEhhbWJ1cmcsIFBsYXR6IGRlciBEZXV0c2NoZW4gRWluaGVpdCwgSGFtYnVyZyJ9fQ=="
	ok, _ := checkLongUrl(longUrl)
	if !ok {
		t.Errorf("Expected true, got false")
	}
}

func TestCheckLongUrlMissingImport(t *testing.T) {
	longUrl := "http://example.com"
	ok, _ := checkLongUrl(longUrl)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckLongUrlMissingBase64(t *testing.T) {
	longUrl := "http://example.com/import/"
	ok, _ := checkLongUrl(longUrl)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckLongUrlInvalidBase64(t *testing.T) {
	longUrl := "http://example.com/import/123"
	ok, _ := checkLongUrl(longUrl)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckLongUrlInvalidJson(t *testing.T) {
	longUrl := "http://example.com/import/aGFsbG8gd2VsdA=="
	ok, _ := checkLongUrl(longUrl)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckShortcutValid(t *testing.T) {
	shortcutLocation := map[string]interface{}{
		"type": "ShortcutLocation",
		"id":   "[#5e639]",
		"name": "",
		"waypoint": map[string]interface{}{
			"lat":     53.5415701077766,
			"lon":     9.984275605794686,
			"address": "Elbphilharmonie Hamburg, Platz der Deutschen Einheit, Hamburg",
		},
	}
	ok, shortcutType := checkShortcut(shortcutLocation)
	if !ok {
		t.Errorf("Expected true, got false")
	}

	if shortcutType != "ShortcutLocation" {
		t.Errorf("Expected ShortcutLocation, got %s", shortcutType)
	}

	shortcutRoute := map[string]interface{}{
		"type": "ShortcutRoute",
		"id":   "[#dd9f9]",
		"name": "",
		"waypoints": [2]map[string]interface{}{
			{
				"lat":     53.5522524,
				"lon":     9.9313068,
				"address": "Altona-Altstadt, 22767, Hamburg, Deutschland",
			},
			{
				"lat":     53.5536507,
				"lon":     9.9893664,
				"address": "Jungfernstieg, Altstadt, 20095, Hamburg, Deutschland",
			},
		},
		"routeLengthText": "4.8 km",
		"routeTimeText":   "17 Min.",
	}
	ok, shortcutType = checkShortcut(shortcutRoute)
	if !ok {
		t.Errorf("Expected true, got false")
	}

	if shortcutType != "ShortcutRoute" {
		t.Errorf("Expected ShortcutRoute, got %s", shortcutType)
	}
}

func TestCheckShortcutJsonNil(t *testing.T) {
	ok, _ := checkShortcut(nil)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckShortcutTypeKeyMissing(t *testing.T) {
	shortcutLocation := map[string]interface{}{
		"id":   "[#5e639]",
		"name": "",
		"waypoint": map[string]interface{}{
			"lat":     53.5415701077766,
			"lon":     9.984275605794686,
			"address": "Elbphilharmonie Hamburg, Platz der Deutschen Einheit, Hamburg",
		},
	}
	ok, _ := checkShortcut(shortcutLocation)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckShortcutTypeKeyInvalid(t *testing.T) {
	shortcutLocation := map[string]interface{}{
		"type": "Shortcut",
		"id":   "[#5e639]",
		"name": "",
		"waypoint": map[string]interface{}{
			"lat":     53.5415701077766,
			"lon":     9.984275605794686,
			"address": "Elbphilharmonie Hamburg, Platz der Deutschen Einheit, Hamburg",
		},
	}
	ok, _ := checkShortcut(shortcutLocation)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckShortcutIdKeyMissing(t *testing.T) {
	shortcutLocation := map[string]interface{}{
		"type": "ShortcutLocation",
		"name": "",
		"waypoint": map[string]interface{}{
			"lat":     53.5415701077766,
			"lon":     9.984275605794686,
			"address": "Elbphilharmonie Hamburg, Platz der Deutschen Einheit, Hamburg",
		},
	}
	ok, _ := checkShortcut(shortcutLocation)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckShortcutNameKeyMissing(t *testing.T) {
	shortcutLocation := map[string]interface{}{
		"type": "ShortcutLocation",
		"id":   "[#5e639]",
		"waypoint": map[string]interface{}{
			"lat":     53.5415701077766,
			"lon":     9.984275605794686,
			"address": "Elbphilharmonie Hamburg, Platz der Deutschen Einheit, Hamburg",
		},
	}
	ok, _ := checkShortcut(shortcutLocation)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckLocationShortcutValid(t *testing.T) {
	shortcutLocation := map[string]interface{}{
		"type": "ShortcutLocation",
		"id":   "[#5e639]",
		"name": "",
		"waypoint": map[string]interface{}{
			"lat":     53.5415701077766,
			"lon":     9.984275605794686,
			"address": "Elbphilharmonie Hamburg, Platz der Deutschen Einheit, Hamburg",
		},
	}
	ok := checkLocationShortcut(shortcutLocation)
	if !ok {
		t.Errorf("Expected true, got false")
	}
}

func TestCheckLocationShortcutTooManyKeys(t *testing.T) {
	shortcutLocation := map[string]interface{}{
		"type": "ShortcutLocation",
		"id":   "[#5e639]",
		"name": "",
		"waypoint": map[string]interface{}{
			"lat":     53.5415701077766,
			"lon":     9.984275605794686,
			"address": "Elbphilharmonie Hamburg, Platz der Deutschen Einheit, Hamburg",
		},
		"extra": "extra",
	}
	ok := checkLocationShortcut(shortcutLocation)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckLocationShortcutWaypointKeyMissing(t *testing.T) {
	shortcutLocation := map[string]interface{}{
		"type": "ShortcutLocation",
		"id":   "[#5e639]",
		"name": "",
	}
	ok := checkLocationShortcut(shortcutLocation)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckRouteShortcutValid(t *testing.T) {
	shortcutRoute := map[string]interface{}{
		"type": "ShortcutRoute",
		"id":   "[#dd9f9]",
		"name": "",
		"waypoints": [2]map[string]interface{}{
			{
				"lat":     53.5522524,
				"lon":     9.9313068,
				"address": "Altona-Altstadt, 22767, Hamburg, Deutschland",
			},
			{
				"lat":     53.5536507,
				"lon":     9.9893664,
				"address": "Jungfernstieg, Altstadt, 20095, Hamburg, Deutschland",
			},
		},
		"routeLengthText": "4.8 km",
		"routeTimeText":   "17 Min.",
	}
	ok := checkRouteShortcut(shortcutRoute)
	if !ok {
		t.Errorf("Expected true, got false")
	}
}

func TestCheckRouteShortcutTooManyKeys(t *testing.T) {
	shortcutRoute := map[string]interface{}{
		"type": "ShortcutRoute",
		"id":   "[#dd9f9]",
		"name": "",
		"waypoints": [2]map[string]interface{}{
			{
				"lat":     53.5522524,
				"lon":     9.9313068,
				"address": "Altona-Altstadt, 22767, Hamburg, Deutschland",
			},
			{
				"lat":     53.5536507,
				"lon":     9.9893664,
				"address": "Jungfernstieg, Altstadt, 20095, Hamburg, Deutschland",
			},
		},
		"routeLengthText": "4.8 km",
		"routeTimeText":   "17 Min.",
		"extra":           "extra",
	}
	ok := checkRouteShortcut(shortcutRoute)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckRouteShortcutWaypointsKeyMissing(t *testing.T) {
	shortcutRoute := map[string]interface{}{
		"type": "ShortcutRoute",
		"id":   "[#dd9f9]",
		"name": "",
	}
	ok := checkRouteShortcut(shortcutRoute)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckRouteShortcutRouteTimeTextKeyMissing(t *testing.T) {
	shortcutRoute := map[string]interface{}{
		"type": "ShortcutRoute",
		"id":   "[#dd9f9]",
		"name": "",
		"waypoints": [2]map[string]interface{}{
			{
				"lat":     53.5522524,
				"lon":     9.9313068,
				"address": "Altona-Altstadt, 22767, Hamburg, Deutschland",
			},
			{
				"lat":     53.5536507,
				"lon":     9.9893664,
				"address": "Jungfernstieg, Altstadt, 20095, Hamburg, Deutschland",
			},
		},
		"routeLengthText": "4.8 km",
	}
	ok := checkRouteShortcut(shortcutRoute)
	if ok {
		t.Errorf("Expected false, got true")
	}
}

func TestCheckRouteShortcutRouteLengthTextKeyMissing(t *testing.T) {
	shortcutRoute := map[string]interface{}{
		"type": "ShortcutRoute",
		"id":   "[#dd9f9]",
		"name": "",
		"waypoints": [2]map[string]interface{}{
			{
				"lat":     53.5522524,
				"lon":     9.9313068,
				"address": "Altona-Altstadt, 22767, Hamburg, Deutschland",
			},
			{
				"lat":     53.5536507,
				"lon":     9.9893664,
				"address": "Jungfernstieg, Altstadt, 20095, Hamburg, Deutschland",
			},
		},
		"routeTimeText": "17 Min.",
	}
	ok := checkRouteShortcut(shortcutRoute)
	if ok {
		t.Errorf("Expected false, got true")
	}
}
