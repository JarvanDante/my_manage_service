package logic

import (
	"strings"
	"testing"
)

func TestRemoveTopLevelYAMLKeyPreservesOrderAndComments(t *testing.T) {
	in := `# 头部注释
server:
  address: ":8000"

# front
frontapi:
  address: ":8001"

database:
  default:
    link: "old"
    debug: false

redis:
  default:
    address: "127.0.0.1:6379"
`
	got := removeTopLevelYAMLKey(in, "database")
	if strings.Contains(got, "database:") {
		t.Fatalf("database block should be removed:\n%s", got)
	}
	if !strings.Contains(got, "# 头部注释") || !strings.Contains(got, "# front") {
		t.Fatalf("comments should be preserved:\n%s", got)
	}
}

func TestReplaceOrInsertDatabaseLinkKeepsDebugFalse(t *testing.T) {
	in := `logger:
  level: "all"

database:
  default:
    debug: false
    link: "pgsql:postgres:old@tcp(127.0.0.1:5432)/my?sslmode=disable"

redis:
  default:
    address: "127.0.0.1:6379"
`
	got := replaceOrInsertDatabaseLink(in, "pgsql:postgres:new@tcp(127.0.0.1:5432)/my?sslmode=disable")
	if !strings.Contains(got, "debug: false") {
		t.Fatalf("debug:false must stay:\n%s", got)
	}
	if strings.Contains(got, "debug: true") {
		t.Fatalf("must not introduce debug:true:\n%s", got)
	}
	if !strings.Contains(got, "postgres:new@") {
		t.Fatalf("link not updated:\n%s", got)
	}
	if !strings.Contains(got, `level: "all"`) || !strings.Contains(got, "redis:") {
		t.Fatalf("other sections changed:\n%s", got)
	}
	// debug 行应仍在 database 段内、且相对顺序在 link 前(与原文一致)
	debugAt := strings.Index(got, "debug: false")
	linkAt := strings.Index(got, "link: ")
	if debugAt < 0 || linkAt < 0 || debugAt > linkAt {
		t.Fatalf("debug/link order broken:\n%s", got)
	}
}

func TestReplaceOrInsertDatabaseLinkInsertsWhenMissing(t *testing.T) {
	in := `database:
  default:
    debug: false
`
	got := replaceOrInsertDatabaseLink(in, "pgsql:u:p@tcp(h:1)/db")
	if !strings.Contains(got, "debug: false") {
		t.Fatalf("debug lost:\n%s", got)
	}
	if !strings.Contains(got, "link:") {
		t.Fatalf("link not inserted:\n%s", got)
	}
}

func TestReplaceOrInsertPaasCredsKeepsAdBase(t *testing.T) {
	in := `logger:
  level: "all"

paas:
  app_key: "old_key"
  app_secret: "old_secret"
  media_base: "http://my_media:8004"
  storage_base: "http://my_storage:8005"
  ad_base: "http://my_ad:8006"

redis:
  default:
    address: "127.0.0.1:6379"
`
	got := replaceOrInsertPaasCreds(in, "ak_new", "secret_new")
	if !strings.Contains(got, `app_key: "ak_new"`) {
		t.Fatalf("app_key not updated:\n%s", got)
	}
	if !strings.Contains(got, `app_secret: "secret_new"`) {
		t.Fatalf("app_secret not updated:\n%s", got)
	}
	if !strings.Contains(got, `ad_base: "http://my_ad:8006"`) ||
		!strings.Contains(got, `media_base: "http://my_media:8004"`) ||
		!strings.Contains(got, `storage_base: "http://my_storage:8005"`) {
		t.Fatalf("paas extras must stay:\n%s", got)
	}
	if strings.Contains(got, "old_key") || strings.Contains(got, "old_secret") {
		t.Fatalf("old creds still present:\n%s", got)
	}
	if !strings.Contains(got, `level: "all"`) || !strings.Contains(got, "redis:") {
		t.Fatalf("other sections changed:\n%s", got)
	}
}

func TestReplaceOrInsertPaasCredsInsertsWhenMissing(t *testing.T) {
	in := `paas:
  ad_base: "http://my_ad:8006"
`
	got := replaceOrInsertPaasCreds(in, "ak_1", "sec_1")
	if !strings.Contains(got, `app_key: "ak_1"`) || !strings.Contains(got, `app_secret: "sec_1"`) {
		t.Fatalf("creds not inserted:\n%s", got)
	}
	if !strings.Contains(got, `ad_base: "http://my_ad:8006"`) {
		t.Fatalf("ad_base lost:\n%s", got)
	}
}
