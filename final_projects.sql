-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Waktu pembuatan: 15 Bulan Mei 2025 pada 08.31
-- Versi server: 8.0.30
-- Versi PHP: 8.3.14

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `final_projects`
--

-- --------------------------------------------------------

--
-- Struktur dari tabel `activity_logs`
--

CREATE TABLE `activity_logs` (
  `id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED DEFAULT NULL,
  `action` varchar(255) NOT NULL,
  `model` varchar(255) DEFAULT NULL,
  `model_id` bigint UNSIGNED DEFAULT NULL,
  `data` text,
  `ip_address` varchar(45) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `classes`
--

CREATE TABLE `classes` (
  `id` int NOT NULL,
  `class` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `classes`
--

INSERT INTO `classes` (`id`, `class`) VALUES
(1, 10),
(2, 11),
(3, 12);

-- --------------------------------------------------------

--
-- Struktur dari tabel `files_masterpiece`
--

CREATE TABLE `files_masterpiece` (
  `id` bigint UNSIGNED NOT NULL,
  `masterpiece_id` bigint UNSIGNED NOT NULL,
  `file_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `files_masterpiece`
--

INSERT INTO `files_masterpiece` (`id`, `masterpiece_id`, `file_path`) VALUES
(1, 1, 'Cuplikan layar 2025-04-25 090629.png'),
(2, 1, 'Cuplikan layar 2025-04-25 090808.png');

-- --------------------------------------------------------

--
-- Struktur dari tabel `files_thesis`
--

CREATE TABLE `files_thesis` (
  `id` bigint UNSIGNED NOT NULL,
  `publication_id` bigint UNSIGNED NOT NULL,
  `file_patch` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `majors`
--

CREATE TABLE `majors` (
  `id` tinyint UNSIGNED NOT NULL,
  `name` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `majors`
--

INSERT INTO `majors` (`id`, `name`) VALUES
(2, 'Perfilman'),
(1, 'Rekayasa Perangkat Lunak');

-- --------------------------------------------------------

--
-- Struktur dari tabel `masterpieces`
--

CREATE TABLE `masterpieces` (
  `id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `status_id` tinyint UNSIGNED NOT NULL,
  `class_id` int NOT NULL,
  `semester_id` int NOT NULL,
  `publication_date` timestamp NULL DEFAULT NULL,
  `link_github` varchar(255) DEFAULT NULL,
  `viewer_count` int UNSIGNED DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Used to store masterpiece files';

--
-- Dumping data untuk tabel `masterpieces`
--

INSERT INTO `masterpieces` (`id`, `user_id`, `status_id`, `class_id`, `semester_id`, `publication_date`, `link_github`, `viewer_count`, `created_at`, `updated_at`) VALUES
(1, 5, 1, 1, 2, '2025-05-15 08:30:21', 'http:github.example.com', 0, '2025-05-15 08:30:21', '2025-05-15 08:30:21');

-- --------------------------------------------------------

--
-- Struktur dari tabel `masterpiece_statuses`
--

CREATE TABLE `masterpiece_statuses` (
  `id` tinyint UNSIGNED NOT NULL,
  `name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `masterpiece_statuses`
--

INSERT INTO `masterpiece_statuses` (`id`, `name`) VALUES
(2, 'Pengembangan'),
(1, 'Selesai');

-- --------------------------------------------------------

--
-- Struktur dari tabel `publications`
--

CREATE TABLE `publications` (
  `id` bigint UNSIGNED NOT NULL,
  `thesis_title_id` bigint UNSIGNED NOT NULL,
  `status_id` tinyint UNSIGNED NOT NULL,
  `publication_date` timestamp NULL DEFAULT NULL,
  `link_github` varchar(255) DEFAULT NULL,
  `viewer_count` int UNSIGNED DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Used to store publications files';

-- --------------------------------------------------------

--
-- Struktur dari tabel `publication_statuses`
--

CREATE TABLE `publication_statuses` (
  `id` tinyint UNSIGNED NOT NULL,
  `name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `roles`
--

CREATE TABLE `roles` (
  `id` tinyint UNSIGNED NOT NULL,
  `name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `roles`
--

INSERT INTO `roles` (`id`, `name`) VALUES
(2, 'Guru'),
(3, 'Pembimbing'),
(1, 'Siswa'),
(5, 'Super Admin'),
(6, 'Testing');

-- --------------------------------------------------------

--
-- Struktur dari tabel `semesters`
--

CREATE TABLE `semesters` (
  `id` int NOT NULL,
  `name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `semesters`
--

INSERT INTO `semesters` (`id`, `name`) VALUES
(1, 'Ganjil'),
(2, 'Genap');

-- --------------------------------------------------------

--
-- Struktur dari tabel `supervision_sessions`
--

CREATE TABLE `supervision_sessions` (
  `id` bigint UNSIGNED NOT NULL,
  `thesis_title_id` bigint UNSIGNED NOT NULL,
  `session_date` datetime NOT NULL,
  `progress` tinyint UNSIGNED NOT NULL,
  `notes` text,
  `document` varchar(255) DEFAULT NULL,
  `status_id` tinyint UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `supervision_statuses`
--

CREATE TABLE `supervision_statuses` (
  `id` tinyint UNSIGNED NOT NULL,
  `name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `task_statuses`
--

CREATE TABLE `task_statuses` (
  `id` tinyint UNSIGNED NOT NULL,
  `name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `thesis_statuses`
--

CREATE TABLE `thesis_statuses` (
  `id` tinyint UNSIGNED NOT NULL,
  `name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `thesis_tasks`
--

CREATE TABLE `thesis_tasks` (
  `id` bigint UNSIGNED NOT NULL,
  `thesis_title_id` bigint UNSIGNED NOT NULL,
  `task_description` text NOT NULL,
  `deadline` date DEFAULT NULL,
  `status_id` tinyint UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `thesis_titles`
--

CREATE TABLE `thesis_titles` (
  `id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `supervisor_id` bigint UNSIGNED DEFAULT NULL,
  `title` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `status_id` tinyint UNSIGNED NOT NULL,
  `submission_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `disposition_date` timestamp NULL DEFAULT NULL,
  `notes` text,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `users`
--

CREATE TABLE `users` (
  `id` bigint UNSIGNED NOT NULL,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `role_id` tinyint UNSIGNED NOT NULL,
  `major_id` tinyint UNSIGNED NOT NULL,
  `batch` year NOT NULL,
  `photo` varchar(255) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data untuk tabel `users`
--

INSERT INTO `users` (`id`, `username`, `email`, `password`, `role_id`, `major_id`, `batch`, `photo`, `created_at`, `updated_at`) VALUES
(1, 'Aditya Maulana Zidqy', 'adityamaullana234@siswa.smktiannajiyah.sch.id', '$2a$10$LElGdmlCdOm.JWABY7k9lu2p6qAcnUg5iypc4Ie91pk2f5UUYhQOG', 1, 1, '2020', 'zidqy.jpg', '2025-05-10 10:13:56', '2025-05-10 10:13:56'),
(2, 'Dliani albana', 'dliani@siswa.smktiannajiyah.sch.id', '$2a$10$9D1kbnxxFWly7VXzm6kIr.yaqNWl15Tx3r0G4E5l0H2FbjdGZubCW', 1, 1, '2022', 'dliani.jpg', '2025-05-10 10:22:51', '2025-05-10 10:22:51'),
(3, 'Gabrielle Echols', 'gabrielle@siswa.smktiannajiyah.sch.id', '$2a$10$HuJN8MazsyFFR56RZBZax.m5niEI8Yt76GcFt6mKnkQgpAvNgB4hW', 1, 1, '2021', 'gabrielle.jpg', '2025-05-10 10:25:30', '2025-05-10 10:25:30'),
(4, 'Sri Handayani', 'handayani@siswa.smktiannajiyah.sch.id', '$2a$10$k6X5C0I/MDMQBVkD0aBMWesDtEstKs7wNtBQ3enGYaqxQAlwyH2Yi', 1, 1, '2025', 'handayani.jpg', '2025-05-10 10:27:06', '2025-05-10 10:27:06'),
(5, 'Zidny Isyah', 'zidny@siswa.smktiannajiyah.sch.id', '$2a$10$UmbuYtsV7xN7aOmOt8zZMecHGgw.N/kkv2c2dM54QUIuWqR.SAraW', 2, 1, '2020', 'zidny.jpg', '2025-05-10 18:00:01', '2025-05-10 18:00:01'),
(15, 'Anjasmara', 'anjasmara@siswa.smktiannajiyah.sch.id', '$2a$10$KpPJxzWFTQKcYBt9IVadF.1ed2/1mlMSI6o1sDyWRNjZ/go4hbEpe', 2, 2, '2020', '1746994325107028800_420310333_2657596564425156_7930308480479293830_n.jpg', '2025-05-11 20:12:05', '2025-05-11 20:12:05'),
(16, 'Krisdayanti', 'krisdayanti@siswa.smktiannajiyah.sch.id', '$2a$10$JsXoTCHG85zNKnv3VWsbEeqT/pKOaNqY0gI.DF3BhUMmkeI1QU9BC', 1, 2, '2019', '1746994505810256400_NakNyamplungCoffe.jpg', '2025-05-11 20:15:06', '2025-05-11 20:15:06'),
(17, 'Lidyawati', 'lidyawati@siswa.smktiannajiyah.sch.id', '$2a$10$JJriZIrkAVAzHrcuT1ZSUec/DhoutdNhm/2OWJOkRvCjnTcHADr1a', 2, 1, '2022', '1746995496558197200_WIN_20240620_19_49_34_Pro.jpg', '2025-05-11 20:31:37', '2025-05-11 20:31:37'),
(18, 'Cindy', 'cindy@siswa.smktiannajiyah.sch.id', '$2a$10$zW.0tKhdUhjC1YbKvKfxI.FDPUGUi28obdOIji6KwdRlW47eAPAHe', 3, 2, '2024', '1747294952179553400_Screenshot (11).png', '2025-05-15 07:42:32', '2025-05-15 07:42:32'),
(19, 'Laraswati', 'laraswati@siswa.smktiannajiyah.sch.id', '$2a$10$ES3n26ltPpzqeH06oUNn7O5u80RUHA5SW0EMmEIrGHF1zNab9wc1W', 3, 2, '2024', '1747295225119752000_Screenshot (11).png', '2025-05-15 07:47:05', '2025-05-15 07:47:05'),
(20, 'Albar Wicaksono', 'albar@siswa.smktiannajiyah.sch.id', '$2a$10$EzHd6hLFhkWhnSxvuyzouOOlZxMahIO5L7JI15VPSzfVD.Ef6j.76', 1, 1, '2020', '1747295645977712100_Screenshot (11).png', '2025-05-15 07:54:06', '2025-05-15 07:54:06'),
(21, 'Suci handayani', 'suci@siswa.smktiannajiyah.sch.id', '$2a$10$rLjaz4i0yZrTjsJmuFkGruYl2cRWMDoRgg4sSSMW5lYA4zc6iHH6u', 1, 1, '2020', '1747296119493900700_Screenshot (11).png', '2025-05-15 08:01:59', '2025-05-15 08:01:59');

--
-- Indexes for dumped tables
--

--
-- Indeks untuk tabel `activity_logs`
--
ALTER TABLE `activity_logs`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`);

--
-- Indeks untuk tabel `classes`
--
ALTER TABLE `classes`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `class` (`class`);

--
-- Indeks untuk tabel `files_masterpiece`
--
ALTER TABLE `files_masterpiece`
  ADD PRIMARY KEY (`id`),
  ADD KEY `masterpiece_id_ibfk_1` (`masterpiece_id`);

--
-- Indeks untuk tabel `files_thesis`
--
ALTER TABLE `files_thesis`
  ADD PRIMARY KEY (`id`),
  ADD KEY `publication_id` (`publication_id`);

--
-- Indeks untuk tabel `majors`
--
ALTER TABLE `majors`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indeks untuk tabel `masterpieces`
--
ALTER TABLE `masterpieces`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`),
  ADD KEY `status_id` (`status_id`),
  ADD KEY `masterpieces_ibfk_3` (`class_id`),
  ADD KEY `masterpieces_ibfk_4` (`semester_id`);

--
-- Indeks untuk tabel `masterpiece_statuses`
--
ALTER TABLE `masterpiece_statuses`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indeks untuk tabel `publications`
--
ALTER TABLE `publications`
  ADD PRIMARY KEY (`id`),
  ADD KEY `thesis_title_id` (`thesis_title_id`),
  ADD KEY `status_id` (`status_id`);

--
-- Indeks untuk tabel `publication_statuses`
--
ALTER TABLE `publication_statuses`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indeks untuk tabel `roles`
--
ALTER TABLE `roles`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indeks untuk tabel `semesters`
--
ALTER TABLE `semesters`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indeks untuk tabel `supervision_sessions`
--
ALTER TABLE `supervision_sessions`
  ADD PRIMARY KEY (`id`),
  ADD KEY `thesis_title_id` (`thesis_title_id`),
  ADD KEY `status_id` (`status_id`);

--
-- Indeks untuk tabel `supervision_statuses`
--
ALTER TABLE `supervision_statuses`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indeks untuk tabel `task_statuses`
--
ALTER TABLE `task_statuses`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indeks untuk tabel `thesis_statuses`
--
ALTER TABLE `thesis_statuses`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indeks untuk tabel `thesis_tasks`
--
ALTER TABLE `thesis_tasks`
  ADD PRIMARY KEY (`id`),
  ADD KEY `thesis_title_id` (`thesis_title_id`),
  ADD KEY `status_id` (`status_id`);

--
-- Indeks untuk tabel `thesis_titles`
--
ALTER TABLE `thesis_titles`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`),
  ADD KEY `supervisor_id` (`supervisor_id`),
  ADD KEY `status_id` (`status_id`);

--
-- Indeks untuk tabel `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`),
  ADD KEY `role_id` (`role_id`),
  ADD KEY `major_id` (`major_id`);

--
-- AUTO_INCREMENT untuk tabel yang dibuang
--

--
-- AUTO_INCREMENT untuk tabel `activity_logs`
--
ALTER TABLE `activity_logs`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `classes`
--
ALTER TABLE `classes`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT untuk tabel `files_masterpiece`
--
ALTER TABLE `files_masterpiece`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT untuk tabel `files_thesis`
--
ALTER TABLE `files_thesis`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `majors`
--
ALTER TABLE `majors`
  MODIFY `id` tinyint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT untuk tabel `masterpieces`
--
ALTER TABLE `masterpieces`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT untuk tabel `masterpiece_statuses`
--
ALTER TABLE `masterpiece_statuses`
  MODIFY `id` tinyint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT untuk tabel `publications`
--
ALTER TABLE `publications`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `publication_statuses`
--
ALTER TABLE `publication_statuses`
  MODIFY `id` tinyint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `roles`
--
ALTER TABLE `roles`
  MODIFY `id` tinyint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT untuk tabel `semesters`
--
ALTER TABLE `semesters`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT untuk tabel `supervision_sessions`
--
ALTER TABLE `supervision_sessions`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `supervision_statuses`
--
ALTER TABLE `supervision_statuses`
  MODIFY `id` tinyint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `task_statuses`
--
ALTER TABLE `task_statuses`
  MODIFY `id` tinyint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `thesis_statuses`
--
ALTER TABLE `thesis_statuses`
  MODIFY `id` tinyint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `thesis_tasks`
--
ALTER TABLE `thesis_tasks`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `thesis_titles`
--
ALTER TABLE `thesis_titles`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `users`
--
ALTER TABLE `users`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=22;

--
-- Ketidakleluasaan untuk tabel pelimpahan (Dumped Tables)
--

--
-- Ketidakleluasaan untuk tabel `activity_logs`
--
ALTER TABLE `activity_logs`
  ADD CONSTRAINT `activity_logs_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL;

--
-- Ketidakleluasaan untuk tabel `files_masterpiece`
--
ALTER TABLE `files_masterpiece`
  ADD CONSTRAINT `masterpiece_id_ibfk_1` FOREIGN KEY (`masterpiece_id`) REFERENCES `masterpieces` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT;

--
-- Ketidakleluasaan untuk tabel `files_thesis`
--
ALTER TABLE `files_thesis`
  ADD CONSTRAINT `files_thesis_ibfk_1` FOREIGN KEY (`publication_id`) REFERENCES `publications` (`id`);

--
-- Ketidakleluasaan untuk tabel `masterpieces`
--
ALTER TABLE `masterpieces`
  ADD CONSTRAINT `masterpieces_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  ADD CONSTRAINT `masterpieces_ibfk_2` FOREIGN KEY (`status_id`) REFERENCES `masterpiece_statuses` (`id`),
  ADD CONSTRAINT `masterpieces_ibfk_3` FOREIGN KEY (`class_id`) REFERENCES `classes` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  ADD CONSTRAINT `masterpieces_ibfk_4` FOREIGN KEY (`semester_id`) REFERENCES `semesters` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT;

--
-- Ketidakleluasaan untuk tabel `publications`
--
ALTER TABLE `publications`
  ADD CONSTRAINT `publications_ibfk_1` FOREIGN KEY (`thesis_title_id`) REFERENCES `thesis_titles` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `publications_ibfk_2` FOREIGN KEY (`status_id`) REFERENCES `publication_statuses` (`id`);

--
-- Ketidakleluasaan untuk tabel `supervision_sessions`
--
ALTER TABLE `supervision_sessions`
  ADD CONSTRAINT `supervision_sessions_ibfk_1` FOREIGN KEY (`thesis_title_id`) REFERENCES `thesis_titles` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `supervision_sessions_ibfk_2` FOREIGN KEY (`status_id`) REFERENCES `supervision_statuses` (`id`);

--
-- Ketidakleluasaan untuk tabel `thesis_tasks`
--
ALTER TABLE `thesis_tasks`
  ADD CONSTRAINT `thesis_tasks_ibfk_1` FOREIGN KEY (`thesis_title_id`) REFERENCES `thesis_titles` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `thesis_tasks_ibfk_2` FOREIGN KEY (`status_id`) REFERENCES `task_statuses` (`id`);

--
-- Ketidakleluasaan untuk tabel `thesis_titles`
--
ALTER TABLE `thesis_titles`
  ADD CONSTRAINT `thesis_titles_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `thesis_titles_ibfk_2` FOREIGN KEY (`supervisor_id`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `thesis_titles_ibfk_3` FOREIGN KEY (`status_id`) REFERENCES `thesis_statuses` (`id`);

--
-- Ketidakleluasaan untuk tabel `users`
--
ALTER TABLE `users`
  ADD CONSTRAINT `users_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
  ADD CONSTRAINT `users_ibfk_2` FOREIGN KEY (`major_id`) REFERENCES `majors` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
