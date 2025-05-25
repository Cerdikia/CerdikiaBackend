CREATE DATABASE IF NOT EXISTS cerdikia;
USE cerdikia;

-- Tabel kelas
CREATE TABLE IF NOT EXISTS kelas (
  id_kelas INT NOT NULL AUTO_INCREMENT,
  kelas VARCHAR(50) NOT NULL,
  PRIMARY KEY (id_kelas)
);

-- Tabel admin
CREATE TABLE IF NOT EXISTS admin (
  email VARCHAR(100) NOT NULL,
  nama VARCHAR(100) NOT NULL,
  keterangan TEXT,
  date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
  image_profile text DEFAULT NULL,
  PRIMARY KEY (email),
  UNIQUE KEY uniq_admin_email (email)
);

-- Tabel siswa
CREATE TABLE IF NOT EXISTS siswa (
  email VARCHAR(100) NOT NULL,
  nama VARCHAR(100) DEFAULT NULL,
  id_kelas INT DEFAULT NULL,
  date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
  image_profile text DEFAULT NULL,
  PRIMARY KEY (email),
  UNIQUE KEY uniq_siswa_email (email),
  KEY fk_siswa_kelas (id_kelas),
  CONSTRAINT fk_siswa_kelas FOREIGN KEY (id_kelas) REFERENCES kelas (id_kelas) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Tabel user_energy
CREATE TABLE IF NOT EXISTS user_energy (
  email VARCHAR(100) NOT NULL,
  energy INT DEFAULT 0,
  last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (email),
  CONSTRAINT fk_user_energy_email FOREIGN KEY (email)
    REFERENCES siswa(email)
    ON DELETE CASCADE ON UPDATE CASCADE
);

-- Tabel barang
CREATE TABLE IF NOT EXISTS barang (
  id_barang INT NOT NULL AUTO_INCREMENT,
  nama_barang TEXT,
  img TEXT,
  description TEXT,
  diamond BIGINT DEFAULT NULL,
  jumlah BIGINT DEFAULT NULL,
  created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (id_barang)
);

-- Tabel logs penukaran poin
CREATE TABLE IF NOT EXISTS logs_penukaran_point (
  id_log INT NOT NULL AUTO_INCREMENT,
  id_barang INT NOT NULL,
  email VARCHAR(100) NOT NULL,
  jumlah INT NOT NULL,
  tanggal_penukaran DATETIME DEFAULT CURRENT_TIMESTAMP,
  kode_penukaran VARCHAR(12) NOT NULL UNIQUE,
  status_penukaran ENUM('menunggu', 'selesai', 'dibatalkan') NOT NULL DEFAULT 'menunggu',
  PRIMARY KEY (id_log),
  KEY fk_logs_barang (id_barang),
  KEY fk_logs_siswa (email),
  CONSTRAINT fk_logs_barang FOREIGN KEY (id_barang) REFERENCES barang (id_barang) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_logs_siswa FOREIGN KEY (email) REFERENCES siswa (email) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Tabel mapel
CREATE TABLE IF NOT EXISTS mapel (
  id_mapel INT NOT NULL AUTO_INCREMENT,
  mapel VARCHAR(255) NOT NULL,
  PRIMARY KEY (id_mapel)
);

-- Tabel guru
CREATE TABLE IF NOT EXISTS guru (
  id int NOT NULL AUTO_INCREMENT,
  email VARCHAR(100) NOT NULL,
  nama VARCHAR(100) NOT NULL,
  jabatan VARCHAR(100) NOT NULL DEFAULT 'guru',
  date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
  image_profile text DEFAULT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS guru_mapel (
  id_guru INT NOT NULL,
  id_mapel INT NOT NULL,
  PRIMARY KEY (id_guru, id_mapel),
  FOREIGN KEY (id_guru) REFERENCES guru(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (id_mapel) REFERENCES mapel(id_mapel) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Tabel modules
CREATE TABLE IF NOT EXISTS modules (
  id_module INT NOT NULL AUTO_INCREMENT,
  id_mapel INT NOT NULL,
  id_kelas INT NOT NULL,
  module INT DEFAULT NULL,
  module_judul TEXT,
  module_deskripsi TEXT,
  is_ready TINYINT(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (id_module),
  KEY fk_module_mapel (id_mapel),
  KEY fk_module_kelas (id_kelas),
  CONSTRAINT fk_module_mapel FOREIGN KEY (id_mapel) REFERENCES mapel (id_mapel) ON DELETE CASCADE,
  CONSTRAINT fk_module_kelas FOREIGN KEY (id_kelas) REFERENCES kelas (id_kelas) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Tabel soal
CREATE TABLE IF NOT EXISTS soal (
  id_soal INT NOT NULL AUTO_INCREMENT,
  id_module INT NOT NULL,
  soal LONGTEXT NOT NULL,
  jenis ENUM('pilihan_ganda', 'essay') NOT NULL DEFAULT 'pilihan_ganda',
  opsi_a LONGTEXT NOT NULL,
  opsi_b LONGTEXT NOT NULL,
  opsi_c LONGTEXT NOT NULL,
  opsi_d LONGTEXT NOT NULL,
  jawaban VARCHAR(10) NOT NULL,
  PRIMARY KEY (id_soal),
  KEY fk_soal_module (id_module),
  CONSTRAINT fk_soal_module FOREIGN KEY (id_module) REFERENCES modules (id_module) ON DELETE CASCADE
);

-- Tabel logs
CREATE TABLE IF NOT EXISTS logs (
  id_logs INT NOT NULL AUTO_INCREMENT,
  email VARCHAR(255) DEFAULT NULL,
  id_kelas INT DEFAULT NULL,
  id_mapel INT DEFAULT NULL,
  id_module INT DEFAULT NULL,
  skor INT DEFAULT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id_logs),
  KEY fk_logs_email (email),
  KEY fk_logs_module (id_module),
  KEY fk_logs_kelas (id_kelas),
  KEY fk_logs_mapel (id_mapel),
  CONSTRAINT fk_logs_email FOREIGN KEY (email) REFERENCES siswa (email) ON DELETE SET NULL,
  CONSTRAINT fk_logs_module FOREIGN KEY (id_module) REFERENCES modules (id_module) ON DELETE SET NULL,
  CONSTRAINT fk_logs_kelas FOREIGN KEY (id_kelas) REFERENCES kelas (id_kelas) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT fk_logs_mapel FOREIGN KEY (id_mapel) REFERENCES mapel (id_mapel) ON DELETE SET NULL ON UPDATE CASCADE
);

-- Tabel verifikasi user
CREATE TABLE IF NOT EXISTS user_verified (
  email VARCHAR(100) NOT NULL,
  verified_status ENUM('accept', 'rejected', 'waiting') DEFAULT 'waiting',
  PRIMARY KEY (email),
  CONSTRAINT fk_verified_siswa FOREIGN KEY (email) REFERENCES siswa (email) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Tabel chat
CREATE TABLE IF NOT EXISTS chat (
  id_message INT AUTO_INCREMENT PRIMARY KEY,
  form VARCHAR(100),
  entity ENUM('personal', 'role') NOT NULL,
  dest VARCHAR(100) NOT NULL, 
  subject VARCHAR(255), 
  message TEXT NOT NULL,
  status ENUM('mengirim', 'terkirim', 'dibaca') DEFAULT 'mengirim',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabel poin user
CREATE TABLE IF NOT EXISTS user_points (
  email VARCHAR(255) NOT NULL,
  diamond INT DEFAULT 0,
  exp INT NOT NULL DEFAULT 0,
  PRIMARY KEY (email),
  CONSTRAINT fk_points_siswa FOREIGN KEY (email) REFERENCES siswa (email) ON DELETE CASCADE
);

-- Tabel data_siswa
CREATE TABLE IF NOT EXISTS data_siswa (
  id_data INT NOT NULL AUTO_INCREMENT,
  email VARCHAR(255) DEFAULT NULL,
  id_kelas INT DEFAULT NULL,
  progres JSON DEFAULT NULL,
  semester ENUM('ganjil', 'genap') NOT NULL,
  tahun_ajaran VARCHAR(50) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id_data),
  KEY fk_data_email (email),
  KEY fk_data_kelas (id_kelas),
  CONSTRAINT fk_data_email FOREIGN KEY (email) REFERENCES siswa (email) ON DELETE SET NULL,
  CONSTRAINT fk_data_kelas FOREIGN KEY (id_kelas) REFERENCES kelas (id_kelas) ON DELETE SET NULL ON UPDATE CASCADE
);

-- -- Tambahkan akun admin default
INSERT INTO admin (email, nama, keterangan)
VALUES 
('mutekinoraffi@gmail.com', 'raffi', 'default admin1'),
('martizasyam22@gmail.com', 'farhan', 'default admin2'),
('fahmihanafi680@gmail.com', 'fahmi', 'default admin3')
ON DUPLICATE KEY UPDATE 
  nama = VALUES(nama), 
  keterangan = VALUES(keterangan);

INSERT INTO mapel (id_mapel, mapel)
VALUES 
(1, 'Bahasa Indonesia'),
(2, 'Matematika'),
(3, 'Bahasa Inggris');

INSERT INTO kelas (id_kelas, kelas)
VALUES (1, '1'), (2, '2'), (3, '3'), (4, '4'), (5, '5'), (6, '6');

-- -- Insert guru (sekali saja)
-- INSERT INTO guru (email, nama, jabatan)
-- VALUES ('guru1@guru.com', 'guru1', 'walikelas 1'), ('guru2@guru.com', 'guru2', 'guru bahasa inggris')
-- ON DUPLICATE KEY UPDATE nama = VALUES(nama), jabatan = VALUES(jabatan);

-- -- Tambahkan relasi ke 2 mapel
-- INSERT INTO guru_mapel (id_guru, id_mapel) VALUES (1, 1), (1, 2), (2, 3)
-- ON DUPLICATE KEY UPDATE id_mapel = id_mapel;