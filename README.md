# Final Project AI-Powered Smart Home Energy Management System RAW

> Studi Independen Bersertifikat — Kampus Merdeka, Kemdikbudristek - Ruangguru

![Kampus Merdeka](https://img.shields.io/badge/Kampus_Merdeka-Studi_Independen-blue)
![Versi](https://img.shields.io/badge/versi-3.0.0-blue)
![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go&logoColor=white)
![AI](https://img.shields.io/badge/AI-TAPAS-orange)
![HuggingFace](https://img.shields.io/badge/HuggingFace-TAPAS-yellow?logo=huggingface&logoColor=black)
![License](https://img.shields.io/badge/license-Academic-lightgrey)

Aplikasi berbasis AI untuk menganalisis konsumsi energi rumah tangga dari data CSV melalui antarmuka CLI interaktif. Dibangun menggunakan Go sebagai backend dan TAPAS (Google, via HuggingFace) sebagai model AI untuk query berbasis tabel. Dikembangkan sebagai bagian dari program **Studi Independen Bersertifikat Kampus Merdeka, Kemdikbudristek**.

## Lihat Repository ReCode#1 disini: https://github.com/Raraeva24/MSIBProject-ReCode_1.git sebagai bagian dari penyempurnaan project ini.


## Daftar Isi

- [Demo](#demo)
- [Fitur](#fitur)
- [Tech Stack](#tech-stack)
- [Cara Kerja](#cara-kerja)
- [Instalasi](#instalasi)
- [Penggunaan](#penggunaan)
- [Format CSV](#format-csv)
- [Contoh Pertanyaan](#contoh-pertanyaan)
- [Struktur Folder](#struktur-folder)
- [Kontak](#kontak)

---

## Demo

**Tampilan menu utama CLI**
```
Welcome to AI Chatbot CLI!
=================================
1. Berapa total penggunaan listrik Refrigerator?
2. Berapa rata-rata penggunaan listrik Washing Machine?
3. Perangkat apa yang paling boros listrik?
4. Berapa total semua penggunaan listrik?
5. Tanya sendiri (ketik pertanyaan bebas)
6. Keluar
=================================
Pilih menu (1-5):
```

**Tampilan saat memilih menu**
```
Pilih menu (1-5): 1
Jawaban:
Total Penggunaan : 25.20 kWh
Cells            : 1.2, 1.2, 1.2, 1.2, 1.2, ...
Coordinates      : [[0 2] [1 2] [2 2] ...]
Aggregator       : SUM
```

**Tampilan pertanyaan bebas**
```
Pilih menu (1-5): 5
Pertanyaanmu: Which appliance has the highest energy consumption?
Jawaban:
Jawaban          : Refrigerator
Cells            : Refrigerator
Coordinates      : [[0 0]]
Aggregator       : NONE
```

---

## Fitur

- **Menu CLI bernomor** — pilih pertanyaan preset dengan menekan angka 1-5
- **Pertanyaan bebas** — menu 5 untuk mengetik pertanyaan sendiri dalam bahasa Inggris
- **Pengolahan hasil otomatis** — SUM dan AVERAGE dihitung ulang dari data Cells secara lokal
- **Output lengkap** — menampilkan semua field Response (hasil, cells, koordinat, aggregator)
- **Fallback handling** — menampilkan pesan error yang informatif jika koneksi ke API gagal

---

## Tech Stack

**AI**

![TAPAS](https://img.shields.io/badge/TAPAS-google%2Ftapas--base--finetuned--wtq-yellow?logo=huggingface&logoColor=black)

**API**

![HuggingFace](https://img.shields.io/badge/HuggingFace-Inference_API-yellow?logo=huggingface&logoColor=black)

**Backend & CLI**

![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)

---

## Cara Kerja

```
User pilih menu (1-6) atau ketik pertanyaan
        │
        ▼
┌──────────────────────────┐
│  Go — buat payload JSON  │  ← query + seluruh tabel CSV
│  {"inputs": {            │
│    "table": {...},       │
│    "query": "..."        │
│  }}                      │
└──────────┬───────────────┘
           │
           ▼
┌──────────────────────────┐
│  TAPAS (HuggingFace)     │  ← model memahami tabel & pertanyaan
│  router.huggingface.co   │
└──────────┬───────────────┘
           │
           ├── Aggregator: SUM     → hitung ulang total dari Cells
           ├── Aggregator: AVERAGE → hitung ulang rata-rata dari Cells
           ├── Aggregator: COUNT   → hitung jumlah item
           └── Aggregator: NONE   → lookup langsung dari tabel
                    │
                    ▼
        ┌─────────────────────────┐
        │  processResponse()      │
        │  format output lengkap  │
        └─────────────────────────┘
                    │
                    ▼
             Tampil di CLI
```

---

## Instalasi

### Prasyarat
- Go `>= 1.21`
- Akun HuggingFace dengan API token aktif

### Langkah

```bash
# 1. Clone repositori
git clone https://github.com/Raraeva24/MSIBProjectRAW.git
cd MSIBProjectRAW

# 2. Install dependensi
go mod tidy

# 3. Buat file .env
cp .env.example .env
```

Isi file `.env`:
```env
HUGGINGFACE_TOKEN=hf_xxxxxxxxxxxxxxxxxxxxxxxx
```

Dapatkan token di: https://huggingface.co/settings/tokens

```bash
# 4. Jalankan aplikasi
go run main.go
```

---

## Format CSV

File CSV yang digunakan harus memiliki kolom berikut:

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `Date` | string | Tanggal pencatatan |
| `Time` | string | Waktu pencatatan |
| `Appliance` | string | Nama perangkat (AC, TV, Refrigerator, dll.) |
| `Energy_Consumption` | float | Konsumsi energi dalam kWh |
| `Room` | string | Ruangan tempat perangkat berada |
| `Status` | string | `On` atau `Off` |

### Contoh Isi CSV

```csv
Date,Time,Appliance,Energy_Consumption,Room,Status
2022-01-01,00:00,Refrigerator,1.2,Kitchen,On
2022-01-01,01:00,Refrigerator,1.2,Kitchen,On
2022-01-01,08:00,TV,0.8,Living Room,Off
2022-01-01,09:00,TV,0.8,Living Room,On
2022-01-01,10:00,Washing Machine,0.9,Bathroom,On
```

---

## Contoh Pertanyaan

Pertanyaan untuk menu bebas (menu 5) harus menggunakan **bahasa Inggris** dan merujuk pada **nama kolom CSV**:

```
What is the total energy consumption of the Refrigerator?
What is the average energy consumption of the Washing Machine?
Which appliance has the highest energy consumption?
What is the total energy consumption of all appliances?
Which room has the highest energy consumption?
How many appliances are on?
```

> ⚠️ **Catatan:** TAPAS hanya bisa menjawab pertanyaan berbasis data tabel. Pertanyaan seperti *"Kenapa refrigerator boros?"* tidak dapat dijawab karena jawabannya tidak ada di dalam tabel CSV.

### Contoh Output

```
Pilih menu (1-5): 2
Jawaban:
Rata-rata        : 1.20 kWh
Cells            : 1.2, 1.2, 1.2, 1.2, 1.2, 1.2, 1.2, 1.2, 1.2
Coordinates      : [[0 2] [16 2] [17 2] [18 2] [19 2] [20 2] [21 2] [22 2] [23 2]]
Aggregator       : AVERAGE
```

---

## Struktur Folder

```
MSIBProject_ReCode_2/
├── main.go            # Entry point — semua logika CLI & AI
├── data-series.csv    # Contoh file CSV
├── .env               # API keys (tidak di-commit)
├── .env.example       # Template .env
├── go.mod
├── go.sum
└── README.md
```

---

## Rencana Re-Code #1

Versi berikutnya akan mengatasi keterbatasan project ini dengan:

> - Menggunakan TAPAS untuk query data tabel
> - Berbasis Web Sederhana
> - Pertanyaan Rekomendasi dapat di jawab berdasarkan aggregasi konsumsi energi per perangkat, status, dan ruangan
> - Dapat mengUpload file CSV terpisah
---

## 👤 Kontak

**Rara Eva Maharani**
- Email: rarevamaharani@gmail.com
- GitHub: [@Raraeva24](https://github.com/Raraeva24)

---

*Repositori ini dikembangkan sebagai bagian dari program Studi Independen Bersertifikat Kampus Merdeka, Kemdikbudristek — 2025.*



<br>
<br>
<br>
<br>

# README LAMA BERISI FLOW PENGERJAAN DARI RUANGGURU



# Artificial Intelligence menggunakan Golang

## Final Project AI-Powered Smart Home Energy Management System

### Description

Kamu akan mengembangkan Sistem Manajemen Energi Rumah Pintar menggunakan Golang dan [model AI Tapas](https://huggingface.co/google/tapas-base-finetuned-wtq) dari Huggingface Model Hub. Sistem ini akan memprediksi dan mengelola penggunaan energi dalam lingkungan rumah pintar. Aplikasi ini akan menerima data tentang penggunaan energi rumah dan memberikan wawasan dan rekomendasi tentang cara mengoptimalkan konsumsi energi.

Fitur:

- Prediksi Konsumsi Energi: Sistem ini akan memprediksi konsumsi energi rumah berdasarkan data historis.

- Rekomendasi Penghematan Energi: Sistem ini akan memberikan rekomendasi tentang cara menghemat energi berdasarkan konsumsi energi yang diprediksi.

Data input dalam bentuk format CSV dengan kolom berikut:

- Date: Tanggal data penggunaan energi.
- Time: Waktu data penggunaan energi.
- Appliance: Nama alat.
- Energy_Consumption: Konsumsi energi alat dalam kWh.
- Room: Ruang tempat alat berada.
- Status: Status alat (On/Off).

Contoh:

```txt
Date,Time,Appliance,Energy_Consumption,Room,Status
2022-01-01,00:00,Refrigerator,1.2,Kitchen,On
2022-01-01,01:00,Refrigerator,1.2,Kitchen,On
...
2022-01-01,08:00,TV,0.8,Living Room,Off
2022-01-01,09:00,TV,0.8,Living Room,On
2022-01-01,10:00,TV,0.8,Living Room,On
...
```

Untuk contoh, kalian bisa menggunakan file yang telah disiapkan `data-series.csv`.

#### Penggunaan Model AI:

Model AI Tapas `tapas-base-finetuned-wtq` akan digunakan untuk memahami data tabel dan membuat prediksi tentang konsumsi energi masa depan. Model ini akan menerima file CSV sebagai input dan menghasilkan prediksi untuk total konsumsi energi hari berikutnya.

Buatlah interface untuk aplikasi ini, bisa berupa aplikasi CLI maupun Web Application. Silahkan dikembangkan sehingga mirip dengan chatbot dimana user bisa bertanya mengenai data-data yang ada di file input.

Silahkan menggunakan model AI lainnya dari Hugging Face Hub untuk membuat aplikasi ini lebih menarik, misal-nya dengan menambahkan model AI `openai-community/gpt2` agar bisa memberikan rekomendasi tentang alat apa yang bisa digunakan lebih sedikit untuk menghemat energi.

### Constraints

Function `CsvToSlice` dan `ConnectAIModel` sudah diberikan dan wajib kalian gunakan. Silahkan membuat function-function lain yang kalian perlukan.

### Test Case Examples

#### Test Case CsvToSlice

**Input**:

```txt
"Name,Age\nJohn,30\nDoe,40"
```

**Expected Output / Behavior**:

{
    "Name": ["John", "Doe"],
    "Age": ["30", "40"]
}

**Explanation**:

Fungsi CsvToSlice menerima string dari file CSV sebagai input dan mengembalikan `map` di mana `key`-nya adalah header kolom dan `value`nya adalah data untuk setiap kolom. Dalam hal ini, string CSV input memiliki dua kolom "Name" dan "Age", dan dua baris data "John, 30" dan "Doe, 40". Fungsi ini harus mengembalikan map dengan dua data. Data pertama harus memiliki key "Name" dan value ["John", "Doe"], dan key kedua harus memiliki key "Age" dan value ["30", "40"].

#### Test Case 2

**Input**:

```txt
Payload: {
    "table": {
        "Name": ["John", "Doe"],
        "Age": ["30", "40"]
    },
    "query": "What is the age of John?"
}
```

**Expected Output / Behavior**:

```txt
{
    "answer": "30",
    "coordinates": [[0, 1]],
    "cells": ["30"],
    "aggregator": ""
}
```

**Explanation**:

Fungsi ConnectAIModel menerima payload dan Huggingface Token sebagai input dan mengembalikan struktur Response. Payload adalah struktur yang berisi `Table` dan `Query`. `Tabel` adalah sebuah map di mana `key`-nya adalah header kolom dan `value`-nya adalah irisan yang berisi data untuk setiap kolom. `Query` adalah string yang mewakili pertanyaan tentang data di tabel. Dalam hal ini, querynya adalah "Berapa umur John?". Fungsi ini harus mengembalikan struktur Response dengan jawaban "30", koordinat [[0, 1]], sel ["30"], dan aggregator.

Happy Coding!