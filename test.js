 function a() {
     process.nextTick(() => {
         for (let i = 0; i <= 1000000000; i++) {
             i + 1;
         }
         console.log(1);
     });
     console.log(2);
     return 3;
 }

 function b() {
     console.log(a());
     return 8;
 }
 console.log(b());
 console.log(9);