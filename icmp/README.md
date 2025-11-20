# 🎯 Go Ping tool (uping)

Bu dastur Go tilida yozilgan bo'lib, standart tarmoq vositasi (`ping`) bilan bir xil vazifani bajaradi. U ICMP so'rovlarini IPv4 va IPv6 protokollari orqali yuborish va javoblarni tahlil qilish imkonini beradi.

## ✨ Xususiyatlari

* ✅ **Universal Protokol:** IPv4 va IPv6 manzilini avtomatik aniqlash va ishlatish.
* ✅ **TTL/Hop Limit:** Qabul qilingan paketdagi TTL (IPv4) yoki Hop Limit (IPv6) qiymatini chiqarish.
* ✅ **CLI Flaglar:** Standart `ping` kabi buyruqlar qatori opsiyalarini qo'llab-quvvatlash.
* ✅ **Xavfsiz Ishlash:** Linuxda **`sudo`** huquqisiz ishlash imkoniyati (`setcap` orqali).

## ⚙️ Talablar

* **Go tili** (1.18 yoki undan yuqori)
* Linux tizimida `setcap` buyrug‘i (odatda tizimda mavjud).

## 🚀 O'rnatish va Sozlash

Dastur tarmoqning "xom (raw) socket" darajasida ishlaydi, shuning uchun Linuxda `sudo` huquqisiz ishlashi uchun maxsus sozlash talab qilinadi.

### 1. Kodni Kompilyatsiya qilish

Kod joylashgan papkada (masalan, `/data/projects/go`) turib, quyidagi buyruqni bering. Bu `uping` nomli bajariladigan faylni yaratadi:

```bash
go build -o uping main.go
```
###  Xavfsizlik Huquqini Berish

```bash
sudo setcap cap_net_raw+ep ./uping
```

###  Foydalanish

```bash
./uping [opsiyalar] <manzil_yoki_domen>
```