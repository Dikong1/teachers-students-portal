// This is a placeholder function. In a real app, you would implement the actual logic.
function chooseToGift(childId) {
    console.log("Gift chosen for child with ID:", childId);
    // Here you would likely make an AJAX call to your backend
}

// Add event listeners to gift buttons
document.addEventListener("DOMContentLoaded", function() {
    var giftButtons = document.querySelectorAll(".gift-button");
    giftButtons.forEach(function(button) {
        button.addEventListener("click", function() {
            var childId = this.closest('.child-wish').getAttribute('data-child-id');
            chooseToGift(childId);
        });
    });
});
// Example JavaScript for client-side interactivity
document.addEventListener("DOMContentLoaded", function() {
    var updateButton = document.querySelector("button[type='submit']");
    updateButton.addEventListener("click", function(event) {
        // You can add validation or other interactivity here
        console.log("Wishes updated!");
    });
});

document.addEventListener('DOMContentLoaded', function() {
    var form = document.querySelector('form');
    form.addEventListener('submit', function(event) {
        event.preventDefault(); // Prevent the default form submission

        var wishes = document.getElementById('wishes').value;
        var data = {
            wishes: wishes
        };

        var jsonData = JSON.stringify(data);

        console.log(jsonData); // For debugging

        // Here you can send jsonData to your server
        // Example using fetch API
        fetch('/update-wishes', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: jsonData
        })
        .then(response => response.json())
        .then(data => {
            console.log('Success:', data);
            // Handle success here (e.g., showing a success message)
        })
        .catch((error) => {
            console.error('Error:', error);
            // Handle errors here
        });
    });
});

function submitRegistration() {
    var formData = {
        name: document.getElementById("firstName").value,
        lastName: document.getElementById("lastName").value,
        email: document.getElementById("email").value,
        phone: document.getElementById("phone").value,
        password: document.getElementById("password").value,
        child: null
    };

    fetch('/submit-volunteer-registration', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.text();
        })
        .then(data => {
            console.log(data);
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });

    // Предотвращаем отправку формы
    return false;
}

document.addEventListener("DOMContentLoaded", function () {
    var form = document.getElementById("loginForm");

    form.addEventListener("submit", function (event) {
        event.preventDefault(); // Предотвратить стандартное поведение формы

        // Получить данные формы
        var formData = {
            phone: document.getElementById("phone").value,
            password: document.getElementById("password").value,
        };

        // Отправить POST-запрос на сервер
        fetch('/submit-volunteer-login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams(formData),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.text();
            })
            .then(data => {
                console.log(data);
                // Обработка успешного ответа, если нужно
            })
            .catch(error => {
                console.error('There has been a problem with your fetch operation:', error);
                // Обработка ошибок
            });
    });
});

