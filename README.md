[//]: # (1 bu pryect tolaqonli taskdagi shartlar bajarildi qoshimcha shartlari ham bajarildi)
[//]: # (2 bu qoshimcha tarizda ratelimter ham qoshildi)
[//]: # (3 doctor,appointment,auth,user-profile qisimlari boldi)
[//]: # (RATELIMTER -bu bir nechta request lar kelishida serverga huchumlar qanchadir soniyada 
api lar request yuborilsa avto block qilib qoyadi)

[//]: # (4 task da user appointment biron qilganda unga sms yoki email orqali yuborish kerak deyailgan
YECHIM: men schedule dan foydalandim bunda vaqtini qisqa qoydim yani 5-minut yani appointment qilganda user
notification va ochiriladi lekin send_at =null bolip turadi va schudel shu vaqtdan hisoblab di
5- minutda keyin hozirgi vaqtdan NOW>=apointmentTime ni qidirib olib keladi va email sms yuboradi
va notification toliq tugallaydi sended_at,status,message larni to'ldiradi
)

[//]: # (sms tarizda yuboriladigan qilsa boladi lekin tekin yoq bolgani sabali email ni tanladim
 sms bilan yuborish qilsa boladi unda eskiz.uz yoki boshqa saytlar dan tel sms code yuboradigan 
 lar sotib olish kerak boladi)

[//]: # (testing repository larga testing yozilgan  testing yana router larga yozmoqchi edim vaqti yetmadi!!!!)

[//]: # (Postgres,pgxPool dan foydalanilgan bu dan foydalanishim sabab ozi postgres ni yopadi)

[//]: # (sqlc dan foydalanganman)

[//]: # (pdocker ,docker-compose yani docker ham toliq yozilgan)
[//]: # (ci-cd deploy ham toliq yozilgan)

[//]: # (pnpm-lock-yam)