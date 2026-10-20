# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Sumber daya dan peminjaman dasar (B1-B3)
- Prompt: "buatkan kode Go net/http modular untuk B1 sampai B3 pinjam kampus, pisah store.go dan api.go"
- Dipakai: Struktur file `store.go`, `api.go`, dan `main.go` untuk pengelolaan sumber daya dan peminjaman.
- Diubah: Menambahkan field `paket_id` ke dalam struct `Peminjaman` agar siap mendukung fitur khas persona Panitia Acara.