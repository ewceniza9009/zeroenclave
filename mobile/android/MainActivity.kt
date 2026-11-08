package com.onevault.app

import android.annotation.SuppressLint
import android.content.Context
import android.content.Intent
import android.graphics.Color
import android.net.Uri
import android.os.Bundle
import android.view.WindowManager
import android.webkit.JavascriptInterface
import android.webkit.PermissionRequest
import android.webkit.ValueCallback
import android.webkit.WebChromeClient
import android.webkit.WebResourceRequest
import android.webkit.WebSettings
import android.webkit.WebView
import android.webkit.WebViewClient
import android.widget.FrameLayout
import android.widget.Toast
import androidx.activity.OnBackPressedCallback
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AppCompatActivity
import androidx.biometric.BiometricManager
import androidx.biometric.BiometricPrompt
import androidx.core.content.ContextCompat
import java.util.concurrent.Executor

class MainActivity : AppCompatActivity() {

    private lateinit var webView: WebView
    private var serverPort: Int = -1

    private var fileUploadCallback: ValueCallback<Array<Uri>>? = null

    private val fileChooserLauncher = registerForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) { result ->
        if (fileUploadCallback == null) return@registerForActivityResult

        val uris: Array<Uri>? = if (result.resultCode == RESULT_OK) {
            val intent = result.data
            val clipData = intent?.clipData
            if (clipData != null) {
                Array(clipData.itemCount) { i -> clipData.getItemAt(i).uri }
            } else if (intent?.data != null) {
                arrayOf(intent.data!!)
            } else {
                null
            }
        } else {
            null
        }

        fileUploadCallback?.onReceiveValue(uris)
        fileUploadCallback = null
    }

    companion object {
        init {
            System.loadLibrary("onevault")
        }
    }

    private external fun startVaultServer(vaultDir: String): Int
    private external fun stopVaultServer()

    @SuppressLint("SetJavaScriptEnabled")
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        // Set status and navigation bar color to match 1Vault dark theme (#131417)
        window.apply {
            clearFlags(WindowManager.LayoutParams.FLAG_TRANSLUCENT_STATUS)
            addFlags(WindowManager.LayoutParams.FLAG_DRAWS_SYSTEM_BAR_BACKGROUNDS)
            statusBarColor = Color.parseColor("#131417")
            navigationBarColor = Color.parseColor("#131417")
        }

        // Start internal Go vault server
        val vaultDir = filesDir.absolutePath
        serverPort = startVaultServer(vaultDir)

        if (serverPort <= 0) {
            Toast.makeText(this, "Failed to start ZeroEnclave secure engine", Toast.LENGTH_LONG).show()
            finish()
            return
        }

        if (ContextCompat.checkSelfPermission(this, android.Manifest.permission.CAMERA) != android.content.pm.PackageManager.PERMISSION_GRANTED) {
            requestPermissions(arrayOf(android.Manifest.permission.CAMERA), 101)
        }

        // Initialize WebView
        webView = WebView(this).apply {
            setBackgroundColor(Color.parseColor("#131417"))
            layoutParams = FrameLayout.LayoutParams(
                FrameLayout.LayoutParams.MATCH_PARENT,
                FrameLayout.LayoutParams.MATCH_PARENT
            )

            settings.apply {
                javaScriptEnabled = true
                domStorageEnabled = true
                databaseEnabled = true
                useWideViewPort = true
                loadWithOverviewMode = true
                setSupportZoom(false)
                displayZoomControls = false
                cacheMode = WebSettings.LOAD_NO_CACHE
                mixedContentMode = WebSettings.MIXED_CONTENT_NEVER_ALLOW
            }
            clearCache(true)

            addJavascriptInterface(AndroidBridge(), "AndroidBridge")

            webChromeClient = object : WebChromeClient() {
                override fun onShowFileChooser(
                    view: WebView?,
                    filePathCallback: ValueCallback<Array<Uri>>?,
                    fileChooserParams: FileChooserParams?
                ): Boolean {
                    fileUploadCallback?.onReceiveValue(null)
                    fileUploadCallback = filePathCallback

                    val intent = fileChooserParams?.createIntent() ?: Intent(Intent.ACTION_GET_CONTENT).apply {
                        type = "*/*"
                        addCategory(Intent.CATEGORY_OPENABLE)
                    }

                    try {
                        fileChooserLauncher.launch(intent)
                    } catch (e: Exception) {
                        fileUploadCallback?.onReceiveValue(null)
                        fileUploadCallback = null
                        return false
                    }
                    return true
                }

                override fun onPermissionRequest(request: PermissionRequest?) {
                    runOnUiThread {
                        request?.grant(request.resources)
                    }
                }
            }

            webViewClient = object : WebViewClient() {
                override fun shouldOverrideUrlLoading(view: WebView?, request: WebResourceRequest?): Boolean {
                    val url = request?.url?.toString() ?: return false
                    val uri = request?.url ?: return false
                    val host = uri.host ?: ""
                    // Keep local server traffic inside WebView unless it's a domain path like /gmail.com
                    if (host == "127.0.0.1" || host == "localhost") {
                        val path = uri.path ?: ""
                        val trimmed = path.trimStart('/')
                        if (trimmed.isNotEmpty() && !trimmed.startsWith("api/") && trimmed.contains(".") &&
                            !trimmed.endsWith(".js") && !trimmed.endsWith(".css") && !trimmed.endsWith(".html") &&
                            !trimmed.endsWith(".png") && !trimmed.endsWith(".ico") && !trimmed.endsWith(".svg") &&
                            !trimmed.endsWith(".woff2") && !trimmed.endsWith(".json")) {
                            try {
                                val target = "https://$trimmed" + (if (uri.query != null) "?${uri.query}" else "")
                                val intent = Intent(Intent.ACTION_VIEW, Uri.parse(target)).apply {
                                    addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                                }
                                startActivity(intent)
                                return true
                            } catch (e: Exception) {
                                return false
                            }
                        }
                        return false
                    }
                    try {
                        var target = url.trim()
                        if (!target.startsWith("http://") && !target.startsWith("https://")) {
                            target = "https://$target"
                        }
                        val intent = Intent(Intent.ACTION_VIEW, Uri.parse(target)).apply {
                            addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                        }
                        startActivity(intent)
                        return true
                    } catch (e: Exception) {
                        return false
                    }
                }
            }
        }

        setContentView(webView)
        webView.loadUrl("http://localhost:$serverPort/")

        // Handle native back button navigation
        onBackPressedDispatcher.addCallback(this, object : OnBackPressedCallback(true) {
            override fun handleOnBackPressed() {
                if (::webView.isInitialized) {
                    webView.evaluateJavascript("window.handleAndroidBack ? window.handleAndroidBack() : false") { result ->
                        if (result != "true") {
                            if (webView.canGoBack()) {
                                webView.goBack()
                            } else {
                                isEnabled = false
                                onBackPressedDispatcher.onBackPressed()
                            }
                        }
                    }
                } else {
                    isEnabled = false
                    onBackPressedDispatcher.onBackPressed()
                }
            }
        })
    }

    override fun onDestroy() {
        super.onDestroy()
        stopVaultServer()
    }

    inner class AndroidBridge {
        private val prefs by lazy {
            getSharedPreferences("onevault_secure_prefs", Context.MODE_PRIVATE)
        }

        @JavascriptInterface
        fun isBiometricAvailable(): Boolean {
            val biometricManager = BiometricManager.from(this@MainActivity)
            return biometricManager.canAuthenticate(
                BiometricManager.Authenticators.BIOMETRIC_STRONG or BiometricManager.Authenticators.BIOMETRIC_WEAK
            ) == BiometricManager.BIOMETRIC_SUCCESS
        }

        @JavascriptInterface
        fun openExternal(rawUrl: String) {
            runOnUiThread {
                try {
                    var target = rawUrl.trim()
                    if (!target.startsWith("http://") && !target.startsWith("https://")) {
                        target = "https://$target"
                    }
                    val intent = Intent(Intent.ACTION_VIEW, Uri.parse(target)).apply {
                        addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                    }
                    startActivity(intent)
                } catch (e: Exception) {
                    Toast.makeText(this@MainActivity, "Could not open browser", Toast.LENGTH_SHORT).show()
                }
            }
        }

        @JavascriptInterface
        fun showToast(message: String) {
            runOnUiThread {
                Toast.makeText(this@MainActivity, message, Toast.LENGTH_SHORT).show()
            }
        }

        @JavascriptInterface
        fun getStoredSecretKey(): String {
            return prefs.getString("secret_key", "") ?: ""
        }

        @JavascriptInterface
        fun saveStoredSecretKey(key: String) {
            if (key.isNotBlank()) {
                prefs.edit().putString("secret_key", key.trim()).apply()
            }
        }

        @JavascriptInterface
        fun hasStoredSecretKey(): Boolean {
            val k = prefs.getString("secret_key", "") ?: ""
            return k.isNotBlank()
        }

        @JavascriptInterface
        fun clearStoredSecretKey() {
            prefs.edit().remove("secret_key").apply()
        }

        @JavascriptInterface
        fun saveBiometricPassword(password: String) {
            if (password.isNotBlank()) {
                prefs.edit().putString("bio_pass", password).apply()
            }
        }

        @JavascriptInterface
        fun hasBiometricPassword(): Boolean {
            val p = prefs.getString("bio_pass", "") ?: ""
            return p.isNotBlank()
        }

        @JavascriptInterface
        fun clearBiometricPassword() {
            prefs.edit().remove("bio_pass").apply()
        }

        @JavascriptInterface
        fun triggerBiometricUnlock() {
            runOnUiThread {
                promptBiometricUnlock()
            }
        }
    }

    private fun promptBiometricUnlock() {
        val prefs = getSharedPreferences("onevault_secure_prefs", Context.MODE_PRIVATE)
        val savedPass = prefs.getString("bio_pass", "") ?: ""
        val savedKey = prefs.getString("secret_key", "") ?: ""

        val executor: Executor = ContextCompat.getMainExecutor(this)
        val prompt = BiometricPrompt(this, executor, object : BiometricPrompt.AuthenticationCallback() {
            override fun onAuthenticationSucceeded(result: BiometricPrompt.AuthenticationResult) {
                super.onAuthenticationSucceeded(result)
                val escapedPass = savedPass.replace("\\", "\\\\").replace("'", "\\'")
                val escapedKey = savedKey.replace("\\", "\\\\").replace("'", "\\'")
                webView.evaluateJavascript("if (window.onBiometricUnlocked) { window.onBiometricUnlocked('$escapedPass', '$escapedKey'); }", null)
            }

            override fun onAuthenticationError(errorCode: Int, errString: CharSequence) {
                super.onAuthenticationError(errorCode, errString)
                if (errorCode != BiometricPrompt.ERROR_USER_CANCELED && errorCode != BiometricPrompt.ERROR_NEGATIVE_BUTTON) {
                    webView.evaluateJavascript("if (window.showToast) { window.showToast('$errString'); }", null)
                }
            }

            override fun onAuthenticationFailed() {
                super.onAuthenticationFailed()
                webView.evaluateJavascript("if (window.showToast) { window.showToast('Fingerprint not recognized'); }", null)
            }
        })

        val promptInfo = BiometricPrompt.PromptInfo.Builder()
            .setTitle("Unlock ZeroEnclave")
            .setSubtitle("Confirm your fingerprint or biometric credential")
            .setNegativeButtonText("Use Master Password")
            .build()

        prompt.authenticate(promptInfo)
    }
}
