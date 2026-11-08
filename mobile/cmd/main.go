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
*/
import "C"
import (
	"context"
	"net"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"1vault/internal/crypto"
	"1vault/internal/server"
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

func main() {}
