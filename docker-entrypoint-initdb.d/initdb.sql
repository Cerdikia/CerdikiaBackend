CREATE DATABASE IF NOT EXISTS cerdikia;

USE cerdikia;

CREATE TABLE users (
    email VARCHAR(255) NOT NULL,
    nama VARCHAR(255) NOT NULL,
    kelas VARCHAR(50) NOT NULL,
    date_created TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    role ENUM('siswa','admin','guru') NOT NULL DEFAULT 'siswa',
    PRIMARY KEY (email)
);

-- Buat tabel mapel
CREATE TABLE IF NOT EXISTS mapel (
  id_mapel INT NOT NULL AUTO_INCREMENT,
  mapel VARCHAR(255) NOT NULL,
  PRIMARY KEY (id_mapel)
);

-- Buat tabel soal
CREATE TABLE IF NOT EXISTS soal (
  id_soal INT NOT NULL AUTO_INCREMENT,
  id_module INT NOT NULL,
  soal TEXT NOT NULL,
  jenis ENUM('pilihan_ganda','essay') NOT NULL,
  jawaban TEXT,
  opsi_jawaban JSON DEFAULT NULL,
  PRIMARY KEY (id_soal, id_module),
  KEY id_module (id_module)
);

-- Buat tabel user_points
CREATE TABLE IF NOT EXISTS user_points (
  email VARCHAR(255) NOT NULL,
  point INT DEFAULT '0',
  PRIMARY KEY (email),
  CONSTRAINT user_points_ibfk_1 FOREIGN KEY (email) REFERENCES users (email) ON DELETE CASCADE
);

-- Buat tabel user_progres
CREATE TABLE IF NOT EXISTS user_progres (
  email VARCHAR(255) NOT NULL,
  id_mapel INT NOT NULL,
  progres JSON DEFAULT NULL,
  PRIMARY KEY (email, id_mapel),
  KEY id_mapel (id_mapel),
  CONSTRAINT user_progres_ibfk_1 FOREIGN KEY (email) REFERENCES users (email) ON DELETE CASCADE,
  CONSTRAINT user_progres_ibfk_2 FOREIGN KEY (id_mapel) REFERENCES mapel (id_mapel) ON DELETE CASCADE
);

-- Buat tabel modules
CREATE TABLE IF NOT EXISTS modules (
  id_module INT NOT NULL AUTO_INCREMENT,
  id_mapel INT NOT NULL,
  kelas VARCHAR(50) NOT NULL,
  module INT DEFAULT NULL,
  module_judul TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  module_deskripsi TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  PRIMARY KEY (id_module),
  KEY id_module (module),
  KEY id_mapel (id_mapel),
  CONSTRAINT modules_ibfk_2 FOREIGN KEY (module) REFERENCES soal (id_module) ON DELETE CASCADE,
  CONSTRAINT modules_ibfk_3 FOREIGN KEY (id_mapel) REFERENCES mapel (id_mapel) ON DELETE CASCADE
);