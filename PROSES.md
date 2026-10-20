# PROSES

## 1. Persona dan masalahnya
Panitia acara sering kali membutuhkan aula, proyektor, dan sound system secara bersamaan untuk suatu kegiatan. Jika salah satu sumber daya tidak tersedia, acara tidak dapat berjalan sehingga panitia memerlukan mekanisme pemesanan paket serentak (semua atau tidak sama sekali) untuk menghindari pemesanan yang terpecah dan tidak sinkron.

## 2. Rancangan struktur data
- Struct `SumberDaya`: Menyimpan `id` (int), `nama` (string), `jenis` (string), dan `stok` (int).
- Struct `Peminjaman`: Menyimpan `id`, `sumber_daya_id`, `peminjam`, `tanggal`, `slot_mulai`, `slot_selesai`, `jumlah`, `status`, serta `paket_id` (int) khusus persona Panitia Acara.
- Penyimpanan menggunakan slice in-memory yang dilindungi oleh `sync.Mutex` untuk menjamin keamanan akses konkuren.

## 3. Satu keputusan desain
Memisahkan logika penyimpanan data (`store.go`) dan router/handler API (`api.go`) untuk menjaga keterbacaan kode (*separation of concerns*) sekaligus menggunakan satu mutex terpusat dalam objek store agar pengelolaan state aman terhadap permintaan bersamaan.

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala: Tanggal fiktif seperti 30 Februari diterima oleh parser tanggal standar.
- Cara menemukan: Pengujian unit dengan input tanggal tidak valid.
- Penyebab dan solusi: `time.Parse` melakukan normalisasi otomatis. Solusi dengan membandingkan kembali hasil format tanggal setelah di-parsing.

### Bug 2
- Gejala: Panjang karakter nama melebihi batas 50 dihitung berdasarkan byte, bukan jumlah rune Unicode.
- Cara menemukan: Validasi karakter khusus atau multibyte.
- Penyebab dan solusi: Menggunakan `len([]rune(s))` untuk menghitung jumlah karakter secara akurat.

## 5. Change request
(Akan diisi setelah menerima lembar perubahan dari dosen.)