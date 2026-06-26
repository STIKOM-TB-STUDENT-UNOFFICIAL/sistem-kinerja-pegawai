#!/bin/bash
set -e

# ============================================================
#  Sistem Kinerja Pegawai — Updater untuk Linux
#  Cukup jalankan:  sudo ./update.sh
# ============================================================

# ---------- Warna untuk output ----------
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# ---------- Variabel ----------
APP_NAME="kinerja-pegawai"
INSTALL_DIR="/opt/${APP_NAME}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# ---------- Fungsi bantuan ----------
info()    { echo -e "${GREEN}[INFO]${NC}    $1"; }
warn()    { echo -e "${YELLOW}[WARN]${NC}    $1"; }
error()   { echo -e "${RED}[ERROR]${NC}   $1"; }
step()    { echo -e "\n${CYAN}▶ $1${NC}"; }

check_root() {
    if [ "$EUID" -ne 0 ]; then
        error "Script ini harus dijalankan sebagai root."
        echo -e "  Gunakan: ${YELLOW}sudo ./update.sh${NC}"
        exit 1
    fi
}

check_root

CHANGED=false

step "Memeriksa perubahan aplikasi..."

# 1. Update main
if [ -f "${SCRIPT_DIR}/main" ]; then
    if [ ! -f "${INSTALL_DIR}/main" ] || ! cmp -s "${SCRIPT_DIR}/main" "${INSTALL_DIR}/main"; then
        info "Mendeteksi perubahan pada 'main'. Memperbarui..."
        cp "${SCRIPT_DIR}/main" "${INSTALL_DIR}/main"
        chmod +x "${INSTALL_DIR}/main"
        CHANGED=true
    else
        info "'main' tidak berubah."
    fi
else
    warn "'main' tidak ditemukan di ${SCRIPT_DIR}, dilewati."
fi

# 2. Update .env
if [ -f "${SCRIPT_DIR}/.env" ]; then
    if [ ! -f "${INSTALL_DIR}/.env" ] || ! cmp -s "${SCRIPT_DIR}/.env" "${INSTALL_DIR}/.env"; then
        info "Mendeteksi perubahan pada '.env'. Memperbarui..."
        cp "${SCRIPT_DIR}/.env" "${INSTALL_DIR}/.env"
        chmod 600 "${INSTALL_DIR}/.env"
        CHANGED=true
    else
        info "'.env' tidak berubah."
    fi
else
    warn "'.env' tidak ditemukan di ${SCRIPT_DIR}, dilewati."
fi

# 3. Update public
if [ -d "${SCRIPT_DIR}/public" ]; then
    if [ ! -d "${INSTALL_DIR}/public" ] || ! diff -r "${SCRIPT_DIR}/public" "${INSTALL_DIR}/public" >/dev/null 2>&1; then
        info "Mendeteksi perubahan pada folder 'public'. Memperbarui..."
        if command -v rsync >/dev/null 2>&1; then
            rsync -ac --delete "${SCRIPT_DIR}/public/" "${INSTALL_DIR}/public/"
        else
            rm -rf "${INSTALL_DIR}/public"
            cp -r "${SCRIPT_DIR}/public" "${INSTALL_DIR}/public"
        fi
        CHANGED=true
    else
        info "Folder 'public' tidak berubah."
    fi
else
    warn "Folder 'public' tidak ditemukan di ${SCRIPT_DIR}, dilewati."
fi

# 4. Update views
if [ -d "${SCRIPT_DIR}/views" ]; then
    if [ ! -d "${INSTALL_DIR}/views" ] || ! diff -r "${SCRIPT_DIR}/views" "${INSTALL_DIR}/views" >/dev/null 2>&1; then
        info "Mendeteksi perubahan pada folder 'views'. Memperbarui..."
        if command -v rsync >/dev/null 2>&1; then
            rsync -ac --delete "${SCRIPT_DIR}/views/" "${INSTALL_DIR}/views/"
        else
            rm -rf "${INSTALL_DIR}/views"
            cp -r "${SCRIPT_DIR}/views" "${INSTALL_DIR}/views"
        fi
        CHANGED=true
    else
        info "Folder 'views' tidak berubah."
    fi
else
    warn "Folder 'views' tidak ditemukan di ${SCRIPT_DIR}, dilewati."
fi

# ---------- Set kepemilikan jika ada perubahan ----------
if [ "$CHANGED" = true ]; then
    step "Mengatur kepemilikan file..."
    chown -R www-data:www-data "${INSTALL_DIR}" 2>/dev/null || \
    chown -R apache:apache "${INSTALL_DIR}" 2>/dev/null || \
    chown -R nobody:nobody "${INSTALL_DIR}" 2>/dev/null || true
    
    step "Merestart systemd service..."
    if systemctl is-active --quiet "${APP_NAME}.service"; then
        systemctl restart "${APP_NAME}.service"
        info "Service '${APP_NAME}' berhasil direstart."
    else
        systemctl start "${APP_NAME}.service"
        info "Service '${APP_NAME}' berhasil dijalankan."
    fi
else
    step "Tidak ada perubahan terdeteksi. Service tidak perlu direstart."
fi

echo -e "\n${GREEN}✅ Selesai!${NC}\n"
