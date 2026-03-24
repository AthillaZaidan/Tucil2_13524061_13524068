# OctoVox - Tucil 2 IF2211 Strategi Algoritma

OctoVox adalah program untuk melakukan voxelization model 3D berformat OBJ menggunakan struktur data octree berbasis Divide and Conquer. Program ini menyediakan visualisasi interaktif, pengaturan max-depth, statistik hasil voxelization, serta ekspor hasil voxelization ke file OBJ baru.

## Deskripsi Singkat Program

Fitur utama:

- Parsing file OBJ (vertex dan face segitiga).
- Perhitungan bounding box model.
- Pembangunan octree sampai kedalaman tertentu (max-depth).
- Uji intersection triangle-box untuk menentukan voxel aktif.
- Visualisasi voxel 3D menggunakan G3N engine.
- Ekspor hasil voxelization ke file OBJ output.
- Pengujian pada kasus normal dan edge case.

## Requirement

Minimum untuk menjalankan program:

- OS: Windows 10/11 64-bit.
- Go: 1.26.1 atau lebih baru.
- C/C++ toolchain untuk cgo (disarankan MinGW-w64/GCC tersedia di PATH).
- GPU/driver dengan dukungan OpenGL (minimum OpenGL 3.x).
- RAM minimal 4 GB (disarankan 8 GB untuk model besar).
- Ruang disk kosong minimal 500 MB.

Catatan:

- Saat build melalui Makefile, dependency DLL audio akan disalin otomatis ke folder `bin`.

## How to Run

### Opsi 1 - Menggunakan Makefile (disarankan)

1. Pastikan berada di root project.
2. Build program:

```bash
make build
```

3. Jalankan program:

```bash
make run
```

4. Bersihkan binary dan DLL hasil build:

```bash
make clean
```

### Opsi 2 - Langsung lewat Go

```bash
cd src
go run .
```

## Format .obj

Parser saat ini mendukung subset format OBJ berikut:

- Baris vertex: `v x y z`
- Baris face segitiga: `f i j k`

Aturan penting:

- Indeks face harus integer positif dan mulai dari 1.
- Face harus segitiga (tepat 3 indeks).
- Referensi indeks face tidak boleh melebihi jumlah vertex.
- Baris kosong, komentar (`# ...`), dan token lain akan diabaikan.

Contoh valid:

```obj
v 0 0 0
v 1 0 0
v 0 1 0
f 1 2 3
```

## Struktur Project

```text
.
|- Makefile
|- go.mod
|- README.md
|- bin/
|- obj/
|- src/
|  |- main.go
|  |- packages/
|     |- intersect/
|     |- octree/
|     |- parser/
|     |- viewer/
|- docs/
	|- main.tex
	|- sections/
	|- public/
```

## Author

- Muhammad Aufar Rizqi Kusuma (13524061)
- Athilla Zaidan Zidna Fann (13524068)

## How to Contribute

Kontribusi dipersilakan melalui alur berikut:

1. Fork repository.
2. Buat branch baru dari `main`.
3. Lakukan perubahan dengan commit yang jelas.
4. Jalankan pengujian/build lokal terlebih dahulu.
5. Buka Pull Request berisi:
	- ringkasan perubahan,
	- alasan perubahan,
	- dampak ke fitur/performa,
	- bukti uji (jika ada).

Panduan tambahan:

- Hindari commit file generated yang tidak perlu (binary, output besar, dll).
- Pastikan perubahan tidak merusak alur `make build` dan `make run`.
