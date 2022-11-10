console.log('Hello World!')

document.getElementById('form-submit-btn').addEventListener('click', function(e) {
    e.preventDefault()
    if (isValidInputs()) {
        console.log('Submit!')
        submit(function() {
            clearForm()
        })
    } else {
        alert('Invalid Input')
    }
})

function isValidInputs() {
    if (isValidFirstName() &&
        isValidLastName() &&
        isValidPhone() &&
        isValidEmail()) {
        return true
    }
    return false
}

function isValidFirstName() {
    let firstName = document.getElementById('firstName').value
    if (firstName.length > 0 && firstName.length <= 32) {
        return true
    }
    return false
}

function isValidLastName() {
    let firstName = document.getElementById('lastName').value
    if (firstName.length > 0 && firstName.length <= 32) {
        return true
    }
    return false
}

function isValidEmail() {
    let regexp = /(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])/
    let email = document.getElementById('email').value
    if (email.match(regexp)) {
        return true
    }
    return false
}

function isValidPhone() {
    let regexp = /^([0-9]{3})([0-9]{3})([0-9]{4})$/
    let phoneNumber = document.getElementById('phoneNumber').value
    if (phoneNumber.match(regexp)) {
        return true
    }
    return false
}

function getFormJSON() {
    var firstName = document.getElementById('firstName').value
    var lastName = document.getElementById('lastName').value
    var phoneNumber = document.getElementById('phoneNumber').value
    var email = document.getElementById('email').value
    return JSON.stringify({
        firstName: firstName,
        lastName: lastName,
        phoneNumber: phoneNumber,
        email: email
    })
}

function clearForm() {
    document.getElementById('firstName').value = ''
    document.getElementById('lastName').value = ''
    document.getElementById('phoneNumber').value = ''
    document.getElementById('email').value = ''
}

function submit(done) {
    let xhr = new XMLHttpRequest()
    xhr.open("POST", "/register", true)

    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            // Print received data from server
            console.log(this.responseText)
            done()
        }
    };

    var data = getFormJSON()
 
    // Sending data with the request
    xhr.send(data);
}