const API_BASE_URL = 'http://localhost:8080/api';

// Helper to make authenticated requests
async function apiCall(endpoint, method = 'GET', body = null) {
    const headers = {
        'Content-Type': 'application/json'
    };

    const token = localStorage.getItem('jwt_token');
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const config = {
        method: method,
        headers: headers
    };

    if (body) {
        config.body = JSON.stringify(body);
    }

    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, config);
        
        let data;
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.includes("application/json")) {
            data = await response.json();
        } else {
            const text = await response.text();
            data = { message: text };
        }
        
        if (!response.ok) {
            // Clean up common gRPC errors
            let msg = data.message || 'API request failed';
            if (msg.includes("desc = ")) {
                msg = msg.split("desc = ")[1];
            }
            throw new Error(msg);
        }
        
        return data;
    } catch (error) {
        console.error('API Error:', error);
        throw error;
    }
}

// --- Auth API ---
async function apiRegister(username, email, password) {
    const data = await apiCall('/auth/register', 'POST', { username, email, password });
    localStorage.setItem('jwt_token', data.token);
    localStorage.setItem('user_id', data.user_id);
    return data;
}

async function apiLogin(email, password) {
    const data = await apiCall('/auth/login', 'POST', { email, password });
    localStorage.setItem('jwt_token', data.token);
    localStorage.setItem('user_id', data.user_id);
    return data;
}

// --- User API ---
async function apiGetProfile() {
    const userId = localStorage.getItem('user_id');
    return await apiCall(`/user/profile?user_id=${userId}`, 'GET');
}

// --- Gacha API ---
async function apiRollGacha(bannerId) {
    const userId = localStorage.getItem('user_id');
    return await apiCall('/gacha/roll', 'POST', { user_id: userId, banner_id: bannerId });
}

// Example usage to replace local storage saveState:
/*
async function syncStateToBackend() {
    // If you want to keep the current index.html mostly intact, 
    // you can call these api functions inside the existing click handlers 
    // instead of rewriting the entire game loop.
}
*/
