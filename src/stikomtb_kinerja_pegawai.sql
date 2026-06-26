-- phpMyAdmin SQL Dump
-- version 4.7.4
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: 23 Jun 2026 pada 19.20
-- Versi Server: 10.1.29-MariaDB
-- PHP Version: 7.2.0

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `stikomtb_kinerja_pegawai`
--

-- --------------------------------------------------------

--
-- Struktur dari tabel `aktivitas`
--

CREATE TABLE `aktivitas` (
  `id` int(11) NOT NULL,
  `userid` varchar(50) NOT NULL,
  `id_tupoksi` int(11) NOT NULL,
  `tanggal` date NOT NULL,
  `catatan` text NOT NULL,
  `mulai` datetime NOT NULL,
  `selesai` datetime DEFAULT NULL,
  `status` enum('Diproses','Diterima','Ditolak','Dalam Pengerjaan') NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data untuk tabel `aktivitas`
--

INSERT INTO `aktivitas` (`id`, `userid`, `id_tupoksi`, `tanggal`, `catatan`, `mulai`, `selesai`, `status`) VALUES
(22, '0118018901', 5, '2026-06-15', 'sfdzsfsaf', '2026-06-15 11:43:15', '2026-06-15 11:43:17', 'Diproses'),
(23, '0118018901', 5, '2026-06-15', 'sadsad', '2026-06-15 11:46:26', '2026-06-19 10:37:26', 'Diproses'),
(37, '0118088603', 1, '2026-06-01', '', '2026-06-01 19:25:00', '2026-06-02 19:25:00', 'Diproses');

-- --------------------------------------------------------

--
-- Struktur dari tabel `anggota`
--

CREATE TABLE `anggota` (
  `id` int(11) NOT NULL,
  `userid_parent` varchar(50) NOT NULL,
  `userid_children` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Struktur dari tabel `detail_pegawai`
--

CREATE TABLE `detail_pegawai` (
  `id` int(11) NOT NULL,
  `userid` varchar(50) NOT NULL,
  `nip` varchar(50) NOT NULL,
  `jabatan` varchar(60) NOT NULL,
  `departemen` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data untuk tabel `detail_pegawai`
--

INSERT INTO `detail_pegawai` (`id`, `userid`, `nip`, `jabatan`, `departemen`) VALUES
(1, '0118088603', '08', 'Ketua', 'Program Studi Sistem Informasi'),
(4, '0103038202', '0103038202', 'Anggota', 'Pendidikan');

-- --------------------------------------------------------

--
-- Struktur dari tabel `file_bukti`
--

CREATE TABLE `file_bukti` (
  `id` int(11) NOT NULL,
  `aktivitas_id` int(11) NOT NULL,
  `nama_file` varchar(64) NOT NULL,
  `lokasi_file` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Struktur dari tabel `komentar_aktivitas`
--

CREATE TABLE `komentar_aktivitas` (
  `id` int(11) NOT NULL,
  `id_aktivitas` int(11) NOT NULL,
  `userid` varchar(50) NOT NULL,
  `komentar` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data untuk tabel `komentar_aktivitas`
--

INSERT INTO `komentar_aktivitas` (`id`, `id_aktivitas`, `userid`, `komentar`) VALUES
(6, 22, '0118088603', ''),
(7, 37, '0118088603', 'Test'),
(9, 23, '0118018901', 'Test'),
(12, 23, '0118088603', 'Test');

-- --------------------------------------------------------

--
-- Struktur dari tabel `login_system`
--

CREATE TABLE `login_system` (
  `userid` varchar(50) NOT NULL,
  `password` varchar(255) CHARACTER SET latin1 NOT NULL,
  `nama_lengkap` varchar(100) CHARACTER SET latin1 NOT NULL,
  `foto` varchar(100) CHARACTER SET latin1 DEFAULT NULL,
  `level` varchar(20) CHARACTER SET latin1 NOT NULL DEFAULT 'mahasiswa',
  `blokir` enum('Y','N') CHARACTER SET latin1 NOT NULL DEFAULT 'N',
  `last_login` datetime DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;

--
-- Dumping data untuk tabel `login_system`
--

INSERT INTO `login_system` (`userid`, `password`, `nama_lengkap`, `foto`, `level`, `blokir`, `last_login`) VALUES
('admin', 'a4616464048c2c9e48a43550b0b67a33', 'Admin System', 'kosong.png', 'admin', 'N', NULL),
('0103038202', '7f4377a680ac32539d9e8828cf8441c0', 'Susiani, M.Pd', 'kosong.png', 'dosen', 'N', '2026-01-23 02:34:00'),
('nurul', 'f6cab623d64aa761deef18e482171958', 'Nurul Syakirah', 'kosong.png', 'pegawai', 'N', NULL),
('Syam', '49717f1e68e4ca8fa6360f333416211a', 'Syam Aryamin', 'kosong.png', 'pegawai', 'N', NULL),
('imam', '2df1b0a118aa58054c9315cd923c0133', 'Imam Akbari', 'kosong.png', 'pegawai', 'N', NULL),
('sandy', '31ead0fa4acc1cf0f53d1da64837a3e5', 'Sandy Putra Siregar', 'kosong.png', 'pegawai', 'N', NULL),
('Irfan', '21232f297a57a5a743894a0e4a801fc3', 'Irfan Sudahri Damanik', 'kosong.png', 'admin', 'N', '2026-02-08 22:18:56'),
('0104029101', '9b73c0dd4d23ecc87075e331d9012bde', 'Zulaini Masruro Nasution, S.Pd.I., M. Pd', '0104029101.jpg', 'dosen', 'N', '2026-02-10 04:21:41'),
('0104028203', '3531355e116e51e476b52377e0cbbbda', 'Ilham Syahputra Saragih, S.Sos., M.M', 'kosong.png', 'dosen', 'N', '2026-02-10 04:20:58'),
('0104018304', 'd124b17456d9ee69aa60bba3a123e165', 'Dr. Solikhun, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-19 18:48:32'),
('0007087406', '4c2d169ab4f74b0fe61537ca7b3eeb88', 'Victor Asido Elyakim P, M.Kom', '0007087406.jpeg', 'dosen', 'N', NULL),
('0023056503', 'f4bbb0996859fa9e9af3b8603f4a816c', 'Rahmat W. Sembiring, M.Sc.IT, Ph.D', 'kosong.png', 'dosen', 'N', '2026-02-09 07:17:05'),
('0101048301', '285812601434475b1ad435573234f652', 'Zulia Almaida Siregar, S.Pd, M.M', 'kosong.png', 'dosen', 'N', '2026-02-09 19:07:52'),
('0101048902', '8d46930247ba220138fd12872dd200fa', 'Widodo Saputra, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-10 04:19:27'),
('0102088407', '06813c2d8ecb473198bf77aabf2d53af', 'Surya Darma, M.Kom', 'kosong.png', 'dosen', 'N', NULL),
('0102108904', '3e2f6115cfb59dd7c88e4c76c5c8bed5', 'Ika Okta Kirana, M.Pd', 'kosong.png', 'dosen', 'N', NULL),
('0103029401', 'c9afed9d1aaa7c2d35dc31f3798e8bcf', 'Ika Purnama Sari, A.Md., M.Pd', 'kosong.png', 'dosen', 'N', NULL),
('0104068101', 'da656ca77a3ecb66d94763992bed366c', 'Dr. Poningsih, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-03 10:23:28'),
('0104098503', 'd290a5cdca186f1908006063789e633e', 'Riki Winanjaya, M. Kom', 'Foto_0104098503_3eff587be5a7b6d6.jpeg', 'dosen', 'N', '2026-02-05 22:46:51'),
('0105088403', 'cd8b45a2c6aa0ea1e2260ec19787755d', 'Fitri Rizki, S.Pd., M.Pd', 'Foto_0105088403_0558afaf9d1c9d08.jpg', 'dosen', 'N', '2026-02-09 00:15:38'),
('0106049201', '3cea27458434057ffd3ac5541ee1ee42', 'Fitri Anggraini, M.Pd', '0106049201.jpg', 'dosen', 'N', '2026-02-06 20:08:45'),
('0107018401', 'df946244e2f00f06cb4bffb649430bed', 'Iin Parlina, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-09 22:28:49'),
('0107057903', 'eed5e5f192f96449937b4332f2648103', 'Saifullah, M.Kom', 'kosong.png', 'dosen', 'N', NULL),
('0108097802', 'f9fb931cfaf02a5319fd1887109f3d82', 'Dr. Ahmad Fithrianto, S.Ag, M.A', 'kosong.png', 'dosen', 'N', NULL),
('0109029502', 'f8963f9e9ef1069e8b4838e9fb0dd518', 'Rizki Alfadillah Nasution, S.P., M.M', 'kosong.png', 'dosen', 'N', '2026-02-03 02:15:20'),
('0109079701', '9f6eb93ffc4ed9d4600eb71cff524ef0', 'P.P.P.A.N.W.Fikrul Ilmi R.H.Zer, M.Kom', '0109079701.jpg', 'dosen', 'N', '2026-02-06 04:17:09'),
('0110057302', '8e4dff47caef0f7bd262c935c867e06b', 'Bahrudi Efendi Damanik, S.E., M.M', 'kosong.png', 'dosen', 'N', '2026-02-08 08:25:10'),
('0110068001', 'e28502c4b6031dcd25e4ee226a6246c9', 'Dr. M. Safii, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-06 19:56:40'),
('0111065301', 'e997914e4c16d2734c0641acbd473935', 'Jalaluddin, S.Pd., M.Pd', 'kosong.png', 'dosen', 'N', '2026-01-19 06:46:47'),
('0111107301', '8887e77a835e85b7f13d67332dc4f6f3', 'Dr. Dedy Hartama, S.T, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-06 19:59:16'),
('0112088002', '74c146b74edcd8fe04b95a7e27b40072', 'Fahmi Firzada, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-10 04:20:37'),
('0113067703', '5ceb143f0532f3ab4d1bd04ddff58da1', 'Handrizal, S.Si., M.Comp.Sc', 'kosong.png', 'dosen', 'N', NULL),
('0113087701', '646cd0548faef12af1f7ffab5fbcc270', 'Asmarani Nasution, M.S.I', 'kosong.png', 'dosen', 'N', '2026-02-06 04:06:25'),
('0114017002', '9a2c95a097983589e4e7560e0cc84f74', 'Irawan, S.E., M.M', 'kosong.png', 'dosen', 'N', '2026-02-06 20:19:32'),
('0114028503', '41ac1de60537a8af8a7d1551adc4bd23', 'Dr. Anjar Wanto, M.Kom', '0114028503.jpg', 'dosen', 'N', '2026-02-19 20:35:36'),
('0114108402', '64b6023d793088c4c10f99d037e44351', 'Gali Rakasiwi, S.Kom', 'kosong.png', 'dosen', 'N', NULL),
('0116038402', '2598e856f4e00023e3e418bd5e5e0cd1', 'Sumarno, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-06 20:10:30'),
('0116128004', '99010aa67bfc6b4134c5d0a3f0b02be7', 'Dedi Suhendro, S.E., M.Si', '0116128004.jpeg', 'dosen', 'N', '2026-02-10 04:21:14'),
('0117049701', '801313b9550c2ccba02380e665dee4b1', 'Putrama Alkhairi, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-10 18:42:02'),
('0118018901', '2e185486841f88c35acebcb2a917867e', 'Eka Irawan, M.Kom', 'kosong.png', 'dosen', 'N', '2026-03-27 04:06:53'),
('0118028303', 'e9bf39aecc750eb51968ce8e8706299d', 'Fetrisia Iriani Ningrum, S.Pd', 'kosong.png', 'dosen', 'N', NULL),
('0118088603', '0192023a7bbd73250516f069df18b500', 'Irfan Sudahri Damanik, M.Kom', '0118088603.jpg', 'dosen', 'N', '2026-05-15 19:09:57'),
('0119116701', '30436896270b3f3f593f9fe8d6b755d8', 'Heri Santoso, M.Kom', 'kosong.png', 'dosen', 'N', NULL),
('0120069202', '95055222666d88852b6ce1836043dc61', 'Wendi Robiansyah, M.Kom', '0120069202.jpg', 'dosen', 'N', '2026-02-03 02:16:54'),
('0122039401', 'a5b48372c71f71f4fa8b6ccd4009fd34', 'Yuegilion Pranayama Purba, S.Si, M.Si', 'kosong.png', 'dosen', 'N', '2026-02-03 03:32:16'),
('0122067502', '765e32b405189b77ddba23a43606e0bc', 'Hendry Qurniawan, S.Kom., M.M', 'kosong.png', 'dosen', 'N', NULL),
('0122129301', 'b749c9f968ca871480d395363b98ef13', 'Ela Roza Batubara, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-03 02:18:28'),
('0126039601', '813b9911497a28cb92a26bbbfb09f246', 'Abdi Rahim Damanik, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-08 09:04:49'),
('0126099401', '28a840ce818f82b0280570237867f0fb', 'Rizky Khairunnisa Sormin, S.S., M.Pd', 'kosong.png', 'dosen', 'N', '2026-02-10 04:20:08'),
('0126118601', '3716394279eceb4bcb984d2536d95da1', 'Muhammad Ridwan Lubis, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-10 04:20:02'),
('0127108703', '453f26ee99a619948590f180de0924a1', 'Heru Satria Tambunan, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-03 19:04:50'),
('0127118402', '26b01d6a02e1e665f9c074cb252be14e', 'Rafiqa Dewi, M.Kom', 'kosong.png', 'dosen', 'N', '2026-01-20 05:09:32'),
('0128109101', 'e085ba9b6fcd8974da286676b6ea0be7', 'Harly Okprana, M. Kom', 'kosong.png', 'dosen', 'N', NULL),
('0129098502', 'dcc7d6d80f6530036ee5c2b9485c56a6', 'Dr. Sundari Retno Andani, S.T., M.Kom', 'kosong.png', 'dosen', 'N', '2026-01-29 21:21:47'),
('0129118801', '2aabb1b6c3c62c492b6441a5b536c718', 'Indra Gunawan, M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-03 18:48:15'),
('0130088603', 'a4d4864102bd9e6d4349d054c0c870aa', 'Dr. Agus Perdana Windarto, M.Kom', '0130088603.JPEG', 'dosen', 'N', '2026-01-20 23:12:07'),
('0131109101', 'bd78618858fe0c2171d0b52e617683ff', 'Nani Hidayati, M.Kom', '0131109101.jpg', 'dosen', 'N', '2026-02-08 07:18:48'),
('0131129601', '3cc4cfa6c9b8dc0d4b47aaf962d1066b', 'Fazli Nugraha Tambunan, S.P., M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-03 19:05:40'),
('egi', 'ae790caa06a60bd61a7dab21301c7e51', 'Egi Batubara', 'kosong.png', 'pegawai', 'N', NULL),
('0130056737', 'f30ed694f4df8e4ab7eb62eab13b7173', 'Dr. Muhammad Zein, S.Pd.I, M.Pd.I', 'kosong.png', 'dosen', 'N', '2026-02-04 04:56:48'),
('0103079801', 'dacce2306360188406c29cfe0438b724', 'Egi Batubara, S.Kom', 'kosong.png', 'dosen', 'N', NULL),
('0124099901', 'd4164414fb79c11e0a0cad84cc401aac', 'Muhammad Trihardyansyah, S.Kom', 'kosong.png', 'dosen', 'N', '2026-02-03 20:18:51'),
('0104049001', '51b87ab06f185a250c0ac6722f86ddb2', 'Ade Ismiaty Ramadhona Ht. Barat, S.Kom', 'kosong.png', 'dosen', 'N', '2026-02-05 20:55:13'),
('0114089101', '749a92d488ffa9ad80d5e0ae82ff8346', 'Imam Akbari, S.Kom', 'kosong.png', 'dosen', 'N', '2026-02-02 22:14:29'),
('0104049701', 'feb7b760e93068cf1442fe722cf73b66', 'Sandy Putra Siregar, S.Kom', 'kosong.png', 'dosen', 'N', '2026-01-30 19:29:41'),
('9362751652130083', '93d4ee396889af6e66c262336cc18b96', 'Syawaluddin Kadafi Parinduri. M.Kom', 'kosong.png', 'dosen', 'N', '2026-02-08 21:12:06'),
('Hardy', 'c71728ba1b6011c40140f2b056864930', 'Muhammad Trihardyansyah', 'kosong.png', 'pegawai', 'N', NULL),
('0111119601', 'b9930dd00b4161020c5cb16ecd58adac', 'Abdullah Ahmad, S.Kom', 'kosong.png', 'dosen', 'N', '2026-02-02 22:18:21'),
('0125117101', 'cd30414cfe91faeccdb269b753f4979b', 'Husnul Arifin S.Ag, S.Pd.I, M.Pd', 'kosong.png', 'dosen', 'N', '2026-02-03 10:23:51'),
('irfan_admin', '0192023a7bbd73250516f069df18b500', 'Irfan Admin', NULL, 'admin', 'N', '2026-05-13 22:04:57'),
('0130056731', '138c3abde84842d856c827e5edf7086a', 'Sahmi Purba, M.PdK', NULL, 'dosen', 'N', '2026-02-09 20:52:33'),
('9999998', 'facf1f36bbcf6f3cefd7429e61999d26', 'Drs. Nesar Achmad Khan, M.Pd', NULL, 'dosen', 'N', '2026-02-09 00:14:44'),
('ridwan_admin', 'a63fd014db1848536b9d441e464ae968', 'Muhammad Ridwan Lubis, M.Kom', NULL, 'admin', 'N', '2026-02-10 04:03:45'),
('0108109901', '4974488adb9d5ad29659ee534bc6c1b5', 'Anan Wibowo, S.Kom', 'kosong.png', 'dosen', 'N', NULL),
('0113089801', '7f363d5723662d9a03a21849a7ff7ebd', 'Khairun Nisa Arifin Nur, S.Kom', 'kosong.png', 'dosen', 'N', '2026-02-05 21:57:04'),
('0103050001', 'd51b26e77e9f07e752f75a404e5953c0', 'Nazlina Izmi Addyna, S.Kom', 'kosong.png', 'dosen', 'N', '2026-02-06 02:03:22'),
('0119099501', '29e85eeeec29082b7ff07af0acf111b7', 'Rahmat Zulpani, S.Kom', 'kosong.png', 'dosen', 'N', NULL),
('0123018301', '80ba6e81c5f5532e782ba1639dc79217', 'Wiwik Sri Astuti, S.Kom', 'kosong.png', 'dosen', 'N', '2026-02-05 22:26:47');

-- --------------------------------------------------------

--
-- Struktur dari tabel `periode`
--

CREATE TABLE `periode` (
  `id` int(11) NOT NULL,
  `nama_periode` varchar(64) NOT NULL,
  `start_date` date NOT NULL,
  `end_date` date NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data untuk tabel `periode`
--

INSERT INTO `periode` (`id`, `nama_periode`, `start_date`, `end_date`) VALUES
(2, 'Juni 2026', '2026-06-01', '2026-06-30');

-- --------------------------------------------------------

--
-- Struktur dari tabel `tupoksi`
--

CREATE TABLE `tupoksi` (
  `id` int(11) NOT NULL,
  `userid` varchar(50) NOT NULL,
  `nama` text NOT NULL,
  `deskripsi` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data untuk tabel `tupoksi`
--

INSERT INTO `tupoksi` (`id`, `userid`, `nama`, `deskripsi`) VALUES
(1, '0118088603', 'Mengisi KRS', 'Gtw'),
(2, '0118088603', 'Membuat KHS', 'Gtw'),
(4, '0118088603', 'zdzgg', 'gfgfdg'),
(5, '0118018901', 'Mengerjakan Ardiva', 'ABC');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `aktivitas`
--
ALTER TABLE `aktivitas`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `anggota`
--
ALTER TABLE `anggota`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `detail_pegawai`
--
ALTER TABLE `detail_pegawai`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `userid` (`userid`);

--
-- Indexes for table `file_bukti`
--
ALTER TABLE `file_bukti`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `komentar_aktivitas`
--
ALTER TABLE `komentar_aktivitas`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `login_system`
--
ALTER TABLE `login_system`
  ADD PRIMARY KEY (`userid`);

--
-- Indexes for table `periode`
--
ALTER TABLE `periode`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `tupoksi`
--
ALTER TABLE `tupoksi`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `aktivitas`
--
ALTER TABLE `aktivitas`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=39;

--
-- AUTO_INCREMENT for table `anggota`
--
ALTER TABLE `anggota`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `detail_pegawai`
--
ALTER TABLE `detail_pegawai`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT for table `file_bukti`
--
ALTER TABLE `file_bukti`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `komentar_aktivitas`
--
ALTER TABLE `komentar_aktivitas`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT for table `periode`
--
ALTER TABLE `periode`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `tupoksi`
--
ALTER TABLE `tupoksi`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
