//go:build android

package main

/*
#include <jni.h>
#include <stdlib.h>

static char* get_string_utf_chars(JNIEnv* env, jstring str) {
    return (char*)(*env)->GetStringUTFChars(env, str, NULL);
}

static void release_string_utf_chars(JNIEnv* env, jstring str, const char* chars) {
    (*env)->ReleaseStringUTFChars(env, str, chars);
}

static jstring new_string_utf(JNIEnv* env, const char* chars) {
    return (*env)->NewStringUTF(env, chars);
}
*/
import "C"
import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unsafe"

	"1vault/internal/crypto"
	"1vault/internal/server"
	vaultsync "1vault/internal/sync"
	"1vault/internal/vault"
	"1vault/web"
)

var (
	mu          sync.Mutex
	httpServer  *http.Server
	globalStore *vault.Store
	serverPort  int
)

//export Java_com_onevault_app_MainActivity_startVaultServer
func Java_com_onevault_app_MainActivity_startVaultServer(env *C.JNIEnv, thiz C.jobject, jVaultDir C.jstring) C.jint {
	mu.Lock()
	defer mu.Unlock()

	if httpServer != nil && serverPort > 0 {
		return C.jint(serverPort)
	}

	cStr := C.get_string_utf_chars(env, jVaultDir)
	vaultDir := C.GoString(cStr)
	C.release_string_utf_chars(env, jVaultDir, cStr)

	// Persist secret key directly in app private storage
	crypto.SetKeyStorageDir(vaultDir)

	vaultPath := filepath.Join(vaultDir, "vault.1v")

	// Initialize vault store (15 minutes autolock)
	globalStore = vault.NewStore(vaultPath, 15)

	// Initialize embedded web server
	srvHandler := server.NewServer(globalStore, web.StaticFS())

	// Bind to collision-free port (preferred 51892, with dynamic OS fallback)
	listener, err := net.Listen("tcp", "127.0.0.1:51892")
	if err != nil {
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return -1
		}
	}

	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return -1
	}
	serverPort = tcpAddr.Port

	httpServer = &http.Server{
		Handler:      srvHandler.Routes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		_ = httpServer.Serve(listener)
	}()

	// Start LAN peer discovery service
	_ = vaultsync.GetDiscoveryService(serverPort)

	return C.jint(serverPort)
}

//export Java_com_onevault_app_MainActivity_stopVaultServer
func Java_com_onevault_app_MainActivity_stopVaultServer(env *C.JNIEnv, thiz C.jobject) {
	mu.Lock()
	defer mu.Unlock()

	if globalStore != nil {
		globalStore.Lock()
	}
	if httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(ctx)
		httpServer = nil
		serverPort = 0
	}
}

func jstringToGo(env *C.JNIEnv, jstr C.jstring) string {
	if jstr == 0 {
		return ""
	}
	cStr := C.get_string_utf_chars(env, jstr)
	if cStr == nil {
		return ""
	}
	res := C.GoString(cStr)
	C.release_string_utf_chars(env, jstr, cStr)
	return res
}

func goToJstring(env *C.JNIEnv, s string) C.jstring {
	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))
	return C.new_string_utf(env, cStr)
}

func ensureUnlockedStore(vaultDir, password, secretKey string) *vault.Store {
	crypto.SetKeyStorageDir(vaultDir)
	if globalStore == nil {
		vaultPath := filepath.Join(vaultDir, "vault.1v")
		globalStore = vault.NewStore(vaultPath, 15)
	}
	if !globalStore.IsUnlocked() && password != "" {
		if secretKey == "" {
			if sk, err := crypto.LoadSecretKey(); err == nil {
				secretKey = sk
			}
		}
		if secretKey != "" {
			_ = globalStore.Unlock(password, secretKey)
		}
	}
	return globalStore
}

//export Java_com_onevault_app_OneVaultAutofillService_queryAutofillItems
func Java_com_onevault_app_OneVaultAutofillService_queryAutofillItems(env *C.JNIEnv, thiz C.jobject, jVaultDir C.jstring, jQuery C.jstring, jPassword C.jstring, jSecretKey C.jstring) C.jstring {
	mu.Lock()
	defer mu.Unlock()

	vaultDir := jstringToGo(env, jVaultDir)
	query := strings.ToLower(strings.TrimSpace(jstringToGo(env, jQuery)))
	password := jstringToGo(env, jPassword)
	secretKey := jstringToGo(env, jSecretKey)

	store := ensureUnlockedStore(vaultDir, password, secretKey)
	if store == nil || !store.IsUnlocked() {
		return goToJstring(env, `{"unlocked":false,"items":[]}`)
	}

	items, err := store.ListItems()
	if err != nil {
		return goToJstring(env, `{"unlocked":true,"items":[]}`)
	}

	// Clean query (could be web domain like "github.com" or Android package like "com.github.android")
	queryTokens := strings.FieldsFunc(query, func(r rune) bool {
		return r == '.' || r == '/' || r == ':' || r == '-' || r == '_'
	})
	ignoreTokens := map[string]bool{
		"com": true, "org": true, "net": true, "io": true, "co": true,
		"app": true, "android": true, "www": true, "http": true, "https": true,
		"mobile": true, "auth": true, "login": true,
	}
	var meaningfulTokens []string
	for _, tk := range queryTokens {
		if len(tk) >= 3 && !ignoreTokens[tk] {
			meaningfulTokens = append(meaningfulTokens, tk)
		}
	}

	type autofillDTO struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Username string `json:"username"`
		Password string `json:"password"`
		URL      string `json:"url"`
	}

	var matched []autofillDTO
	var fallbacks []autofillDTO

	for _, it := range items {
		if it.Trashed || (it.Password == "" && it.Username == "") {
			continue
		}
		dto := autofillDTO{
			ID:       it.ID,
			Title:    it.Title,
			Username: it.Username,
			Password: it.Password,
			URL:      it.URL,
		}
		fallbacks = append(fallbacks, dto)

		if query == "" {
			continue
		}
		uLower := strings.ToLower(it.URL)
		tLower := strings.ToLower(it.Title)
		if strings.Contains(uLower, query) || strings.Contains(tLower, query) {
			matched = append(matched, dto)
			continue
		}
		for _, tk := range meaningfulTokens {
			if strings.Contains(uLower, tk) || strings.Contains(tLower, tk) {
				matched = append(matched, dto)
				break
			}
		}
	}

	if len(matched) == 0 && len(fallbacks) > 0 {
		// Provide top 8 logins if no exact domain/package match so the user can still pick from Autofill
		limit := 8
		if len(fallbacks) < limit {
			limit = len(fallbacks)
		}
		matched = fallbacks[:limit]
	}

	out, _ := json.Marshal(map[string]interface{}{
		"unlocked": true,
		"items":    matched,
	})
	return goToJstring(env, string(out))
}

//export Java_com_onevault_app_OneVaultAutofillService_saveAutofillItem
func Java_com_onevault_app_OneVaultAutofillService_saveAutofillItem(env *C.JNIEnv, thiz C.jobject, jVaultDir C.jstring, jTitle C.jstring, jUsername C.jstring, jPassword C.jstring, jUrl C.jstring, jMasterPass C.jstring, jSecretKey C.jstring) C.jboolean {
	mu.Lock()
	defer mu.Unlock()

	vaultDir := jstringToGo(env, jVaultDir)
	title := strings.TrimSpace(jstringToGo(env, jTitle))
	username := strings.TrimSpace(jstringToGo(env, jUsername))
	password := jstringToGo(env, jPassword)
	urlStr := strings.TrimSpace(jstringToGo(env, jUrl))
	masterPass := jstringToGo(env, jMasterPass)
	secretKey := jstringToGo(env, jSecretKey)

	if password == "" {
		return C.JNI_FALSE
	}
	if title == "" {
		title = urlStr
	}
	if title == "" {
		title = "Android Login"
	}

	store := ensureUnlockedStore(vaultDir, masterPass, secretKey)
	if store == nil || !store.IsUnlocked() {
		return C.JNI_FALSE
	}

	_, err := store.AddItem(vault.VaultItem{
		Title:    title,
		Category: vault.CategoryLogin,
		Username: username,
		Password: password,
		URL:      urlStr,
		Tags:     []string{"android-autofill"},
	})
	if err != nil {
		return C.JNI_FALSE
	}
	return C.JNI_TRUE
}

func main() {}
