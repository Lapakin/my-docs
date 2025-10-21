// Auth Form Handlers
document.addEventListener('DOMContentLoaded', function() {
    // Login Form Handler
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();

            const login = document.getElementById('login').value;
            const password = document.getElementById('password').value;
            const errorMessage = document.getElementById('errorMessage');

            try {
                const response = await fetch('/api/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ login, password }),
                });

                const data = await response.json();

                if (response.ok) {
                    window.location.href = '/dashboard';
                } else {
                    errorMessage.textContent = data.error || 'Login failed';
                    errorMessage.classList.add('show');
                }
            } catch (error) {
                errorMessage.textContent = 'Network error. Please try again.';
                errorMessage.classList.add('show');
            }
        });
    }

    // Register Form Handler
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();

            const username = document.getElementById('fullName').value;
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const errorMessage = document.getElementById('errorMessage');

            try {
                const response = await fetch('/api/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ username, email, password }),
                });

                const data = await response.json();

                if (response.ok) {
                    window.location.href = '/login?registered=true';
                } else {
                    errorMessage.textContent = data.error || 'Registration failed';
                    errorMessage.classList.add('show');
                }
            } catch (error) {
                errorMessage.textContent = 'Network error. Please try again.';
                errorMessage.classList.add('show');
            }
        });
    }
});

