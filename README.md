<p align="center">
  <img src="web/static/logo.svg" width="96" height="96" alt="ZeroEnclave Logo">
</p>

# ZeroEnclave — Zero-Knowledge Password & Media Enclave (Windows & Android)

**ZeroEnclave** (`v1.4.3`) is a standalone, local-first, zero-knowledge password manager, encrypted media vault, and cross-device sync enclave engineered in pure **Go**. It delivers full cross-platform support for **Windows Desktop** and **Android**.

## What's New in v1.4.3

* **Relay Auto-Discovery (IP changes no longer break sync):**
  * When the saved relay address stops responding (e.g. your router gave the relay PC a new IP), the app searches your LAN, finds the relay again, verifies it, and updates the saved Relay URL automatically on Windows and Android.
  * Devices of the same vault advertise a hosted relay through encrypted-bucket-scoped UDP beacons; the relay host also self-heals to its own current LAN address.
* **Mermaid Diagrams in Notes:** write a ```mermaid``` block (or use the new toolbar button) to render flowcharts, sequence/ER/state diagrams, bundled offline and theme-aware.
* **Fixes:** duplicate items created on save, light-mode highlight/separator visibility, consistent Folder/Category combobox, Save/Cancel moved to the modal header, Android portrait viewport.

<details>
<summary>Previous release: v1.4.1</summary>

* **Android Autofill Service:**
  * Native `OneVaultAutofillService` fills usernames and passwords directly inside other Android apps and Chrome/browsers.
  * Captures and saves new logins (`onSaveRequest`) straight into your encrypted vault.
  * 1-tap **Android Autofill** activation shortcut inside the mobile sidebar under Tools.
* **WebAuthn Passkey Storage & Browser Extension Bridge:**
  * Generate, store, and sync hardware-grade **ECDSA P-256 (`ES256`) WebAuthn Passkeys** inside your encrypted vault (`CTAP2` CBOR `COSE_Key`, `authenticatorData`, `attestationObject`, and signature counter).
  * Dedicated **Passkeys** view in the sidebar, `+ Passkey` generator button on login items, and a page-level `navigator.credentials.create` / `navigator.credentials.get` passkey bridge in the Chrome/Edge/Brave/Vivaldi extension.
* **Item History & Trash Safety Net:**
  * Deleting an item now soft-deletes it to **Trash** first, with 1-click **Restore** or **Empty Trash** / **Delete Forever**.
  * Every item detail card includes a **History** button (`25` snapshots per item) showing password, username, URL, folder/tag changes, and line-by-line note diffs with 1-click version rollback or fork.
* **Custom Tags, Folders & Pinned Favorites:**
  * Assign items to a **Folder** (`Work`, `Personal`, `Finance`, `DevOps`) and multiple comma-separated **#tags** with quick-pick chips.
  * Dynamic **Folders & Tags** section in the sidebar, **By Tag / Folder** option in the Group By menu, and a **Pinned Favorites** section at the top of the item list.

</details>

<details>
<summary>Previous release: v1.3.1</summary>

* **Instant Live List Refresh (No More F5):**
  * New logins, SSH keys, notes, and cards appear in the list immediately, including items saved from the browser extension.
  * A lightweight background revision poll keeps the list in sync on desktop and Android, where there is no F5.
  * After saving, the list switches to **All** and scrolls the new item into view.
* **At-a-Glance Item Identity:**
  * Each type has its own avatar: terminal for SSH, document for Notes, card for Cards, and key or website favicon for Logins.
  * Recently created items carry a glowing **✦ NEW** badge; items saved via the extension show an **Ext** badge.
* **Alphabetized Lists & Group By:**
  * Items sort naturally A to Z.
  * New **Group By** menu: Category, Source (Extension or App), Date Added, and First Letter, with sticky group headers.
* **Command Palette (Ctrl+K):** Improved search with match highlighting, category filters, quick actions, and full keyboard navigation (Up/Down/Enter/Esc).
* **Decluttered Toolbar & Rows:**
  * Compact icon buttons for Drive, Sync, Group By, and New Item (`+`); the item count stays on one line (`42 items`).
  * Removed repetitive `[LOGIN]`/`[NOTE]` chips from row titles; the avatar conveys the type.
* **Android Fix:** Group By dropdown and command palette now use solid opaque surfaces, so they are readable on Android WebView.

</details>

<details>
<summary>Previous release: v1.2.9</summary>

* **Appearance & Theme System (Midnight OLED, Deep Slate & Clean Light):**
  * Added **Midnight OLED** mode (`#000000` pure black canvas) for battery-saving high contrast on OLED phones, tablets, and monitors.
  * Added **Deep Slate** (refined neutral dark theme) and **Clean Light** mode.
  * Added 6 vivid accent palettes: **Electric Blue** (Default), **Emerald Crypt**, **Amethyst Vault**, **Cyber Amber**, **Crimson Sentinel**, and **Titanium Slate**.
  * Instant pre-render hydration eliminates flash of unstyled theme on boot.
* **Clean Light Mode Overhaul & High-Contrast Typography:**
  * **Zero Black Card Artifacts:** Converted all credential panels, security health cards, modal dialogs, and Watchtower sections to true daylight `#FFFFFF` surfaces with `#CBD5E1` borders and clean elevation.
  * **Form Helper Micro-Actions & Button Normalization:** Added refined micro-action styling for `.form-helper-link` (`Generate Ed25519 Key` and `Generate Password`) with smooth border radii, ambient hover elevation, and daylight high contrast, alongside a global `<button>` CSS reset that completely eliminates raw unstyled 90s browser bevels and harsh black borders across all dialogs.
  * **Secondary Button & Hover Contrast:** Fixed `.action-btn-secondary` and `.tool-icon-btn` text and SVG strokes to remain crisp and dark (`#0F172A`) both at rest and during `:hover`, completely resolving the white-on-white vanishing bug.
  * **Primary Action Button Contrast:** Enhanced high-contrast typography and SVG icon styling on primary action buttons (e.g. `Copy Password`, `Save Item`) across all accent color variations.
  * **Mobile & Tablet Detail View Contrast:** Fixed `.mobile-detail-back-btn` (`← Items`) and `.mobile-hero-title` to enforce crisp `#0F172A` text, borders, and vector icons in Clean Light Mode on tablets and mobile screens.
  * **Markdown Studio Toolbar:** Clean daylight theme for the secure notes markdown editor toolbar with high-contrast tool buttons, active states, and crisp borders.
  * **Clean Light Unlock & Login Screen:** Theme preference automatically hydrates into the lock and login screen, providing a consistent daylight aesthetic without dark card anomalies.
  * **Clean Favicon & Avatar Containers:** Eliminated hardcoded `#181E2B` black avatar backgrounds; all website favicons render inside crisp `#FFFFFF` boxes with `#CBD5E1` borders.
  * **Watchtower Security Audit Legibility:** Color-coded alert headers, recommendation pills, and count badges feature WCAG-compliant tones (`#DC2626`, `#C2410C`, `#B45309`) with high contrast against light surfaces.
  * **Notes & Markdown Readability:** Formatted markdown bodies and raw editor text to `#0F172A` in light mode.
* **Watchtower Security Audit Layout Polish:**
  * Redesigned score meters, recommendation cards, and credential pill envelopes with generous padding, margins, and crisp tabular numbers.
  * Optimized responsive density across both desktop monitors and mobile/tablet touchscreens.
* **Complete Emoji Sweep & Precision Vector Iconography:**
  * Purged all remaining emoji symbols and unicode glyphs from buttons, sidebars, relay monitors, and alerts.
  * Replaced with high-precision Lucide/Feather vector SVGs adhering to strict UI Craft design standards.
* **Native Browser Launching & Navigation Fix:**
  * Resolved internal WebView relative URL navigation errors (such as `127.0.0.1:51892/gmail.com 404`).
  * Non-protocol URLs automatically normalize to `https://` and route directly to the user's default OS browser.
* **Mobile & Tablet Drawer Spacing:**
  * Fixed close button alignment, margin, and touch target on mobile and tablet screens.

</details>

---

## Core Capabilities

* **Zero-Knowledge Multi-Device Sync:**
  * **Built-in Native Relay Server (Default & Zero Setup):** Boots automatically on port `8443` whenever you launch ZeroEnclave on desktop. Ultra-lightweight (~5MB RAM), zero dependencies, no Docker required, and instant local LAN pairing. All relay data is stored securely in your vault's `relay_data/` directory.
  * **Live Sync Progress & Percentage Tracking:** Displays a real-time progress bar (`0% – 100%`) with dynamic stage diagnostics (`Connecting`, `Merging remote delta X of Y`, `Uploading local changes`, `Reconciling attachments and media`).
  * **Zero-Knowledge Docker Relay Server (Optional 24/7 Daemon):** Self-hosted blind relay server container (`cmd/relay`) for users who want 24/7 background syncing even when the desktop app window is closed. The server only ever handles blind Base64 ciphertext deltas and encrypted `1VBLOB` chunks sealed under your derived bucket key (`SHA-256("zeroenclave:sync:bucket:" + SecretKey)`). The server has zero access to passwords, filenames, metadata, or encryption keys.
  * **In-App Relay Assistant:** View live relay connectivity status, toggle between the built-in relay and Docker container, and inspect LAN IP addresses for phone/tablet pairing directly in the **Cloud & Sync** dialog.
  * **Local Wi-Fi Mesh Discovery Fast-Path:** When two or more devices (PC, laptop, smartphone) share the same local Wi-Fi or LAN, background UDP beacons on port `8444` discover peer nodes within milliseconds and sync peer-to-peer at gigabit LAN speed, automatically bypassing cloud/internet relays.
  * **Instant Camera QR Code Pairing:** Generate an offline standard ISO/IEC 18004 QR code on your desktop screen and scan it in seconds with your phone's camera. Secret keys and relay configs synchronize instantly with zero typing, and newly paired secondary devices automatically initialize and sync on first unlock.
  * **Git-Style Revision History & Zero-Data-Loss Offline Conflict Resolution:** Every edit generates a deterministic SHA-256 content commit (`RevisionID`) linked to its parent (`ParentRevID`) with up to 25 encrypted snapshots per item. Clean sequential edits fast-forward automatically; concurrent offline edits on multiple devices preserve both versions in the item's encrypted history (`reason: "sync-conflict"`), alert you with a `⚠️ Conflict` badge, and let you inspect a line-by-line red/green diff, restore any snapshot, or fork both versions into separate items.
* **Markdown Studio & Interactive Secure Notes:**
  * **Write / Split / Preview Modes:** Author rich encrypted documentation with live side-by-side preview.
  * **Full Formatting Toolbar:** One-click controls for Headings (`H1`–`H3`), Bold/Italic/Strikethrough/Highlight (`==text==`), fenced code blocks with copy buttons, click-to-reveal redacted spoilers (`||secret||`), interactive task checklists (`- [ ]` / `- [x]`), GitHub-style Callouts (`[!NOTE]`, `[!TIP]`, `[!WARNING]`, `[!CAUTION]`), encrypted item wiki-links (`[[Item Title]]`), and an automatic **Table Structuring Builder** (`2x2`, `3x3`, `4x4`, `Key/Value`, `+ Add Row`, `+ Add Column`, and `Auto-Align Columns`).
* **Developer & DevOps Power Tools:**
  * **Built-in Ed25519 SSH Keypair & Fingerprint Generator:** Generate cryptographic `Ed25519` OpenSSH private/public keypairs and `SHA256:...` fingerprints inside the item editor with one click.
  * **Runtime `.env` Secret Injector CLI (`ZeroEnclave run`):** Inject secrets and `KEY=VALUE` variables from any vault item directly into a child process's memory (`ZeroEnclave.exe run --item "Prod API" -- npm start`) without ever writing `.env` files to disk.
  * **Self-Destructing 1-Time E2EE Drop Links:** Share credentials with a burn-after-reading URL (`/drop/<id>#<key>`) where the AES-256-GCM decryption key stays strictly in the URL `#fragment` (never sent to the relay) and the ciphertext is permanently shredded after the first view.
* **Active Defense & Anti-Coercion Security:**
  * **Anti-Coercion Duress / Decoy Vault:** Configure a secondary Duress Password that unlocks an isolated decoy vault with harmless dummy credentials while keeping your real vault completely hidden.
  * **2-of-3 Shamir Secret Sharing (`GF(256)`):** Split your Emergency Kit Secret Key into 3 cryptographic shares where any 2 reconstruct your key and 1 share alone reveals zero information.
  * **Hardware Screen-Capture Shield (`SetWindowDisplayAffinity`):** One-click toggle in the Emergency Kit to block screenshots, OBS, and remote desktop tools from capturing the ZeroEnclave window.
* **Mountable Virtual Encrypted Drive (`Z:\`):**
  * Mount your encrypted vault as a standard Windows Virtual Drive (`Z:\`) without third-party kernel drivers, Dokany, WinFsp, or Administrator privileges.
  * Open, edit, and save vault files directly in Windows Explorer, Notepad, VS Code, or Photoshop.
  * **Anti-Forensic Panic Shred:** One-click unmount securely wipes drive mappings, scrubs temporary files with random byte overwrites, and terminates open handles.
* **Passkeys & Windows Hello Biometrics:**
  * Passwordless login via WebAuthn hardware security keys (YubiKey), Windows Hello PIN/fingerprint/facial recognition, and Android BiometricPrompt.
* **Native Desktop Window:** Opens as a true native Windows desktop application (Edge WebView2 + Win32 dark frame) — **not** inside a web browser tab.
* **No Console Window:** Compiled with `-H windowsgui` so launching `ZeroEnclave.exe` / `1vault.exe` runs silently without popping up a terminal window.
* **Streaming Encrypted Blobs (`1VBLOB` Engine):**
  * Vault **photos, 4K videos, mobile/desktop apps (`.apk`, `.exe`), archives, and documents** of any size without increasing vault unlock time or exhausting device RAM.
  * Every file is encrypted with a unique 256-bit **File Encryption Key (FEK)** (envelope-encrypted inside your vault index) and streamed in **64 KB AES-256-GCM chunks** with positional AAD protection.
  * **$\mathcal{O}(1)$ Video Scrubbing:** Supports HTTP `Range` requests (`206 Partial Content`) by seeking directly to any 64 KB encrypted chunk on disk (`19 + i × 65552`), allowing instant playback and scrubbing of multi-gigabyte videos inside the app.
* **Modern UI Craft Design:**
  * **Brand Favicon Resolution:** Automatically displays website favicons for accounts (Google, GitHub, Snaphunt, etc.) with gradient monogram fallbacks.
  * **2-Column Desktop Grid:** 65% credential & media management + 35% live Security Intel sidebar (entropy calculation, breach checks, stream blob metrics, sync status).
  * **Quick Actions:** One-click copy buttons (`Copy Password`, `Copy Username`, `Open Website`, `+ Vault File / Media / App`) with instant visual checkmark feedback.
  * **Fluid Mobile Experience:** Automatic smartphone detection with off-canvas category drawers and slide-over details.

---

## Security Architecture

* **Key Derivation (KDF):** **Argon2id** (memory: 64 MB, time: 3 iterations, 4 parallelism threads) to defeat GPU/ASIC brute-force attacks.
* **Vault Index Encryption (AEAD):** **AES-256-GCM** authenticated encryption with random 96-bit nonces per save.
* **Streaming Media Blob Encryption (`1VBLOB\x01`):**
  * **Per-File Envelope Encryption:** Random 256-bit FEK generated via `crypto/rand` per attachment.
  * **Chunked AEAD:** 64 KB ($65,536\text{ B}$) chunks sealed with AES-256-GCM (`12B` nonce = `8B prefix || 4B chunk index`; `9B` AAD = `8B chunk index || 1B final-chunk flag`) to prevent chunk reordering or truncation attacks.
* **Dual-Layer Protection (Secret Key):** Combines Master Password with a generated 128-bit Secret Key (`A3-XXXXXX-XXXXXX-...`).
* **Zero-Knowledge Sync:** Relay servers only process encrypted blobs. Delta payloads are sealed with AES-256-GCM using keys derived from your vault's Secret Key.
* **In-Memory Hygiene:** Plaintext secrets, FEKs, and derived keys are zeroed in memory upon locking or program exit.
* **Breach Protection:** Integrated k-anonymity checks against the HaveIBeenPwned API (only sends the first 5 SHA-1 hash characters, never the password).
* **Watchtower Security Audit:** Live auditing of weak passwords, reused credentials, identical username/password pairs, and redundant accounts.

---

## Device Linking & Pairing

To synchronize your Desktop, Smartphone, and Tablet into one shared vault:

1. **On your primary device (e.g. Desktop):**
   - Unlock your vault and open **Cloud & Sync**.
   - Your unique pairing **QR Code** will be displayed.
2. **On your secondary device (Smartphone or Tablet):**
   - On the welcome / unlock screen, tap **"📷 Scan QR to pair this device"**.
   - Point your camera at the desktop screen.
   - The Secret Key and sync settings will populate automatically.
   - Enter your Master Password and tap **Unlock**.
3. **Done:** All devices now share the same encrypted vault. Any change made on one device automatically syncs to the others.

---

## Disaster Recovery (If Your Device is Lost or Destroyed)

Because ZeroEnclave is zero-knowledge and end-to-end encrypted, you are never locked out as long as you keep your credentials:

1. **Install ZeroEnclave on your new computer.**
2. **Enter your Master Password and 128-bit Secret Key** (from your printed Emergency Kit).
3. **Connect to your Sync source:**
   - **Via Relay Server:** Enter your Relay Server URL under **Cloud & Sync** and click **Sync Now**.
   - **Via Paired Phone:** On the same local Wi-Fi, the phone and new computer will discover each other automatically and restore all passwords and media peer-to-peer.
4. Your complete vault index and all streaming media attachments will be **100% restored in seconds**.

---

## Sync Relay Options: Built-in vs. Docker

ZeroEnclave gives you two flexible options for sync relaying:

### Option 1: Built-in Native Relay Server (Default & Zero Setup)
* **Zero Configuration:** Automatically launches directly inside `ZeroEnclave.exe` on port `8443` when you open the desktop app.
* **Ultra Lightweight:** Consumes only ~5 MB of RAM and requires zero third-party software, dependencies, or Docker installations.
* **Seamless Local Sync:** Your phone and other devices can immediately connect over your local Wi-Fi or LAN via the auto-generated pairing QR code.
* **Lifecycle:** Runs as long as the desktop application window is open.

### Option 2: Docker Container Relay (Optional 24/7 Daemon)
* **24/7 Always-On:** If you want your phone, tablet, and browser extension to sync in the background even when your desktop app window is closed, you can run the relay as a background Docker container (`zeroenclave-sync-relay`).
* **Easy In-App Launch:** Open **Cloud & Sync**, expand **⚡ Built-in & Docker Relays**, and click **"▶ Start Local Docker Relay"**.
* **Or via Docker Compose:**
  ```bash
  git clone https://github.com/ewceniza9009/1Vault.git
  cd 1Vault
  docker compose up -d
  ```

---

## Downloads & Installation (v1.4.3)

Pre-built binaries and packages are available under [Releases](https://github.com/ewceniza9009/1Vault/releases):

### 1. Windows Installer (Recommended)
Download [`ZeroEnclave-Setup-v1.4.3.exe`](https://github.com/ewceniza9009/1Vault/releases/latest):
* **One-Click Install:** Installs without requiring Administrator privileges and preserves existing vault data seamlessly.
* **Shortcuts:** Automatically creates Start Menu and Desktop shortcuts (`ZeroEnclave`) with the hexagonal ZeroEnclave emblem.
* **Clean Uninstall:** Registered in Windows **Settings > Apps & Features**.

### 2. Windows Portable Executable
Download [`ZeroEnclave.exe`](https://github.com/ewceniza9009/1Vault/releases/latest) to run directly without installing.

### 3. Android Application (APK)
Download [`ZeroEnclave-v1.4.3.apk`](https://github.com/ewceniza9009/1Vault/releases/latest):
* Production signed Release APK built with full optimization.
* Powered by the embedded Go cryptographic & `1VBLOB` streaming core compiled via Android NDK CGO.
* Supports Android 8.0+ (`arm64-v8a` and `x86_64`).
* Camera QR code pairing, local Wi-Fi auto-discovery, and biometric unlock integration.

---

## Building from Source

### Prerequisites
* Go 1.22+
* Windows: MinGW (`gcc`, `windres`) and [NSIS](https://nsis.sourceforge.io/) (for creating the installer).
* Android: Android SDK, NDK (r26+), and JDK 17+.

### Building Windows Desktop & Installer
```powershell
# 1. Run unit tests
go test -v ./...

# 2. Build desktop executable with embedded ZeroEnclave icon and metadata
windres -O coff -o 1vault_windows_amd64.syso 1vault.rc
go build -ldflags="-H windowsgui -s -w" -o ZeroEnclave.exe .

# 3. Compile Windows Installer with NSIS
& "C:\Program Files (x86)\NSIS\makensis.exe" installer.nsi
```

### Building Android APK
```powershell
# 1. Compile native Go shared libraries using Android NDK Clang
$ndkBin = "$env:LOCALAPPDATA\Android\Sdk\ndk\27.0.12077973\toolchains\llvm\prebuilt\windows-x86_64\bin"

# ARM64 (arm64-v8a)
$env:CC = "$ndkBin\aarch64-linux-android28-clang.cmd"
$env:GOOS = "android"; $env:GOARCH = "arm64"; $env:CGO_ENABLED = "1"
go build -buildmode=c-shared -ldflags="-s -w" -o "android/app/src/main/jniLibs/arm64-v8a/libonevault.so" ./cmd/mobile
Copy-Item "android/app/src/main/jniLibs/arm64-v8a/libonevault.so" "android/app/src/main/jniLibs/arm64-v8a/lib1vault.so" -Force

# x86_64
$env:CC = "$ndkBin\x86_64-linux-android28-clang.cmd"
$env:GOOS = "android"; $env:GOARCH = "amd64"; $env:CGO_ENABLED = "1"
go build -buildmode=c-shared -ldflags="-s -w" -o "android/app/src/main/jniLibs/x86_64/libonevault.so" ./cmd/mobile
Copy-Item "android/app/src/main/jniLibs/x86_64/libonevault.so" "android/app/src/main/jniLibs/x86_64/lib1vault.so" -Force

# 2. Build Release APK with Gradle
cmd /c "cd android && gradlew.bat assembleRelease"
# Output artifact: android/app/build/outputs/apk/release/ZeroEnclave-v1.4.3.apk
```
