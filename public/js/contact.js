function sendMail(){
    let name = document.getElementById('name').value;
    let email = document.getElementById('email').value;
    let phone = document.getElementById('phone').value;
    let subject = document.getElementById('subject').value;
    let message = document.getElementById('message').value;

    if (name == "" || email == "" || phone == "" || subject == "" || message == ""){
        alert("Semua Field Harus Diisi")
    }
    else {
    const defaultEmail = "rizkiramdhani2610@gmail.com";

    let mailTo = document.createElement('a')
    mailTo.href = `mailto:${defaultEmail}?subject=${subject}&body=Halo, Nama saya ${name}, ${message}, Kamu bisa hubungi Saya di ${phone}.`
    mailTo.target = "_blank"
    mailTo.click()
    }
}