#!/bin/bash
set -e

# ============================================================
#  Sistem Kinerja Pegawai — Installer untuk Linux
#  Cukup jalankan:  sudo ./install.sh
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
SERVICE_FILE="/etc/systemd/system/${APP_NAME}.service"
APACHE_CONF="/etc/apache2/sites-available/${APP_NAME}.conf"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

DB_NAME="stikomtb_kinerja_pegawai"
SQL_FILE="${SCRIPT_DIR}/stikomtb_kinerja_pegawai.sql"

# ---------- Fungsi bantuan ----------
print_banner() {
    echo ""
    echo -e "${CYAN}╔══════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║   ${GREEN}Installer — Sistem Kinerja Pegawai${CYAN}             ║${NC}"
    echo -e "${CYAN}║   ${NC}STIKOM Tunas Bangsa${CYAN}                             ║${NC}"
    echo -e "${CYAN}╚══════════════════════════════════════════════════╝${NC}"
    echo ""
}

info()    { echo -e "${GREEN}[INFO]${NC}    $1"; }
warn()    { echo -e "${YELLOW}[WARN]${NC}    $1"; }
error()   { echo -e "${RED}[ERROR]${NC}   $1"; }
step()    { echo -e "\n${CYAN}▶ $1${NC}"; }

check_root() {
    if [ "$EUID" -ne 0 ]; then
        error "Script ini harus dijalankan sebagai root."
        echo -e "  Gunakan: ${YELLOW}sudo ./install.sh${NC}"
        exit 1
    fi
}

# ---------- Deteksi package manager ----------
detect_pkg_manager() {
    if command -v apt-get &>/dev/null; then
        PKG_MANAGER="apt"
    elif command -v dnf &>/dev/null; then
        PKG_MANAGER="dnf"
    elif command -v yum &>/dev/null; then
        PKG_MANAGER="yum"
    else
        error "Package manager tidak ditemukan (apt/dnf/yum)."
        exit 1
    fi
    info "Package manager terdeteksi: ${PKG_MANAGER}"
}

# ---------- Install dependensi ----------
install_dependencies() {
    step "Menginstall dependensi sistem..."

    case $PKG_MANAGER in
        apt)
            apt-get update -qq
            apt-get install -y -qq mariadb-server apache2 >/dev/null 2>&1
            ;;
        dnf)
            dnf install -y -q mariadb-server httpd >/dev/null 2>&1
            ;;
        yum)
            yum install -y -q mariadb-server httpd >/dev/null 2>&1
            ;;
    esac

    info "Dependensi berhasil diinstall."
}

# ---------- Setup MariaDB ----------
setup_database() {
    step "Mengkonfigurasi database MariaDB..."

    # Pastikan MariaDB berjalan
    systemctl start mariadb 2>/dev/null || systemctl start mysql 2>/dev/null || true
    systemctl enable mariadb 2>/dev/null || systemctl enable mysql 2>/dev/null || true

    # Baca konfigurasi DB dari .env
    if [ -f "${SCRIPT_DIR}/.env" ]; then
        source <(grep -E '^DB_' "${SCRIPT_DIR}/.env" | sed 's/\r//g')
    fi

    DB_USER="${DB_USER:-root}"
    DB_PASS="${DB_PASS:-}"
    DB_HOST="${DB_HOST:-localhost}"
    DB_PORT="${DB_PORT:-3306}"

    # Buat database jika belum ada
    if [ -n "$DB_PASS" ]; then
        MYSQL_CMD="mysql -u${DB_USER} -p${DB_PASS} -h${DB_HOST} -P${DB_PORT}"
    else
        MYSQL_CMD="mysql -u${DB_USER} -h${DB_HOST} -P${DB_PORT}"
    fi

    info "Membuat database '${DB_NAME}'..."
    $MYSQL_CMD -e "CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;" 2>/dev/null

    # Import SQL jika file ada
    if [ -f "$SQL_FILE" ]; then
        info "Mengimport data dari ${SQL_FILE}..."
        $MYSQL_CMD "$DB_NAME" < "$SQL_FILE" 2>/dev/null
        info "Database berhasil diimport."
    else
        warn "File SQL tidak ditemukan: ${SQL_FILE}"
        warn "Anda perlu mengimport database secara manual."
    fi
}

# ---------- Install aplikasi ----------
install_application() {
    step "Menginstall aplikasi ke ${INSTALL_DIR}..."

    # Buat direktori instalasi
    mkdir -p "${INSTALL_DIR}"

    # Copy binary utama
    if [ -f "${SCRIPT_DIR}/main" ]; then
        cp "${SCRIPT_DIR}/main" "${INSTALL_DIR}/main"
        chmod +x "${INSTALL_DIR}/main"
        info "Binary aplikasi berhasil disalin."
    else
        error "Binary 'main' tidak ditemukan di ${SCRIPT_DIR}!"
        exit 1
    fi

    # Copy file .env
    if [ -f "${SCRIPT_DIR}/.env" ]; then
        cp "${SCRIPT_DIR}/.env" "${INSTALL_DIR}/.env"
        chmod 600 "${INSTALL_DIR}/.env"
        info "File konfigurasi .env berhasil disalin."
    fi

    # Copy folder public (assets, uploads)
    if [ -d "${SCRIPT_DIR}/public" ]; then
        cp -r "${SCRIPT_DIR}/public" "${INSTALL_DIR}/public"
        info "Folder public berhasil disalin."
    fi

    # Copy folder views (templates)
    if [ -d "${SCRIPT_DIR}/views" ]; then
        cp -r "${SCRIPT_DIR}/views" "${INSTALL_DIR}/views"
        info "Folder views berhasil disalin."
    fi

    # Set kepemilikan
    chown -R www-data:www-data "${INSTALL_DIR}" 2>/dev/null || \
    chown -R apache:apache "${INSTALL_DIR}" 2>/dev/null || \
    chown -R nobody:nobody "${INSTALL_DIR}" 2>/dev/null || true

    info "Aplikasi berhasil diinstall di ${INSTALL_DIR}"
}

# ---------- Setup systemd service ----------
setup_service() {
    step "Mengkonfigurasi systemd service..."

    if [ -f "${SCRIPT_DIR}/kinerja-pegawai.service" ]; then
        cp "${SCRIPT_DIR}/kinerja-pegawai.service" "${SERVICE_FILE}"
    else
        # Buat file service secara otomatis
        cat > "${SERVICE_FILE}" <<EOF
[Unit]
Description=Sistem Kinerja Pegawai
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data

WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/main

Restart=always
RestartSec=5

Environment=GIN_MODE=release
EnvironmentFile=${INSTALL_DIR}/.env

[Install]
WantedBy=multi-user.target
EOF
    fi

    systemctl daemon-reload
    systemctl enable "${APP_NAME}.service"
    systemctl start "${APP_NAME}.service"

    info "Service '${APP_NAME}' berhasil diaktifkan dan dijalankan."
}

# ---------- Setup Apache reverse proxy ----------
setup_apache() {
    step "Mengkonfigurasi Apache reverse proxy..."

    # Cek apakah Apache (apache2/httpd) tersedia
    if command -v a2enmod &>/dev/null; then
        # Debian/Ubuntu — apache2
        a2enmod proxy proxy_http >/dev/null 2>&1

        if [ -f "${SCRIPT_DIR}/apache-kinerja-pegawai.conf" ]; then
            cp "${SCRIPT_DIR}/apache-kinerja-pegawai.conf" "${APACHE_CONF}"
        else
            cat > "${APACHE_CONF}" <<EOF
<VirtualHost *:80>
    ServerName skp

    ProxyPreserveHost On

    ProxyPass / http://127.0.0.1:3000/
    ProxyPassReverse / http://127.0.0.1:3000/

    ErrorLog \${APACHE_LOG_DIR}/kinerja-pegawai-error.log
    CustomLog \${APACHE_LOG_DIR}/kinerja-pegawai-access.log combined
</VirtualHost>
EOF
        fi

        a2ensite "${APP_NAME}.conf" >/dev/null 2>&1
        systemctl reload apache2

        info "Apache virtual host berhasil dikonfigurasi."
    elif command -v httpd &>/dev/null; then
        # CentOS/RHEL — httpd
        HTTPD_CONF_DIR="/etc/httpd/conf.d"
        mkdir -p "$HTTPD_CONF_DIR"

        if [ -f "${SCRIPT_DIR}/apache-kinerja-pegawai.conf" ]; then
            cp "${SCRIPT_DIR}/apache-kinerja-pegawai.conf" "${HTTPD_CONF_DIR}/${APP_NAME}.conf"
        fi

        systemctl enable httpd >/dev/null 2>&1
        systemctl restart httpd

        info "HTTPD virtual host berhasil dikonfigurasi."
    else
        warn "Apache/HTTPD tidak terdeteksi. Lewati konfigurasi reverse proxy."
    fi
}

# ---------- Ringkasan ----------
print_summary() {
    # Ambil port dari .env
    APP_PORT="3000"
    if [ -f "${INSTALL_DIR}/.env" ]; then
        PORT_VAL=$(grep -E '^PORT=' "${INSTALL_DIR}/.env" | cut -d'=' -f2 | tr -d '"' | tr -d '\r')
        [ -n "$PORT_VAL" ] && APP_PORT="$PORT_VAL"
    fi

    echo ""
    echo -e "${GREEN}╔══════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║          ✅  INSTALASI BERHASIL!                 ║${NC}"
    echo -e "${GREEN}╚══════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "  📁  Direktori instalasi : ${CYAN}${INSTALL_DIR}${NC}"
    echo -e "  🌐  Aplikasi berjalan   : ${CYAN}http://localhost:${APP_PORT}${NC}"
    echo -e "  🗄️  Database            : ${CYAN}${DB_NAME}${NC}"
    echo -e "  ⚙️  Service             : ${CYAN}${APP_NAME}.service${NC}"
    echo ""
    echo -e "  ${YELLOW}Perintah berguna:${NC}"
    echo -e "    Cek status  : ${CYAN}sudo systemctl status ${APP_NAME}${NC}"
    echo -e "    Restart     : ${CYAN}sudo systemctl restart ${APP_NAME}${NC}"
    echo -e "    Lihat log   : ${CYAN}sudo journalctl -u ${APP_NAME} -f${NC}"
    echo -e "    Stop        : ${CYAN}sudo systemctl stop ${APP_NAME}${NC}"
    echo ""
}

# ============================================================
#  MAIN
# ============================================================
print_banner
check_root
detect_pkg_manager

echo ""
echo -e "${YELLOW}Installer akan melakukan:${NC}"
echo -e "  1. Install MariaDB & Apache2"
echo -e "  2. Buat & import database"
echo -e "  3. Install aplikasi ke /opt/${APP_NAME}"
echo -e "  4. Setup systemd service"
echo -e "  5. Konfigurasi Apache reverse proxy"
echo ""
read -rp "Lanjutkan instalasi? [Y/n] " confirm
confirm="${confirm:-Y}"

if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    info "Instalasi dibatalkan."
    exit 0
fi

install_dependencies
setup_database
install_application
setup_service
setup_apache
print_summary
