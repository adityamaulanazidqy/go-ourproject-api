-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Waktu pembuatan: 10 Bulan Mei 2025 pada 18.55
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
-- Struktur dari tabel `files_masterpiece`
--

CREATE TABLE `files_masterpiece` (
  `id` bigint UNSIGNED NOT NULL,
  `masterpiece_id` bigint UNSIGNED NOT NULL,
  `file_patch` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

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
-- Struktur dari tabel `masterpiece`
--

CREATE TABLE `masterpiece` (
  `id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED NOT NULL,
  `status_id` tinyint UNSIGNED NOT NULL,
  `publication_date` timestamp NULL DEFAULT NULL,
  `link_github` varchar(255) DEFAULT NULL,
  `viewer_count` int UNSIGNED DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Used to store masterpiece files';

-- --------------------------------------------------------

--
-- Struktur dari tabel `masterpiece_statuses`
--

CREATE TABLE `masterpiece_statuses` (
  `id` tinyint UNSIGNED NOT NULL,
  `name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

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
(1, 'Siswa');

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
(5, 'Zidny Isyah', 'zidny@siswa.smktiannajiyah.sch.id', '$2a$10$UmbuYtsV7xN7aOmOt8zZMecHGgw.N/kkv2c2dM54QUIuWqR.SAraW', 2, 1, '2020', 'zidny.jpg', '2025-05-10 18:00:01', '2025-05-10 18:00:01');

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
-- Indeks untuk tabel `files_masterpiece`
--
ALTER TABLE `files_masterpiece`
  ADD PRIMARY KEY (`id`),
  ADD KEY `masterpiece_id` (`masterpiece_id`);

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
-- Indeks untuk tabel `masterpiece`
--
ALTER TABLE `masterpiece`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`),
  ADD KEY `status_id` (`status_id`);

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
-- AUTO_INCREMENT untuk tabel `files_masterpiece`
--
ALTER TABLE `files_masterpiece`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

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
-- AUTO_INCREMENT untuk tabel `masterpiece`
--
ALTER TABLE `masterpiece`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `masterpiece_statuses`
--
ALTER TABLE `masterpiece_statuses`
  MODIFY `id` tinyint UNSIGNED NOT NULL AUTO_INCREMENT;

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
  MODIFY `id` tinyint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

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
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

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
  ADD CONSTRAINT `files_masterpiece_ibfk_1` FOREIGN KEY (`masterpiece_id`) REFERENCES `masterpiece` (`id`);

--
-- Ketidakleluasaan untuk tabel `files_thesis`
--
ALTER TABLE `files_thesis`
  ADD CONSTRAINT `files_thesis_ibfk_1` FOREIGN KEY (`publication_id`) REFERENCES `publications` (`id`);

--
-- Ketidakleluasaan untuk tabel `masterpiece`
--
ALTER TABLE `masterpiece`
  ADD CONSTRAINT `masterpiece_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  ADD CONSTRAINT `masterpiece_ibfk_2` FOREIGN KEY (`status_id`) REFERENCES `masterpiece_statuses` (`id`);

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
